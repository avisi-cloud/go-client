package acloudapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestGetClusterResponseFields(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/orgs/org1/clusters/env1/cluster1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"slug": "cluster1",
			"environmentSlug": "env1",
			"serviceSubnet": "10.96.0.0/12",
			"clusterPodSubnet": "10.244.0.0/16",
			"actualState": {"status": "running", "version": "v1.33.1", "collectedAt": "2026-10-16T08:59:36.000+00:00"},
			"cloudAccount": {
				"identity": "account1",
				"displayName": "Account 1",
				"primaryCloudCredentials": {"identity": "credentials1", "displayName": "Credentials 1", "isPrimary": true}
			},
			"cloudCredentials": {"identity": "credentials1", "displayName": "Credentials 1", "isPrimary": true},
			"upstreamCluster": {"identity": "upstream1", "slug": "upstream", "environmentSlug": "env1", "version": "v1.33.1"},
			"downstreamClusters": [{"identity": "downstream1", "slug": "downstream", "environmentSlug": "env2", "version": "v1.32.5"}]
		}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	c := &clientImpl{RestyClient: NewRestyClient(nil, ClientOpts{APIUrl: server.URL})}
	cluster, err := c.GetCluster(context.Background(), "org1", "env1", "cluster1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cluster.ServiceSubnet != "10.96.0.0/12" {
		t.Fatalf("unexpected serviceSubnet in response: %v", cluster.ServiceSubnet)
	}
	if cluster.ClusterPodSubnet != "10.244.0.0/16" {
		t.Fatalf("unexpected clusterPodSubnet in response: %v", cluster.ClusterPodSubnet)
	}
	wantActualState := &ActualClusterState{
		Status:      "running",
		Version:     "v1.33.1",
		CollectedAt: time.Date(2026, 10, 16, 8, 59, 36, 0, time.UTC),
	}
	if cluster.ActualState == nil || cluster.ActualState.Status != wantActualState.Status ||
		cluster.ActualState.Version != wantActualState.Version || !cluster.ActualState.CollectedAt.Equal(wantActualState.CollectedAt) {
		t.Fatalf("unexpected actualState in response: %+v", cluster.ActualState)
	}
	wantCredentials := &CloudCredentialsReference{Identity: "credentials1", DisplayName: "Credentials 1", IsPrimary: true}
	if cluster.CloudAccount == nil || !reflect.DeepEqual(cluster.CloudAccount.PrimaryCloudCredentials, wantCredentials) {
		t.Fatalf("unexpected cloudAccount in response: %+v", cluster.CloudAccount)
	}
	if !reflect.DeepEqual(cluster.CloudCredentials, wantCredentials) {
		t.Fatalf("unexpected cloudCredentials in response: %+v", cluster.CloudCredentials)
	}
	wantUpstream := &UpstreamCluster{Identity: "upstream1", Slug: "upstream", EnvironmentSlug: "env1", Version: "v1.33.1"}
	if !reflect.DeepEqual(cluster.UpstreamCluster, wantUpstream) {
		t.Fatalf("unexpected upstreamCluster in response: %+v", cluster.UpstreamCluster)
	}
	wantDownstream := []ClusterReference{{Identity: "downstream1", Slug: "downstream", EnvironmentSlug: "env2", Version: "v1.32.5"}}
	if !reflect.DeepEqual(cluster.DownstreamClusters, wantDownstream) {
		t.Fatalf("unexpected downstreamClusters in response: %+v", cluster.DownstreamClusters)
	}
}

func TestUpdateClusterRequestBody(t *testing.T) {
	environmentSlug := "env2"
	cloudCredentialsIdentity := "credentials2"
	rebootMaintenanceScheduleIdentity := "schedule2"
	tests := []struct {
		name   string
		update UpdateCluster
		want   map[string]any
	}{
		{
			name:   "unset fields are omitted from the request body",
			update: UpdateCluster{},
			want:   map[string]any{},
		},
		{
			name: "ipWhitelist, environmentSlug and cloudCredentialsIdentity",
			update: UpdateCluster{
				IPWhitelist: []IPWhitelistEntry{
					{Cidr: "192.0.2.0/24", Description: "office"},
				},
				EnvironmentSlug:          &environmentSlug,
				CloudCredentialsIdentity: &cloudCredentialsIdentity,
			},
			want: map[string]any{
				"ipWhitelist": []any{
					map[string]any{"cidr": "192.0.2.0/24", "description": "office"},
				},
				"environmentSlug":          "env2",
				"cloudCredentialsIdentity": "credentials2",
			},
		},
		{
			name: "rebootMaintenanceScheduleIdentity",
			update: UpdateCluster{
				RebootMaintenanceScheduleIdentity: &rebootMaintenanceScheduleIdentity,
			},
			want: map[string]any{
				"rebootMaintenanceScheduleIdentity": "schedule2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestBody map[string]any
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v1/orgs/org1/clusters/env1/cluster1", func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Errorf("failed to read request body: %v", err)
				}
				if err := json.Unmarshal(body, &requestBody); err != nil {
					t.Errorf("failed to unmarshal request body: %v", err)
				}
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"slug":"cluster1","environmentSlug":"env1"}`)
			})
			server := httptest.NewServer(mux)
			defer server.Close()

			c := &clientImpl{RestyClient: NewRestyClient(nil, ClientOpts{APIUrl: server.URL})}
			if _, err := c.UpdateCluster(context.Background(), "org1", "env1", "cluster1", tt.update); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(requestBody, tt.want) {
				t.Fatalf("request body = %v, want %v", requestBody, tt.want)
			}
		})
	}
}
