package acloudapi

import (
	"context"
	"sort"
)

func (c *clientImpl) GetClusterVersions(ctx context.Context) ([]ClusterVersion, error) {
	var clusterVersions []ClusterVersion
	response, err := c.R().
		SetContext(ctx).
		SetResult(&clusterVersions).
		Get("/api/v1/cluster-versions")
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}

	// sort result
	sort.SliceStable(clusterVersions, func(i, j int) bool {
		return clusterVersions[i].Version < clusterVersions[j].Version
	})
	return clusterVersions, nil
}
