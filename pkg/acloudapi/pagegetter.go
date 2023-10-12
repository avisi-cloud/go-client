package acloudapi

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	MAX_PAGING_LOOPS = 100
)

var (
	ErrMaximumPagingLoopsExceeded = fmt.Errorf("exceeded maximum paging loops")
)

type PageGetter interface {
	Get(ctx context.Context, url string, page int) (PagedResult, error)
}

type RestyPageGetter struct {
	client *RestyClient
}

func (rg *RestyPageGetter) Get(ctx context.Context, url string, page int) (PagedResult, error) {
	pagedResult := PagedResult{}

	response, err := rg.client.Resty().R().
		SetContext(ctx).
		SetQueryParam("page", strconv.Itoa(page)).
		SetResult(&pagedResult).
		Get(url)
	if err := rg.client.CheckResponse(response, err); err != nil {
		return pagedResult, err
	}

	return pagedResult, nil
}

func MarshalPagedResultContent[T any](pagedResult PagedResult) ([]T, error) {
	marshal, err := json.Marshal(pagedResult.Content)
	if err != nil {
		return nil, err
	}

	var content []T
	err = json.Unmarshal(marshal, &content)
	if err != nil {
		return nil, err
	}

	return content, nil
}
