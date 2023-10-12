package acloudapi

import "context"

func (c *clientImpl) GetServiceLevelAgreements(ctx context.Context) ([]ServiceLevelAgreement, error) {
	var slas []ServiceLevelAgreement
	response, err := c.R().
		SetContext(ctx).
		SetResult(&slas).
		Get("/api/v1/service-level-agreement")
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return slas, nil
}
