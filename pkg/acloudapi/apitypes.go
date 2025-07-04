package acloudapi

import (
	"fmt"
	"time"
)

type CloudProviderType string

const (
	CloudProviderEmpty        CloudProviderType = ""
	CloudProviderNone         CloudProviderType = "none"
	CloudProviderByon         CloudProviderType = "byon"
	CloudProviderAWS          CloudProviderType = "aws"
	CloudProviderAzure        CloudProviderType = "azure"
	CloudProviderMultipass    CloudProviderType = "multipass"
	CloudProviderOpenstack    CloudProviderType = "openstack"
	CloudProviderDigitalOcean CloudProviderType = "digitalocean"
	CloudProviderHetzner      CloudProviderType = "hetzner"
	CloudProviderVSphere      CloudProviderType = "vsphere"
)

// Cluster represents the Cluster resource in the Avisi Cloud API
type Cluster struct {
	Name                    string                     `json:"name" yaml:"Name"`
	Identity                string                     `json:"identity" yaml:"Identity"`
	EnvironmentIdentity     string                     `json:"environmentIdentity" yaml:"EnvironmentIdentity"`
	EnvironmentSlug         string                     `json:"environmentSlug" yaml:"EnvironmentSlug"`
	CustomerIdentity        string                     `json:"customerIdentity" yaml:"CustomerIdentity"`
	CustomerSlug            string                     `json:"customerSlug" yaml:"CustomerSlug"`
	Slug                    string                     `json:"slug" yaml:"Slug"`
	CNI                     string                     `json:"cni" yaml:"CNI"`
	Description             string                     `json:"description" yaml:"Description,omitempty"`
	CloudProvider           string                     `json:"cloudProvider" yaml:"CloudProvider"`
	CloudAccount            *CloudAccountReference     `json:"cloudAccount" yaml:"CloudAccount,omitempty"`
	CloudCredentials        *CloudCredentialsReference `json:"cloudCredentials" yaml:"CloudCredentials,omitempty"`
	Region                  string                     `json:"region" yaml:"Region"`
	Version                 string                     `json:"version" yaml:"Version"`
	UpdateChannel           *UpdateChannelResponse     `json:"updateChannel" yaml:"UpdateChannel,omitempty"`
	AutoUpgrade             bool                       `json:"autoUpgrade" yaml:"AutoUpgrade"`
	HighlyAvailable         bool                       `json:"highlyAvailable" yaml:"HighlyAvailable"`
	EnableNetworkEncryption bool                       `json:"enableNetworkEncryption" yaml:"EnableNetworkEncryption"`
	// Deprecated: replaced by PodSecurityStandardsProfile which offers support for selecting a specific default PSS profile. This setting does not do anything since Kubernetes v1.23
	EnablePodSecurityStandards   bool                   `json:"enablePodSecurityStandards" yaml:"EnablePodSecurityStandards"`
	PodSecurityStandardsProfile  string                 `json:"podSecurityStandardsProfile" yaml:"PodSecurityStandardsProfile"`
	EnableMultiAvailAbilityZones bool                   `json:"enableMultiAvailabilityZones" yaml:"EnableMultiAvailabilityZones"`
	EnableNATGateway             bool                   `json:"enableNATGateway" yaml:"EnableNATGateway"`
	Status                       string                 `json:"status" yaml:"Status,omitempty"`
	DesiredStatus                string                 `json:"desiredStatus" yaml:"-"` // TODO: hidden for now in yaml since it is confusing
	ProvisionStatus              ClusterProvisionStatus `json:"provisionStatus" yaml:"ProvisionStatus"`
	CreatedAt                    time.Time              `json:"createdAt" yaml:"CreatedAt"`
	ModifiedAt                   time.Time              `json:"modifiedAt" yaml:"ModifiedAt"`
	DeletedAt                    *time.Time             `json:"deletedAt" yaml:"DeletedAt,omitempty"`
	Memory                       int                    `json:"memory" yaml:"Memory"`
	CPU                          int                    `json:"cpu" yaml:"CPU"`
	IPWhitelist                  []IpWhitelistResponse  `json:"ipWhitelist" yaml:"IPWhitelist,omitempty"`
	AmeOIDC                      bool                   `json:"ameOIDC" yaml:"AmeOIDC"`
	DeleteProtection             bool                   `json:"deleteProtection" yaml:"DeleteProtection"`
	Addons                       map[string]APIAddon    `json:"addons" yaml:"Addons,omitempty"`
	ObservabilityTenant          *ObservabilityTenant   `json:"observabilityTenant,omitempty" yaml:"ObservabilityTenant,omitempty"`
	EnvironmentPrometheusRules   bool                   `json:"environmentPrometheusRules" yaml:"EnvironmentPrometheusRules"`
	MaintenanceSchedule          *MaintenanceSchedule   `json:"maintenanceSchedule,omitempty" yaml:"MaintenanceSchedule,omitempty"`
	AutoscalerSettings           *AutoscalerSettings    `json:"clusterAutoscalerSettings,omitempty" yaml:"ClusterAutoscalerSettings,omitempty"`
	UpstreamCluster              *UpstreamCluster       `json:"upstreamCluster,omitempty" yaml:"UpstreamCluster,omitempty"`
}

