package acloudapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func (c *clientImpl) GetSilences(ctx context.Context, org, observabilityTenantSlug string) ([]Silence, error) {
	silences := []Silence{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&silences).
		Get(fmt.Sprintf("/api/v1/orgs/%s/observability/%s/silences", org, observabilityTenantSlug))
	if err := c.CheckResponse(response, err); err != nil {
		if response.StatusCode() == http.StatusNotFound {
			return silences, nil
		}
		return nil, err
	}
	return silences, nil
}

func (c *clientImpl) CreateSilence(ctx context.Context, createSilence CreateSilence, org, observabilityTenantSlug string) (*Silence, error) {
	if strings.TrimSpace(createSilence.Comment) == "" {
		return nil, fmt.Errorf("comment is required")
	}
	if len(createSilence.Matchers) == 0 {
		return nil, fmt.Errorf("matchers must not be empty")
	}
	if createSilence.EndsAt.Before(createSilence.StartsAt) {
		return nil, fmt.Errorf("endsAt must be after startAt")
	}

	silence := Silence{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&silence).
		SetBody(&createSilence).
		Post(fmt.Sprintf("/api/v1/orgs/%s/observability/%s/silences", org, observabilityTenantSlug))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return &silence, nil
}

func (c *clientImpl) ExpireSilence(ctx context.Context, org, observabilityTenantSlug, silenceID string) error {
	response, err := c.R().
		SetContext(ctx).
		Delete(fmt.Sprintf("/api/v1/orgs/%s/observability/%s/silences/%s", org, observabilityTenantSlug, silenceID))
	if err := c.CheckResponse(response, err); err != nil {
		return err
	}
	return nil
}
