package kubelinter.objectkinds

# Any represents the ObjectKind that matches any object
is_any := true

# Ingress represents Kubernetes Ingress objects
is_ingress if {
	input.kind == "Ingress"
	input.apiVersion == "networking.k8s.io/v1"
}

# Service represents Kubernetes Service objects
is_service if {
	input.kind == "Service"
	input.apiVersion == "v1"
}

# ServiceAccount represents Kubernetes ServiceAccount objects
is_serviceaccount if {
	input.kind == "ServiceAccount"
	input.apiVersion == "v1"
}

# ClusterRole represents Kubernetes ClusterRole objects
is_clusterrole if {
	input.kind == "ClusterRole"
	input.apiVersion == "rbac.authorization.k8s.io/v1"
}

# ClusterRoleBinding represents Kubernetes ClusterRoleBinding objects
is_clusterrolebinding if {
	input.kind == "ClusterRoleBinding"
	input.apiVersion == "rbac.authorization.k8s.io/v1"
}

# Role represents Kubernetes Role objects
is_role if {
	input.kind == "Role"
	input.apiVersion == "rbac.authorization.k8s.io/v1"
}

# RoleBinding represents Kubernetes RoleBinding objects
is_rolebinding if {
	input.kind == "RoleBinding"
	input.apiVersion == "rbac.authorization.k8s.io/v1"
}

# NetworkPolicy represents Kubernetes NetworkPolicy objects
is_networkpolicy if {
	input.kind == "NetworkPolicy"
	input.apiVersion == "networking.k8s.io/v1"
}

# PodDisruptionBudget represents Kubernetes PodDisruptionBudget objects
is_poddisruptionbudget if {
	input.kind == "PodDisruptionBudget"
	input.apiVersion == "policy/v1"
}

# HorizontalPodAutoscaler represents Kubernetes HorizontalPodAutoscaler objects
is_horizontalpodautoscaler if {
	input.kind == "HorizontalPodAutoscaler"
	input.apiVersion == "autoscaling/v1"
}

is_horizontalpodautoscaler if {
	input.kind == "HorizontalPodAutoscaler"
	input.apiVersion == "autoscaling/v2"
}

is_horizontalpodautoscaler if {
	input.kind == "HorizontalPodAutoscaler"
	input.apiVersion == "autoscaling/v2beta1"
}

is_horizontalpodautoscaler if {
	input.kind == "HorizontalPodAutoscaler"
	input.apiVersion == "autoscaling/v2beta2"
}

# SecurityContextConstraints represents OpenShift SecurityContextConstraints objects
is_securitycontextconstraints if {
	input.kind == "SecurityContextConstraints"
	input.apiVersion == "security.openshift.io/v1"
}

# ServiceMonitor represents Prometheus Service Monitor objects
is_servicemonitor if {
	input.kind == "ServiceMonitor"
	input.apiVersion == "monitoring.coreos.com/v1"
}

# ScaledObject represents KEDA ScaledObject objects
is_scaledobject if {
	input.kind == "ScaledObject"
	input.apiVersion == "keda.sh/v1alpha1"
}

# JobLike matches Job and CronJob objects
is_job_like if {
	input.kind == "Job"
	input.apiVersion == "batch/v1"
}

is_job_like if {
	input.kind == "CronJob"
	input.apiVersion == "batch/v1"
}

# DeploymentLike matches various deployment-like objects
is_deployment_like if {
	input.kind == "Deployment"
	input.apiVersion == "apps/v1"
}

is_deployment_like if {
	input.kind == "DaemonSet"
	input.apiVersion == "apps/v1"
}

is_deployment_like if {
	input.kind == "StatefulSet"
	input.apiVersion == "apps/v1"
}

is_deployment_like if {
	input.kind == "ReplicaSet"
	input.apiVersion == "apps/v1"
}

is_deployment_like if {
	input.kind == "Pod"
	input.apiVersion == "v1"
}

is_deployment_like if {
	input.kind == "ReplicationController"
	input.apiVersion == "v1"
}

is_deployment_like if {
	input.kind == "Job"
	input.apiVersion == "batch/v1"
}

is_deployment_like if {
	input.kind == "CronJob"
	input.apiVersion == "batch/v1"
}

# OpenShift specific deployment-like objects
is_deployment_like if {
	input.kind == "DeploymentConfig"
	input.apiVersion == "apps.openshift.io/v1"
}
