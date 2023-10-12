package acloudapi

import (
	"context"
	"fmt"
	"time"
)

func (c *adminClientImpl) GetOrganisation(ctx context.Context, organisationIdentity string) (*AdminOrganisation, error) {
	organisation := AdminOrganisation{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&organisation).
		Get(fmt.Sprintf("/admin/v1/orgs/%s", organisationIdentity))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return &organisation, nil
}

type AdminOrganisation struct {
	ID                                  string    `json:"id"`
	Name                                string    `json:"name"`
	VatCode                             string    `json:"vatCode"`
	ContactEmail                        string    `json:"contactEmail"`
	BillingEmail                        string    `json:"billingEmail"`
	VatCodeValidated                    bool      `json:"vatCodeValidated"`
	VatCodeValidatedAt                  time.Time `json:"vatCodeValidatedAt"`
	PhoneNumber                         string    `json:"phoneNumber"`
	CreatedAt                           time.Time `json:"createdAt"`
	AcceptedTerms                       bool      `json:"acceptedTerms"`
	AcceptedTermsAt                     time.Time `json:"acceptedTermsAt"`
	RestrictedToAvailableCloudProviders bool      `json:"restrictedToAvailableCloudProviders"`
	StripeCustomer                      string    `json:"stripeCustomer"`
	Slug                                string    `json:"slug"`
}
