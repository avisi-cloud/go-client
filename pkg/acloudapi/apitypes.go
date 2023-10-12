package acloudapi

import (
	"fmt"
	"time"
)

type Cluster struct {
	Name                         string                     `json:"name" yaml:"Name"`
	Identity                     string                     `json:"identity" yaml:"Identity"`
	EnvironmentIdentity          string                     `json:"environmentIdentity" yaml:"EnvironmentIdentity"`
	EnvironmentSlug              string                     `json:"environmentSlug" yaml:"EnvironmentSlug"`
	CustomerIdentity             string                     `json:"customerIdentity" yaml:"CustomerIdentity"`
	CustomerSlug                 string                     `json:"customerSlug" yaml:"CustomerSlug"`
	Slug                         string                     `json:"slug" yaml:"Slug"`
	Description                  string                     `json:"description" yaml:"Description,omitempty"`
	CloudProvider                string                     `json:"cloudProvider" yaml:"CloudProvider"`
	CloudAccount                 *CloudAccountReference     `json:"cloudAccount" yaml:"CloudAccount,omitempty"`
	CloudCredentials             *CloudCredentialsReference `json:"cloudCredentials" yaml:"CloudCredentials,omitempty"`
	Region                       string                     `json:"region" yaml:"Region"`
	Version                      string                     `json:"version" yaml:"Version"`
	UpdateChannel                *UpdateChannelResponse     `json:"updateChannel" yaml:"UpdateChannel,omitempty"`
	AutoUpgrade                  bool                       `json:"autoUpgrade" yaml:"AutoUpgrade"`
	HighlyAvailable              bool                       `json:"highlyAvailable" yaml:"HighlyAvailable"`
	EnableNetworkEncryption      bool                       `json:"enableNetworkEncryption" yaml:"EnableNetworkEncryption"`
	EnablePodSecurityStandards   bool                       `json:"enablePodSecurityStandards" yaml:"EnablePodSecurityStandards"`
	EnableMultiAvailAbilityZones bool                       `json:"enableMultiAvailabilityZones" yaml:"EnableMultiAvailabilityZones"`
	EnableNATGateway             bool                       `json:"enableNATGateway" yaml:"EnableNATGateway"`
	Status                       string                     `json:"status" yaml:"Status,omitempty"`
	DesiredStatus                string                     `json:"desiredStatus" yaml:"-"` // TODO: hidden for now in yaml since it is confusing
	ProvisionStatus              ClusterProvisionStatus     `json:"provisionStatus" yaml:"ProvisionStatus"`
	CreatedAt                    time.Time                  `json:"createdAt" yaml:"CreatedAt"`
	ModifiedAt                   time.Time                  `json:"modifiedAt" yaml:"ModifiedAt"`
	DeletedAt                    *time.Time                 `json:"deletedAt" yaml:"DeletedAt,omitempty"`
	Memory                       int                        `json:"memory" yaml:"Memory"`
	CPU                          int                        `json:"cpu" yaml:"CPU"`
	SLA                          string                     `json:"sla" yaml:"SLA"`
	IPWhitelist                  []IpWhitelistResponse      `json:"ipWhitelist" yaml:"IPWhitelist,omitempty"`
	AmeOIDC                      bool                       `json:"ameOIDC" yaml:"AmeOIDC"`
	DeleteProtection             bool                       `json:"deleteProtection" yaml:"DeleteProtection"`
	Addons                       map[string]APIAddon        `json:"addons" yaml:"Addons,omitempty"`
	ObservabilityTenant          *ObservabilityTenant       `json:"observabilityTenant,omitempty" yaml:"ObservabilityTenant,omitempty"`
	EnvironmentPrometheusRules   bool                       `json:"environmentPrometheusRules" yaml:"EnvironmentPrometheusRules"`
}

type IpWhitelistResponse struct {
	Cidr        string `json:"cidr" yaml:"Cidr"`
	Description string `json:"description" yaml:"Description"`
}

