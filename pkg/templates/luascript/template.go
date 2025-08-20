package luascript

import (
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/luaengine"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/luascript/internal/params"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Lua Script",
		Key:         "lua-script",
		Description: "Run custom Lua script checks against Kubernetes objects",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.Any},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			script, err := loadScript(p)
			if err != nil {
				return nil, err
			}

			timeout := time.Duration(p.Timeout) * time.Second
			if p.Timeout == 0 {
				timeout = 5 * time.Second // default timeout
			}

			engine := luaengine.New(script, timeout)

			return func(lintCtx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				diagnostics, err := engine.ExecuteCheck(lintCtx, object)
				if err != nil {
					// Return the error as a diagnostic
					return []diagnostic.Diagnostic{{
						Message: "Lua script error: " + err.Error(),
					}}
				}
				return diagnostics
			}, nil
		}),
	})
}

// loadScript loads the Lua script from file or inline content
func loadScript(p params.Params) (string, error) {
	if p.Inline != "" && p.Script != "" {
		return "", errors.New("cannot specify both 'script' and 'inline' parameters")
	}

	if p.Inline != "" {
		return p.Inline, nil
	}

	if p.Script == "" {
		return "", errors.New("must specify either 'script' or 'inline' parameter")
	}

	// Load script from file
	scriptPath := p.Script
	if !filepath.IsAbs(scriptPath) {
		// Make relative paths relative to current working directory
		wd, err := os.Getwd()
		if err != nil {
			return "", errors.Wrap(err, "getting working directory")
		}
		scriptPath = filepath.Join(wd, scriptPath)
	}

	content, err := os.ReadFile(scriptPath) // #nosec G304 - script path comes from user configuration
	if err != nil {
		return "", errors.Wrapf(err, "reading script file %s", scriptPath)
	}

	return string(content), nil
}
