package acloudapi

import (
	"context"
	"fmt"
)

type AddUserResponse struct {
	OidcId                string `json:"oidcId" yaml:"OidcId"`
	Identity              string `json:"identity" yaml:"Identity"`
	Username              string `json:"username" yaml:"Username"`
	Email                 string `json:"email" yaml:"Email"`
	CreatedAt             string `json:"createdAt" yaml:"CreatedAt"`
	CanCreateOrganisation bool   `json:"canCreateOrganisation" yaml:"CanCreateOrganisation"`
}

// AddUser associates your existing OpenID Connect user ID with your Avisi Cloud organization,
// linking your account without creating a new user.
func (c *clientImpl) AddUser(ctx context.Context) (AddUserResponse, error) {
	result := AddUserResponse{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&result).
		Post(fmt.Sprintf("/api/v1/user/add"))
	if err := c.CheckResponse(response, err); err != nil {
		return AddUserResponse{}, err
	}
	return result, nil
}