type UpstreamCluster struct {
	Identity        string `json:"identity,omitempty" yaml:"Identity,omitempty"`
	Slug            string `json:"slug,omitempty" yaml:"Slug,omitempty"`
	EnvironmentSlug string `json:"environmentSlug,omitempty" yaml:"EnvironmentSlug,omitempty"`
	Version         string `json:"version,omitempty" yaml:"Version,omitempty"`
}

type MaintenanceSchedule struct {
	Identity           string              `json:"identity" yaml:"Identity"`
	Name               string              `json:"name" yaml:"Name"`
	MaintenanceWindows []MaintenanceWindow `json:"windows" yaml:"MaintenanceWindows"`
}

type MaintenanceWindow struct {
	Day       string `json:"day" yaml:"Day"`
	StartTime string `json:"startTime" yaml:"StartTime"`
	Duration  int    `json:"duration" yaml:"duration"`
}

func (m MaintenanceWindow) String() string {
	return fmt.Sprintf("%s %s %d minutes", m.Day, m.StartTime, m.Duration)
}

// IpWhitelistResponse represents the response structure for IP whitelist information.
type IpWhitelistResponse struct {
	Cidr        string `json:"cidr" yaml:"Cidr"`
	Description string `json:"description" yaml:"Description"`
}

// CloudAccountReference represents a reference to a cloud account.
type CloudAccountReference struct {
	Identity    string `json:"identity" yaml:"Identity"`
	DisplayName string `json:"displayName" yaml:"DisplayName"`
}

// CloudCredentialsReference represents a reference to cloud credentials.
type CloudCredentialsReference struct {
	Identity    string `json:"identity" yaml:"Identity"`
	DisplayName string `json:"displayName" yaml:"DisplayName"`
}

// Identifier returns the cluster identifier in the form of {organisation-slug}/{environment-slug}/{cluster-slug}
func (c Cluster) Identifier() string {
	return fmt.Sprintf("%s/%s/%s", c.CustomerSlug, c.EnvironmentSlug, c.Slug)
}

// FullIdentifier returns the cluster identifier including its cluster-identity
// in the form of {organisation-slug}/{environment-slug}/{cluster-slug} ({cluster-identity})
func (c Cluster) FullIdentifier() string {
	return fmt.Sprintf("%s/%s/%s (%s)", c.CustomerSlug, c.EnvironmentSlug, c.Slug, c.Identity)
}

// ClusterProvisionStatus represents the status of a cluster provision process.
type ClusterProvisionStatus string

const (
	UNKNOWN                        ClusterProvisionStatus = "UNKNOWN"
	ACCEPTED                       ClusterProvisionStatus = "ACCEPTED"
	OIDC_PROVISIONED               ClusterProvisionStatus = "OIDC_PROVISIONED"
	CLUSTER_PROVISIONED            ClusterProvisionStatus = "CLUSTER_PROVISIONED"
	INITIAL_NODE_POOLS_PROVISIONED ClusterProvisionStatus = "INITIAL_NODE_POOLS_PROVISIONED"
	INITIAL_ADDONS_PROVISIONED     ClusterProvisionStatus = "INITIAL_ADDONS_PROVISIONED"
	DONE                           ClusterProvisionStatus = "DONE"
)

