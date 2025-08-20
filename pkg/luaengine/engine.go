package luaengine

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	lua "github.com/yuin/gopher-lua"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/k8sutil"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
)

const (
	defaultTimeout = 5 * time.Second
	maxTimeout     = 30 * time.Second
)

// Engine represents a Lua script execution engine for kube-linter checks
type Engine struct {
	script  string
	timeout time.Duration
}

// New creates a new Lua engine with the given script content and timeout
func New(script string, timeout time.Duration) *Engine {
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	if timeout > maxTimeout {
		timeout = maxTimeout
	}

	return &Engine{
		script:  script,
		timeout: timeout,
	}
}

// ExecuteCheck runs the Lua script against the given object and returns diagnostics
func (e *Engine) ExecuteCheck(lintCtx lintcontext.LintContext, object lintcontext.Object) ([]diagnostic.Diagnostic, error) {
	L := lua.NewState()
	defer L.Close()

	// Register extract functions
	e.registerExtractAPI(L, object.K8sObject)

	// Register diagnostic creation function
	e.registerDiagnosticAPI(L)

	// Execute the script with timeout
	done := make(chan error, 1)
	var result []diagnostic.Diagnostic

	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- errors.Errorf("lua script panicked: %v", r)
				return
			}
		}()

		if err := L.DoString(e.script); err != nil {
			done <- errors.Wrap(err, "executing lua script")
			return
		}

		// Call the check function
		if err := L.CallByParam(lua.P{
			Fn:      L.GetGlobal("check"),
			NRet:    1,
			Protect: true,
		}); err != nil {
			done <- errors.Wrap(err, "calling check function")
			return
		}

		// Get the result
		ret := L.Get(-1)
		L.Pop(1)

		var err error
		result, err = e.parseDiagnostics(ret)
		done <- err
	}()

	select {
	case err := <-done:
		return result, err
	case <-time.After(e.timeout):
		return nil, errors.New("lua script execution timeout")
	}
}

// registerExtractAPI registers K8s object extraction functions in Lua
func (e *Engine) registerExtractAPI(l *lua.LState, obj k8sutil.Object) {
	extractTable := l.NewTable()

	// Register podSpec extraction
	extractTable.RawSetString("podSpec", l.NewFunction(func(state *lua.LState) int {
		podSpec, found := extract.PodSpec(obj)
		if !found {
			state.Push(lua.LNil)
			return 1
		}

		// Convert to Lua table
		luaValue, err := e.goToLua(state, podSpec)
		if err != nil {
			state.RaiseError("failed to convert podSpec: %v", err)
			return 0
		}
		state.Push(luaValue)
		return 1
	}))

	// Register labels extraction
	extractTable.RawSetString("labels", l.NewFunction(func(state *lua.LState) int {
		labels := extract.Labels(obj)
		luaValue, err := e.goToLua(state, labels)
		if err != nil {
			state.RaiseError("failed to convert labels: %v", err)
			return 0
		}
		state.Push(luaValue)
		return 1
	}))

	// Register annotations extraction
	extractTable.RawSetString("annotations", l.NewFunction(func(state *lua.LState) int {
		annotations := extract.Annotations(obj)
		luaValue, err := e.goToLua(state, annotations)
		if err != nil {
			state.RaiseError("failed to convert annotations: %v", err)
			return 0
		}
		state.Push(luaValue)
		return 1
	}))

	// Register GVK extraction
	extractTable.RawSetString("gvk", l.NewFunction(func(state *lua.LState) int {
		gvk := extract.GVK(obj)
		luaValue, err := e.goToLua(state, gvk)
		if err != nil {
			state.RaiseError("failed to convert gvk: %v", err)
			return 0
		}
		state.Push(luaValue)
		return 1
	}))

	l.SetGlobal("extract", extractTable)
}

// registerDiagnosticAPI registers diagnostic creation functions
func (e *Engine) registerDiagnosticAPI(l *lua.LState) {
	// Register diagnostic creation function
	l.SetGlobal("diagnostic", l.NewFunction(func(state *lua.LState) int {
		message := state.CheckString(1)

		table := state.NewTable()
		table.RawSetString("message", lua.LString(message))

		state.Push(table)
		return 1
	}))
}

// parseDiagnostics converts Lua return value to Go diagnostics
func (e *Engine) parseDiagnostics(lv lua.LValue) ([]diagnostic.Diagnostic, error) {
	if lv.Type() == lua.LTNil {
		return nil, nil
	}

	if lv.Type() != lua.LTTable {
		return nil, errors.New("check function must return a table (array) of diagnostics")
	}

	table := lv.(*lua.LTable)
	var diagnostics []diagnostic.Diagnostic

	table.ForEach(func(key, value lua.LValue) {
		if value.Type() != lua.LTTable {
			return
		}

		diagTable := value.(*lua.LTable)
		messageValue := diagTable.RawGetString("message")
		if messageValue.Type() != lua.LTString {
			return
		}

		diagnostics = append(diagnostics, diagnostic.Diagnostic{
			Message: string(messageValue.(lua.LString)),
		})
	})

	return diagnostics, nil
}

// goToLua converts Go values to Lua values using JSON marshaling
func (e *Engine) goToLua(l *lua.LState, value interface{}) (lua.LValue, error) {
	// Convert to JSON first, then to Lua
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	var jsonValue interface{}
	if err := json.Unmarshal(jsonBytes, &jsonValue); err != nil {
		return nil, err
	}

	return e.jsonToLua(l, jsonValue), nil
}

// jsonToLua converts JSON values to Lua values recursively
func (e *Engine) jsonToLua(l *lua.LState, value interface{}) lua.LValue {
	switch v := value.(type) {
	case nil:
		return lua.LNil
	case bool:
		return lua.LBool(v)
	case float64:
		return lua.LNumber(v)
	case string:
		return lua.LString(v)
	case []interface{}:
		table := l.NewTable()
		for i, item := range v {
			table.RawSetInt(i+1, e.jsonToLua(l, item))
		}
		return table
	case map[string]interface{}:
		table := l.NewTable()
		for key, val := range v {
			table.RawSetString(key, e.jsonToLua(l, val))
		}
		return table
	default:
		return lua.LString("")
	}
}
