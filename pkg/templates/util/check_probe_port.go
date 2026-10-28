package util

import (
	"fmt"
	"strconv"
	"strings"

	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var sentinel = struct{}{}

func CheckProbePort(container *v1.Container, probe *v1.Probe) []diagnostic.Diagnostic {
	if probe == nil {
		return nil
	}

	ports := map[intstr.IntOrString]struct{}{}
	for _, port := range container.Ports {
		if port.Protocol != "" && port.Protocol != v1.ProtocolTCP {
			continue
		}
		ports[intstr.FromInt32(port.ContainerPort)] = sentinel
		ports[intstr.FromString(port.Name)] = sentinel
	}

	if httpProbe := probe.HTTPGet; httpProbe != nil {
		if _, ok := ports[httpProbe.Port]; !ok && !probePortInArgs(container, httpProbe.Port) {
			return []diagnostic.Diagnostic{{
				Message: fmt.Sprintf("container %q does not expose port %s for the HTTPGet", container.Name, httpProbe.Port.String()),
			}}
		}
	}

	if tcpProbe := probe.TCPSocket; tcpProbe != nil {
		if _, ok := ports[tcpProbe.Port]; !ok && !probePortInArgs(container, tcpProbe.Port) {
			return []diagnostic.Diagnostic{{
				Message: fmt.Sprintf("container %q does not expose port %s for the TCPSocket", container.Name, tcpProbe.Port.String()),
			}}
		}
	}

	if grpcProbe := probe.GRPC; grpcProbe != nil {
		if _, ok := ports[intstr.FromInt32(grpcProbe.Port)]; !ok && !probePortInArgs(container, intstr.FromInt32(grpcProbe.Port)) {
			return []diagnostic.Diagnostic{{
				Message: fmt.Sprintf("container %q does not expose port %d for the GRPC check", container.Name, grpcProbe.Port),
			}}
		}
	}
	return nil
}

// probePortInArgs reports whether the probe port is wired up through the
// container args or command rather than declared as a containerPort. A
// containerPort entry is informational only: a process can listen on any port
// regardless of whether it is declared. Charts such as opentelemetry-operator
// pass the health-probe address purely through a flag (for example
// "--health-probe-addr=:8081"), so the declared ports do not include it and the
// probe-port checks would otherwise report a false positive (see issue #1086).
//
// Only a numeric probe port can be matched this way; a named port (string)
// still has to resolve against a declared containerPort, so it is left to the
// caller's normal lookup.
func probePortInArgs(container *v1.Container, port intstr.IntOrString) bool {
	if port.Type != intstr.Int {
		return false
	}
	portNum := port.IntValue()
	if portNum <= 0 {
		return false
	}
	needle := strconv.Itoa(portNum)
	for _, arg := range container.Args {
		if argContainsPort(arg, needle) {
			return true
		}
	}
	for _, cmd := range container.Command {
		if argContainsPort(cmd, needle) {
			return true
		}
	}
	return false
}

// argContainsPort reports whether needle (the port number rendered as a string)
// appears in arg as a standalone integer token, so that searching for "8081"
// does not match a larger number such as "18081" or "80818" (or a port that
// happens to be a substring of an image tag or version).
func argContainsPort(arg, needle string) bool {
	for from := 0; ; {
		idx := strings.Index(arg[from:], needle)
		if idx < 0 {
			return false
		}
		start := from + idx
		end := start + len(needle)
		beforeOK := start == 0 || !isDigit(arg[start-1])
		afterOK := end == len(arg) || !isDigit(arg[end])
		if beforeOK && afterOK {
			return true
		}
		from = start + 1
	}
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
