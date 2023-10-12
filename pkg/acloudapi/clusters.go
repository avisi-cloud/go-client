package acloudapi

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
)

type clustersResponse struct {
	clusters []Cluster
	err      error
}

func (c *clientImpl) GetClusters(ctx context.Context, opts ...GetClusterOpts) ([]Cluster, error) {
	memberships, err := c.GetMemberships(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get memberships/organisations: %v", err)
	}

	numberOfRequests := len(memberships)
	responseChan := make(chan clustersResponse, numberOfRequests)
	wg := sync.WaitGroup{}
	for _, membership := range memberships {
		wg.Add(1)
		go func(mbs Membership) {
			clustersByOrg, err := c.GetClustersByOrg(ctx, mbs.Slug, opts...)
			resp := clustersResponse{
				clusters: clustersByOrg,
			}
			if err != nil {
				resp.err = fmt.Errorf("failed to get node-pools by cluster %s: %w", mbs.Slug, err)
			}
			responseChan <- resp
			wg.Done()
		}(membership)
	}

	wg.Wait()
	close(responseChan)
	clusters := make([]Cluster, 0)
	for i := 0; i < numberOfRequests; i++ {
		response := <-responseChan
		if response.err != nil {
			return nil, response.err // returning first error
		}

		clusters = append(clusters, response.clusters...)
	}
	SortClusters(clusters)

	return clusters, nil
}

func (c *clientImpl) GetClustersByOrg(ctx context.Context, org string, opts ...GetClusterOpts) ([]Cluster, error) {
	queryParams := GetClusterOptsToQueryParams(opts, GetClusterOpts{ShowCompute: True()})
	pagedResult, err := c.GetPaged(ctx, fmt.Sprintf("/api/v1/orgs/%s/clusters%s", org, OptionalQueryParams(queryParams)))
	if err != nil {
		return nil, err
	}
	return c.marshalClustersFromPagedResult(org, pagedResult)
}

func (c *clientImpl) GetClustersByOrgAndEnv(ctx context.Context, org, env string, opts ...GetClusterOpts) ([]Cluster, error) {
	queryParams := GetClusterOptsToQueryParams(opts, GetClusterOpts{ShowCompute: True()})
	pagedResult, err := c.GetPaged(ctx, fmt.Sprintf("/api/v1/orgs/%s/clusters/%s%s", org, env, OptionalQueryParams(queryParams)))
	if err != nil {
		return nil, err
	}
	return c.marshalClustersFromPagedResult(org, pagedResult)
}

func (c *clientImpl) GetCluster(ctx context.Context, org, env, clusterSlug string, opts ...GetClusterOpts) (*Cluster, error) {
	queryParams := GetClusterOptsToQueryParams(opts, GetClusterOpts{ShowCompute: True()})
	cluster := Cluster{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&cluster).
		Get(fmt.Sprintf("/api/v1/orgs/%s/clusters/%s/%s%s", org, env, clusterSlug, OptionalQueryParams(queryParams)))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	FixCluster(&cluster, org)
	return &cluster, nil
}

func (c *clientImpl) GetClusterOIDCConfig(ctx context.Context, org, env, cluster string) (*ClusterMetadataResponse, error) {
	result := ClusterMetadataResponse{}

	response, err := c.R().
		SetContext(ctx).
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/orgs/%s/clusters/%s/%s/oidc-config", org, env, cluster))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *clientImpl) CreateCluster(ctx context.Context, organisationSlug, environmentSlug string, create CreateCluster) (*Cluster, error) {
	cluster := Cluster{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&cluster).
		SetBody(&create).
		Post(fmt.Sprintf("/api/v1/orgs/%s/clusters/%s", organisationSlug, environmentSlug))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	FixCluster(&cluster, organisationSlug)
	return &cluster, nil
}

