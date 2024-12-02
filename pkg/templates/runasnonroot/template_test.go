package runasnonroot

import (
	"testing"

	"golang.stackrox.io/kube-linter/internal/pointers"
	v1 "k8s.io/api/core/v1"
)

func TestEffectiveRunAsNonRoot(t *testing.T) {
	tests := []struct {
		name        string
		podSC       *v1.PodSecurityContext
		containerSC *v1.SecurityContext
		expected    bool
	}{
		{
			name:        "both nil",
			podSC:       nil,
			containerSC: nil,
			expected:    false,
		},
		{
			name: "podSC set",
			podSC: &v1.PodSecurityContext{
				RunAsNonRoot: pointers.Bool(true),
			},
			containerSC: nil,
			expected:    true,
		},
		{
			name:  "containerSC set",
			podSC: nil,
			containerSC: &v1.SecurityContext{
				RunAsNonRoot: pointers.Bool(true),
			},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := effectiveRunAsNonRoot(test.podSC, test.containerSC)
			if result != test.expected {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestEffectiveRunAsUser(t *testing.T) {
	user := pointers.Int64(1000)
	tests := []struct {
		name        string
		podSC       *v1.PodSecurityContext
		containerSC *v1.SecurityContext
		expected    *int64
	}{
		{
			name:        "both nil",
			podSC:       nil,
			containerSC: nil,
			expected:    nil,
		},
		{
			name: "podSC set",
			podSC: &v1.PodSecurityContext{
				RunAsUser: user,
			},
			containerSC: nil,
			expected:    user,
		},
		{
			name:  "containerSC set",
			podSC: nil,
			containerSC: &v1.SecurityContext{
				RunAsUser: user,
			},
			expected: user,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := effectiveRunAsUser(test.podSC, test.containerSC)
			if result != test.expected {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestEffectiveRunAsGroup(t *testing.T) {
	group := pointers.Int64(1000)
	tests := []struct {
		name        string
		podSC       *v1.PodSecurityContext
		containerSC *v1.SecurityContext
		expected    *int64
	}{
		{
			name:        "both nil",
			podSC:       nil,
			containerSC: nil,
			expected:    nil,
		},
		{
			name: "podSC set",
			podSC: &v1.PodSecurityContext{
				RunAsGroup: group,
			},
			containerSC: nil,
			expected:    group,
		},
		{
			name:  "containerSC set",
			podSC: nil,
			containerSC: &v1.SecurityContext{
				RunAsGroup: group,
			},
			expected: group,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := effectiveRunAsGroup(test.podSC, test.containerSC)
			if result != test.expected {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestIsNonZero(t *testing.T) {
	tests := []struct {
		name     string
		number   *int64
		expected bool
	}{
		{
			name:     "Nil case",
			number:   nil,
			expected: false,
		},
		{
			name:     "zero case",
			number:   pointers.Int64(0),
			expected: false,
		},
		{
			name:     "non-zero case",
			number:   pointers.Int64(1000),
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isNonZero(test.number)
			if result != test.expected {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		})
	}
}

// func TestRunAsNonRoot(t *testing.T) {
// 	suite.Run(t, new(RunAsNonRootTestSuite))
// }
//
// type RunAsNonRootTestSuite struct {
// 	templates.TemplateTestSuite
// 	ctx *mocks.MockLintContext
// }
//
// func (s *RunAsNonRootTestSuite) SetupTest() {
// 	s.Init(templateKey)
// 	s.ctx = mocks.NewMockContext()
// }
//
// func (s *RunAsNonRootTestSuite) TestEffectiveRunAsGroup() {
// 	const targetName = "deployment01"
// 	testCases := []struct {
// 		name        string
// 		podSC       *v1.PodSecurityContext
// 		containerSC *v1.SecurityContext
// 		expected    map[string][]diagnostic.Diagnostic
// 	}{
// 		{
// 			name: "both nil",
// 			podSC: &v1.PodSecurityContext{
// 				RunAsGroup: pointers.Int64(1000),
// 			},
// 			containerSC: nil,
// 			expected:    nil,
// 		},
// 		{
// 			name:        "both nil",
// 			podSC:       nil,
// 			containerSC: nil,
// 			expected:    nil,
// 		},
// 		{
// 			name:        "both nil",
// 			podSC:       nil,
// 			containerSC: nil,
// 			expected:    nil,
// 		},
// 	}
//
// 	for _, tc := range testCases {
// 		s.Run(tc.name, func() {
// 			s.ctx.AddMockDeployment(s.T(), targetName)
// 			s.Validate(s.ctx, []templates.TestCase{{
// 				Diagnostics: tc.expected,
// 			}})
// 		})
// 	}
// }
//
// // [[ "${message1}" == "Deployment: container \"app\" has runAsGroup set to 0" ]]
// // [[ "${message2}" == "Deployment: container \"app\" is not set to runAsNonRoot" ]]
// // [[ "${message3}" == "DeploymentConfig: container \"app2\" has runAsGroup set to 0" ]]
// // [[ "${message4}" == "DeploymentConfig: container \"app2\" is not set to runAsNonRoot" ]]
//
// func (suite *RunAsNonRootTestSuite) TestIsNonZero() {
// 	tests := []struct {
// 		name     string
// 		number   *int64
// 		expected bool
// 	}{
// 		{"Nil number", nil, false},
// 		{"Zero number", pointers.Int64(0), false},
// 		{"Non-zero number", pointers.Int64(1), true},
// 	}
//
// 	for _, test := range tests {
// 		suite.Run(test.name, func() {
// 			result := isNonZero(test.number)
// 			suite.Equal(test.expected, result)
// 		})
// 	}
// }
