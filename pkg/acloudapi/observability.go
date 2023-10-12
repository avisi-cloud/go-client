package acloudapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type ObservabilityTenant struct {
	Identity     string     `json:"identity" yaml:"Identity"`
	CustomerSlug string     `json:"customer" yaml:"Customer,omitempty"`
	Name         string     `json:"name" yaml:"Name"`
	Slug         string     `json:"slug" yaml:"Slug"`
	Available    bool       `json:"available" yaml:"Available"`
	IpWhiteList  string     `json:"ipWhiteList" yaml:"IPWhiteList,omitempty"`
	CreatedAt    time.Time  `json:"createdAt" yaml:"CreatedAt"`
	ModifiedAt   time.Time  `json:"modifiedAt" yaml:"ModifiedAt"`
	DeletedAt    *time.Time `json:"DeletedAt" yaml:"DeletedAt,omitempty"`
}

type ObservabilityAlert struct {
	Labels      map[string]string `json:"labels" yaml:"Labels"`
	Annotations map[string]string `json:"annotations" yaml:"Annotations"`
	ActiveAt    time.Time         `json:"activeAt" yaml:"ActiveAt"`
	State       string            `json:"state" yaml:"State"`
	Value       string            `json:"value" yaml:"Value"`
}

type ObservabilityAlertmanagerAndPrometheusrulesResponse struct {
	Rules              map[string]string `json:"ruleFiles" yaml:"Rules"`
	Templates          map[string]string `json:"templateFiles" yaml:"Templates"`
	AlertManagerConfig string            `json:"alertManagerConfig" yaml:"Alertmanager"`
}

type ObservabilityAlertmanager struct {
	Rules              map[string]string `json:"ruleFiles" yaml:"Rules"`
	Templates          map[string]string `json:"templateFiles" yaml:"Templates"`
	AlertManagerConfig string            `json:"alertManagerConfig" yaml:"Alertmanager"`
}

type PrometheusRules struct {
	Name    string
	Content []byte
}

func (c *clientImpl) GetObservabilityTenants(ctx context.Context, org string) ([]ObservabilityTenant, error) {
	pagedResult, err := c.GetPaged(ctx, fmt.Sprintf("/api/v1/orgs/%s/monitoring", org))
	if err != nil {
		return nil, err
	}
	return MarshalPagedResultContent[ObservabilityTenant](pagedResult)
}

func (c *clientImpl) GetObservabilityTenantBySlug(ctx context.Context, org, slug string) (*ObservabilityTenant, error) {
	cluster := ObservabilityTenant{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&cluster).
		Get(fmt.Sprintf("/api/v1/orgs/%s/monitoring/%s", org, slug))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return &cluster, nil
}

func (c *clientImpl) GetObservabilityTenantAlertmanagerConfiguration(ctx context.Context, org, slug string) (*ObservabilityAlertmanager, error) {
	alertmanagerConfigResponse, getResponse, err := getAlertManagerConfigResponse(ctx, org, slug, c)
	if err := c.CheckResponse(getResponse, err); err != nil {
		if getResponse.StatusCode() == http.StatusNotFound {
			return &ObservabilityAlertmanager{
				Rules:              alertmanagerConfigResponse.Rules,
				Templates:          alertmanagerConfigResponse.Templates,
				AlertManagerConfig: alertmanagerConfigResponse.AlertManagerConfig,
			}, nil
		}
		return nil, err
	}
	return &ObservabilityAlertmanager{
		Rules:              alertmanagerConfigResponse.Rules,
		Templates:          alertmanagerConfigResponse.Templates,
		AlertManagerConfig: alertmanagerConfigResponse.AlertManagerConfig,
	}, nil
}

