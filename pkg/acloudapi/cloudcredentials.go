package acloudapi

import (
	"context"
	"fmt"
	"time"
)

type CloudCredential struct {
	Identity             string            `json:"identity"`
	CloudAccountIdentity string            `json:"cloudAccountIdentity"`
	CloudType            CloudProviderType `json:"cloudType"`
	DisplayName          string            `json:"displayName"`
	Metadata             map[string]string `json:"metadata"`
	IsPrimary            bool              `json:"isPrimary"`
	CreatedAt            time.Time         `json:"createdAt"`
}

type CloudCredentialCredentials interface {
	GetCloudProviderType() CloudProviderType
}

type CloudCredentialAWS struct {
	AccessKeyID     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
}

func (c *CloudCredentialAWS) GetCloudProviderType() CloudProviderType {
	return CloudProviderAWS
}

type CloudCredentialAzure struct {
	ClientID       string `json:"clientId"`
	ClientSecret   string `json:"clientSecret"`
	SubscriptionID string `json:"subscriptionId"`
	TenantID       string `json:"tenantId"`
}

func (c *CloudCredentialAzure) GetCloudProviderType() CloudProviderType {
	return CloudProviderAzure
}

type CloudCredentialDigitalOcean struct {
	APIKey string `json:"apiKey"`
}

func (c *CloudCredentialDigitalOcean) GetCloudProviderType() CloudProviderType {
	return CloudProviderDigitalOcean
}

type CloudCredentialHetzner struct {
	APIKey string `json:"apiKey"`
}

func (c *CloudCredentialHetzner) GetCloudProviderType() CloudProviderType {
	return CloudProviderHetzner
}

type CloudCredentialOpenStack struct {
	CredentialID     string `json:"credentialId"`
	CredentialSecret string `json:"credential"`
}

func (c *CloudCredentialOpenStack) GetCloudProviderType() CloudProviderType {
	return CloudProviderOpenstack
}

type CloudCredentialVSphere struct {
	Username          string `json:"username"`
	Password          string `json:"password"`
	ProvisionUsername string `json:"provisionUsername,omitempty"`
	ProvisionPassword string `json:"provisionPassword,omitempty"`
}

func (c *CloudCredentialVSphere) GetCloudProviderType() CloudProviderType {
	return CloudProviderVSphere
}

// CreateCloudCredential is the struct for creating a new cloud credential
type CreateCloudCredential struct {
	// DisplayName is the name of the cloud account
	DisplayName string `json:"displayName"`
	// Credentials holds the cloud account credentials struct for the cloud type
	Credentials CloudCredentialCredentials `json:"credentials"`
}

func (c *CreateCloudCredential) Validate() error {
	if c.DisplayName == "" {
		return fmt.Errorf("displayName is missing")
	}
	if c.Credentials == nil {
		return fmt.Errorf("credentials are missing")
	}
	if c.Credentials.GetCloudProviderType() == "" {
		return fmt.Errorf("cloud provider type is missing")
	}
	return nil
}

func (c *clientImpl) GetCloudCredentials(ctx context.Context, org, cloudAccountIdentity string) ([]CloudCredential, error) {
	pagedResult, err := c.GetPaged(ctx, fmt.Sprintf("/api/v1/orgs/%s/cloud-accounts/%s/credentials", org, cloudAccountIdentity))
	if err != nil {
		return nil, err
	}
	return MarshalPagedResultContent[CloudCredential](pagedResult)
}

// CreateCloudCredential creates a new cloud credential within a cloud account. The Cloud Account must exist and the cloud type must be valid
// The cloud account identity is the unique identifier of the cloud account.
//
// Ensure that the cloud account exists before calling this function, and that the cloud type is valid.
// The cloud type is the type of the cloud account, e.g. aws, azure, digitalocean, hetzner, openstack, vsphere
//
// The create struct must contain the displayName and the credentials for the cloud type. The credentials struct must match the cloud type.
func (c *clientImpl) CreateCloudCredential(ctx context.Context, org string, cloudAccount CloudAccount, create CreateCloudCredential) (*CloudCredential, error) {
	cloudType := cloudAccount.CloudProfile.Type

	if cloudAccount.Identity == "" {
		return nil, fmt.Errorf("cloud account identity is missing")
	}
	if cloudType == "" {
		return nil, fmt.Errorf("cloud account %q has no cloud type", cloudAccount.Identity)
	}
	if err := create.Validate(); err != nil {
		return nil, err
	}

	if string(create.Credentials.GetCloudProviderType()) != cloudType {
		return nil, fmt.Errorf("cloud provider type %q does not match cloud account type %q", create.Credentials.GetCloudProviderType(), cloudType)
	}

	credentials := CloudCredential{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&cloudAccount).
		SetBody(&create).
		Post(fmt.Sprintf("/api/v1/orgs/%s/cloud-accounts/%s/credentials/%s", org, cloudAccount.Identity, cloudType))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return &credentials, nil
}

// DeleteCloudCredential deletes the cloud account with the given identity
func (c *clientImpl) DeleteCloudCredential(ctx context.Context, org, cloudAccountIdentity, credentialsIdentity string) error {
	response, err := c.R().
		SetContext(ctx).
		Delete(fmt.Sprintf("/api/v1/orgs/%s/cloud-accounts/%s/credentials/%s", org, cloudAccountIdentity, credentialsIdentity))
	if err := c.CheckResponse(response, err); err != nil {
		return err
	}
	return nil
}