type CloudAccountReference struct {
	Identity    string `json:"identity" yaml:"Identity"`
	DisplayName string `json:"displayName" yaml:"DisplayName"`
}

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

type APIAddon struct {
	Enabled      bool              `json:"enabled" yaml:"Enabled"`
	CustomValues map[string]string `json:"customValues,omitempty" yaml:"CustomValues,omitempty"`
}

type CreateCluster struct {
	Name                 string `json:"name"`
	EnvironmentID        string `json:"environmentId"`
	Description          string `json:"description,omitempty"`
	CloudAccountIdentity string `json:"cloudAccountIdentity"`
	Region               string `json:"region"`

	Version       string `json:"version,omitempty"`
	UpdateChannel string `json:"updateChannel,omitempty"`

	EnableNATGateway             bool `json:"enableNATGateway"`
	EnableNetworkEncryption      bool `json:"enableNetworkEncryption"`
	EnablePodSecurityStandards   bool `json:"enablePodSecurityStandards"`
	EnableMultiAvailabilityZones bool `json:"enableMultiAvailabilityZones"`
	EnableAutoUpgrade            bool `json:"enableAutoUpgrade"`
	EnableHighAvailability       bool `json:"enableHighAvailability"`

	SLA string `json:"sla,omitempty"`

	ServiceSubnet    string `json:"serviceSubnet,omitempty"`
	ClusterPodSubnet string `json:"clusterPodSubnet,omitempty"`

	NodePools   []NodePools         `json:"nodePools"`
	IPWhitelist []IPWhitelistEntry  `json:"ipWhitelist,omitempty"`
	Addons      map[string]APIAddon `json:"addons,omitempty"`
}

type IPWhitelistEntry struct {
	Cidr        string `json:"cidr"`
	Description string `json:"description,omitempty"`
}

type UpdateCluster struct {
	Status                     string              `json:"status,omitempty"`
	Version                    string              `json:"version,omitempty"`
	UpdateChannel              string              `json:"updateChannel,omitempty"`
	EnableNetworkProxy         *bool               `json:"enableNetworkProxy,omitempty"`
	EnableNetworkEncryption    *bool               `json:"enableNetworkEncryption,omitempty"`
	EnableAutoUpgrade          *bool               `json:"enableAutoUpgrade,omitempty"`
	EnableHighAvailability     *bool               `json:"enableHighAvailability,omitempty"`
	EnablePodSecurityStandards *bool               `json:"enablePodSecurityStandards,omitempty"`
	IPWhitelist                []string            `json:"ipWhitelist,omitempty"`
	Addons                     map[string]APIAddon `json:"addons,omitempty"`
}