func (c *clientImpl) AddObservabilityTenantPrometheusRules(ctx context.Context, org, slug string, rules []PrometheusRules, force bool) error {
	alertmanagerConfigResponse, getResponse, err := getAlertManagerConfigResponse(ctx, org, slug, c)
	if err := c.CheckResponse(getResponse, err); err != nil {
		return err
	}
	var existingRules []string
	for _, promRule := range rules {
		if !force {
			for currentRuleName := range alertmanagerConfigResponse.Rules {
				if promRule.Name == currentRuleName {
					existingRules = append(existingRules, promRule.Name)
				}
			}
		}
		alertmanagerConfigResponse.Rules[promRule.Name] = string(promRule.Content)
	}
	if !force && len(existingRules) > 0 {
		return fmt.Errorf("one or multiple rules already exists: %q\nuse --force if you want to go through", strings.Join(existingRules, ", "))
	}
	err = postNewAlertManagerConfig(ctx, org, slug, alertmanagerConfigResponse, c)
	if err != nil {
		return err
	}
	return nil
}

func (c *clientImpl) DeleteObservabilityTenantPrometheusRules(ctx context.Context, org, slug string, names []string) error {
	alertmanagerConfigResponse, getResponse, err := getAlertManagerConfigResponse(ctx, org, slug, c)
	if err := c.CheckResponse(getResponse, err); err != nil {
		return err
	}

	for _, name := range names {
		delete(alertmanagerConfigResponse.Rules, name)
	}

	err = postNewAlertManagerConfig(ctx, org, slug, alertmanagerConfigResponse, c)
	if err != nil {
		return err
	}
	return nil
}

func (c *clientImpl) OverwriteObservabilityTenantPrometheusRules(ctx context.Context, org, slug string, rules []PrometheusRules) error {
	alertmanagerConfigResponse, getResponse, err := getAlertManagerConfigResponse(ctx, org, slug, c)
	if err := c.CheckResponse(getResponse, err); err != nil {
		return err
	}
	// removes all the current rules
	alertmanagerConfigResponse.Rules = make(map[string]string)
	// sets all the new rules
	for _, promRule := range rules {
		alertmanagerConfigResponse.Rules[promRule.Name] = string(promRule.Content)
	}

	err = postNewAlertManagerConfig(ctx, org, slug, alertmanagerConfigResponse, c)
	if err != nil {
		return err
	}
	return nil
}

func (c *clientImpl) GetObservabilityOrganisationAlerts(ctx context.Context, org string) ([]ObservabilityAlert, error) {
	alerts := []ObservabilityAlert{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&alerts).
		Get(fmt.Sprintf("/api/v1/orgs/%s/alerts", org))
	if err := c.CheckResponse(response, err); err != nil {
		if response.StatusCode() == http.StatusNotFound {
			return alerts, nil
		}
		return nil, err
	}
	return alerts, nil
}

func (c *clientImpl) GetObservabilityTenantAlerts(ctx context.Context, org string, slug string) ([]ObservabilityAlert, error) {
	alerts := []ObservabilityAlert{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&alerts).
		Get(fmt.Sprintf("/api/v1/orgs/%s/alerts/%s", org, slug))
	if err := c.CheckResponse(response, err); err != nil {
		if response.StatusCode() == http.StatusNotFound {
			return alerts, nil
		}
		return nil, err
	}
	return alerts, nil
}

func getAlertManagerConfigResponse(ctx context.Context, org string, slug string, c *clientImpl) (ObservabilityAlertmanagerAndPrometheusrulesResponse, *resty.Response, error) {
	alertmanagerConfigResponse := ObservabilityAlertmanagerAndPrometheusrulesResponse{}
	getResponse, err := c.R().
		SetContext(ctx).
		SetResult(&alertmanagerConfigResponse).
		Get(fmt.Sprintf("/api/v1/orgs/%s/monitoring/%s/alertmanager", org, slug))
	return alertmanagerConfigResponse, getResponse, err
}

func postNewAlertManagerConfig(ctx context.Context, org string, slug string, alertmanagerConfigResponse ObservabilityAlertmanagerAndPrometheusrulesResponse, c *clientImpl) error {
	jsonConfig, err := json.Marshal(alertmanagerConfigResponse)
	if err != nil {
		return err
	}
	response, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonConfig).
		SetContext(ctx).
		Post(fmt.Sprintf("/api/v1/orgs/%s/monitoring/%s/alertmanager", org, slug))
	if err := c.CheckResponse(response, err); err != nil {
		return err
	}
	return nil
}
