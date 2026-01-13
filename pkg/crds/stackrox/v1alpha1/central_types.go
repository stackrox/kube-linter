/*
Copyright 2021.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

// Important: Run "make generate manifests" to regenerate code after modifying this file

// -------------------------------------------------------------
// Spec

// CentralSpec defines the desired state of Central
type CentralSpec struct {
	// Settings for the Central component, which is responsible for all user interaction.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1,displayName="Central Component Settings"
	Central *CentralComponentSpec `json:"central,omitempty"`

	// Settings for the Scanner component, which is responsible for vulnerability scanning of container
	// images.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=2,displayName="Scanner Component Settings"
	Scanner *ScannerComponentSpec `json:"scanner,omitempty"`

	// Settings for the Scanner V4 component, which can run in addition to the previously existing Scanner components
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=3,displayName="Scanner V4 Component Settings"
	ScannerV4 *ScannerV4Spec `json:"scannerV4,omitempty"`

	// Settings related to outgoing network traffic.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=4
	Egress *Egress `json:"egress,omitempty"`

	// Settings related to Transport Layer Security, such as Certificate Authorities.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=5
	TLS *TLSConfig `json:"tls,omitempty"`

	// Additional image pull secrets to be taken into account for pulling images.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Image Pull Secrets",order=6,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:advanced"}
	ImagePullSecrets []LocalSecretReference `json:"imagePullSecrets,omitempty"`

	// Customizations to apply on all Central Services components.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName=Customizations,order=7,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:advanced"}
	Customize *CustomizeSpec `json:"customize,omitempty"`

	// Deprecated field. This field will be removed in a future release.
	// Miscellaneous settings.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName=Miscellaneous,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	Misc *MiscSpec `json:"misc,omitempty" deprecated:"true"`

	// Overlays
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName=Overlays,order=8,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	Overlays []*K8sObjectOverlay `json:"overlays,omitempty"`

	// Monitoring configuration.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=9,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:advanced"}
	Monitoring *GlobalMonitoring `json:"monitoring,omitempty"`

	// Network configuration.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName=Network,order=10,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:advanced"}
	Network *GlobalNetworkSpec `json:"network,omitempty"`

	// Config-as-Code configuration.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName=Config-as-Code,order=11,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:advanced"}
	ConfigAsCode *ConfigAsCodeSpec `json:"configAsCode,omitempty"`
}

// Egress defines settings related to outgoing network traffic.
type Egress struct {
	// Configures whether Red Hat Advanced Cluster Security should run in online or offline (disconnected) mode.
	// In offline mode, automatic updates of vulnerability definitions and kernel modules are disabled.
	// The default is: Online.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName=Connectivity Policy,order=1
	ConnectivityPolicy *ConnectivityPolicy `json:"connectivityPolicy,omitempty"`
}

// ConnectivityPolicy is a type for values of spec.egress.connectivityPolicy.
// +kubebuilder:validation:Enum=Online;Offline
type ConnectivityPolicy string

const (
	// ConnectivityOnline means that Central is allowed to make outbound connections to the Internet.
	ConnectivityOnline ConnectivityPolicy = "Online"
	// ConnectivityOffline means that Central must not make outbound connections to the Internet.
	ConnectivityOffline ConnectivityPolicy = "Offline"
)

func (p ConnectivityPolicy) Pointer() *ConnectivityPolicy {
	return &p
}

// CentralComponentSpec defines settings for the "central" component.
type CentralComponentSpec struct {
	// Specify a secret that contains the administrator password in the "password" data item.
	// If omitted, the operator will auto-generate a password and store it in the "password" item
	// in the "central-htpasswd" secret.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Administrator Password",order=1
	AdminPasswordSecret *LocalSecretReference `json:"adminPasswordSecret,omitempty"`

	// Disable admin password generation. Do not use this for first-time installations,
	// as you will have no way to perform initial setup and configuration of alternative authentication methods.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	AdminPasswordGenerationDisabled *bool `json:"adminPasswordGenerationDisabled,omitempty"`

	// Here you can configure if you want to expose central through a node port, a load balancer, or an OpenShift
	// route.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=2
	Exposure *Exposure `json:"exposure,omitempty"`

	// By default, Central will only serve an internal TLS certificate, which means that you will
	// need to handle TLS termination at the ingress or load balancer level.
	// If you want to terminate TLS in Central and serve a custom server certificate, you can specify
	// a secret containing the certificate and private key here.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="User-facing TLS certificate secret",order=3
	DefaultTLSSecret *LocalSecretReference `json:"defaultTLSSecret,omitempty"`

	// Configures monitoring endpoint for Central. The monitoring endpoint
	// allows other services to collect metrics from Central, provided in
	// Prometheus compatible format.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=4
	Monitoring *Monitoring `json:"monitoring,omitempty"`

	// Unused field. This field exists solely for backward compatibility starting from version v4.6.0.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	Persistence *ObsoletePersistence `json:"persistence,omitempty"`

	// Settings for Central DB, which is responsible for data persistence.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=5,displayName="Central DB Settings"
	DB *CentralDBSpec `json:"db,omitempty"`

	// Configures telemetry settings for Central. If enabled, Central transmits telemetry and diagnostic
	// data to a remote storage backend.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=6,displayName="Telemetry",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	Telemetry *Telemetry `json:"telemetry,omitempty"`

	// Configures resources within Central in a declarative manner.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=7,displayName="Declarative Configuration",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	DeclarativeConfiguration *DeclarativeConfiguration `json:"declarativeConfiguration,omitempty"`

	// Configures the encryption of notifier secrets stored in the Central DB.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=8,displayName="Notifier Secrets Encryption",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	NotifierSecretsEncryption *NotifierSecretsEncryption `json:"notifierSecretsEncryption,omitempty"`

	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=99
	DeploymentSpec `json:",inline"`
}

// GetDB returns Central's db config
func (c *CentralComponentSpec) GetDB() *CentralDBSpec {
	if c == nil {
		return nil
	}
	return c.DB
}

// GetAdminPasswordSecret provides a way to retrieve the admin password that is safe to use on a nil receiver object.
func (c *CentralComponentSpec) GetAdminPasswordSecret() *LocalSecretReference {
	if c == nil {
		return nil
	}
	return c.AdminPasswordSecret
}

// GetAdminPasswordGenerationDisabled provides a way to retrieve the AdminPasswordEnabled setting that is safe to use on a nil receiver object.
func (c *CentralComponentSpec) GetAdminPasswordGenerationDisabled() bool {
	if c == nil {
		return false
	}
	return pointer.BoolDeref(c.AdminPasswordGenerationDisabled, false)
}

// ShouldManageDB returns true if central DB should be managed by the Operator.
func (c *CentralComponentSpec) ShouldManageDB() bool {
	return c == nil || c.DB == nil || c.DB.ConnectionStringOverride == nil
}

// GetNotifierSecretsEncryptionEnabled provides a way to retrieve the NotifierSecretsEncryption.Enabled setting that is safe to use on a nil receiver object.
func (c *CentralComponentSpec) GetNotifierSecretsEncryptionEnabled() bool {
	if c == nil || c.NotifierSecretsEncryption == nil {
		return false
	}
	return pointer.BoolDeref(c.NotifierSecretsEncryption.Enabled, false)
}

// DeclarativeConfiguration defines settings for adding resources in a declarative manner.
type DeclarativeConfiguration struct {
	// List of config maps containing declarative configuration.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Config maps containing declarative configuration"
	ConfigMaps []LocalConfigMapReference `json:"configMaps,omitempty"`

	// List of secrets containing declarative configuration.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Secrets containing declarative configuration"
	Secrets []LocalSecretReference `json:"secrets,omitempty"`
}

// NotifierSecretsEncryption defines settings for encrypting notifier secrets in the Central DB.
type NotifierSecretsEncryption struct {
	// Enables the encryption of notifier secrets stored in the Central DB.
	// The default is: false.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1
	Enabled *bool `json:"enabled,omitempty"`
}

// CentralDBSpec defines settings for the "central db" component.
// TODO(ROX-14395): drop `IsEnabled` field when bumping API version.
// isEnabled is effectively no-op starting from the version 3.74.0. It should be removed when we
// bump API version of ACS custom resources. Removing it before is unsafe and may break compatibility.
type CentralDBSpec struct {
	// Obsolete field.
	// This field will be removed in a future release.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	IsEnabled *CentralDBEnabled `json:"isEnabled,omitempty"`

	// Specify a secret that contains the password in the "password" data item. This can only be used when
	// specifying a connection string manually.
	// When omitted, the operator will auto-generate a DB password and store it in the "password" item
	// in the "central-db-password" secret.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Administrator Password",order=1
	PasswordSecret *LocalSecretReference `json:"passwordSecret,omitempty"`

	// Specify a connection string that corresponds to a database managed elsewhere. If set, the operator will not manage the Central DB.
	// When using this option, you must explicitly set a password secret; automatically generating a password will not
	// be supported.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=2,displayName="Connection String"
	ConnectionStringOverride *string `json:"connectionString,omitempty"`

	// Configures how Central DB should store its persistent data. You can choose between using a persistent
	// volume claim (recommended default), and a host path.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=3
	Persistence *DBPersistence `json:"persistence,omitempty"`

	// Config map containing postgresql.conf and pg_hba.conf that will be used if modifications need to be applied.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=4,displayName="Config map that will override postgresql.conf and pg_hba.conf"
	ConfigOverride *LocalConfigMapReference `json:"configOverride,omitempty"`

	// Configures the database connection pool size.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=5,displayName="Database Connection Pool Size Settings"
	ConnectionPoolSize *DBConnectionPoolSize `json:"connectionPoolSize,omitempty"`

	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=99
	DeploymentSpec `json:",inline"`
}

// CentralDBEnabled is a type for values of spec.central.db.enabled.
// +kubebuilder:validation:Enum=Default;Enabled
type CentralDBEnabled string

const (
	// CentralDBEnabledDefault configures the central to use PostgreSQL database.
	// Deprecated const.
	CentralDBEnabledDefault CentralDBEnabled = "Default"

	// CentralDBEnabledTrue configures the central to use a PostgreSQL database.
	// Deprecated const.
	CentralDBEnabledTrue CentralDBEnabled = "Enabled"
)

// DBConnectionPoolSize configures the database connection pool size.
type DBConnectionPoolSize struct {
	// Minimum number of connections in the connection pool.
	// The default is: 10.
	//+kubebuilder:validation:Minimum=1
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Minimum Connections"
	MinConnections *int32 `json:"minConnections"`

	// Maximum number of connections in the connection pool.
	// The default is: 90.
	//+kubebuilder:validation:Minimum=1
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Maximum Connections"
	MaxConnections *int32 `json:"maxConnections"`
}

// GetPasswordSecret provides a way to retrieve the admin password that is safe to use on a nil receiver object.
func (c *CentralDBSpec) GetPasswordSecret() *LocalSecretReference {
	if c == nil {
		return nil
	}
	return c.PasswordSecret
}

// GetPersistence returns the persistence for Central DB
func (c *CentralDBSpec) GetPersistence() *DBPersistence {
	if c == nil {
		return nil
	}
	return c.Persistence
}

// ObsoletePersistence contains obsolete persistence settings for central.
type ObsoletePersistence struct {
	// Obsolete unused field.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	PersistentVolumeClaim *ObsoletePersistentVolumeClaim `json:"persistentVolumeClaim,omitempty"`

	// Obsolete unused field.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	HostPath *HostPathSpec `json:"hostPath,omitempty"`
}

// HostPathSpec defines settings for host path config.
type HostPathSpec struct {
	// The path on the host running Central.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=99
	Path *string `json:"path,omitempty"`
}

// ObsoletePersistentVolumeClaim contains obsolete PVC-based persistence settings.
type ObsoletePersistentVolumeClaim struct {
	// Obsolete unused field.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	ClaimName *string `json:"claimName,omitempty"`

	// Obsolete unused field.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	Size *string `json:"size,omitempty"`

	// Obsolete unused field.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:hidden"}
	StorageClassName *string `json:"storageClassName,omitempty"`
}

// DBPersistence defines persistence settings for Central DB.
type DBPersistence struct {
	// Uses a Kubernetes persistent volume claim (PVC) to manage the storage location of persistent data.
	// Recommended for most users.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Persistent volume claim",order=1
	PersistentVolumeClaim *DBPersistentVolumeClaim `json:"persistentVolumeClaim,omitempty"`

	// Stores persistent data in a directory on the host. This is not recommended, and should only
	// be used together with a node selector (only available in YAML view).
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Host path",order=99
	HostPath *HostPathSpec `json:"hostPath,omitempty"`
}

// GetPersistentVolumeClaim returns the configured PVC
func (p *DBPersistence) GetPersistentVolumeClaim() *DBPersistentVolumeClaim {
	if p == nil {
		return nil
	}
	return p.PersistentVolumeClaim
}

// GetHostPath returns the configured host path
func (p *DBPersistence) GetHostPath() string {
	if p == nil {
		return ""
	}
	if p.HostPath == nil {
		return ""
	}

	return pointer.StringDeref(p.HostPath.Path, "")
}

// DBPersistentVolumeClaim defines PVC-based persistence settings for Central DB.
type DBPersistentVolumeClaim struct {
	// The name of the PVC to manage persistent data. If no PVC with the given name exists, it will be
	// created.
	// The default is: central-db.
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

// Exposure defines how Central is exposed.
type Exposure struct {
	// Expose Central through an OpenShift route.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1,displayName="Route"
	Route *ExposureRoute `json:"route,omitempty"`

	// Expose Central through a load balancer service.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=2,displayName="Load Balancer"
	LoadBalancer *ExposureLoadBalancer `json:"loadBalancer,omitempty"`

	// Expose Central through a node port.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=3,displayName="Node Port"
	NodePort *ExposureNodePort `json:"nodePort,omitempty"`
}

// ExposureLoadBalancer defines settings for exposing central via a LoadBalancer.
type ExposureLoadBalancer struct {
	// The default is: false.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1
	Enabled *bool `json:"enabled,omitempty"`

	// The default is: 443.
	//+kubebuilder:validation:Minimum=1
	//+kubebuilder:validation:Maximum=65535
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=2,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:fieldDependency:.enabled:true"}
	Port *int32 `json:"port,omitempty"`

	// If you have a static IP address reserved for your load balancer, you can enter it here.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=3,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:fieldDependency:.enabled:true"}
	IP *string `json:"ip,omitempty"`
}

// ExposureNodePort defines settings for exposing central via a NodePort.
type ExposureNodePort struct {
	// The default is: false.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1
	Enabled *bool `json:"enabled,omitempty"`

	// Use this to specify an explicit node port. Most users should leave this empty.
	//+kubebuilder:validation:Minimum=1
	//+kubebuilder:validation:Maximum=65535
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=2,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:fieldDependency:.enabled:true"}
	Port *int32 `json:"port,omitempty"`
}

// ExposureRoute defines settings for exposing Central via a Route.
type ExposureRoute struct {
	// Expose Central with a passthrough route.
	// The default is: false.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1
	Enabled *bool `json:"enabled,omitempty"`

	// Specify a custom hostname for the Central route.
	// If unspecified, an appropriate default value will be automatically chosen by the OpenShift route operator.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=2
	Host *string `json:"host,omitempty"`

	// Set up a Central route with reencrypt TLS termination.
	// For reencrypt routes, the request is terminated on the OpenShift router with a custom certificate.
	// The request is then reencrypted by the OpenShift router and sent to Central.
	// [user] --TLS--> [OpenShift router] --TLS--> [Central]
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=3,displayName="Re-Encrypt Route"
	Reencrypt *ExposureRouteReencrypt `json:"reencrypt,omitempty"`
}

// ExposureRouteReencrypt defines settings for exposing Central via a reencrypt Route.
type ExposureRouteReencrypt struct {
	// Expose Central with a reencrypt route.
	// Should not be used for sensor communication.
	// The default is: false.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1
	Enabled *bool `json:"enabled,omitempty"`

	// Specify a custom hostname for the Central reencrypt route.
	// If unspecified, an appropriate default value will be automatically chosen by the OpenShift route operator.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=2
	Host *string `json:"host,omitempty"`

	// TLS settings for exposing Central via a reencrypt Route.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=3,displayName="TLS Settings"
	TLS *ExposureRouteReencryptTLS `json:"tls,omitempty"`
}

// ExposureRouteReencryptTLS defines TLS settings for exposing Central via a reencrypt Route.
type ExposureRouteReencryptTLS struct {
	// The PEM encoded certificate chain that may be used to establish a complete chain of trust.
	// Defaults to the OpenShift certificate authority.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1,displayName="CA Certificate"
	CaCertificate *string `json:"caCertificate,omitempty"`

	// The PEM encoded certificate that is served on the route. Must be a single serving
	// certificate instead of a certificate chain.
	// Defaults to a certificate signed by the OpenShift certificate authority.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=2,displayName="Certificate"
	Certificate *string `json:"certificate,omitempty"`

	// The CA certificate of the final destination, i.e. of Central.
	// Used by the OpenShift router for health checks on the secure connection.
	// Defaults to the Central certificate authority.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=3,displayName="Destination CA Certificate"
	DestinationCACertificate *string `json:"destinationCACertificate,omitempty"`

	// The PEM encoded private key of the certificate that is served on the route.
	// Defaults to a certificate signed by the OpenShift certificate authority.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=4,displayName="Private Key"
	Key *string `json:"key,omitempty"`
}

// Telemetry defines telemetry settings for Central.
type Telemetry struct {
	// Specifies whether Telemetry is enabled.
	// The default is: true.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	Enabled *bool `json:"enabled,omitempty"`

	// Defines the telemetry storage backend for Central.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=2,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:fieldDependency:.enabled:true"}
	Storage *TelemetryStorage `json:"storage,omitempty"`
}

// TelemetryStorage defines the telemetry storage backend for Central.
type TelemetryStorage struct {
	// Storage API endpoint.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1
	Endpoint *string `json:"endpoint,omitempty"`

	// Storage API key. If not set, telemetry is disabled.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=2
	Key *string `json:"key,omitempty"`
}

// Note the following struct should mostly match LocalScannerComponentSpec for the SecuredCluster type. Different Scanner
// types struct are maintained because of UI exposed documentation differences.

// ScannerComponentSpec defines settings for the central "scanner" component.
type ScannerComponentSpec struct {
	// If you do not want to deploy the Red Hat Advanced Cluster Security Scanner, you can disable it here
	// (not recommended). By default, the scanner is enabled.
	// If you do so, all the settings in this section will have no effect.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Scanner Component",order=1
	ScannerComponent *ScannerComponentPolicy `json:"scannerComponent,omitempty"`

	// Settings pertaining to the analyzer deployment, such as for autoscaling.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=2,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:fieldDependency:.scannerComponent:Enabled"}
	Analyzer *ScannerAnalyzerComponent `json:"analyzer,omitempty"`

	// Settings pertaining to the database used by the Red Hat Advanced Cluster Security Scanner.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=3,displayName="DB",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:fieldDependency:.scannerComponent:Enabled"}
	DB *DeploymentSpec `json:"db,omitempty"`

	// Configures monitoring endpoint for Scanner. The monitoring endpoint
	// allows other services to collect metrics from Scanner, provided in
	// Prometheus compatible format.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=4
	Monitoring *Monitoring `json:"monitoring,omitempty"`
}

// ScannerV4Spec defines settings for the central "Scanner V4" component.
type ScannerV4Spec struct {
	// Can be specified as "Enabled" or "Disabled".
	// If this field is not specified, the following defaulting takes place:
	// * for new installations, Scanner V4 is enabled starting with ACS 4.8;
	// * for upgrades to 4.8 from previous releases, Scanner V4 is disabled.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1,displayName="Scanner V4 component"
	ScannerComponent *ScannerV4ComponentPolicy `json:"scannerComponent,omitempty"`

	// Settings pertaining to the indexer deployment.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=2,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:fieldDependency:.scannerComponent:Enabled"}
	Indexer *ScannerV4Component `json:"indexer,omitempty"`

	// Settings pertaining to the matcher deployment.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=3,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:fieldDependency:.scannerComponent:Enabled"}
	Matcher *ScannerV4Component `json:"matcher,omitempty"`

	// Settings pertaining to the DB deployment.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=4,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:fieldDependency:.scannerComponent:Enabled"}
	DB *ScannerV4DB `json:"db,omitempty"`

	// Configures monitoring endpoint for Scanner V4. The monitoring endpoint
	// allows other services to collect metrics from Scanner V4, provided in
	// Prometheus compatible format.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=5,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:fieldDependency:.scannerComponent:Enabled"}
	Monitoring *Monitoring `json:"monitoring,omitempty"`
}

// IsEnabled checks whether Scanner V4 is enabled. This method is safe to be used with nil receivers.
func (s *ScannerV4Spec) IsEnabled() bool {
	if s == nil || s.ScannerComponent == nil {
		return false // disabled by default
	}
	return *s.ScannerComponent == ScannerV4ComponentEnabled
}

// GetAnalyzer returns the analyzer component even if receiver is nil
func (s *ScannerComponentSpec) GetAnalyzer() *ScannerAnalyzerComponent {
	if s == nil {
		return nil
	}
	return s.Analyzer
}

// IsEnabled checks whether scanner is enabled. This method is safe to be used with nil receivers.
func (s *ScannerComponentSpec) IsEnabled() bool {
	if s == nil || s.ScannerComponent == nil {
		return true // enabled by default
	}
	return *s.ScannerComponent == ScannerComponentEnabled
}

// ScannerComponentPolicy is a type for values of spec.scanner.scannerComponent.
// +kubebuilder:validation:Enum=Enabled;Disabled
type ScannerComponentPolicy string

const (
	// ScannerComponentEnabled means that scanner should be installed.
	ScannerComponentEnabled ScannerComponentPolicy = "Enabled"
	// ScannerComponentDisabled means that scanner should not be installed.
	ScannerComponentDisabled ScannerComponentPolicy = "Disabled"
)

// ScannerV4ComponentPolicy is a type for values of spec.scannerV4.scannerComponent
// +kubebuilder:validation:Enum=Default;Enabled;Disabled
type ScannerV4ComponentPolicy string

const (
	// Keep this for compatibility and potentially for reasoning about expected defaults.
	ScannerV4ComponentDefault ScannerV4ComponentPolicy = "Default"
	// ScannerV4ComponentEnabled explicitly enables the Scanner V4 component.
	ScannerV4ComponentEnabled ScannerV4ComponentPolicy = "Enabled"
	// ScannerV4ComponentDisabled explicitly disables the Scanner V4 component.
	ScannerV4ComponentDisabled ScannerV4ComponentPolicy = "Disabled"
)

type ConfigAsCodeSpec struct {
	// If you want to deploy the Config as Code component, set this to "Enabled"
	// The default is: Enabled.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=1,displayName="Config as Code component"
	ComponentPolicy *ConfigAsCodeComponentPolicy `json:"configAsCodeComponent,omitempty"`

	//+operator-sdk:csv:customresourcedefinitions:type=spec,order=99
	DeploymentSpec `json:",inline"`
}

// ConfigAsCodeComponentPolicy is a type for values of spec.configAsCode.configAsCodeComponent
// +kubebuilder:validation:Enum=Enabled;Disabled
type ConfigAsCodeComponentPolicy string

const (
	// ConfigAsCodeComponentEnabled explicitly enables the Config as Code component.
	ConfigAsCodeComponentEnabled ConfigAsCodeComponentPolicy = "Enabled"
	// ConfigAsCodeComponentDisabled explicitly disables the Config as Code component.
	ConfigAsCodeComponentDisabled ConfigAsCodeComponentPolicy = "Disabled"
)

// -------------------------------------------------------------
// Status

// CentralStatus defines the observed state of Central.
type CentralStatus struct {
	Conditions      []StackRoxCondition `json:"conditions"`
	DeployedRelease *StackRoxRelease    `json:"deployedRelease,omitempty"`

	// The deployed version of the product.
	//+operator-sdk:csv:customresourcedefinitions:type=status,order=1
	ProductVersion string `json:"productVersion,omitempty"`
	//+operator-sdk:csv:customresourcedefinitions:type=status,order=2
	Central *CentralComponentStatus `json:"central,omitempty"`

	// ObservedGeneration is the generation most recently observed by the controller.
	//+operator-sdk:csv:customresourcedefinitions:type=status,order=4
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// AdminPasswordStatus shows status related to the admin password.
type AdminPasswordStatus struct {
	// Info stores information on how to obtain the admin password.
	//+operator-sdk:csv:customresourcedefinitions:type=status,order=1,displayName="Admin Credentials Info"
	Info string `json:"info,omitempty"`

	// AdminPasswordSecretReference contains reference for the admin password
	//+operator-sdk:csv:customresourcedefinitions:type=status,order=2,displayName="Admin Password Secret Reference",xDescriptors={"urn:alm:descriptor:io.kubernetes:Secret"}
	SecretReference *string `json:"adminPasswordSecretReference,omitempty"`
}

// CentralComponentStatus describes status specific to the central component.
type CentralComponentStatus struct {
	// AdminPassword stores information related to the auto-generated admin password.
	AdminPassword *AdminPasswordStatus `json:"adminPassword,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+operator-sdk:csv:customresourcedefinitions:resources={{Deployment,v1,""},{Secret,v1,""},{Service,v1,""},{Route,v1,""}}
//+kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.status.productVersion`
//+kubebuilder:printcolumn:name="AdminPassword",type=string,JSONPath=`.status.central.adminPassword.adminPasswordSecretReference`
//+kubebuilder:printcolumn:name="Message",type=string,JSONPath=`.status.conditions[?(@.type=="Deployed")].message`
//+kubebuilder:printcolumn:name="Progressing",type=string,JSONPath=`.status.conditions[?(@.type=="Progressing")].status`
//+kubebuilder:printcolumn:name="Available",type=string,JSONPath=`.status.conditions[?(@.type=="Available")].status`
//+genclient

// Central is the configuration template for the central services. This includes the API server, persistent storage,
// and the web UI, as well as the image scanner.
type Central struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CentralSpec   `json:"spec,omitempty"`
	Status CentralStatus `json:"status,omitempty"`

	// This field will never be serialized, it is used for attaching defaulting decisions to a Central struct during reconciliation.
	Defaults CentralSpec `json:"-"`
}

// GetCondition returns a specific condition by type, or nil if not found.
func (c *Central) GetCondition(condType ConditionType) *StackRoxCondition {
	return getCondition(c.Status.Conditions, condType)
}

// SetCondition updates or adds a condition. Returns true if the condition changed.
func (c *Central) SetCondition(updatedCond StackRoxCondition) bool {
	var updated bool
	c.Status.Conditions, updated = updateCondition(c.Status.Conditions, updatedCond)
	return updated
}

// GetGeneration returns the metadata.generation of the Central resource.
func (c *Central) GetGeneration() int64 {
	return c.ObjectMeta.GetGeneration()
}

// GetObservedGeneration returns the observedGeneration of the Central status sub-resource.
func (c *Central) GetObservedGeneration() int64 {
	return c.Status.ObservedGeneration
}

//+kubebuilder:object:root=true

// CentralList contains a list of Central
type CentralList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Central `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Central{}, &CentralList{})
}

var (
	// CentralGVK is the GVK for the Central type.
	CentralGVK = GroupVersion.WithKind("Central")

	ScannerV4Default  = ScannerV4ComponentDefault
	ScannerV4Enabled  = ScannerV4ComponentEnabled
	ScannerV4Disabled = ScannerV4ComponentDisabled
)

// IsScannerEnabled returns true if scanner is enabled.
func (c *Central) IsScannerEnabled() bool {
	return c.Spec.Scanner.IsEnabled()
}

// IsScannerV4Enabled returns true if Scanner V4 is enabled.
func (c *Central) IsScannerV4Enabled() bool {
	return c.Spec.ScannerV4.IsEnabled()
}
