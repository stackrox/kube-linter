package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestCheckProbePort(t *testing.T) {
	httpProbe := func(port intstr.IntOrString) *v1.Probe {
		return &v1.Probe{ProbeHandler: v1.ProbeHandler{HTTPGet: &v1.HTTPGetAction{Port: port}}}
	}

	testCases := []struct {
		name      string
		container v1.Container
		probe     *v1.Probe
		wantDiag  bool
	}{
		{
			name:      "nil probe is ignored",
			container: v1.Container{Name: "c"},
			probe:     nil,
			wantDiag:  false,
		},
		{
			name: "declared container port matches",
			container: v1.Container{
				Name:  "c",
				Ports: []v1.ContainerPort{{ContainerPort: 8080}},
			},
			probe:    httpProbe(intstr.FromInt(8080)),
			wantDiag: false,
		},
		{
			name: "undeclared port not present anywhere is flagged",
			container: v1.Container{
				Name:  "c",
				Ports: []v1.ContainerPort{{ContainerPort: 8080}},
			},
			probe:    httpProbe(intstr.FromInt(8081)),
			wantDiag: true,
		},
		{
			name: "port exposed via args is accepted (issue #1086)",
			container: v1.Container{
				Name:  "manager",
				Args:  []string{"--metrics-addr=0.0.0.0:8080", "--health-probe-addr=:8081"},
				Ports: []v1.ContainerPort{{Name: "metrics", ContainerPort: 8080, Protocol: v1.ProtocolTCP}},
			},
			probe:    httpProbe(intstr.FromInt(8081)),
			wantDiag: false,
		},
		{
			name: "port exposed via command is accepted",
			container: v1.Container{
				Name:    "manager",
				Command: []string{"/manager", "--health-probe-addr=:8081"},
			},
			probe:    httpProbe(intstr.FromInt(8081)),
			wantDiag: false,
		},
		{
			name: "port that is only a substring of a larger number is still flagged",
			container: v1.Container{
				Name: "manager",
				Args: []string{"--metrics-addr=0.0.0.0:18081"},
			},
			probe:    httpProbe(intstr.FromInt(8081)),
			wantDiag: true,
		},
		{
			name: "trailing-digit superset of the port is still flagged",
			container: v1.Container{
				Name: "manager",
				Args: []string{"--addr=:80818"},
			},
			probe:    httpProbe(intstr.FromInt(8081)),
			wantDiag: true,
		},
		{
			name: "named string probe port is not matched against args",
			container: v1.Container{
				Name: "manager",
				Args: []string{"--health-probe-name=healthz"},
			},
			probe:    httpProbe(intstr.FromString("healthz")),
			wantDiag: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.container
			diags := CheckProbePort(&c, tc.probe)
			if tc.wantDiag {
				assert.NotEmpty(t, diags, "expected a diagnostic")
			} else {
				assert.Empty(t, diags, "expected no diagnostic")
			}
		})
	}
}

func TestArgContainsPort(t *testing.T) {
	cases := []struct {
		arg    string
		needle string
		want   bool
	}{
		{"--health-probe-addr=:8081", "8081", true},
		{"--addr=0.0.0.0:8081", "8081", true},
		{"8081", "8081", true},
		{"--addr=:18081", "8081", false},
		{"--addr=:80818", "8081", false},
		{"--addr=:808", "8081", false},
		{"--flag", "8081", false},
		{"--addr=:8081 --other=:8081", "8081", true},
	}
	for _, c := range cases {
		assert.Equalf(t, c.want, argContainsPort(c.arg, c.needle), "argContainsPort(%q, %q)", c.arg, c.needle)
	}
}
