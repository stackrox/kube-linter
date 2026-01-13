/*
Copyright 2021 Red Hat.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Important: Run "make generate manifests" to regenerate code after modifying this file

// -------------------------------------------------------------
// Spec

// SecuredClusterSpec defines the desired configuration state of a secured cluster.
type SecuredClusterSpec struct {
	// The unique name of this cluster, as it will be shown in the Red Hat Advanced Cluster Security UI.
	// Note: Once a name is set here, you will not be able to change it again. You will need to delete
	// and re-create this object in order to register a cluster with a new name.
	//+kubebuilder:validation:Required
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1
	ClusterName *string `json:"clusterName"`

	// Custom labels associated with a secured cluster in Red Hat Advanced Cluster Security.
	ClusterLabels map[string]string `json:"clusterLabels,omitempty"`

	// The endpoint of the Red Hat Advanced Cluster Security Central instance to connect to,
	// including the port number. If no port is specified and the endpoint contains an https://
	// protocol specification, then the port 443 is implicitly assumed.
	// If using a non-gRPC capable load balancer, use the WebSocket protocol by prefixing the endpoint
	// address with wss://.
	// Note: when leaving this blank, Sensor will attempt to connect to a Central instance running in the same
	// namespace.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=2
	CentralEndpoint *string `json:"centralEndpoint,omitempty"`

	// Settings for the Sensor component.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=3,displayName="Sensor Settings"
	Sensor *SensorComponentSpec `json:"sensor,omitempty"`

	// Settings for the Admission Control component, which is necessary for preventive policy enforcement,
	// and for Kubernetes event monitoring.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=4,displayName="Admission Control Settings"
	AdmissionControl *AdmissionControlComponentSpec `json:"admissionControl,omitempty"`

	// Settings for the components running on each node in the cluster (Collector and Compliance).
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=5,displayName="Per Node Settings"
	PerNode *PerNodeSpec `json:"perNode,omitempty"`

	// Settings relating to the ingestion of Kubernetes audit logs.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=6,displayName="Kubernetes Audit Logs Ingestion Settings"
	AuditLogs *AuditLogsSpec `json:"auditLogs,omitempty"`

	// Settings relating to process baselines.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=7,displayName="Process Baselines Settings"
	ProcessBaselines *ProcessBaselinesSpec `json:"processBaselines,omitempty"`

	// Settings for the Scanner component, which is responsible for vulnerability scanning of container
	// images stored in a cluster-local image repository.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=8,displayName="Scanner Component Settings"
	Scanner *LocalScannerComponentSpec `json:"scanner,omitempty"`

	// Settings for the Scanner V4 components, which can run in addition to the previously existing Scanner components
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=9,displayName="Scanner V4 Component Settings"
	ScannerV4 *LocalScannerV4ComponentSpec `json:"scannerV4,omitempty"`
	// Above default is necessary to make the nested default work see: https://github.com/kubernetes-sigs/controller-tools/issues/622

	// Settings related to Transport Layer Security, such as Certificate Authorities.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=10
	TLS *TLSConfig `json:"tls,omitempty"`

	// Additional image pull secrets to be taken into account for pulling images.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Image Pull Secrets",order=11,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:advanced"}
	ImagePullSecrets []LocalSecretReference `json:"imagePullSecrets,omitempty"`

	// Customizations to apply on all Secured Cluster Services components.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName=Customizations,order=12,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:advanced"}
	Customize *CustomizeSpec `json:"customize,omitempty"`

	// Deprecated field. This field will be removed in a future release.
	// Miscellaneous settings.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName=Miscellaneous,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	Misc *MiscSpec `json:"misc,omitempty" deprecated:"true"`

	// Overlays
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName=Overlays,order=13,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	Overlays []*K8sObjectOverlay `json:"overlays,omitempty"`

	// Monitoring configuration.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=14,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:advanced"}
	Monitoring *GlobalMonitoring `json:"monitoring,omitempty"`

	// Set this parameter to override the default registry in images. For example, nginx:latest -> <registry override>/library/nginx:latest
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Custom Default Image Registry",order=15,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:advanced"}
	RegistryOverride *string `json:"registryOverride,omitempty"`

	// Network configuration.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName=Network,order=16,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:advanced"}
	Network *GlobalNetworkSpec `json:"network,omitempty"`
}

// ProcessBaselinesAutoLockMode is a type for values of spec.processBaselineAutoLockMode.
// +kubebuilder:validation:Enum=Enabled;Disabled
type ProcessBaselinesAutoLockMode string

const (
	// ProcessBaselineLockModeEnabled means: Process baseline auto-locking will be enabled
	ProcessBaselinesAutoLockModeEnabled ProcessBaselinesAutoLockMode = "Enabled"
	// ProcessBaselineLockModeDisabled means: Process baseline auto-locking will be disabled
	ProcessBaselinesAutoLockModeDisabled ProcessBaselinesAutoLockMode = "Disabled"
)

// Pointer returns the given ProcessBaselineAutoLockMode as a pointer, needed in k8s resource structs.
func (p ProcessBaselinesAutoLockMode) Pointer() *ProcessBaselinesAutoLockMode {
	return &p
}

// ProcessBaselinesSpec defines settings for the process baseline auto-locking feature.
type ProcessBaselinesSpec struct {
	// Should process baselines be automatically locked when the observation period (1 hour by default) ends.
	// The default is: Disabled.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:select:Enabled", "urn:alm:descriptor:com.tectonic.ui:select:Disabled"}
	AutoLock *ProcessBaselinesAutoLockMode `json:"autoLock,omitempty"`
}

// SensorComponentSpec defines settings for sensor.
type SensorComponentSpec struct {
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1
	DeploymentSpec `json:",inline"`
}

// AdmissionControlComponentSpec defines settings for the admission controller configuration.
type AdmissionControlComponentSpec struct {
	// Deprecated field. This field will be removed in a future release.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	ListenOnCreates *bool `json:"listenOnCreates,omitempty"`

	// Deprecated field. This field will be removed in a future release.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	ListenOnUpdates *bool `json:"listenOnUpdates,omitempty"`

	// Deprecated field. This field will be removed in a future release.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	ListenOnEvents *bool `json:"listenOnEvents,omitempty" deprecated:"true"`

	// Set to Disabled to disable policy enforcement for the admission controller. This is not recommended.
	// On new deployments starting with version 4.9, defaults to Enabled.
	// On old deployments, defaults to Enabled if at least one of listenOnCreates or listenOnUpdates is true.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1
	Enforcement *PolicyEnforcement `json:"enforcement,omitempty"`

	// Deprecated field. This field will be removed in a future release.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	ContactImageScanners *ImageScanPolicy `json:"contactImageScanners,omitempty" deprecated:"true"`

	// Deprecated field. This field will be removed in a future release.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	TimeoutSeconds *int32 `json:"timeoutSeconds,omitempty" deprecated:"true"`

	// Enables teams to bypass admission control in a monitored manner in the event of an emergency.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=2
	// The default is: BreakGlassAnnotation.
	Bypass *BypassPolicy `json:"bypass,omitempty"`

	// If set to "Fail", the admission controller's webhooks are configured to fail-closed in case admission controller
	// fails to respond in time. A failure policy "Ignore" configures the webhooks to fail-open.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=3
	// The default is: Ignore.
	FailurePolicy *FailurePolicy `json:"failurePolicy,omitempty"`

	// Settings pertaining to the Admission Controller running on a Secured Cluster.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=4
	DeploymentSpec `json:",inline"`

	// The number of replicas of the admission control pod.
	//+kubebuilder:validation:Minimum=1
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Replicas",order=5
	// The default is: 3.
	Replicas *int32 `json:"replicas,omitempty"`
}

// PolicyEnforcement defines whether policy enforcement is enabled or disabled.
// +kubebuilder:validation:Enum=Enabled;Disabled
type PolicyEnforcement string

const (
	// PolicyEnforcementEnabled means: policy enforcement is enabled.
	PolicyEnforcementEnabled PolicyEnforcement = "Enabled"
	// PolicyEnforcementDisabled means: policy enforcement is disabled.
	PolicyEnforcementDisabled PolicyEnforcement = "Disabled"
)

// ImageScanPolicy defines whether images should be scanned at admission control time.
// +kubebuilder:validation:Enum=ScanIfMissing;DoNotScanInline
type ImageScanPolicy string

const (
	// ScanIfMissing means that images which do not have a known scan result should be scanned as part of an admission request.
	ScanIfMissing ImageScanPolicy = "ScanIfMissing"
	// DoNotScanInline means that images which do not have a known scan result will not be scanned when processing an admission request.
	DoNotScanInline ImageScanPolicy = "DoNotScanInline"
)

// Pointer returns the given ImageScanPolicy as a pointer, needed in k8s resource structs.
func (p ImageScanPolicy) Pointer() *ImageScanPolicy {
	return &p
}

// BypassPolicy defines whether admission controller can be bypassed.
// +kubebuilder:validation:Enum=BreakGlassAnnotation;Disabled
type BypassPolicy string

const (
	// BypassBreakGlassAnnotation means that the admission controller can be bypassed by adding an admission.stackrox.io/break-glass annotation to a resource.
	// Bypassing the admission controller triggers a policy violation which includes deployment details.
	// We recommend providing an issue-tracker link or some other reference as the value of this annotation so that others can understand why you bypassed the admission controller.
	BypassBreakGlassAnnotation BypassPolicy = "BreakGlassAnnotation"
	// BypassDisabled means that the admission controller cannot be bypassed.
	BypassDisabled BypassPolicy = "Disabled"
)

// Pointer returns the given BypassPolicy as a pointer, needed in k8s resource structs.
func (p BypassPolicy) Pointer() *BypassPolicy {
	return &p
}

// FailurePolicy defines the failure policy for the admission controller webhooks, i.e. if a webhook request
// shall fail in case the webhook does not respond in time (fail-closed) or if the request shall be allowed
// in such a scenario (fail-open).
// +kubebuilder:validation:Enum=Ignore;Fail
type FailurePolicy string

const (
	// FailurePolicyFail instructs the admission controller's webhooks to fail-closed.
	FailurePolicyFail FailurePolicy = "Fail"
	// FailurePolicyIgnore instructs the admission controller's webhooks to fail-open.
	FailurePolicyIgnore FailurePolicy = "Ignore"
)

// PerNodeSpec declares configuration settings for components which are deployed to all nodes.
type PerNodeSpec struct {
	// Settings for the Collector container, which is responsible for collecting process and networking
	// activity at the host level.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1,displayName="Collector Settings"
	Collector *CollectorContainerSpec `json:"collector,omitempty"`

	// Settings for the Compliance container, which is responsible for checking host-level configurations.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=2,displayName="Compliance Settings"
	Compliance *ContainerSpec `json:"compliance,omitempty"`

	// Settings for the Node-Inventory container, which is responsible for scanning the Nodes' filesystem.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=3,displayName="Node Scanning Settings"
	NodeInventory *ContainerSpec `json:"nodeInventory,omitempty"`

	// Settings for the Sensitive File Activity container, which is responsible for file activity monitoring on the Node.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=4,displayName="SFA"
	SFA *SFAContainerSpec `json:"sfa,omitempty"`

	// To ensure comprehensive monitoring of your cluster activity, Red Hat Advanced Cluster Security
	// will run services on every node in the cluster, including tainted nodes by default. If you do
	// not want this behavior, please select 'AvoidTaints' here.
	// The default is: TolerateTaints.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=5
	TaintToleration *TaintTolerationPolicy `json:"taintToleration,omitempty"`

	// HostAliases allows configuring additional hostnames to resolve in the pod's hosts file.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=6,displayName="Host Aliases"
	HostAliases []corev1.HostAlias `json:"hostAliases,omitempty"`
}

// CollectionMethod defines the method of collection used by collector. Options are 'EBPF', 'CORE_BPF', 'NoCollection', or 'KernelModule'. Note that the collection method will be switched to CORE_BPF if KernelModule or EBPF is used.
// +kubebuilder:validation:Enum=EBPF;CORE_BPF;NoCollection;KernelModule
type CollectionMethod string

const (
	// CollectionEBPF means: use EBPF collection.
	CollectionEBPF CollectionMethod = "EBPF"
	// CollectionCOREBPF means: use CORE_BPF collection.
	CollectionCOREBPF CollectionMethod = "CORE_BPF"
	// CollectionNone means: NO_COLLECTION.
	CollectionNone CollectionMethod = "NoCollection"
	// CollectionKernelModule means: use KERNEL_MODULE collection.
	CollectionKernelModule CollectionMethod = "KernelModule"
)

// Pointer returns the given CollectionMethod as a pointer, needed in k8s resource structs.
func (c CollectionMethod) Pointer() *CollectionMethod {
	return &c
}

// AuditLogsSpec configures settings related to audit log ingestion.
type AuditLogsSpec struct {
	// Whether collection of Kubernetes audit logs should be enabled or disabled. Currently, this is only
	// supported on OpenShift 4, and trying to enable it on non-OpenShift 4 clusters will result in an error.
	// Use the 'Auto' setting to enable it on compatible environments, and disable it elsewhere.
	// The default is: Auto.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1
	Collection *AuditLogsCollectionSetting `json:"collection,omitempty"`
}

// AuditLogsCollectionSetting determines if audit log collection is enabled.
// +kubebuilder:validation:Enum=Auto;Disabled;Enabled
type AuditLogsCollectionSetting string

const (
	// AuditLogsCollectionAuto means to configure audit logs collection according to the environment (enable on
	// OpenShift 4.x, disable on all other environments).
	AuditLogsCollectionAuto AuditLogsCollectionSetting = "Auto"
	// AuditLogsCollectionDisabled means to disable audit logs collection.
	AuditLogsCollectionDisabled AuditLogsCollectionSetting = "Disabled"
	// AuditLogsCollectionEnabled means to enable audit logs collection.
	AuditLogsCollectionEnabled AuditLogsCollectionSetting = "Enabled"
)

// Pointer returns a pointer with the given value.
func (s AuditLogsCollectionSetting) Pointer() *AuditLogsCollectionSetting {
	return &s
}

// TaintTolerationPolicy is a type for values of spec.collector.taintToleration
// +kubebuilder:validation:Enum=TolerateTaints;AvoidTaints
type TaintTolerationPolicy string

const (
	// TaintTolerate means tolerations are applied to collector, and the collector pods can schedule onto all nodes with taints.
	TaintTolerate TaintTolerationPolicy = "TolerateTaints"
	// TaintAvoid means no tolerations are applied, and the collector pods won't schedule onto nodes with taints.
	TaintAvoid TaintTolerationPolicy = "AvoidTaints"
)

// Pointer returns the given TaintTolerationPolicy as a pointer, needed in k8s resource structs.
func (t TaintTolerationPolicy) Pointer() *TaintTolerationPolicy {
	return &t
}

// CollectorContainerSpec defines settings for the collector container.
type CollectorContainerSpec struct {
	// The method for system-level data collection. CORE_BPF is recommended.
	// If you select "NoCollection", you will not be able to see any information about network activity
	// and process executions. The remaining settings in this section will not have any effect.
	// The value is a subject of conversion by the operator if needed, e.g. to
	// remove deprecated methods.
	// The default is: CORE_BPF.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:select:CORE_BPF", "urn:alm:descriptor:com.tectonic.ui:select:NoCollection"}
	Collection *CollectionMethod `json:"collection,omitempty"`

	// Obsolete field.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=2,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	ImageFlavor *CollectorImageFlavor `json:"imageFlavor,omitempty"`

	// Obsolete field. This field will be removed in a future release.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=3,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	ForceCollection *bool `json:"forceCollection,omitempty"`

	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=4
	ContainerSpec `json:",inline"`
}

// SFAContainerSpec defines settings for the Sensitive File Activity agent container.
type SFAContainerSpec struct {
	// Specifies whether Sensitive File Activity agent is deployed.
	// The default is: Disabled.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1,displayName="SFA Agent"
	Agent *DeploySFAAgent `json:"agent,omitempty"`

	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=2
	ContainerSpec `json:",inline"`
}

// DeploySFAAgent is a type for values of spec.perNode.sfa.agent
// +kubebuilder:validation:Enum=Enabled;Disabled
type DeploySFAAgent string

const (
	SFAAgentEnabled  DeploySFAAgent = "Enabled"
	SFAAgentDisabled DeploySFAAgent = "Disabled"
)

// Pointer returns the given DeploySFAAgent value as a pointer, needed in k8s resource structs.
func (v DeploySFAAgent) Pointer() *DeploySFAAgent {
	return &v
}

// ContainerSpec defines container settings.
type ContainerSpec struct {
	// Allows overriding the default resource settings for this component. Please consult the documentation
	// for an overview of default resource requirements and a sizing guide.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:resourceRequirements"},order=100
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

// CollectorImageFlavor is a type for values of spec.collector.collector.imageFlavor
// +kubebuilder:validation:Enum=Regular;Slim
type CollectorImageFlavor string

const (
	// ImageFlavorRegular means to use regular collector images.
	ImageFlavorRegular CollectorImageFlavor = "Regular"
	// ImageFlavorSlim means to use slim collector images.
	ImageFlavorSlim CollectorImageFlavor = "Slim"
)

// Note the following struct should mostly match ScannerComponentSpec for the Central's type. Different Scanner
// types struct are maintained because of UI exposed documentation differences.

// LocalScannerComponentSpec defines settings for the "scanner" component.
type LocalScannerComponentSpec struct {
	// If you do not want to deploy the Red Hat Advanced Cluster Security Scanner, you can disable it here
	// (not recommended).
	// If you do so, all the settings in this section will have no effect.
	// The default is: AutoSense.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Scanner Component",order=1
	ScannerComponent *LocalScannerComponentPolicy `json:"scannerComponent,omitempty"`

	// Settings pertaining to the analyzer deployment, such as for autoscaling.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=2
	Analyzer *ScannerAnalyzerComponent `json:"analyzer,omitempty"`

	// Settings pertaining to the database used by the Red Hat Advanced Cluster Security Scanner.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=3,displayName="DB"
	DB *DeploymentSpec `json:"db,omitempty"`
}

// LocalScannerV4ComponentSpec defines settings for the "Scanner V4" component in SecuredClusters
type LocalScannerV4ComponentSpec struct {
	// If you want to enable the Scanner V4 component set this to "AutoSense"
	// If this field is not specified or set to "Default", the following defaulting takes place:
	// * for new installations, Scanner V4 is enabled starting with ACS 4.8;
	// * for upgrades to 4.8 from previous releases, Scanner V4 is disabled.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Scanner V4 component",order=1
	ScannerComponent *LocalScannerV4ComponentPolicy `json:"scannerComponent,omitempty"`

	// Settings pertaining to the indexer deployment.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=2,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:fieldDependency:.scannerComponent:AutoSense"}
	Indexer *ScannerV4Component `json:"indexer,omitempty"`

	// Settings pertaining to the DB deployment.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=3,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:fieldDependency:.scannerComponent:AutoSense"}
	DB *ScannerV4DB `json:"db,omitempty"`

	// Configures monitoring endpoint for Scanner V4. The monitoring endpoint
	// allows other services to collect metrics from Scanner V4, provided in
	// Prometheus compatible format.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=4,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:fieldDependency:.scannerComponent:AutoSense"}
	Monitoring *Monitoring `json:"monitoring,omitempty"`
}

// LocalScannerComponentPolicy is a type for values of spec.scanner.scannerComponent.
// +kubebuilder:validation:Enum=AutoSense;Disabled
type LocalScannerComponentPolicy string

const (
	// LocalScannerComponentAutoSense means that scanner should be installed,
	// unless there is a Central resource in the same namespace.
	// In that case typically a central scanner will be deployed as a component of Central.
	LocalScannerComponentAutoSense LocalScannerComponentPolicy = "AutoSense"
	// LocalScannerComponentDisabled means that scanner should not be installed.
	LocalScannerComponentDisabled LocalScannerComponentPolicy = "Disabled"
)

// Pointer returns the pointer of the policy.
func (l LocalScannerComponentPolicy) Pointer() *LocalScannerComponentPolicy {
	return &l
}

// LocalScannerV4ComponentPolicy is a type for values of spec.scannerV4.scannerComponent
// +kubebuilder:validation:Enum=Default;AutoSense;Disabled
type LocalScannerV4ComponentPolicy string

const (
	// LocalScannerV4ComponentDefault means that local Scanner V4 will use the default semantics
	// to determine whether scannerV4 components should be used.
	// Currently this defaults to "Disabled" semantics.
	// TODO: change default to "AutoSense" semantics with version 4.5
	LocalScannerV4ComponentDefault LocalScannerV4ComponentPolicy = "Default"
	// LocalScannerV4ComponentAutoSense means that Scanner V4 should be installed,
	// unless there is a Central resource in the same namespace.
	// In that case typically a central Scanner V4 will be deployed as a component of Central
	LocalScannerV4ComponentAutoSense LocalScannerV4ComponentPolicy = "AutoSense"
	// LocalScannerV4ComponentDisabled means that scanner should not be installed.
	LocalScannerV4ComponentDisabled LocalScannerV4ComponentPolicy = "Disabled"
)

// Pointer returns the pointer of the policy.
func (l LocalScannerV4ComponentPolicy) Pointer() *LocalScannerV4ComponentPolicy {
	return &l
}

// -------------------------------------------------------------
// Status

// SecuredClusterStatus defines the observed state of SecuredCluster
type SecuredClusterStatus struct {
	Conditions      []StackRoxCondition `json:"conditions"`
	DeployedRelease *StackRoxRelease    `json:"deployedRelease,omitempty"`

	// The deployed version of the product.
	//+operator-sdk:csv:customresourcedefinitions:type=status,order=1
	ProductVersion string `json:"productVersion,omitempty"`

	// The assigned cluster name per the spec. This cannot be changed afterwards. If you need to change the
	// cluster name, please delete and recreate this resource.
	//+operator-sdk:csv:customresourcedefinitions:type=status,displayName="Cluster Name",order=2
	ClusterName string `json:"clusterName,omitempty"`

	// ObservedGeneration is the generation most recently observed by the controller.
	//+operator-sdk:csv:customresourcedefinitions:type=status,order=3
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+operator-sdk:csv:customresourcedefinitions:resources={{Deployment,v1,""},{DaemonSet,v1,""}}
//+kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.status.productVersion`
//+kubebuilder:printcolumn:name="Message",type=string,JSONPath=`.status.conditions[?(@.type=="Deployed")].message`
//+kubebuilder:printcolumn:name="Progressing",type=string,JSONPath=`.status.conditions[?(@.type=="Progressing")].status`
//+kubebuilder:printcolumn:name="Available",type=string,JSONPath=`.status.conditions[?(@.type=="Available")].status`
//+genclient

// SecuredCluster is the configuration template for the secured cluster services. These include Sensor, which is
// responsible for the connection to Central, and Collector, which performs host-level collection of process and
// network events.<p>
// **Important:** Please see the _Installation Prerequisites_ on the main RHACS operator page before deploying, or
// consult the RHACS documentation on creating cluster init bundles.
type SecuredCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecuredClusterSpec   `json:"spec,omitempty"`
	Status SecuredClusterStatus `json:"status,omitempty"`

	// This field will never be serialized, it is used for attaching defaulting decisions to a SecuredCluster struct during reconciliation.
	Defaults SecuredClusterSpec `json:"-"`
}

//+kubebuilder:object:root=true

// SecuredClusterList contains a list of SecuredCluster
type SecuredClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecuredCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecuredCluster{}, &SecuredClusterList{})
}

var (
	// SecuredClusterGVK is the GVK for the SecuredCluster type.
	SecuredClusterGVK = GroupVersion.WithKind("SecuredCluster")

	LocalScannerV4AutoSense = LocalScannerV4ComponentAutoSense
	LocalScannerV4Disabled  = LocalScannerV4ComponentDisabled
)

// GetCondition returns a specific condition by type, or nil if not found.
func (c *SecuredCluster) GetCondition(condType ConditionType) *StackRoxCondition {
	return getCondition(c.Status.Conditions, condType)
}

// SetCondition updates or adds a condition. Returns true if the condition changed.
func (c *SecuredCluster) SetCondition(updatedCond StackRoxCondition) bool {
	var updated bool
	c.Status.Conditions, updated = updateCondition(c.Status.Conditions, updatedCond)
	return updated
}

// GetGeneration returns the metadata.generation of the Central resource.
func (c *SecuredCluster) GetGeneration() int64 {
	return c.ObjectMeta.GetGeneration()
}

// GetObservedGeneration returns the observedGeneration of the Central status sub-resource.
func (c *SecuredCluster) GetObservedGeneration() int64 {
	return c.Status.ObservedGeneration
}
