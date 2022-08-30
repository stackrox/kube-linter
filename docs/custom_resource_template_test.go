package certlifetime

import (
	"fmt"
	"time"

	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/k8sutil"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Define a constant that can be used to reference the objectKind in the check
const (
	Certificate = cmv1.CertificateKind
)

// Define the GVK of the object that we want to match
var (
	certificateGVK = cmv1.SchemeGroupVersion.WithKind(Certificate)
)

func init() {
	// Register our matcher and objectkind with the global matcher registry.
	// This function can contain any arbitrary logic to match objects we want to check with this lint check
	objectkinds.RegisterObjectKind(Certificate, objectkinds.MatcherFunc(func(gvk schema.GroupVersionKind) bool {
		return gvk == certificateGVK
	}))

	templates.Register(CertLifetime)
}

var CertLifetime = check.Template{
	HumanName:   "Certificate Lifetime",
	Key:         "cert-lifetime",
	Description: "Flag certificates lasting longer than 1 year",
	SupportedObjectKinds: config.ObjectKindsDesc{
		ObjectKinds: []string{Certificate},
	},
	Parameters:             nil,
	ParseAndValidateParams: func(params map[string]interface{}) (interface{}, error) { return nil, nil },
	Instantiate: func(_ interface{}) (check.Func, error) {
		return func(_ lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
			cert, ok := object.K8sObject.(*cmv1.Certificate)
			if !ok {
				return []diagnostic.Diagnostic{{Message: "Invalid certificate"}}
			}

			if cert.Spec.Duration.Duration.Hours() > 8760 {
				return []diagnostic.Diagnostic{{Message: "Certificates with lifetimes longer than one year are not allowed"}}
			}
			return nil
		}, nil
	},
}

func ExampleCertLifetime() {
	dur, _ := time.ParseDuration("9001h")
	cert := &cmv1.Certificate{
		ObjectMeta: v1.ObjectMeta{
			Name: "certificate-long",
		},
		Spec: cmv1.CertificateSpec{
			Duration: &v1.Duration{Duration: dur},
		},
	}

	checker, _ := CertLifetime.Instantiate(nil)

	asK8sObj, _ := runtime.Object(cert).(k8sutil.Object)

	result := checker(nil, lintcontext.Object{K8sObject: asK8sObj})

	fmt.Println(result[0].Message)
	// Output: Certificates with lifetimes longer than one year are not allowed
}
