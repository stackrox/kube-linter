//+kubebuilder:object:generate=true

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// MiscSpec defines miscellaneous settings for custom resources.
type MiscSpec struct {
	// Deprecated field. This field will be removed in a future release.
	// Set this to true to have the operator create SecurityContextConstraints (SCCs) for the operands. This
	// isn't usually needed, and may interfere with other workloads.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1,displayName="Create SecurityContextConstraints for Operand",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	CreateSCCs *bool `json:"createSCCs,omitempty" deprecated:"true"`
}

// PinToNodesPolicy defines preset node pinning configurations.
// Use this for common scenarios like pinning to OpenShift infrastructure nodes.
// +kubebuilder:validation:Enum=None;InfraRole
type PinToNodesPolicy string

const (
	// PinToNodesNone does not apply any node scheduling constraints.
	PinToNodesNone PinToNodesPolicy = "None"
	// PinToNodesInfraRole pins deployments to OpenShift infrastructure nodes by setting
	// nodeSelector to "node-role.kubernetes.io/infra" and adding the corresponding tolerations.
	PinToNodesInfraRole PinToNodesPolicy = "InfraRole"
)

// DeploymentDefaultsSpec defines default scheduling constraints for Deployment-based components.
type DeploymentDefaultsSpec struct {
	// Pin all Deployment-based components to specific node types. This is a convenience setting
	// that automatically configures both nodeSelector and tolerations with predefined values.
	// Use this for common scenarios like running on OpenShift infrastructure nodes.
	// For custom node selection, use the explicit nodeSelector and tolerations fields instead.
	// Cannot be used together with nodeSelector or tolerations fields.
	// The default is: None.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1
	PinToNodes *PinToNodesPolicy `json:"pinToNodes,omitempty"`

	// Default nodeSelector applied to all Deployment-based components. Use this for custom node
	// selection criteria.
	// Cannot be used together with pinToNodes.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Node Selector",order=2
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Default tolerations applied to all Deployment-based components. Use this when your target
	// nodes have custom taints that pods must tolerate.
	// Cannot be used together with pinToNodes.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=3
	Tolerations []*corev1.Toleration `json:"tolerations,omitempty"`
}

// CustomizeSpec defines customizations to apply.
type CustomizeSpec struct {
	// Custom labels to set on all managed objects.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1
	Labels map[string]string `json:"labels,omitempty"`
	// Custom annotations to set on all managed objects.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=2
	Annotations map[string]string `json:"annotations,omitempty"`

	// Custom environment variables to set on managed pods' containers.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=3,displayName="Environment Variables"
	EnvVars []corev1.EnvVar `json:"envVars,omitempty"`

	// Global nodeSelector and tolerations for Deployment-based components. DaemonSets (Collector) are not affected.
	// Component-level nodeSelector and tolerations settings override these defaults on a field-by-field basis.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=4,displayName="Deployment Defaults"
	DeploymentDefaults *DeploymentDefaultsSpec `json:"deploymentDefaults,omitempty"`
}

// DeploymentSpec defines settings that affect a deployment.
type DeploymentSpec struct {
	// Allows overriding the default resource settings for this component. Please consult the documentation
	// for an overview of default resource requirements and a sizing guide.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:resourceRequirements"},order=100
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// If you want this component to only run on specific nodes, you can configure a node selector here.
	// This setting overrides spec.customize.deploymentDefaults.nodeSelector.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Node Selector",order=101
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// If you want this component to only run on specific nodes, you can configure tolerations of tainted nodes.
	// This setting overrides spec.customize.deploymentDefaults.tolerations.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:tolerations"},order=102
	Tolerations []*corev1.Toleration `json:"tolerations,omitempty"`

	// HostAliases allows configuring additional hostnames to resolve in the pod's hosts file.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=103,displayName="Host Aliases"
	HostAliases []corev1.HostAlias `json:"hostAliases,omitempty"`
}

// StackRoxCondition defines a condition for a StackRox custom resource.
type StackRoxCondition struct {
	Type    ConditionType   `json:"type"`
	Status  ConditionStatus `json:"status"`
	Reason  ConditionReason `json:"reason,omitempty"`
	Message string          `json:"message,omitempty"`

	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
}

