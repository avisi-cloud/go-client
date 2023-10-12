package acloudapi

import (
	"context"
	"fmt"
)

func (c *clientImpl) GetMemberships(ctx context.Context) ([]Membership, error) {
	pagedResult, err := c.GetPaged(ctx, fmt.Sprintf("/api/v1/memberships"))
	if err != nil {
		return nil, err
	}

	return MarshalPagedResultContent[Membership](pagedResult)
}
