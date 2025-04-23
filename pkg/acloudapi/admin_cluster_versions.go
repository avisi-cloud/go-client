package acloudapi

import (
	"context"
	"time"
)

const (
	ClusterVersionsURL = "/admin/v1/cluster-versions"
)

func (c *adminClientImpl) ListClusterVersions(ctx context.Context) ([]AdminClusterVersion, error) {
	all, err := c.GetPaged(ctx, ClusterVersionsURL)
	if err != nil {
		return nil, err
	}
	return MarshalPagedResultContent[AdminClusterVersion](all)
}

func (c *adminClientImpl) ListAvailableClusterVersions(ctx context.Context) ([]AdminClusterVersion, error) {
	available, err := c.GetPaged(ctx, ClusterVersionsURL+"/available")
	if err != nil {
		return nil, err
	}
	return MarshalPagedResultContent[AdminClusterVersion](available)
}

func (c *adminClientImpl) ListHistoryClusterVersions(ctx context.Context) ([]AdminClusterVersion, error) {
	history, err := c.GetPaged(ctx, ClusterVersionsURL+"/history")
	if err != nil {
		return nil, err
	}
	return MarshalPagedResultContent[AdminClusterVersion](history)
}

func (c *adminClientImpl) GetClusterVersion(ctx context.Context, version string) (*AdminClusterVersion, error) {
	clusterVersion := AdminClusterVersion{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&clusterVersion).
		Get(ClusterVersionsURL + "/" + version)
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return &clusterVersion, nil
}

func (c *adminClientImpl) UpdateClusterVersion(ctx context.Context, version string, request AdminUpdateClusterVersionRequest) (*AdminClusterVersion, error) {
	clusterVersion := AdminClusterVersion{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&clusterVersion).
		SetBody(&request).
		Put(ClusterVersionsURL + "/" + version)
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return &clusterVersion, nil
}

func (c *adminClientImpl) CreateClusterVersion(ctx context.Context, request AdminCreateClusterVersionRequest) (*AdminClusterVersion, error) {
	clusterVersion := AdminClusterVersion{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&clusterVersion).
		SetBody(&request).
		Post(ClusterVersionsURL)
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return &clusterVersion, nil
}

func (c *adminClientImpl) DeleteClusterVersion(ctx context.Context, version string) error {
	response, err := c.R().
		SetContext(ctx).
		Delete(ClusterVersionsURL + "/" + version)
	if err := c.CheckResponse(response, err); err != nil {
		return err
	}
	return nil
}

type AdminClusterVersion struct {
	Version                  string     `json:"version"`
	KubernetesVersion        string     `json:"kubernetesVersion"`
	ClusterControllerVersion string     `json:"clusterControllerVersion"`
	AddonControllerVersion   string     `json:"addonControllerVersion"`
	Available                bool       `json:"available"`
	CreatedAt                time.Time  `json:"createdAt"`
	ModifiedAt               time.Time  `json:"modifiedAt"`
	DeletedAt                *time.Time `json:"deletedAt,omitempty"`
	Note                     string     `json:"note"`
	ClusterCount             int64      `json:"clusterCount"`
}

type AdminCreateClusterVersionRequest struct {
	Version                  string `json:"version"`
	KubernetesVersion        string `json:"kubernetesVersion"`
	ClusterControllerVersion string `json:"clusterControllerVersion"`
	AddonControllerVersion   string `json:"addonControllerVersion,omitempty"`
	Available                bool   `json:"available"`
	Note                     string `json:"note,omitempty"`
}

type AdminUpdateClusterVersionRequest struct {
	Available bool `json:"available"`
}
