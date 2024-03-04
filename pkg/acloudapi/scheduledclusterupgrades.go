package acloudapi

import (
	"fmt"
	"strings"
)

func ListScheduledClusterUpgradesOptsToQueryParams(opts []ListScheduledClusterUpgradesOpts, defaults ListScheduledClusterUpgradesOpts) string {
	return toQueryParamsListScheduledClusterUpgradesOpts(mergeListScheduledClusterUpgradesOpts(opts, defaults))
}

func mergeListScheduledClusterUpgradesOpts(opts []ListScheduledClusterUpgradesOpts, defaults ListScheduledClusterUpgradesOpts) ListScheduledClusterUpgradesOpts {
	merged := ListScheduledClusterUpgradesOpts{}
	for _, opt := range opts {
		if opt.ClusterIdentities != nil {
			merged.ClusterIdentities = opt.ClusterIdentities
		}
		if opt.Statuses != nil {
			merged.Statuses = opt.Statuses
		}
	}
	return setDefaultsListScheduledClusterUpgradesOpts(merged, defaults)
}

func setDefaultsListScheduledClusterUpgradesOpts(merged ListScheduledClusterUpgradesOpts, defaults ListScheduledClusterUpgradesOpts) ListScheduledClusterUpgradesOpts {
	if merged.ClusterIdentities == nil {
		merged.ClusterIdentities = defaults.ClusterIdentities
	}
	if merged.Statuses == nil {
		merged.Statuses = defaults.Statuses
	}
	return merged
}

func toQueryParamsListScheduledClusterUpgradesOpts(opts ListScheduledClusterUpgradesOpts) string {
	url := ""
	if opts.ClusterIdentities != nil {
		url = fmt.Sprintf("clusterIdentities=%s", strings.Join(opts.ClusterIdentities, ","))
	}
	if opts.Statuses != nil {
		if len(url) > 0 {
			url = fmt.Sprintf("%s&", url)
		}
		stringSlice := make([]string, len(opts.Statuses))
		for i, item := range opts.Statuses {
			stringSlice[i] = string(item)
		}
		url = fmt.Sprintf("%sstatuses=%s", url, strings.Join(stringSlice, ","))
	}
	return url
}
