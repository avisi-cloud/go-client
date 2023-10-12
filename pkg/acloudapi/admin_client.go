package acloudapi

import (
	"context"

	"github.com/go-resty/resty/v2"
)

type AdminClusterAPI interface {
	GetCluster(ctx context.Context, clusterIdentity string, opts ...GetClusterOpts) (*Cluster, error)
}

type AdminOrganisationAPI interface {
	GetOrganisation(ctx context.Context, organisationIdentity string) (*AdminOrganisation, error)
}

type AdminClient interface {
	AdminClusterAPI
	AdminOrganisationAPI

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
