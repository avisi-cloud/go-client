package acloudapi

import (
	"context"
	"fmt"
)

type CloudAccount struct {
	Identity                        string               `json:"identity"`
	DisplayName                     string               `json:"displayName"`
	Metadata                        CloudAccountMetadata `json:"metadata"`
	CloudProfile                    CloudProfile         `json:"cloudProfile"`
	Enabled                         bool                 `json:"enabled"`
	PrimaryCloudCredentialsIdentity string               `json:"primaryCloudCredentialsIdentity"`
}

type CloudAccountMetadata struct {
	// VsphereParentResourcePool is the parent resource pool for the vSphere cloud account. Required for Vsphere
	VSphereParentResourcePool *string `json:"parentResourcePool,omitempty"`
	// VSphereParentFolder is the parent folder for the vSphere cloud account. Required for Vsphere
	VsphereParentFolder *string `json:"parentFolder,omitempty"`

	// OpenStackTenantID is the ID of the OpenStack tenant. Optional
	OpenStackTenantID *string `json:"tenantId,omitempty"`
}

type CreateCloudAccount struct {
	// DisplayName is the name of the cloud account
	DisplayName string `json:"displayName"`
	// CloudProfile is the identity of the cloud profile to use
	CloudProfile string `json:"cloudProfile"`

	// Metadata is a map of additional information for the cloud account
	// See https://docs.avisi.cloud for more information
	Metadata CloudAccountMetadata `json:"metadata"`
}

type UpdateCloudAccount struct {
	// DisplayName is the name of the cloud account
	DisplayName string `json:"displayName"`

	// Enabled is a flag to enable or disable the cloud account
	Enabled bool `json:"enabled"`

	// PrimaryCloudCredentials is the identity of the primary cloud credentials to use
	PrimaryCloudCredentials string `json:"primaryCloudCredentials"`
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
		return nil, fmt.Errorf("cloud account %q for cloud provider %q not found", name, cloudProvider)
	}
	result := cloudAccounts[0]

	return &result, nil
}

func (c *clientImpl) CreateCloudAccount(ctx context.Context, org string, create CreateCloudAccount) (*CloudAccount, error) {
	cloudAccount := CloudAccount{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&cloudAccount).
		SetBody(&create).
		Post(fmt.Sprintf("/api/v1/orgs/%s/cloud-accounts", org))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return &cloudAccount, nil
}

// UpdateCloudAccount updates the cloud account with the given identity
func (c *clientImpl) UpdateCloudAccount(ctx context.Context, org, identity string, update UpdateCloudAccount) (*CloudAccount, error) {
	cloudAccount := CloudAccount{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&cloudAccount).
		SetBody(&update).
		Patch(fmt.Sprintf("/api/v1/orgs/%s/cloud-accounts/%s", org, identity))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return &cloudAccount, nil
}

// DeleteCloudAccount deletes the cloud account with the given identity
func (c *clientImpl) DeleteCloudAccount(ctx context.Context, org, identity string) error {
	response, err := c.R().
		SetContext(ctx).
		Delete(fmt.Sprintf("/api/v1/orgs/%s/cloud-accounts/%s", org, identity))
	if err := c.CheckResponse(response, err); err != nil {
		return err
	}
	return nil
}

// GetCloudProfiles returns a list of cloud profiles for the given organization
// The cloud profiles are used to create cloud accounts
// The cloud profiles contain information about the cloud provider and the regions available
func (c *clientImpl) GetCloudProfiles(ctx context.Context, org string) ([]CloudProfile, error) {
	cloudProfiles := []CloudProfile{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&cloudProfiles).
		Get(fmt.Sprintf("/api/v1/orgs/%s/cloud-profiles", org))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return cloudProfiles, nil
}