// APIAddon represents an API addon.
type APIAddon struct {
	Enabled      bool              `json:"enabled" yaml:"Enabled"`
	CustomValues map[string]string `json:"customValues,omitempty" yaml:"CustomValues,omitempty"`
}

// CreateCluster represents the configuration for creating a cluster.
type CreateCluster struct {
	Name                 string `json:"name" yaml:"Name"`
	EnvironmentID        string `json:"environmentId" yaml:"EnvironmentId"`
	Description          string `json:"description,omitempty" yaml:"Description,omitempty"`
	CloudAccountIdentity string `json:"cloudAccountIdentity" yaml:"CloudAccountIdentity"`
	Region               string `json:"region" yaml:"Region"`

	Version       string `json:"version,omitempty" yaml:"Version,omitempty"`
	UpdateChannel string `json:"updateChannel,omitempty" yaml:"UpdateChannel,omitempty"`
	CNI           string `json:"cni,omitempty" yaml:"CNI,omitempty"`

	EnableNATGateway             bool   `json:"enableNATGateway" yaml:"EnableNATGateway"`
	EnableNetworkEncryption      bool   `json:"enableNetworkEncryption" yaml:"EnableNetworkEncryption"`
	EnablePodSecurityStandards   bool   `json:"enablePodSecurityStandards" yaml:"EnablePodSecurityStandards"`
	PodSecurityStandardsProfile  string `json:"podSecurityStandardsProfile,omitempty" yaml:"PodSecurityStandardsProfile,omitempty"`
	EnableMultiAvailabilityZones bool   `json:"enableMultiAvailabilityZones" yaml:"EnableMultiAvailabilityZones"`
	EnableAutoUpgrade            bool   `json:"enableAutoUpgrade" yaml:"EnableAutoUpgrade"`
	EnableHighAvailability       bool   `json:"enableHighAvailability" yaml:"EnableHighAvailability"`

	ServiceSubnet    string `json:"serviceSubnet,omitempty" yaml:"ServiceSubnet,omitempty"`
	ClusterPodSubnet string `json:"clusterPodSubnet,omitempty" yaml:"ClusterPodSubnet,omitempty"`

	NodePools   []NodePools         `json:"nodePools" yaml:"NodePools"`
	IPWhitelist []IPWhitelistEntry  `json:"ipWhitelist,omitempty" yaml:"IpWhitelist,omitempty"`
	Addons      map[string]APIAddon `json:"addons,omitempty" yaml:"Addons,omitempty"`

	AutoscalerSettings          *AutoscalerSettings     `json:"clusterAutoscalerSettings,omitempty" yaml:"ClusterAutoscalerSettings,omitempty"`
	MaintenanceScheduleIdentity string                  `json:"maintenanceScheduleIdentity,omitempty" yaml:"MaintenanceScheduleIdentity,omitempty"`
	UpstreamClusterIdentity     string                  `json:"upstreamClusterIdentity,omitempty" yaml:"UpstreamClusterIdentity,omitempty"`
	UpstreamCluster             *UpstreamClusterRequest `json:"upstreamCluster,omitempty" yaml:"UpstreamCluster,omitempty"`
}

type UpstreamClusterRequest struct {
	Identity string `json:"identity,omitempty" yaml:"Identity,omitempty"`
}

// IPWhitelistEntry represents an entry in the IP whitelist.
type IPWhitelistEntry struct {
	Cidr        string `json:"cidr" yaml:"Cidr"`
	Description string `json:"description,omitempty" yaml:"Description,omitempty"`
}