// ConditionType is a type of values of condition type.
type ConditionType string

// ConditionStatus is a type of values of condition status.
type ConditionStatus string

// ConditionReason is a type of values of condition reason.
type ConditionReason string

// These are the allowed values for StackRoxCondition fields.
const (
	ConditionInitialized    ConditionType = "Initialized"
	ConditionDeployed       ConditionType = "Deployed"
	ConditionReleaseFailed  ConditionType = "ReleaseFailed"
	ConditionIrreconcilable ConditionType = "Irreconcilable"

	// These are specifically owned by the status controllers.
	ConditionProgressing ConditionType = "Progressing"
	ConditionAvailable   ConditionType = "Available"

	StatusTrue    ConditionStatus = "True"
	StatusFalse   ConditionStatus = "False"
	StatusUnknown ConditionStatus = "Unknown"

	ReasonInstallSuccessful   ConditionReason = "InstallSuccessful"
	ReasonUpgradeSuccessful   ConditionReason = "UpgradeSuccessful"
	ReasonUninstallSuccessful ConditionReason = "UninstallSuccessful"
	ReasonInstallError        ConditionReason = "InstallError"
	ReasonUpgradeError        ConditionReason = "UpgradeError"
	ReasonReconcileError      ConditionReason = "ReconcileError"
	ReasonUninstallError      ConditionReason = "UninstallError"
)

// StackRoxRelease describes the Helm "release" that was most recently applied.
type StackRoxRelease struct {
	Version string `json:"version,omitempty"`
}

// AdditionalCA defines a certificate for an additional Certificate Authority.
type AdditionalCA struct {
	// Must be a valid file basename
	Name string `json:"name"`
	// PEM format
	Content string `json:"content"`
}

// TLSConfig defines common TLS-related settings for all components.
type TLSConfig struct {
	// Allows you to specify additional trusted Root CAs.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Additional CAs"
	AdditionalCAs []AdditionalCA `json:"additionalCAs,omitempty"`
}

// LocalSecretReference is a reference to a secret within the same namespace.
type LocalSecretReference struct {
	// The name of the referenced secret.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:io.kubernetes:Secret"}
	Name string `json:"name"`
}

// LocalConfigMapReference is a reference to a config map within the same namespace.
type LocalConfigMapReference struct {
	// The name of the referenced config map.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:io.kubernetes:ConfigMap"}
	Name string `json:"name"`
}

// ScannerAnalyzerComponent describes the analyzer component
type ScannerAnalyzerComponent struct {
	// Controls the number of analyzer replicas and autoscaling.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1
	Scaling *ScannerComponentScaling `json:"scaling,omitempty"`

	DeploymentSpec `json:",inline"`
}

// GetScaling returns scaling config even if receiver is nil
func (s *ScannerAnalyzerComponent) GetScaling() *ScannerComponentScaling {
	if s == nil {
		return nil
	}
	return s.Scaling
}

// ScannerV4Component defines common configuration for Scanner V4 indexer and matcher components.
type ScannerV4Component struct {
	// Controls the number of replicas and autoscaling for this component.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1
	Scaling        *ScannerComponentScaling `json:"scaling,omitempty"`
	DeploymentSpec `json:",inline"`
}

// ScannerV4DB defines configuration for the Scanner V4 database component.
type ScannerV4DB struct {
	// Configures how Scanner V4 should store its persistent data.
	// You can use a persistent volume claim (the recommended default), a host path,
	// or an emptyDir volume if Scanner V4 is running on a secured cluster without default StorageClass.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1
	Persistence    *ScannerV4Persistence `json:"persistence,omitempty"`
	DeploymentSpec `json:",inline"`
}

// GetPersistence returns the configured PVC
func (sdb *ScannerV4DB) GetPersistence() *ScannerV4Persistence {
	if sdb == nil {
		return nil
	}
	return sdb.Persistence
}

// ScannerV4Persistence defines persistence settings for Scanner V4.
type ScannerV4Persistence struct {
	// Uses a Kubernetes persistent volume claim (PVC) to manage the storage location of persistent data.
	// Recommended for most users.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Persistent volume claim",order=1
	PersistentVolumeClaim *ScannerV4PersistentVolumeClaim `json:"persistentVolumeClaim,omitempty"`

	// Stores persistent data in a directory on the host. This is not recommended, and should only
	// be used together with a node selector (only available in YAML view).
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Host path",order=99
	HostPath *HostPathSpec `json:"hostPath,omitempty"`
}

