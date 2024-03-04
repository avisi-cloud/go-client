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

func (c *adminClientImpl) ListClusters(ctx context.Context, opts ...GetClusterOpts) ([]Cluster, error) {
	queryParams := GetClusterOptsToQueryParams(opts, GetClusterOpts{IncludeDetails: True(), ShowCompute: True(), HideDeleted: True()})
	all, err := c.GetPaged(ctx, fmt.Sprintf("/admin/v1/clusters%s", OptionalQueryParams(queryParams)))
	if err != nil {
		return nil, err
	}
	content, err := MarshalPagedResultContent[Cluster](all)
	if err != nil {
		return nil, err
	}
	for i, _ := range content {
		item := &content[i]
		FixCluster(item, item.CustomerSlug)
	}
	return content, nil
}

func (c *adminClientImpl) UpdateCluster(ctx context.Context, request AdminUpdateClusterRequest) (*Cluster, error) {
	cluster := Cluster{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&cluster).
		SetBody(&request).
		Put(fmt.Sprintf("/admin/v1/clusters/%s", request.ClusterIdentity))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	FixCluster(&cluster, cluster.CustomerSlug)
	return &cluster, nil
}

type AdminUpdateClusterRequest struct {
	ClusterIdentity string `json:"clusterIdentity"`
	Version         string `json:"version"`
}