// UpdateCluster represents the data structure for updating a cluster.
type UpdateCluster struct {
	Status                  *string `json:"status,omitempty" yaml:"Status,omitempty"`
	Version                 *string `json:"version,omitempty" yaml:"Version,omitempty"`
	CNI                     *string `json:"cni,omitempty" yaml:"CNI,omitempty"`
	UpdateChannel           *string `json:"updateChannel,omitempty" yaml:"UpdateChannel,omitempty"`
	EnableNetworkProxy      *bool   `json:"enableNetworkProxy,omitempty" yaml:"EnableNetworkProxy,omitempty"`
	EnableNetworkEncryption *bool   `json:"enableNetworkEncryption,omitempty" yaml:"EnableNetworkEncryption,omitempty"`
	EnableAutoUpgrade       *bool   `json:"enableAutoUpgrade,omitempty" yaml:"EnableAutoUpgrade,omitempty"`
	EnableHighAvailability  *bool   `json:"enableHighAvailability,omitempty" yaml:"EnableHighAvailability,omitempty"`
	// Deprecated: replaced by PodSecurityStandardsProfile which offers support for selecting a specific default PSS profile. This setting does not do anything since Kubernetes v1.23
	EnablePodSecurityStandards  *bool                   `json:"enablePodSecurityStandards,omitempty" yaml:"EnablePodSecurityStandards,omitempty"`
	PodSecurityStandardsProfile *string                 `json:"podSecurityStandardsProfile,omitempty" yaml:"PodSecurityStandardsProfile,omitempty"`
	DeleteProtection            *bool                   `json:"deleteProtection,omitempty" yaml:"DeleteProtection,omitempty"`
	IPWhitelist                 []string                `json:"ipWhitelist,omitempty" yaml:"IpWhitelist,omitempty"`
	Addons                      map[string]APIAddon     `json:"addons,omitempty" yaml:"Addons,omitempty"`
	AutoscalerSettings          *AutoscalerSettings     `json:"clusterAutoscalerSettings,omitempty" yaml:"ClusterAutoscalerSettings,omitempty"`
	MaintenanceScheduleIdentity *string                 `json:"maintenanceScheduleIdentity,omitempty" yaml:"MaintenanceScheduleIdentity,omitempty"`
	UpstreamCluster             *UpstreamClusterRequest `json:"upstreamCluster,omitempty" yaml:"UpstreamCluster,omitempty"`
}

// NodePools is used by CreateCluster
type NodePools struct {
	Name             string                  `json:"name" yaml:"Name"`
	Identity         string                  `json:"identity" yaml:"Identity"`
	AutoScaling      bool                    `json:"autoScaling" yaml:"AutoScaling"`
	MinSize          int                     `json:"minSize" yaml:"MinSize"`
	MaxSize          int                     `json:"maxSize" yaml:"MaxSize"`
	NodeSize         string                  `json:"nodeSize" yaml:"NodeSize"`
	AvailabilityZone string                  `json:"availabilityZone,omitempty" yaml:"AvailabilityZone,omitempty"`
	UpgradeStrategy  NodePoolUpgradeStrategy `json:"upgradeStrategy" yaml:"UpgradeStrategy"`
}

// NodePool represents a pool of nodes in a cluster.
type NodePool struct {
	ID                  int                     `json:"id" yaml:"ID"`
	Identity            string                  `json:"identity" yaml:"Identity"`
	Name                string                  `json:"name" yaml:"Name"`
	AvailabilityZone    string                  `json:"availabilityZone" yaml:"AvailabilityZone,omitempty"`
	NodeSize            string                  `json:"nodeSize" yaml:"NodeSize"`
	AutoScaling         bool                    `json:"autoScaling" yaml:"AutoScaling"`
	MinSize             int                     `json:"minSize" yaml:"MinSize"`
	MaxSize             int                     `json:"maxSize" yaml:"MaxSize"`
	NodeAutoReplacement bool                    `json:"enableNodeAutoReplacement" yaml:"EnableNodeAutoReplacement"`
	EnableNodeReboots   bool                    `json:"enableNodeReboots" yaml:"EnableNodeReboots"`
	UpgradeStrategy     NodePoolUpgradeStrategy `json:"upgradeStrategy" yaml:"UpgradeStrategy"`
	Annotations         map[string]string       `json:"annotations" yaml:"Annotations"`
	Labels              map[string]string       `json:"labels" yaml:"Labels"`
	Taints              []NodeTaint             `json:"taints" yaml:"Taints"`
	ProvisionStatus     string                  `json:"provisionStatus" yaml:"ProvisionStatus,omitempty"`
	CreatedAt           time.Time               `json:"createdAt" yaml:"CreatedAt,omitempty"`
	ModifiedAt          time.Time               `json:"modifiedAt" yaml:"ModifiedAt,omitempty"`

	ClusterIdentity string  `json:"clusterIdentity" yaml:"ClusterIdentity,omitempty"` // adds a reference to Cluster
	Cluster         Cluster `json:"-" yaml:"-"`                                       // adds a reference to Cluster for in-code
}

