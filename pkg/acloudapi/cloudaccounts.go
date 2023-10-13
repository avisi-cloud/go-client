package acloudapi

import (
	"context"
	"fmt"
)

type CloudAccount struct {
	Identity                        string            `json:"identity"`
	DisplayName                     string            `json:"displayName"`
	Metadata                        map[string]string `json:"metadata"`
	CloudProfile                    CloudProfile      `json:"cloudProfile"`
	Enabled                         bool              `json:"enabled"`
	PrimaryCloudCredentialsIdentity string            `json:"primaryCloudCredentialsIdentity"`
}

type CloudProfile struct {
	Identity      string            `json:"identity"`
	DisplayName   string            `json:"displayName"`
	Metadata      map[string]string `json:"metadata"`
	CloudProvider string            `json:"cloudProvider"`
	Regions       []string          `json:"regions"`
	Enabled       bool              `json:"enabled"`
	Public        bool              `json:"public"`
	Type          string            `json:"type"`

	CloudProviderResponse        CloudProvider `json:"cloudProviderResponse"`
	CloudProviderRegionResponses []Region      `json:"cloudProviderRegionResponses"`
}

func (c *clientImpl) GetCloudAccounts(ctx context.Context, org string) ([]CloudAccount, error) {
	pagedResult, err := c.GetPaged(ctx, fmt.Sprintf("/api/v1/orgs/%s/cloud-accounts", org))
	if err != nil {
		return nil, err
	}
	return MarshalPagedResultContent[CloudAccount](pagedResult)
}

func (c *clientImpl) FindCloudAccountByName(ctx context.Context, org, name, cloudProvider string) (*CloudAccount, error) {
	pagedResult, err := c.GetPaged(ctx, fmt.Sprintf("/api/v1/orgs/%s/cloud-accounts?display-name=%s&cloud-provider-slug=%s", org, name, cloudProvider))
	if err != nil {
		return nil, err
	}

	if pagedResult.TotalElements > 1 {
		return nil, fmt.Errorf("ambiguous results, expected 1 cloudaccount, got %d", pagedResult.TotalElements)
	}

	cloudAccounts, err := MarshalPagedResultContent[CloudAccount](pagedResult)
	if err != nil {
		return nil, err
	}
	if len(cloudAccounts) == 0 {
		return nil, fmt.Errorf("cloud account %s for cloud provider %s not found", name, cloudProvider)
	}
	result := cloudAccounts[0]

	return &result, nil
}
