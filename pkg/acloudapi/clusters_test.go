package acloudapi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_toQueryParams(t *testing.T) {
	type args struct {
		mergedGetClusterOpts GetClusterOpts
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "none",
			args: args{
				mergedGetClusterOpts: GetClusterOpts{
					IncludeDetails: nil,
					ShowCompute:    nil,
				},
			},
			want: "",
		},
		{
			name: "includeDetails",
			args: args{
				mergedGetClusterOpts: GetClusterOpts{
					IncludeDetails: True(),
					ShowCompute:    nil,
				},
			},
			want: "includeDetails=true",
		},
		{
			name: "showCompute",
			args: args{
				mergedGetClusterOpts: GetClusterOpts{
					IncludeDetails: nil,
					ShowCompute:    True(),
				},
			},
			want: "show-compute=true",
		},
		{
			name: "both",
			args: args{
				mergedGetClusterOpts: GetClusterOpts{
					IncludeDetails: True(),
					ShowCompute:    True(),
				},
			},
			want: "includeDetails=true&show-compute=true",
		},
		{
			name: "both-false",
			args: args{
				mergedGetClusterOpts: GetClusterOpts{
					IncludeDetails: False(),
					ShowCompute:    False(),
				},
			},
			want: "includeDetails=false&show-compute=false",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toQueryParams(tt.args.mergedGetClusterOpts); got != tt.want {
				t.Errorf("toQueryParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptionalQueryParams(t *testing.T) {
	type args struct {
		queryParams string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{
				queryParams: "",
			},
			want: "",
		},
		{
			name: "set",
			args: args{
				queryParams: "param=1&other=2",
			},
			want: "?param=1&other=2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OptionalQueryParams(tt.args.queryParams); got != tt.want {
				t.Errorf("OptionalQueryParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetClustersErrorMessage(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/memberships", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"content":[{"slug":"org1"}],"last":true}`)
	})
	mux.HandleFunc("/api/v1/orgs/org1/clusters", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"message":"boom"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	c := &clientImpl{RestyClient: NewRestyClient(nil, ClientOpts{APIUrl: server.URL})}
	ctx := context.Background()
	_, err := c.GetClusters(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
	want := "failed to get clusters for organisation org1"
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("unexpected error: %v", err)
	}
}