// GetPersistentVolumeClaim returns the configured PVC
func (p *ScannerV4Persistence) GetPersistentVolumeClaim() *ScannerV4PersistentVolumeClaim {
	if p == nil {
		return nil
	}
	return p.PersistentVolumeClaim
}

// GetHostPath returns the configured host path
func (p *ScannerV4Persistence) GetHostPath() string {
	if p == nil {
		return ""
	}
	if p.HostPath == nil {
		return ""
	}

	return pointer.StringDeref(p.HostPath.Path, "")
}

// ScannerV4PersistentVolumeClaim defines PVC-based persistence settings for Scanner V4 DB.
type ScannerV4PersistentVolumeClaim struct {
	// The name of the PVC to manage persistent data. If no PVC with the given name exists, it will be
	// created.
	// The default is: scanner-v4-db.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Claim Name",order=1
	ClaimName *string `json:"claimName,omitempty"`

	// The size of the persistent volume when created through the claim. If a claim was automatically created,
	// this can be used after the initial deployment to resize (grow) the volume (only supported by some
	// storage class controllers).
	//+kubebuilder:validation:Pattern=^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Size",order=2,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Size *string `json:"size,omitempty"`

	// The name of the storage class to use for the PVC. If your cluster is not configured with a default storage
	// class, you must select a value here.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Storage Class",order=3,xDescriptors={"urn:alm:descriptor:io.kubernetes:StorageClass"}
	StorageClassName *string `json:"storageClassName,omitempty"`
}

// GetStorageClassName gets the StorageClassName string value
// returns empty string if the object or the StorageClassName pointer is nil
func (s *ScannerV4PersistentVolumeClaim) GetStorageClassName() string {
	if s == nil || s.StorageClassName == nil {
		return ""
	}

	return *s.StorageClassName
}

// GetClaimName gets the ClaimName string value
// returns empty string if the object or the ClaimName pointer is nil
func (s *ScannerV4PersistentVolumeClaim) GetClaimName() string {
	if s == nil || s.ClaimName == nil {
		return ""
	}

	return *s.ClaimName
}

// ScannerComponentScaling defines replication settings of scanner components.
type ScannerComponentScaling struct {
	// When enabled, the number of component replicas is managed dynamically based on the load, within the limits
	// specified below.
	// The default is: Enabled.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Autoscaling",order=1
	AutoScaling *AutoScalingPolicy `json:"autoScaling,omitempty"`

	// When autoscaling is disabled, the number of replicas will always be configured to match this value.
	// The default is: 3.
	//+kubebuilder:validation:Minimum=1
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Default Replicas",order=2
	Replicas *int32 `json:"replicas,omitempty"`

	// The default is: 2.
	//+kubebuilder:validation:Minimum=1
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Autoscaling Minimum Replicas",order=3,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:fieldDependency:.autoScaling:Enabled"}
	MinReplicas *int32 `json:"minReplicas,omitempty"`

	// The default is: 5.
	//+kubebuilder:validation:Minimum=1
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Autoscaling Maximum Replicas",order=4,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:fieldDependency:.autoScaling:Enabled"}
	MaxReplicas *int32 `json:"maxReplicas,omitempty"`
}

// AutoScalingPolicy is a type for values of spec.scanner.analyzer.replicas.autoScaling.
// +kubebuilder:validation:Enum=Enabled;Disabled
type AutoScalingPolicy string

const (
	// ScannerAutoScalingEnabled means that scanner autoscaling should be enabled.
	ScannerAutoScalingEnabled AutoScalingPolicy = "Enabled"
	// ScannerAutoScalingDisabled means that scanner autoscaling should be disabled.
	ScannerAutoScalingDisabled AutoScalingPolicy = "Disabled"
)

// Monitoring defines settings for monitoring endpoint.
type Monitoring struct {
	// Expose the monitoring endpoint. A new service, "monitoring",
	// with port 9090, will be created as well as a network policy allowing
	// inbound connections to the port.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1
	ExposeEndpoint *ExposeEndpoint `json:"exposeEndpoint,omitempty"`
}