// NodePools is used by CreateCluster
type NodePools struct {
	Name             string `json:"name"`
	AutoScaling      bool   `json:"autoScaling"`
	MinSize          int    `json:"minSize"`
	MaxSize          int    `json:"maxSize"`
	NodeSize         string `json:"nodeSize"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
}

type NodePool struct {
	ID                  int               `json:"id" yaml:"ID"`
	Identity            string            `json:"identity" yaml:"Identity"`
	Name                string            `json:"name" yaml:"Name"`
	AvailabilityZone    string            `json:"availabilityZone" yaml:"AvailabilityZone,omitempty"`
	NodeSize            string            `json:"nodeSize" yaml:"NodeSize"`
	AutoScaling         bool              `json:"autoScaling" yaml:"AutoScaling"`
	MinSize             int               `json:"minSize" yaml:"MinSize"`
	MaxSize             int               `json:"maxSize" yaml:"MaxSize"`
	NodeAutoReplacement bool              `json:"enableNodeAutoReplacement" yaml:"EnableNodeAutoReplacement"`
	Annotations         map[string]string `json:"annotations" yaml:"Annotations"`
	Labels              map[string]string `json:"labels" yaml:"Labels"`
	Taints              []NodeTaint       `json:"taints" yaml:"Taints"`
	Status              string            `json:"status" yaml:"Status,omitempty"`
	CreatedAt           time.Time         `json:"createdAt" yaml:"CreatedAt,omitempty"`
	ModifiedAt          time.Time         `json:"modifiedAt" yaml:"ModifiedAt,omitempty"`
	DeletedAt           *time.Time        `json:"deletedAt" yaml:"DeletedAt,omitempty"`

	ClusterIdentity string  `json:"clusterIdentity" yaml:"ClusterIdentity,omitempty"` // adds a reference to Cluster
	Cluster         Cluster `json:"-" yaml:"-"`                                       // adds a reference to Cluster for in-code
}

type NodeTaint struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Effect string `json:"effect"`
}

func (n NodePool) FullIdentifier() string {
	c := n.Cluster
	return fmt.Sprintf("%s/%s/%s/%s (%d)", c.CustomerSlug, c.EnvironmentSlug, c.Slug, n.Name, n.ID)
}

type CreateNodePool struct {
	Name                string            `json:"name"`
	AvailabilityZone    string            `json:"availabilityZone,omitempty"`
	NodeSize            string            `json:"nodeSize"`
	MinSize             int               `json:"minSize"`
	MaxSize             int               `json:"maxSize"`
	AutoScaling         bool              `json:"autoScaling"`
	NodeAutoReplacement bool              `json:"enableNodeAutoReplacement"`
	Annotations         map[string]string `json:"annotations"`
	Labels              map[string]string `json:"labels"`
	Taints              []NodeTaint       `json:"taints"`
}

type ClusterMetadataResponse struct {
	Endpoint      string `json:"endpoint"`
	CACertificate string `json:"caCertificate"`
	ClientID      string `json:"clientId"`
	ClientSecret  string `json:"clientSecret"`
	IssuerUrl     string `json:"issuerUrl"`
}

type ClusterVersion struct {
	Version string `json:"version"`
}

type Membership struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
}

type CloudProvider struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Slug       string    `json:"slug"`
	Available  bool      `json:"available"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
}

type Region struct {
	ID         int       `json:"id"`
	Provider   string    `json:"provider"`
	Name       string    `json:"name"`
	Slug       string    `json:"slug"`
	Available  bool      `json:"available"`
	Restricted bool      `json:"restricted"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
}

type AvailabilityZone struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Slug       string    `json:"slug"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
}

type ServiceLevelAgreement struct {
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	Value       int       `json:"value"`
	AutoUpgrade bool      `json:"autoUpgrade"`
	CreatedAt   time.Time `json:"createdAt"`
	ModifiedAt  time.Time `json:"modifiedAt"`
}

type Environment struct {
	ID               int        `json:"id"`
	Name             string     `json:"name"`
	Purpose          string     `json:"purpose"`
	Type             string     `json:"type"`
	Description      string     `json:"description"`
	CreatedAt        time.Time  `json:"createdAt"`
	ModifiedAt       time.Time  `json:"modifiedAt"`
	DeletedAt        *time.Time `json:"deletedAt"`
	TotalClusters    int        `json:"totalClusters"`
	TotalCPU         int        `json:"totalCpu"`
	TotalMemory      int        `json:"totalMemory"`
	Slug             string     `json:"slug"`
	OrganisationSlug string     `json:"organisationSlug"`
}

type CreateEnvironment struct {
	Name        string `json:"name"`
	Purpose     string `json:"purpose"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type UpdateEnvironment struct {
	Name        string `json:"name"`
	Purpose     string `json:"purpose"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type NodeType struct {
	Type   string `json:"type"`
	CPU    int    `json:"cpu"`
	Memory int    `json:"memory"`
}

type UpdateChannelResponse struct {
	Name                     string `json:"name" yaml:"Name"`
	Available                bool   `json:"available" yaml:"Available"`
	KubernetesClusterVersion string `json:"kubernetesClusterVersion" yaml:"Version"`
}
