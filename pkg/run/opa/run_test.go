package opa

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/open-policy-agent/opa/v1/ast"
	"github.com/open-policy-agent/opa/v1/rego"
	"github.com/open-policy-agent/opa/v1/util"
	"github.com/stretchr/testify/assert"
)

var jsonFile = `{
	   "kind": "AdmissionReview",
	   "request": {
	       "kind": {
	           "kind": "Pod",
	           "version": "v1"
	       },
	       "object": {
	           "metadata": {
	               "name": "myapp",
	               "labels": {
	                   "costcenter": "fakecode"
	               }
	           },
	           "spec": {
	               "containers": [
	                   {
	                       "image": "nginx",
	                       "name": "nginx-frontend"
	                   },
	                   {
	                       "image": "mysql",
	                       "name": "mysql-backend"
	                   }
	               ]
	           }
	       }
	   }
	}`

func TestName(t *testing.T) {
	ctx := context.TODO()
	input := util.MustUnmarshalJSON([]byte(jsonFile))

	// Load policies and data from a folder manually
	//loaded, err := loader.AllRegos([]string{"x:policies"})
	//assert.NoError(t, err)

	modules := []func(*rego.Rego){
		rego.Query("data.abc.deny[_]"),
		rego.SetRegoVersion(ast.RegoV1),
		rego.Load([]string{"policies"}, nil),

		rego.Input(input),
	}
	//for _, v := range loaded.Modules {
	//	spew.Dump(v.Name)
	//	modules = append(modules, rego.ParsedModule(v.Parsed))
	//}

	eval := rego.New(
		modules...,
	)

	rs, err := eval.Eval(ctx)
	assert.NoError(t, err)
	spew.Dump(rs)

}
