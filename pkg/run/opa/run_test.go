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
		rego.Query(`
{
  "xyz": [x | x := data.abc.xyz.deny[_]],
  "abc": [a | a := data.abc.abc.deny[_]]
}
`),
		rego.SetRegoVersion(ast.RegoV1),
		rego.Load([]string{"policies"}, nil),

		rego.Input(input),
	}
	eval := rego.New(
		modules...,
	)

	rs, err := eval.Eval(ctx)
	assert.NoError(t, err)
	spew.Dump(rs)

	messages := []string{}
	for _, result := range rs {
		for _, r := range result.Expressions {
			msgs, ok := r.Value.(map[string]interface{})
			assert.True(t, ok)
			for k, v := range msgs {
				println(k)
				strs, ok := v.([]interface{})
				assert.True(t, ok)
				for _, str := range strs {
					messages = append(messages, str.(string))
				}
			}
		}
	}

	assert.Len(t, messages, 4)
	spew.Dump(messages)

}
