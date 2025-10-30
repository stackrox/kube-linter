package kubeconform

import (
	"fmt"
	"os"

	"github.com/yannh/kubeconform/pkg/resource"
	"github.com/yannh/kubeconform/pkg/validator"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/kubeconform/internal/params"
)

const (
	templateKey = "kubeconform"
)

func init() {
	templates.Register(check.Template{
		HumanName:   templateKey,
		Key:         templateKey,
		Description: "Flag objects that does not match schema using https://github.com/yannh/kubeconform",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.Any},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate:            params.WrapInstantiateFunc(validate),
	})
}

func validate(p params.Params) (check.Func, error) {
	// Create cache directory if it doesn't exist
	if p.Cache != "" {
		if err := os.MkdirAll(p.Cache, 0o750); err != nil {
			return nil, fmt.Errorf("creating cache directory %s: %w", p.Cache, err)
		}
	}

	v, err := validator.New(p.SchemaLocations, validator.Opts{
		Cache:                p.Cache,
		SkipKinds:            sliceToMap(p.SkipKinds),
		RejectKinds:          sliceToMap(p.RejectKinds),
		KubernetesVersion:    p.KubernetesVersion,
		Strict:               p.Strict,
		IgnoreMissingSchemas: p.IgnoreMissingSchemas,
	})
	if err != nil {
		return nil, fmt.Errorf("creating kubeconform validator: %w", err)
	}
	return func(ctx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
		res := v.ValidateResource(resource.Resource{
			Path:  object.Metadata.FilePath,
			Bytes: object.Metadata.Raw,
		})
		if res.Status == validator.Invalid {
			return []diagnostic.Diagnostic{
				{Message: fmt.Sprintf("resource is not valid: %v", res.Err)},
			}
		}
		if res.Status == validator.Error {
			return []diagnostic.Diagnostic{
				{Message: fmt.Sprintf("error while processing resource: %v", res.Err)},
			}
		}
		return nil
	}, nil
}

func sliceToMap(keys []string) map[string]struct{} {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[k] = struct{}{}
	}
	return m
}
