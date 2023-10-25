package acloudapi

import (
	"context"
	"fmt"
)

func (c *clientImpl) GetNodePoolJoinConfig(ctx context.Context, cluster Cluster, nodePool NodePool) (*NodePoolJoinConfig, error) {
	joinConfig := NodePoolJoinConfig{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&joinConfig).
		Post(fmt.Sprintf("/api/v1/orgs/%s/clusters/%s/%s/pools/%s/join-config", cluster.CustomerSlug, cluster.EnvironmentSlug, cluster.Slug, nodePool.Identity))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return &joinConfig, nil
}
