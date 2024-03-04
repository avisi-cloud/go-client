package acloudapi

import (
	"context"
)

func (c *adminClientImpl) ListUpdateChannels(ctx context.Context) ([]UpdateChannelResponse, error) {
	all, err := c.GetPaged(ctx, "/admin/v1/update-channels")
	if err != nil {
		return nil, err
	}
	return MarshalPagedResultContent[UpdateChannelResponse](all)
}
