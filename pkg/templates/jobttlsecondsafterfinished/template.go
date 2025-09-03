package jobttlsecondsafterfinished

import (
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/jobttlsecondsafterfinished/internal/params"
)

const templateKey = "job-ttl-seconds-after-finished"

func init() {
	templates.Register(check.Template{
		HumanName:   "ttlSecondsAfterFinished impact for standalone and managed Job objects",
		Key:         templateKey,
		Description: "Flag standalone Job objects not setting ttlSecondsAfterFinished. Flag CronJob objects setting ttlSecondsAfterFinished",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.JobLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return func(lintCtx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				jobSpec, kind, ok := extract.JobSpec(object.K8sObject)
				if !ok {
					return nil
				}
				switch kind {
				case "Job":
					if jobSpec.TTLSecondsAfterFinished == nil {
						return []diagnostic.Diagnostic{{Message: "Standalone Job does not specify ttlSecondsAfterFinished"}}
					}
				case "CronJob":
					if jobSpec.TTLSecondsAfterFinished != nil {
						return []diagnostic.Diagnostic{{Message: "Managed Job specifies ttlSecondsAfterFinished which might conflict with successfulJobsHistoryLimit and failedJobsHistoryLimit from CronJob that have default values. Final behaviour is determined by the strictest parameter, and therefore, setting ttlSecondsAfterFinished at the job level can result with unexpected behaviour with regard to finished jobs removal"}}
					}
				}
				return nil
			}, nil
		}),
	})
}
