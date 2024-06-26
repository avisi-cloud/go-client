package acloudapi

import (
	"context"

	"github.com/go-resty/resty/v2"
)

type ClusterAPI interface {
	GetClusters(ctx context.Context, opts ...GetClusterOpts) ([]Cluster, error)
	GetClustersByOrg(ctx context.Context, organisationSlug string, opts ...GetClusterOpts) ([]Cluster, error)
	GetClustersByOrgAndEnv(ctx context.Context, organisationSlug, environmentSlug string, opts ...GetClusterOpts) ([]Cluster, error)

	GetCluster(ctx context.Context, organisationSlug, environmentSlug, cluster string, opts ...GetClusterOpts) (*Cluster, error)
	GetClusterOIDCConfig(ctx context.Context, organisationSlug, environmentSlug, clusterSlug string) (*ClusterMetadataResponse, error)

	CreateCluster(ctx context.Context, organisationSlug, environmentSlug string, create CreateCluster) (*Cluster, error)
	UpdateCluster(ctx context.Context, organisationSlug, environmentSlug, clusterSlug string, update UpdateCluster) (*Cluster, error)
	DeleteCluster(ctx context.Context, organisationSlug, environmentSlug, clusterSlug string, update UpdateCluster) error
}

type ClusterVersionAPI interface {
	GetClusterVersions(ctx context.Context) ([]ClusterVersion, error)
}

type CloudAccountAPI interface {
	GetCloudAccounts(ctx context.Context, organisationSlug string) ([]CloudAccount, error)
}

type CloudProvidersAPI interface {
	GetCloudProviders(ctx context.Context, organisationSlug string) ([]CloudProvider, error)
	GetRegions(ctx context.Context, organisationSlug, cloudProviderSlug string) ([]Region, error)
	GetAvailabilityZones(ctx context.Context, organisationSlug, cloudProviderSlug, regionSlug string) ([]AvailabilityZone, error)
	GetNodeTypes(ctx context.Context, cloudProviderSlug string) ([]NodeType, error)
}

type EnvironmentsAPI interface {
	GetEnvironment(ctx context.Context, org, env string) (*Environment, error)
	CreateEnvironment(ctx context.Context, createEnvironment CreateEnvironment, org string) (*Environment, error)
	UpdateEnvironment(ctx context.Context, updateEnvironment UpdateEnvironment, org, env string) (*Environment, error)
	DeleteEnvironment(ctx context.Context, org, env string) error
	GetEnvironments(ctx context.Context, organisationSlug string) ([]Environment, error)
}

type NodePoolsAPI interface {
	GetNodePools(ctx context.Context) ([]NodePool, error)
	GetNodePoolsByOrg(ctx context.Context, organisationSlug string) ([]NodePool, error)
	GetNodePoolsByCluster(ctx context.Context, cluster Cluster) ([]NodePool, error)
	GetNodePoolsByClusters(ctx context.Context, clusters []Cluster) ([]NodePool, error)
	GetNodePoolJoinConfig(ctx context.Context, cluster Cluster, nodePool NodePool) (*NodePoolJoinConfig, error)
	CreateNodePool(ctx context.Context, cluster Cluster, create CreateNodePool) (*NodePool, error)
	UpdateNodePool(ctx context.Context, cluster Cluster, nodePoolID int, update CreateNodePool) (*NodePool, error)
	DeleteNodePool(ctx context.Context, cluster Cluster, nodePoolID int) error
}

type CloudAccountsAPI interface {
	GetCloudAccounts(ctx context.Context, org string) ([]CloudAccount, error)
	CreateCloudAccount(ctx context.Context, org string, createCloudAccount CreateCloudAccount) (*CloudAccount, error)
	UpdateCloudAccount(ctx context.Context, org, cloudAccount string, updateCloudAccount UpdateCloudAccount) (*CloudAccount, error)
	DeleteCloudAccount(ctx context.Context, org, cloudAccount string) error
	FindCloudAccountByName(ctx context.Context, org, name, cloudProvider string) (*CloudAccount, error)
	GetCloudProfiles(ctx context.Context, org string) ([]CloudProfile, error)

	GetCloudCredentials(ctx context.Context, org, cloudAccountIdentity string) ([]CloudCredential, error)
	CreateCloudCredential(ctx context.Context, org string, cloudAccount CloudAccount, create CreateCloudCredential) (*CloudCredential, error)
	DeleteCloudCredential(ctx context.Context, org, cloudAccountIdentity, cloudCredentialIdentity string) error
}

type MembershipsAPI interface {
	GetMemberships(ctx context.Context) ([]Membership, error)
}

type UpdateChannelAPI interface {
	GetUpdateChannels(ctx context.Context, org string) ([]UpdateChannelResponse, error)
}

type ObservabilityAPI interface {
	GetObservabilityTenants(ctx context.Context, org string) ([]ObservabilityTenant, error)
	GetObservabilityTenantBySlug(ctx context.Context, org, slug string) (*ObservabilityTenant, error)
	GetObservabilityOrganisationAlerts(ctx context.Context, org string) ([]ObservabilityAlert, error)
	GetObservabilityTenantAlerts(ctx context.Context, org, slug string) ([]ObservabilityAlert, error)
	GetObservabilityTenantAlertmanagerConfiguration(ctx context.Context, org, slug string) (*ObservabilityAlertmanager, error)
	AddObservabilityTenantPrometheusRules(ctx context.Context, org, slug string, rules []PrometheusRules, force bool) error
	OverwriteObservabilityTenantPrometheusRules(ctx context.Context, org, slug string, rules []PrometheusRules) error
	DeleteObservabilityTenantPrometheusRules(ctx context.Context, org, slug string, names []string) error

	GetSilences(ctx context.Context, org, observabilityTenantSlug string) ([]Silence, error)
	CreateSilence(ctx context.Context, createSilence CreateSilence, org, observabilityTenantSlug string) (*Silence, error)
	ExpireSilence(ctx context.Context, org, observabilityTenantSlug, silenceID string) error
}

type OrganisationAPI interface {
	GetOrganisation(ctx context.Context, organisationSlug string) (*Organisation, error)
}

type Client interface {
	CloudAccountAPI
	CloudProvidersAPI
	ClusterAPI
	ClusterVersionAPI
	EnvironmentsAPI
	NodePoolsAPI
	MembershipsAPI
	UpdateChannelAPI
	ObservabilityAPI
	OrganisationAPI
	CloudAccountsAPI

	Resty() *resty.Client
}

type clientImpl struct {
	*RestyClient
}

func NewClient(authenticator Authenticator, opts ClientOpts) Client {
	return &clientImpl{
		NewRestyClient(authenticator, opts),
	}
}