func (c *clientImpl) UpdateCluster(ctx context.Context, org, env, clusterSlug string, update UpdateCluster) (*Cluster, error) {
	cluster := Cluster{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&cluster).
		SetBody(&update).
		Patch(fmt.Sprintf("/api/v1/orgs/%s/clusters/%s/%s", org, env, clusterSlug))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	FixCluster(&cluster, org)
	return &cluster, nil
}

func (c *clientImpl) DeleteCluster(ctx context.Context, org, env, slug string, update UpdateCluster) error {
	response, err := c.R().
		SetContext(ctx).
		SetBody(&update).
		Patch(fmt.Sprintf("/api/v1/orgs/%s/clusters/%s/%s", org, env, slug))
	if err := c.CheckResponse(response, err); err != nil {
		return err
	}
	return nil
}

func (c *clientImpl) marshalClustersFromPagedResult(org string, result PagedResult) ([]Cluster, error) {
	clusters, err := MarshalPagedResultContent[Cluster](result)
	if err != nil {
		return nil, err
	}

	for i := range clusters {
		FixCluster(&clusters[i], org)
	}
	SortClusters(clusters)
	return clusters, nil
}

// SortClusters sorts a slice of clusters on CustomerSlug, EnvironmentSlug and then ClusterSlug
func SortClusters(clusters []Cluster) {
	sort.SliceStable(clusters, func(i, j int) bool {
		return compareCluster(&clusters[i], &clusters[j])
	})
}

func compareCluster(c1 *Cluster, c2 *Cluster) bool {
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
	return c1.Slug < c2.Slug
}

func FixCluster(cluster *Cluster, org string) {
	cluster.CustomerSlug = org                               // fix missing CustomerSlug/Org
	cluster.Status = FixStatus(cluster.Status)               // temporary map 'started' to 'running'
	cluster.DesiredStatus = FixStatus(cluster.DesiredStatus) // temporary map 'started' to 'running'
}

func FixStatus(status string) string {
	if status == "started" {
		return "running"
	}
	return status
}

type GetClusterOpts struct {
	IncludeDetails *bool
	ShowCompute    *bool
}

func OptionalQueryParams(queryParams string) string {
	if queryParams == "" {
		return ""
	}
	return "?" + queryParams
}

func GetClusterOptsToQueryParams(opts []GetClusterOpts, defaults GetClusterOpts) string {
	return toQueryParams(mergeGetClusterOpts(opts, defaults))
}

func mergeGetClusterOpts(opts []GetClusterOpts, defaults GetClusterOpts) GetClusterOpts {
	merged := GetClusterOpts{}
	for _, opt := range opts {
		if opt.IncludeDetails != nil {
			merged.IncludeDetails = opt.IncludeDetails
		}
		if opt.ShowCompute != nil {
			merged.ShowCompute = opt.ShowCompute
		}
	}
	return setDefaults(merged, defaults)
}

func setDefaults(merged GetClusterOpts, defaults GetClusterOpts) GetClusterOpts {
	if merged.IncludeDetails == nil {
		merged.IncludeDetails = defaults.IncludeDetails
	}
	if merged.ShowCompute == nil {
		merged.ShowCompute = defaults.IncludeDetails
	}
	return merged
}

func toQueryParams(mergedGetClusterOpts GetClusterOpts) string {
	url := ""
	if mergedGetClusterOpts.IncludeDetails != nil {
		url = fmt.Sprintf("includeDetails=%t", *mergedGetClusterOpts.IncludeDetails)
	}
	if mergedGetClusterOpts.ShowCompute != nil {
		if len(url) > 0 {
			url = fmt.Sprintf("%s&", url)
		}
		url = fmt.Sprintf("%sshow-compute=%t", url, *mergedGetClusterOpts.ShowCompute)
	}
	return url
}

func BoolPointer(b bool) *bool {
	return &b
}

func False() *bool {
	return BoolPointer(false)
}

func True() *bool {
	return BoolPointer(true)
}
