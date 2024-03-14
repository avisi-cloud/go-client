package acloudapi

import (
	"context"
	"fmt"
)

func (c *adminClientImpl) ListScheduledClusterUpgrades(ctx context.Context, opts ...ListScheduledClusterUpgradesOpts) ([]ScheduledClusterUpgrade, error) {
	queryParams := ListScheduledClusterUpgradesOptsToQueryParams(opts, ListScheduledClusterUpgradesOpts{})
	all, err := c.GetPaged(ctx, fmt.Sprintf("/admin/v1/scheduled-cluster-upgrades%s", OptionalQueryParams(queryParams)))
	if err != nil {
		return nil, err
	}
	return MarshalPagedResultContent[ScheduledClusterUpgrade](all)
}

func (c *adminClientImpl) GetScheduledClusterUpgrade(ctx context.Context, identity string) (*ScheduledClusterUpgrade, error) {
	scheduledClusterUpgrade := ScheduledClusterUpgrade{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&scheduledClusterUpgrade).
		Get(fmt.Sprintf("/admin/v1/scheduled-cluster-upgrades/%s", identity))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return &scheduledClusterUpgrade, nil
}

func (c *adminClientImpl) CancelScheduledClusterUpgrade(ctx context.Context, identity string) (*ScheduledClusterUpgrade, error) {
	scheduledClusterUpgrade := ScheduledClusterUpgrade{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&scheduledClusterUpgrade).
		Delete(fmt.Sprintf("/admin/v1/scheduled-cluster-upgrades/%s", identity))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return &scheduledClusterUpgrade, nil
}

func (c *adminClientImpl) CreateScheduledClusterUpgrade(ctx context.Context, request CreateScheduledClusterUpgradeRequest) (*ScheduledClusterUpgrade, error) {
	scheduledClusterUpgrade := ScheduledClusterUpgrade{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&scheduledClusterUpgrade).
		SetBody(&request).
		Post("/admin/v1/scheduled-cluster-upgrades")
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return &scheduledClusterUpgrade, nil
}

func (c *adminClientImpl) UpdateScheduledClusterUpgrade(ctx context.Context, request UpdateScheduledClusterUpgradeRequest) (*ScheduledClusterUpgrade, error) {
	scheduledClusterUpgrade := ScheduledClusterUpgrade{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&scheduledClusterUpgrade).
		SetBody(&request).
		Put(fmt.Sprintf("/admin/v1/scheduled-cluster-upgrades/%s", request.Identity))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return &scheduledClusterUpgrade, nil
}
