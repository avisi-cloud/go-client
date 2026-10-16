package acloudapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

func TestGetClusterNodePortAccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/orgs/org1/clusters/env1/cluster1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"slug":"cluster1","environmentSlug":"env1","nodePortAccess":%q}`, NodePortAccessAllowlist)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	c := &clientImpl{RestyClient: NewRestyClient(nil, ClientOpts{APIUrl: server.URL})}
	cluster, err := c.GetCluster(context.Background(), "org1", "env1", "cluster1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cluster.NodePortAccess != NodePortAccessAllowlist {
		t.Fatalf("unexpected nodePortAccess in response: %v", cluster.NodePortAccess)
	}
}

func TestCreateClusterNodePortAccess(t *testing.T) {
	tests := []struct {
		name           string
		nodePortAccess NodePortAccess
		wantInBody     bool
	}{
		{
			name:           "unset is omitted from the request body",
			nodePortAccess: "",
			wantInBody:     false,
		},
		{
			name:           "internal",
			nodePortAccess: NodePortAccessInternal,
			wantInBody:     true,
		},
		{
			name:           "allowlist",
			nodePortAccess: NodePortAccessAllowlist,
			wantInBody:     true,
		},
		{
			name:           "public",
			nodePortAccess: NodePortAccessPublic,
			wantInBody:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestBody map[string]any
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v1/orgs/org1/clusters/env1", func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Errorf("failed to read request body: %v", err)
				}
				if err := json.Unmarshal(body, &requestBody); err != nil {
					t.Errorf("failed to unmarshal request body: %v", err)
				}
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, `{"name":"cluster1","slug":"cluster1","nodePortAccess":%q}`, NodePortAccessPublic)
			})
			server := httptest.NewServer(mux)
			defer server.Close()

			c := &clientImpl{RestyClient: NewRestyClient(nil, ClientOpts{APIUrl: server.URL})}
			cluster, err := c.CreateCluster(context.Background(), "org1", "env1", CreateCluster{
				Name:           "cluster1",
				NodePortAccess: tt.nodePortAccess,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got, inBody := requestBody["nodePortAccess"]
			if inBody != tt.wantInBody {
				t.Fatalf("nodePortAccess in request body = %v, want %v", inBody, tt.wantInBody)
			}
			if tt.wantInBody && got != string(tt.nodePortAccess) {
				t.Fatalf("nodePortAccess in request body = %v, want %v", got, tt.nodePortAccess)
			}
			if cluster.NodePortAccess != NodePortAccessPublic {
				t.Fatalf("unexpected nodePortAccess in response: %v", cluster.NodePortAccess)
			}
		})
	}
}

func TestUpdateClusterNodePortAccess(t *testing.T) {
	internal := NodePortAccessInternal
	allowlist := NodePortAccessAllowlist
	public := NodePortAccessPublic
	tests := []struct {
		name           string
		nodePortAccess *NodePortAccess
		wantInBody     bool
	}{
		{
			name:           "nil is omitted from the request body",
			nodePortAccess: nil,
			wantInBody:     false,
		},
		{
			name:           "internal",
			nodePortAccess: &internal,
			wantInBody:     true,
		},
		{
			name:           "allowlist",
			nodePortAccess: &allowlist,
			wantInBody:     true,
		},
		{
			name:           "public",
			nodePortAccess: &public,
			wantInBody:     true,
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
				fmt.Fprintf(w, `{"name":"cluster1","slug":"cluster1","nodePortAccess":%q}`, NodePortAccessInternal)
			})
			server := httptest.NewServer(mux)
			defer server.Close()

			c := &clientImpl{RestyClient: NewRestyClient(nil, ClientOpts{APIUrl: server.URL})}
			cluster, err := c.UpdateCluster(context.Background(), "org1", "env1", "cluster1", UpdateCluster{
				NodePortAccess: tt.nodePortAccess,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got, inBody := requestBody["nodePortAccess"]
			if inBody != tt.wantInBody {
				t.Fatalf("nodePortAccess in request body = %v, want %v", inBody, tt.wantInBody)
			}
			if tt.wantInBody && got != string(*tt.nodePortAccess) {
				t.Fatalf("nodePortAccess in request body = %v, want %v", got, *tt.nodePortAccess)
			}
			if cluster.NodePortAccess != NodePortAccessInternal {
				t.Fatalf("unexpected nodePortAccess in response: %v", cluster.NodePortAccess)
			}
		})
	}
}

func TestParseNodePortAccess(t *testing.T) {
	for _, nodePortAccess := range AllNodePortAccess {
		t.Run(string(nodePortAccess), func(t *testing.T) {
			got, err := ParseNodePortAccess(string(nodePortAccess))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != nodePortAccess {
				t.Fatalf("ParseNodePortAccess() = %v, want %v", got, nodePortAccess)
			}
		})
	}

	t.Run("invalid", func(t *testing.T) {
		if _, err := ParseNodePortAccess("whitelist"); err == nil {
			t.Fatal("expected error")
		}
	})
}