// IsEnabled checks whether exposing of endpoint is enabled.
// This method is safe to be used with nil receivers.
func (s *Monitoring) IsEnabled() bool {
	if s == nil || s.ExposeEndpoint == nil {
		return false // disabled by default
	}

	return *s.ExposeEndpoint == ExposeEndpointEnabled
}

// ExposeEndpoint is a type for monitoring sub-struct.
// +kubebuilder:validation:Enum=Enabled;Disabled
type ExposeEndpoint string

const (
	// ExposeEndpointEnabled means that component should expose monitoring port.
	ExposeEndpointEnabled ExposeEndpoint = "Enabled"
	// ExposeEndpointDisabled means that component should not expose monitoring port.
	ExposeEndpointDisabled ExposeEndpoint = "Disabled"
)

// GlobalMonitoring defines settings related to global monitoring. Contrary to
// `Monitoring`, the corresponding Helm flag lives in the global scope `.monitoring`.
type GlobalMonitoring struct {
	OpenShiftMonitoring *OpenShiftMonitoring `json:"openshift,omitempty"`
}

// OpenShiftMonitoring defines settings related to OpenShift Monitoring
type OpenShiftMonitoring struct {
	// The default is: true.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	Enabled *bool `json:"enabled"`
}

// IsOpenShiftMonitoringDisabled returns true if OpenShiftMonitoring is disabled.
// This function is nil safe.
func (m *GlobalMonitoring) IsOpenShiftMonitoringDisabled() bool {
	return m != nil && m.OpenShiftMonitoring != nil && m.OpenShiftMonitoring.Enabled != nil && !*m.OpenShiftMonitoring.Enabled
}

// Set this to Disabled to prevent the operator from creating NetworkPolicy objects.
// +kubebuilder:validation:Enum=Enabled;Disabled
type NetworkPolicies string

const (
	// NetworkPoliciesEnabled means that network policies should be created.
	NetworkPoliciesEnabled NetworkPolicies = "Enabled"
	// NetworkPoliciesDisabled means that network policies should not be created.
	NetworkPoliciesDisabled NetworkPolicies = "Disabled"
)

// IsNetworkPoliciesEnabled checks whether network policies are enabled.
// This method is safe to be used with nil receivers.
func (s *GlobalNetworkSpec) IsNetworkPoliciesEnabled() bool {
	if s == nil || s.Policies == nil {
		return true // enabled by default
	}

	return *s.Policies == NetworkPoliciesEnabled
}

// GlobalNetworkSpec defines settings related to Helm chart network parameters. The corresponding Helm flags
// live in the global scope `.network`.
type GlobalNetworkSpec struct {
	// To provide security at the network level, the ACS Operator creates NetworkPolicy resources by default. If you want to manage your own NetworkPolicy objects then set this to "Disabled".
	// The default is: Enabled.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1,displayName="Network Policies"
	Policies *NetworkPolicies `json:"policies,omitempty"`
}

// +kubebuilder:object:generate=false
type ObjectForStatusController interface {
	ctrlClient.Object
	GetCondition(condType ConditionType) *StackRoxCondition
	SetCondition(StackRoxCondition) bool
	GetGeneration() int64
	GetObservedGeneration() int64
}

// getCondition returns a specific condition by type, or nil if not found.
func getCondition(conditions []StackRoxCondition, condType ConditionType) *StackRoxCondition {
	for i := range conditions {
		if conditions[i].Type == condType {
			return &conditions[i]
		}
	}
	return nil
}

// updateCondition updates or adds a condition. Returns the updated condition list alongside
// a boolean indicating if the condition has changed.
func updateCondition(conditions []StackRoxCondition, updatedCond StackRoxCondition) ([]StackRoxCondition, bool) {
	for i, cond := range conditions {
		if cond.Type == updatedCond.Type {
			// Check if update is needed.
			if cond.Status == updatedCond.Status &&
				cond.Reason == updatedCond.Reason &&
				cond.Message == updatedCond.Message {
				return conditions, false
			}
			// Update existing condition.
			conditions[i] = updatedCond
			return conditions, true
		}
	}
	// Condition doesn't exist, add it.
	return append(conditions, updatedCond), true
}
