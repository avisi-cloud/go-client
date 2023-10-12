package acloudapi

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
)

func (c *clientImpl) GetNodePools(ctx context.Context) ([]NodePool, error) {
	clusters, err := c.GetClusters(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get clusters: %w", err)
	}

	return c.GetNodePoolsByClusters(ctx, clusters)
}

func (c *clientImpl) GetNodePoolsByOrg(ctx context.Context, organisationSlug string) ([]NodePool, error) {
	clusters, err := c.GetClustersByOrg(ctx, organisationSlug)
	if err != nil {
		return nil, fmt.Errorf("failed to get clusters by organisation %q: %w", organisationSlug, err)
	}

	return c.GetNodePoolsByClusters(ctx, clusters)
}

type nodePoolResponse struct {
	NodePools []NodePool
	err       error
}

func (c *clientImpl) GetNodePoolsByClusters(ctx context.Context, clusters []Cluster) ([]NodePool, error) {
	numberOfRequests := len(clusters)
	responseChan := make(chan nodePoolResponse, numberOfRequests)
	wg := sync.WaitGroup{}
	for _, cluster := range clusters {
		wg.Add(1)
		go func(cls Cluster) {
			nodePoolsByCluster, err := c.GetNodePoolsByCluster(ctx, cls)
			resp := nodePoolResponse{
				NodePools: nodePoolsByCluster,
			}
			if err != nil {
				resp.err = fmt.Errorf("failed to get node-pools by cluster %s: %w", cls.FullIdentifier(), err)
			}
			responseChan <- resp
			wg.Done()
		}(cluster)
	}

	wg.Wait()
	close(responseChan)
	nodePools := make([]NodePool, 0)
	for i := 0; i < numberOfRequests; i++ {
		response := <-responseChan
		if response.err != nil {
			return nil, response.err // returning first error
		}

		nodePools = append(nodePools, response.NodePools...)
	}
	SortNodePools(nodePools)
	return nodePools, nil
}

func (c *clientImpl) GetNodePoolsByCluster(ctx context.Context, cluster Cluster) ([]NodePool, error) {
	var nodePools []NodePool
	response, err := c.R().
		SetContext(ctx).
		SetResult(&nodePools).
		Get(fmt.Sprintf("/api/v1/orgs/%s/clusters/%s/%s/pools", cluster.CustomerSlug, cluster.EnvironmentSlug, cluster.Slug))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	for i := range nodePools { // we need to modify the nodePool in the Slice
		c.mapNodePool(&nodePools[i], cluster)
	}
	SortNodePools(nodePools)
	return nodePools, nil
}

func (c *clientImpl) mapNodePool(nodePool *NodePool, cluster Cluster) {
	// add a reference to the cluster
	nodePool.Cluster = cluster
	nodePool.ClusterIdentity = cluster.Identity
}

func (c *clientImpl) CreateNodePool(ctx context.Context, cluster Cluster, create CreateNodePool) (*NodePool, error) {
	nodePool := NodePool{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&nodePool).
		SetBody(&create).
		Post(fmt.Sprintf("/api/v1/orgs/%s/clusters/%s/%s/pools", cluster.CustomerSlug, cluster.EnvironmentSlug, cluster.Slug))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	c.mapNodePool(&nodePool, cluster)
	return &nodePool, nil
}

func (c *clientImpl) UpdateNodePool(ctx context.Context, cluster Cluster, nodePoolID int, update CreateNodePool) (*NodePool, error) {
	nodePool := NodePool{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&nodePool).
		SetBody(&update).
		Put(fmt.Sprintf("/api/v1/orgs/%s/clusters/%s/%s/pools/%d", cluster.CustomerSlug, cluster.EnvironmentSlug, cluster.Slug, nodePoolID))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	c.mapNodePool(&nodePool, cluster)
	return &nodePool, nil
}

func (c *clientImpl) DeleteNodePool(ctx context.Context, cluster Cluster, nodePoolID int) error {
	response, err := c.R().
		SetContext(ctx).
		Delete(fmt.Sprintf("/api/v1/orgs/%s/clusters/%s/%s/pools/%d", cluster.CustomerSlug, cluster.EnvironmentSlug, cluster.Slug, nodePoolID))
	if err := c.CheckResponse(response, err); err != nil {
		return err
	}
	return nil
}

// SortNodePools sorts a slice of NodePools by CustomerSlug, EnvironmentSlug, ClusterSlug and then NodePool name
func SortNodePools(nodePools []NodePool) {
	sort.SliceStable(nodePools, func(i, j int) bool {
		return compareNodePool(&nodePools[i], &nodePools[j])
	})
}

func compareNodePool(n1 *NodePool, n2 *NodePool) bool {
	c1 := n1.Cluster
	c2 := n2.Cluster
	switch strings.Compare(c1.CustomerSlug, c2.CustomerSlug) {
	case -1:
		return true
	case 1:
		return false
	}
	switch strings.Compare(c1.EnvironmentSlug, c2.EnvironmentSlug) {
	case -1:
		return true
	case 1:
		return false
	}
	switch strings.Compare(c1.Slug, c2.Slug) {
	case -1:
		return true
	case 1:
		return false
	}
	return n1.Name < n2.Name
}