// NodePoolUpgradeStrategy represents the upgrade strategy for a node pool.
type NodePoolUpgradeStrategy string

const (
	NodePoolUpgradeStrategyReplace                         NodePoolUpgradeStrategy = "REPLACE"
	NodePoolUpgradeStrategyInPlace                         NodePoolUpgradeStrategy = "INPLACE"
	NodePoolUpgradeStrategyInPlaceWithoutDrain             NodePoolUpgradeStrategy = "INPLACE_WITHOUT_DRAIN"
	NodePoolUpgradeStrategyReplaceMinorInPlacePatch        NodePoolUpgradeStrategy = "REPLACE_MINOR_INPLACE_PATCH"
	NodePoolUpgradeStrategyReplaceMinorInPlacePatchNoDrain NodePoolUpgradeStrategy = "REPLACE_MINOR_INPLACE_PATCH_WITHOUT_DRAIN"
)

var AllNodePoolUpgradeStrategies = []NodePoolUpgradeStrategy{
	NodePoolUpgradeStrategyReplace,
	NodePoolUpgradeStrategyInPlace,
	NodePoolUpgradeStrategyInPlaceWithoutDrain,
	NodePoolUpgradeStrategyReplaceMinorInPlacePatch,
	NodePoolUpgradeStrategyReplaceMinorInPlacePatchNoDrain,
}

func ParseNodePoolUpgradeStrategy(s string) (NodePoolUpgradeStrategy, error) {
	for _, strategy := range AllNodePoolUpgradeStrategies {
		if string(strategy) == s {
			return strategy, nil
		}
	}
	return "", fmt.Errorf("unsupported upgrade strategy: %s", s)
}

// NodeTaint represents a taint applied to a Kubernetes node.
type NodeTaint struct {
	Key    string `json:"key" yaml:"Key"`       // Key is the key of the taint.
	Value  string `json:"value" yaml:"Value"`   // Value is the value of the taint.
	Effect string `json:"effect" yaml:"Effect"` // Effect is the effect of the taint.
}

// NodePoolJoinConfig represents the configuration for joining a node pool. Only used for Bring Your Own Node cluster node pools.
type NodePoolJoinConfig struct {
	Versions                NodeJoinConfigVersions `json:"versions" yaml:"Versions"`
	CloudInitUserDataBase64 string                 `json:"cloudInitUserDataBase64" yaml:"CloudInitUserDataBase64"`
	InstallScriptBase64     string                 `json:"installScriptBase64" yaml:"InstallScriptBase64"`
	UpgradeScriptBase64     string                 `json:"upgradeScriptBase64" yaml:"UpgradeScriptBase64"`
	JoinCommand             string                 `json:"joinCommand" yaml:"JoinCommand"`
	KubeletConfigBase64     string                 `json:"kubeletConfigBase64" yaml:"KubeletConfigBase64"`
}

// NodeJoinConfigVersions represents the versions of various components used in the node join configuration. Only used for Bring Your Own Node cluster node pools.
type NodeJoinConfigVersions struct {
	CloudInit  string `json:"cloudInit" yaml:"CloudInit"`
	Kubernetes string `json:"kubernetes" yaml:"Kubernetes"`
	Containerd string `json:"containerd" yaml:"Containerd"`
	Crictl     string `json:"crictl" yaml:"Crictl"`
}

func (n NodePool) FullIdentifier() string {
	c := n.Cluster
	return fmt.Sprintf("%s/%s/%s/%s (%d)", c.CustomerSlug, c.EnvironmentSlug, c.Slug, n.Name, n.ID)
}

