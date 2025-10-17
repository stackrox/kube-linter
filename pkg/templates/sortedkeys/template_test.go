package sortedkeys

import (
	"testing"

	"github.com/stretchr/testify/suite"
	appsV1 "k8s.io/api/apps/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/sortedkeys/internal/params"
)

func TestSortedKeys(t *testing.T) {
	suite.Run(t, new(SortedKeysTestSuite))
}

type SortedKeysTestSuite struct {
	templates.TemplateTestSuite
	ctx *mocks.MockLintContext
}

func (s *SortedKeysTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *SortedKeysTestSuite) TestSortedKeys() {
	const deploymentName = "sorted-deployment"

	// Create a deployment with properly sorted YAML keys
	sortedYAML := []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: sorted-deployment
spec:
  replicas: 3
`)

	deployment := &appsV1.Deployment{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name: deploymentName,
		},
	}

	s.ctx.AddObjectWithRaw(deploymentName, deployment, sortedYAML)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				Recursive: true,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: nil, // No diagnostics expected for sorted keys
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *SortedKeysTestSuite) TestUnsortedTopLevelKeys() {
	const deploymentName = "unsorted-deployment"

	// Create a deployment with unsorted YAML keys (kind comes after spec)
	unsortedYAML := []byte(`apiVersion: apps/v1
metadata:
  name: unsorted-deployment
spec:
  replicas: 1
kind: Deployment
`)

	deployment := &appsV1.Deployment{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name: deploymentName,
		},
	}

	s.ctx.AddObjectWithRaw(deploymentName, deployment, unsortedYAML)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				Recursive: true,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: {
					{Message: "Keys are not sorted at root. Expected order: [apiVersion, kind, metadata, spec], got: [apiVersion, metadata, spec, kind]"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *SortedKeysTestSuite) TestUnsortedNestedKeys() {
	const deploymentName = "unsorted-nested"

	// Create a deployment with unsorted nested keys
	unsortedNestedYAML := []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: unsorted-nested
spec:
  template:
    spec:
      containers:
        - name: app
          image: myapp:latest
    metadata:
      labels:
        app: myapp
`)

	deployment := &appsV1.Deployment{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name: deploymentName,
		},
	}

	s.ctx.AddObjectWithRaw(deploymentName, deployment, unsortedNestedYAML)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				Recursive: true,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: {
					{Message: "Keys are not sorted at spec.template. Expected order: [metadata, spec], got: [spec, metadata]"},
					{Message: "Keys are not sorted at spec.template.spec.containers[0]. Expected order: [image, name], got: [name, image]"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *SortedKeysTestSuite) TestNonRecursiveCheck() {
	const deploymentName = "non-recursive"

	// Create a deployment with unsorted nested keys but recursive=false
	unsortedNestedYAML := []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: non-recursive
spec:
  template:
    spec:
      containers:
        - name: app
    metadata:
      labels:
        app: myapp
`)

	deployment := &appsV1.Deployment{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name: deploymentName,
		},
	}

	s.ctx.AddObjectWithRaw(deploymentName, deployment, unsortedNestedYAML)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				Recursive: false,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: nil, // Top-level keys are sorted, nested unsorted keys ignored
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *SortedKeysTestSuite) TestConfigMapUnsortedData() {
	const configMapName = "unsorted-configmap"

	// ConfigMap with unsorted data keys
	unsortedDataYAML := []byte(`apiVersion: v1
data:
  zebra: "z"
  apple: "a"
kind: ConfigMap
metadata:
  name: unsorted-configmap
`)

	// Note: We need to use k8s types, but for simplicity in this test
	// we'll use a deployment as a stand-in since the check works on raw YAML
	deployment := &appsV1.Deployment{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name: configMapName,
		},
	}

	s.ctx.AddObjectWithRaw(configMapName, deployment, unsortedDataYAML)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				Recursive: true,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				configMapName: {
					// Data keys are unsorted (zebra before apple)
					{Message: "Keys are not sorted at data. Expected order: [apple, zebra], got: [zebra, apple]"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *SortedKeysTestSuite) TestDeeplyNestedStructure() {
	const deploymentName = "deeply-nested"

	// Test deeply nested structure (4 levels) with unsorted keys at level 3
	deeplyNestedYAML := []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    kubernetes.io/change-cause: "updated image"
    app.kubernetes.io/version: "1.0.0"
  labels:
    app: myapp
  name: deeply-nested
spec:
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
        - env:
            - name: LOG_LEVEL
              value: debug
            - name: APP_MODE
              value: production
          image: myapp:latest
          name: main
          resources:
            requests:
              memory: "64Mi"
              cpu: "100m"
            limits:
              memory: "128Mi"
              cpu: "500m"
`)

	deployment := &appsV1.Deployment{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name: deploymentName,
		},
	}

	s.ctx.AddObjectWithRaw(deploymentName, deployment, deeplyNestedYAML)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				Recursive: true,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: {
					// annotations keys unsorted
					{Message: "Keys are not sorted at metadata.annotations. Expected order: [app.kubernetes.io/version, kubernetes.io/change-cause], got: [kubernetes.io/change-cause, app.kubernetes.io/version]"},
					// resources keys unsorted (requests before limits)
					{Message: "Keys are not sorted at spec.template.spec.containers[0].resources. Expected order: [limits, requests], got: [requests, limits]"},
					// resources.requests keys unsorted
					{Message: "Keys are not sorted at spec.template.spec.containers[0].resources.requests. Expected order: [cpu, memory], got: [memory, cpu]"},
					// resources.limits keys unsorted
					{Message: "Keys are not sorted at spec.template.spec.containers[0].resources.limits. Expected order: [cpu, memory], got: [memory, cpu]"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *SortedKeysTestSuite) TestMultipleContainersInArray() {
	const deploymentName = "multi-container"

	// Test multiple containers with various unsorted patterns
	multiContainerYAML := []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: multi-container
spec:
  template:
    metadata:
      labels:
        app: multi
    spec:
      containers:
        - name: nginx
          image: nginx:latest
          ports:
            - protocol: TCP
              containerPort: 80
              name: http
        - name: sidecar
          image: sidecar:latest
          volumeMounts:
            - name: data
              mountPath: /data
        - image: logger:latest
          name: logger
      volumes:
        - name: data
          emptyDir: {}
`)

	deployment := &appsV1.Deployment{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name: deploymentName,
		},
	}

	s.ctx.AddObjectWithRaw(deploymentName, deployment, multiContainerYAML)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				Recursive: true,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: {
					// First container keys unsorted (name before image)
					{Message: "Keys are not sorted at spec.template.spec.containers[0]. Expected order: [image, name, ports], got: [name, image, ports]"},
					// First container ports unsorted
					{Message: "Keys are not sorted at spec.template.spec.containers[0].ports[0]. Expected order: [containerPort, name, protocol], got: [protocol, containerPort, name]"},
					// Second container keys unsorted (name before image)
					{Message: "Keys are not sorted at spec.template.spec.containers[1]. Expected order: [image, name, volumeMounts], got: [name, image, volumeMounts]"},
					// Second container volumeMounts unsorted
					{Message: "Keys are not sorted at spec.template.spec.containers[1].volumeMounts[0]. Expected order: [mountPath, name], got: [name, mountPath]"},
					// volumes unsorted (name before emptyDir)
					{Message: "Keys are not sorted at spec.template.spec.volumes[0]. Expected order: [emptyDir, name], got: [name, emptyDir]"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *SortedKeysTestSuite) TestEdgeCaseEmptyAndSingleKey() {
	const deploymentName = "edge-cases"

	// Test edge cases: empty objects and single key objects
	edgeCaseYAML := []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  annotations: {}
  labels:
    single: value
  name: edge-cases
spec:
  replicas: 1
`)

	deployment := &appsV1.Deployment{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name: deploymentName,
		},
	}

	s.ctx.AddObjectWithRaw(deploymentName, deployment, edgeCaseYAML)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				Recursive: true,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: nil, // No errors - empty and single key objects are fine
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *SortedKeysTestSuite) TestNumericAndSpecialCharKeys() {
	const configMapName = "numeric-keys"

	// Test numeric keys and special characters - should sort lexicographically
	numericKeysYAML := []byte(`apiVersion: v1
data:
  "10-config": value10
  "2-config": value2
  "1-config": value1
  config-z: valuez
  config-a: valuea
kind: ConfigMap
metadata:
  name: numeric-keys
`)

	deployment := &appsV1.Deployment{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name: configMapName,
		},
	}

	s.ctx.AddObjectWithRaw(configMapName, deployment, numericKeysYAML)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				Recursive: true,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				configMapName: {
					// Numeric keys should be sorted lexicographically: "1-config", "10-config", "2-config", "config-a", "config-z"
					{Message: "Keys are not sorted at data. Expected order: [1-config, 10-config, 2-config, config-a, config-z], got: [10-config, 2-config, 1-config, config-z, config-a]"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *SortedKeysTestSuite) TestComplexServiceWithSelectors() {
	const serviceName = "complex-service"

	// Real-world Service manifest with multiple sections
	serviceYAML := []byte(`apiVersion: v1
kind: Service
metadata:
  annotations:
    service.beta.kubernetes.io/aws-load-balancer-type: nlb
  labels:
    app: web
  name: complex-service
spec:
  ports:
    - port: 443
      targetPort: 8443
      protocol: TCP
      name: https
    - name: http
      port: 80
      protocol: TCP
      targetPort: 8080
  selector:
    app: web
  type: LoadBalancer
`)

	deployment := &appsV1.Deployment{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name: serviceName,
		},
	}

	s.ctx.AddObjectWithRaw(serviceName, deployment, serviceYAML)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				Recursive: true,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				serviceName: {
					// First port keys unsorted
					{Message: "Keys are not sorted at spec.ports[0]. Expected order: [name, port, protocol, targetPort], got: [port, targetPort, protocol, name]"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *SortedKeysTestSuite) TestMixedSortedAndUnsorted() {
	const deploymentName = "mixed-sorting"

	// Some levels sorted, some not
	mixedYAML := []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: test
    version: v1
  name: mixed-sorting
spec:
  selector:
    matchLabels:
      version: v1
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
        - image: test:latest
          name: test
`)

	deployment := &appsV1.Deployment{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name: deploymentName,
		},
	}

	s.ctx.AddObjectWithRaw(deploymentName, deployment, mixedYAML)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				Recursive: true,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: {
					// metadata.labels is sorted (app, version) ✓
					// spec.selector.matchLabels is unsorted (version, app) ✗
					{Message: "Keys are not sorted at spec.selector.matchLabels. Expected order: [app, version], got: [version, app]"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *SortedKeysTestSuite) TestComplexPodSpecWithVolumes() {
	const deploymentName = "complex-pod-spec"

	// Complex pod spec with init containers, volumes, security context
	complexPodYAML := []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: complex-pod-spec
spec:
  template:
    metadata:
      labels:
        app: complex
    spec:
      containers:
        - image: app:latest
          name: app
          volumeMounts:
            - mountPath: /data
              name: data-volume
          securityContext:
            runAsNonRoot: true
            allowPrivilegeEscalation: false
      initContainers:
        - command:
            - sh
            - -c
            - echo init
          image: busybox:latest
          name: init
      serviceAccountName: app-sa
      securityContext:
        fsGroup: 1000
        runAsUser: 1000
        runAsGroup: 1000
      volumes:
        - name: data-volume
          persistentVolumeClaim:
            claimName: data-pvc
`)

	deployment := &appsV1.Deployment{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name: deploymentName,
		},
	}

	s.ctx.AddObjectWithRaw(deploymentName, deployment, complexPodYAML)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				Recursive: true,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: {
					// spec.template.spec keys unsorted (serviceAccountName before securityContext)
					{Message: "Keys are not sorted at spec.template.spec. Expected order: [containers, initContainers, securityContext, serviceAccountName, volumes], got: [containers, initContainers, serviceAccountName, securityContext, volumes]"},
					// containers[0] keys unsorted (volumeMounts before securityContext)
					{Message: "Keys are not sorted at spec.template.spec.containers[0]. Expected order: [image, name, securityContext, volumeMounts], got: [image, name, volumeMounts, securityContext]"},
					// containers[0].securityContext unsorted
					{Message: "Keys are not sorted at spec.template.spec.containers[0].securityContext. Expected order: [allowPrivilegeEscalation, runAsNonRoot], got: [runAsNonRoot, allowPrivilegeEscalation]"},
					// spec.template.spec.securityContext unsorted (runAsUser before runAsGroup)
					{Message: "Keys are not sorted at spec.template.spec.securityContext. Expected order: [fsGroup, runAsGroup, runAsUser], got: [fsGroup, runAsUser, runAsGroup]"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *SortedKeysTestSuite) TestAllKeysSortedComplex() {
	const deploymentName = "all-sorted-complex"

	// Complex manifest with everything sorted correctly
	allSortedYAML := []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    app.kubernetes.io/version: "1.0"
    kubernetes.io/description: "test app"
  labels:
    app: test
    environment: prod
  name: all-sorted-complex
spec:
  replicas: 3
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
        - env:
            - name: ENV
              value: prod
          image: test:latest
          name: main
          ports:
            - containerPort: 8080
              name: http
              protocol: TCP
          resources:
            limits:
              cpu: "1"
              memory: 512Mi
            requests:
              cpu: 100m
              memory: 128Mi
      volumes:
        - emptyDir: {}
          name: cache
`)

	deployment := &appsV1.Deployment{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name: deploymentName,
		},
	}

	s.ctx.AddObjectWithRaw(deploymentName, deployment, allSortedYAML)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				Recursive: true,
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				deploymentName: nil, // Everything is properly sorted
			},
			ExpectInstantiationError: false,
		},
	})
}
