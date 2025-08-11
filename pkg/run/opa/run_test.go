package opa

import (
	"context"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/open-policy-agent/opa/v1/ast"
	"github.com/open-policy-agent/opa/v1/rego"
	"github.com/open-policy-agent/opa/v1/util"
	"github.com/stretchr/testify/assert"
)

var jsonFile = `{
  "object": {
    "apiVersion": "apps/v1",
    "kind": "Deployment",
    "metadata": {
      "name": "app"
    },
    "spec": {
      "template": {
        "spec": {
          "containers": [
            {
              "name": "app",
              "image": "app:latest"
            }
          ]
        }
      }
    }
  },
  "objects": [],
  "param": {
    "latesttag": {
      "blockList": [
        ".*:(latest)$",
        "^[^:]*$",
        "(.*/[^:]+)$"
      ]
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
		rego.Query(`data.kubelinter.template.latesttag.deny`),
		rego.SetRegoVersion(ast.RegoV1),
		rego.Load([]string{"/home/janisz/go/src/github.com/stackrox/kube-linter/kubelinter"}, nil),
		rego.Input(input),
	}
	eval := rego.New(
		modules...,
	)

	rs, err := eval.Eval(ctx)
	assert.NoError(t, err)
	spew.Dump(rs)

	messages, err := result(rs)
	assert.NoError(t, err)

	assert.Len(t, messages, 1)
	spew.Dump(messages)

}
