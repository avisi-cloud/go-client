package acloudapi

import (
	"context"
	"fmt"
	"time"
)

type Organisation struct {
	ID                                  string    `json:"id"`
	Name                                string    `json:"name"`
	VatCode                             string    `json:"vatCode"`
	CompanyName                         string    `json:"companyName"`
	CompanyAddress                      string    `json:"companyAddress"`
	ContactEmail                        string    `json:"contactEmail"`
	BillingEmail                        string    `json:"billingEmail"`
	VatCodeValidated                    bool      `json:"vatCodeValidated"`
	VatCodeValidatedAt                  time.Time `json:"vatCodeValidatedAt"`
	PhoneNumber                         string    `json:"phoneNumber"`
	CreatedAt                           time.Time `json:"createdAt"`
	AcceptedTerms                       bool      `json:"acceptedTerms"`
	AcceptedTermsAt                     time.Time `json:"acceptedTermsAt"`
	RestrictedToAvailableCloudProviders bool      `json:"restrictedToAvailableCloudProviders"`
	Type                                string    `json:"type"`
	Slug                                string    `json:"slug"`
}

func (c *clientImpl) GetOrganisation(ctx context.Context, organisationSlug string) (*Organisation, error) {
	organisation := Organisation{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&organisation).
		Get(fmt.Sprintf("/api/v1/organisations/%s", organisationSlug))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return &organisation, nil
}
