package acloudapi

import (
	"context"
	"fmt"
)

func (c *clientImpl) GetCloudProviders(ctx context.Context, organisationSlug string) ([]CloudProvider, error) {
	pagedResult, err := c.GetPaged(ctx, fmt.Sprintf("/api/v1/orgs/%s/cloud-providers", organisationSlug))
	if err != nil {
		return nil, err
	}
	return MarshalPagedResultContent[CloudProvider](pagedResult)
}

func (c *clientImpl) GetRegions(ctx context.Context, organisationSlug, cloudProviderSlug string) ([]Region, error) {
	pagedResult, err := c.GetPaged(ctx, fmt.Sprintf("/api/v1/orgs/%s/cloud-providers/%s/regions", organisationSlug, cloudProviderSlug))
	if err != nil {
		return nil, err
	}
	return MarshalPagedResultContent[Region](pagedResult)
}

func (c *clientImpl) GetAvailabilityZones(ctx context.Context, organisationSlug, cloudProviderSlug, regionSlug string) ([]AvailabilityZone, error) {
	pagedResult, err := c.GetPaged(ctx, fmt.Sprintf("/api/v1/orgs/%s/cloud-providers/%s/regions/%s/availability-zones", organisationSlug, cloudProviderSlug, regionSlug))
	if err != nil {
		return nil, err
	}
	return MarshalPagedResultContent[AvailabilityZone](pagedResult)
}

func (c *clientImpl) GetNodeTypes(ctx context.Context, cloudProviderSlug string) ([]NodeType, error) {
	pagedResult, err := c.GetPaged(ctx, fmt.Sprintf("/api/v1/cloud-providers/%s/nodetypes", cloudProviderSlug))
	if err != nil {
		return nil, err
	}
	return MarshalPagedResultContent[NodeType](pagedResult)
}
