package acloudapi

import (
	"context"
	"fmt"
)

func (c *clientImpl) GetEnvironments(ctx context.Context, org string) ([]Environment, error) {
	pagedResult, err := c.GetPaged(ctx, fmt.Sprintf("/api/v1/orgs/%s/environments?show-compute=true", org))
	if err != nil {
		return nil, err
	}
	return MarshalPagedResultContent[Environment](pagedResult)
}

func (c *clientImpl) GetEnvironment(ctx context.Context, org, env string) (*Environment, error) {
	environment := Environment{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&environment).
		Get(fmt.Sprintf("/api/v1/orgs/%s/environments/%s?show-compute=true", org, env))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return &environment, nil
}

func (c *clientImpl) CreateEnvironment(ctx context.Context, createEnvironment CreateEnvironment, org string) (*Environment, error) {
	environment := Environment{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&environment).
		SetBody(&createEnvironment).
		Post(fmt.Sprintf("/api/v1/orgs/%s/environments", org))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return &environment, nil
}

func (c *clientImpl) UpdateEnvironment(ctx context.Context, updateEnvironment UpdateEnvironment, org, env string) (*Environment, error) {
	environment := Environment{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&environment).
		SetBody(&updateEnvironment).
		Patch(fmt.Sprintf("/api/v1/orgs/%s/environments/%s", org, env))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return &environment, nil
}

func (c *clientImpl) DeleteEnvironment(ctx context.Context, org, env string) error {
	response, err := c.R().
		SetContext(ctx).
		Delete(fmt.Sprintf("/api/v1/orgs/%s/environments/%s", org, env))
	if err := c.CheckResponse(response, err); err != nil {
		return err
	}
	return nil
}