// CreateNodePool represents the configuration for creating a node pool.
type CreateNodePool struct {
	Name                string                  `json:"name" yaml:"Name"`
	AvailabilityZone    string                  `json:"availabilityZone,omitempty" yaml:"AvailabilityZone,omitempty"`
	NodeSize            string                  `json:"nodeSize" yaml:"NodeSize"`
	MinSize             int                     `json:"minSize" yaml:"MinSize"`
	MaxSize             int                     `json:"maxSize" yaml:"MaxSize"`
	AutoScaling         bool                    `json:"autoScaling" yaml:"AutoScaling"`
	NodeAutoReplacement bool                    `json:"enableNodeAutoReplacement" yaml:"EnableNodeAutoReplacement"`
	EnableNodeReboots   bool                    `json:"enableNodeReboots" yaml:"EnableNodeReboots"`
	UpgradeStrategy     NodePoolUpgradeStrategy `json:"upgradeStrategy,omitempty" yaml:"UpgradeStrategy"`
	Annotations         map[string]string       `json:"annotations" yaml:"Annotations"`
	Labels              map[string]string       `json:"labels" yaml:"Labels"`
	Taints              []NodeTaint             `json:"taints" yaml:"Taints"`
}

type ClusterMetadataResponse struct {
	Endpoint      string `json:"endpoint" yaml:"Endpoint"`
	CACertificate string `json:"caCertificate" yaml:"CaCertificate"`
	ClientID      string `json:"clientId" yaml:"ClientId"`
	ClientSecret  string `json:"clientSecret" yaml:"ClientSecret"`
	IssuerUrl     string `json:"issuerUrl" yaml:"IssuerUrl"`
}

type ClusterVersion struct {
	Version string `json:"version" yaml:"Version"`
}

type Membership struct {
	Email string `json:"email" yaml:"Email"`
	ID    string `json:"id" yaml:"Id"`
	Name  string `json:"name" yaml:"Name"`
	Slug  string `json:"slug" yaml:"Slug"`
}

type CloudProvider struct {
	ID         int       `json:"id" yaml:"Id"`
	Name       string    `json:"name" yaml:"Name"`
	Slug       string    `json:"slug" yaml:"Slug"`
	Available  bool      `json:"available" yaml:"Available"`
	CreatedAt  time.Time `json:"createdAt" yaml:"CreatedAt"`
	ModifiedAt time.Time `json:"modifiedAt" yaml:"ModifiedAt"`
}

type Region struct {
	ID         int       `json:"id" yaml:"Id"`
	Provider   string    `json:"provider" yaml:"Provider"`
	Name       string    `json:"name" yaml:"Name"`
	Slug       string    `json:"slug" yaml:"Slug"`
	Available  bool      `json:"available" yaml:"Available"`
	Restricted bool      `json:"restricted" yaml:"Restricted"`
	CreatedAt  time.Time `json:"createdAt" yaml:"CreatedAt"`
	ModifiedAt time.Time `json:"modifiedAt" yaml:"ModifiedAt"`
}

type AvailabilityZone struct {
	ID         int       `json:"id" yaml:"Id"`
	Name       string    `json:"name" yaml:"Name"`
	Slug       string    `json:"slug" yaml:"Slug"`
	CreatedAt  time.Time `json:"createdAt" yaml:"CreatedAt"`
	ModifiedAt time.Time `json:"modifiedAt" yaml:"ModifiedAt"`
}

type ServiceLevelAgreement struct {
	Slug        string    `json:"slug" yaml:"Slug"`
	Name        string    `json:"name" yaml:"Name"`
	Value       int       `json:"value" yaml:"Value"`
	AutoUpgrade bool      `json:"autoUpgrade" yaml:"AutoUpgrade"`
	CreatedAt   time.Time `json:"createdAt" yaml:"CreatedAt"`
	ModifiedAt  time.Time `json:"modifiedAt" yaml:"ModifiedAt"`
}

type Environment struct {
	ID               int        `json:"id" yaml:"Id"`
	Name             string     `json:"name" yaml:"Name"`
	Purpose          string     `json:"purpose" yaml:"Purpose"`
	Type             string     `json:"type" yaml:"Type"`
	Description      string     `json:"description" yaml:"Description"`
	CreatedAt        time.Time  `json:"createdAt" yaml:"CreatedAt"`
	ModifiedAt       time.Time  `json:"modifiedAt" yaml:"ModifiedAt"`
	DeletedAt        *time.Time `json:"deletedAt" yaml:"DeletedAt"`
	TotalClusters    int        `json:"totalClusters" yaml:"TotalClusters"`
	TotalCPU         int        `json:"totalCpu" yaml:"TotalCpu"`
	TotalMemory      int        `json:"totalMemory" yaml:"TotalMemory"`
	Slug             string     `json:"slug" yaml:"Slug"`
	OrganisationSlug string     `json:"organisationSlug" yaml:"OrganisationSlug"`
}

