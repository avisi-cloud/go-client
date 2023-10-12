package acloudapi

import (
	"context"
	"fmt"
)

func (c *adminClientImpl) GetCluster(ctx context.Context, clusterIdentity string, opts ...GetClusterOpts) (*Cluster, error) {
	cluster := Cluster{}
	queryParams := GetClusterOptsToQueryParams(opts, GetClusterOpts{IncludeDetails: True()})
	response, err := c.R().
		SetContext(ctx).
		SetResult(&cluster).
		Get(fmt.Sprintf("/admin/v1/clusters/%s%s", clusterIdentity, OptionalQueryParams(queryParams)))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	FixCluster(&cluster, cluster.CustomerSlug)
	return &cluster, nil
}
