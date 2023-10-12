package acloudapi

import (
	"context"
	"fmt"
	"sort"
)

func (c *clientImpl) GetUpdateChannels(ctx context.Context, org string) ([]UpdateChannelResponse, error) {
	pagedResult, err := c.GetPaged(ctx, fmt.Sprintf("/api/v1/orgs/%s/update-channels", org))
	if err != nil {
		return nil, err
	}

	result, err := MarshalPagedResultContent[UpdateChannelResponse](pagedResult)
	if err != nil {
		return nil, err
	}

	sortUpdateChannel(result)

	return result, nil
}

func sortUpdateChannel(updateChannels []UpdateChannelResponse) {
	sort.SliceStable(updateChannels, func(i, j int) bool {
		left := updateChannels[i]
		right := updateChannels[j]

		return left.KubernetesClusterVersion > right.KubernetesClusterVersion
	})
}