type CreateEnvironment struct {
	Name        string `json:"name" yaml:"Name"`
	Purpose     string `json:"purpose" yaml:"Purpose"`
	Type        string `json:"type" yaml:"Type"`
	Description string `json:"description" yaml:"Description"`
}

type UpdateEnvironment struct {
	Name        string `json:"name" yaml:"Name"`
	Purpose     string `json:"purpose" yaml:"Purpose"`
	Type        string `json:"type" yaml:"Type"`
	Description string `json:"description" yaml:"Description"`
}

type NodeType struct {
	Type   string `json:"type" yaml:"Type"`
	CPU    int    `json:"cpu" yaml:"Cpu"`
	Memory int    `json:"memory" yaml:"Memory"`
}

type UpdateChannelResponse struct {
	Name                     string `json:"name" yaml:"Name"`
	Available                bool   `json:"available" yaml:"Available"`
	KubernetesClusterVersion string `json:"kubernetesClusterVersion" yaml:"Version"`
}

type CreateSilence struct {
	Matchers []SilenceMatcher `json:"matchers" yaml:"Matchers"`
	StartsAt time.Time        `json:"startsAt" yaml:"StartsAt"`
	EndsAt   time.Time        `json:"endsAt" yaml:"EndsAt"`
	Comment  string           `json:"comment" yaml:"Comment"`
}

type Silence struct {
	Id        string           `json:"id" yaml:"Id"`
	Matchers  []SilenceMatcher `json:"matchers" yaml:"Matchers"`
	StartsAt  time.Time        `json:"startsAt" yaml:"StartsAt"`
	EndsAt    time.Time        `json:"endsAt" yaml:"EndsAt"`
	CreatedBy string           `json:"createdBy" yaml:"CreatedBy"`
	Comment   string           `json:"comment" yaml:"Comment"`
	Status    SilenceStatus    `json:"status" yaml:"Status"`
}

type SilenceStatus struct {
	State string `json:"state" yaml:"State"`
}

type SilenceMatcher struct {
	Name    string `json:"name" yaml:"Name"`
	Value   string `json:"value" yaml:"Value"`
	IsRegex bool   `json:"isRegex" yaml:"IsRegex"`
	IsEqual bool   `json:"isEqual" yaml:"IsEqual"`
}

type ScheduledClusterUpgrade struct {
	Identity                  string                        `json:"identity" yaml:"Identity"`
	ClusterIdentity           string                        `json:"clusterIdentity" yaml:"ClusterIdentity"`
	CreatedAt                 time.Time                     `json:"createdAt" yaml:"CreatedAt"`
	ModifiedAt                time.Time                     `json:"modifiedAt" yaml:"ModifiedAt"`
	WindowStart               time.Time                     `json:"windowStart" yaml:"WindowStart"`
	WindowEnd                 time.Time                     `json:"windowEnd" yaml:"WindowEnd"`
	FromClusterVersion        string                        `json:"fromClusterVersion" yaml:"FromClusterVersion"`
	ToClusterVersion          string                        `json:"toClusterVersion" yaml:"ToClusterVersion"`
	Status                    ScheduledClusterUpgradeStatus `json:"status" yaml:"Status"`
	Reason                    string                        `json:"reason" yaml:"Reason"`
	ScheduleRequestDate       time.Time                     `json:"scheduleRequestDate" yaml:"ScheduleRequestDate"`
	IgnoreMaintenanceSchedule bool                          `json:"ignoreMaintenanceSchedule" yaml:"IgnoreMaintenanceSchedule"`
}

type ScheduledClusterUpgradeStatus string

