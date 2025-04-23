package acloudapi

import (
	"context"

	"github.com/go-resty/resty/v2"
)

type AdminClusterAPI interface {
	GetCluster(ctx context.Context, clusterIdentity string, opts ...GetClusterOpts) (*Cluster, error)
	ListClusters(ctx context.Context, opts ...GetClusterOpts) ([]Cluster, error)
	UpdateCluster(ctx context.Context, request AdminUpdateClusterRequest) (*Cluster, error)
}

type AdminOrganisationAPI interface {
	GetOrganisation(ctx context.Context, organisationIdentity string) (*AdminOrganisation, error)
}

type AdminScheduledClusterUpgradesAPI interface {
	ListScheduledClusterUpgrades(ctx context.Context, opts ...ListScheduledClusterUpgradesOpts) ([]ScheduledClusterUpgrade, error)
	GetScheduledClusterUpgrade(ctx context.Context, identity string) (*ScheduledClusterUpgrade, error)
	CancelScheduledClusterUpgrade(ctx context.Context, identity string) (*ScheduledClusterUpgrade, error)
	CreateScheduledClusterUpgrade(ctx context.Context, request CreateScheduledClusterUpgradeRequest) (*ScheduledClusterUpgrade, error)
	UpdateScheduledClusterUpgrade(ctx context.Context, request UpdateScheduledClusterUpgradeRequest) (*ScheduledClusterUpgrade, error)
}

type AdminUpdateChannelsAPI interface {
	ListUpdateChannels(ctx context.Context) ([]UpdateChannelResponse, error)
}

type AdminClusterVersionsAPI interface {
	ListClusterVersions(ctx context.Context) ([]AdminClusterVersion, error)
	ListAvailableClusterVersions(ctx context.Context) ([]AdminClusterVersion, error)
	ListHistoryClusterVersions(ctx context.Context) ([]AdminClusterVersion, error)

	GetClusterVersion(ctx context.Context, version string) (*AdminClusterVersion, error)
	UpdateClusterVersion(ctx context.Context, version string, request AdminUpdateClusterVersionRequest) (*AdminClusterVersion, error)
	CreateClusterVersion(ctx context.Context, request AdminCreateClusterVersionRequest) (*AdminClusterVersion, error)
	DeleteClusterVersion(ctx context.Context, version string) error
}

type AdminClient interface {
	AdminClusterAPI
	AdminOrganisationAPI
	AdminScheduledClusterUpgradesAPI
	AdminUpdateChannelsAPI
	AdminClusterVersionsAPI

	Resty() *resty.Client
}

type adminClientImpl struct {
	*RestyClient
}

func NewAdminClient(authenticator Authenticator, opts ClientOpts) AdminClient {
	return &adminClientImpl{
		NewRestyClient(authenticator, opts),
	}
}
