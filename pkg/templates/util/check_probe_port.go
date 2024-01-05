package util

import (
	"fmt"

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

	if httpProbe := probe.ProbeHandler.HTTPGet; httpProbe != nil {
		if _, ok := ports[httpProbe.Port]; !ok {
			return []diagnostic.Diagnostic{{
				Message: fmt.Sprintf("container %q does not expose port %s for the HTTPGet", container.Name, httpProbe.Port.String()),
			}}
		}
	}

	if tcpProbe := probe.ProbeHandler.TCPSocket; tcpProbe != nil {
		if _, ok := ports[tcpProbe.Port]; !ok {
			return []diagnostic.Diagnostic{{
				Message: fmt.Sprintf("container %q does not expose port %s for the TCPSocket", container.Name, tcpProbe.Port.String()),
			}}
		}
	}

	if grpcProbe := probe.ProbeHandler.GRPC; grpcProbe != nil {
		if _, ok := ports[intstr.FromInt32(grpcProbe.Port)]; !ok {
			return []diagnostic.Diagnostic{{
				Message: fmt.Sprintf("container %q does not expose port %d for the GRPC check", container.Name, grpcProbe.Port),
			}}
		}
	}
	return nil
}