const (
	// Requested represents the status where a cluster upgrade
	// has been explicitly requested but not yet scheduled.
	// This status is used to indicate that a new scheduled
	// cluster upgrade has been initiated by the user or system.
	Requested ScheduledClusterUpgradeStatus = "REQUESTED"

	// Scheduled indicates that this cluster upgrade has been scheduled
	// and is pending execution.
	Scheduled ScheduledClusterUpgradeStatus = "SCHEDULED"

	// ScheduledNotified indicates that this upgrade has been scheduled
	// and the user has been notified of the planned upgrade.
	ScheduledNotified ScheduledClusterUpgradeStatus = "SCHEDULED_NOTIFIED"

	// InProgress signifies that the cluster upgrade is currently in progress
	// but has not yet completed.
	InProgress ScheduledClusterUpgradeStatus = "IN_PROGRESS"

	// Succeeded indicates that the cluster upgrade has been completed successfully.
	Succeeded ScheduledClusterUpgradeStatus = "SUCCEEDED"

	// Superseded means that this scheduled cluster upgrade has been superseded
	// due to external factors, such as a change in the update channel version
	// or a change in the cluster's status.
	Superseded ScheduledClusterUpgradeStatus = "SUPERSEDED"

	// Failed indicates that the scheduled cluster upgrade failed to meet the
	// necessary success criteria within the allotted time frame.
	Failed ScheduledClusterUpgradeStatus = "FAILED"

	// Missed means that the scheduled cluster upgrade did not initiate
	// within its designated time window.
	Missed ScheduledClusterUpgradeStatus = "MISSED"
)

type CreateScheduledClusterUpgradeRequest struct {
	ClusterIdentity    string    `json:"clusterIdentity"`
	WindowStart        time.Time `json:"windowStart"`
	WindowEnd          time.Time `json:"windowEnd"`
	FromClusterVersion string    `json:"fromClusterVersion"`
	ToClusterVersion   string    `json:"toClusterVersion"`
}

type UpdateScheduledClusterUpgradeRequest struct {
	Identity    string                        `json:"identity"`
	Status      ScheduledClusterUpgradeStatus `json:"status,omitempty"`
	Reason      string                        `json:"reason,omitempty"`
	WindowStart *time.Time                    `json:"windowStart,omitempty"`
	WindowEnd   *time.Time                    `json:"windowEnd,omitempty"`
	Version     string                        `json:"version,omitempty"`
}

type ListScheduledClusterUpgradesOpts struct {
	ClusterIdentities []string
	Statuses          []ScheduledClusterUpgradeStatus
}

type AutoscalerSettings struct {
	// ScaleDownUtilizationThreshold specifies the scale down utilization threshold
	ScaleDownUtilizationThreshold *float64 `json:"scale-down-utilization-threshold,omitempty" yaml:"ScaleDownUtilizationThreshold,omitempty"`
	// ScaleDownGpuUtilizationThreshold specifies the scale down GPU utilization threshold
	ScaleDownGpuUtilizationThreshold *float64 `json:"scale-down-gpu-utilization-threshold,omitempty" yaml:"ScaleDownGpuUtilizationThreshold,omitempty"`
	// ScaleDownDelayAfterAdd specifies the scale down delay after add
	ScaleDownDelayAfterAdd *string `json:"scale-down-delay-after-add,omitempty" yaml:"ScaleDownDelayAfterAdd,omitempty"`
	// ScaleDownUnneededTime specifies the scale down unneeded time
	ScaleDownUnneededTime *string `json:"scale-down-unneeded-time,omitempty" yaml:"ScaleDownUnneededTime,omitempty"`
	// ScaleDownUnreadyTime specifies the scale down unready time
	ScaleDownUnreadyTime *string `json:"scale-down-unready-time,omitempty" yaml:"ScaleDownUnreadyTime,omitempty"`
	// MaxNodeProvisionTime specifies the max node provision time
	MaxNodeProvisionTime *string `json:"max-node-provision-time,omitempty" yaml:"MaxNodeProvisionTime,omitempty"`
	// UnremovableNodeRecheckTimeout specifies the unremovable node recheck timeout
	UnremovableNodeRecheckTimeout *string `json:"unremovable-node-recheck-timeout,omitempty" yaml:"UnremovableNodeRecheckTimeout,omitempty"`
}
