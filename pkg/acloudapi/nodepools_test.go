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

func TestCreateNodePoolSecurityUpdatesOnJoin(t *testing.T) {
	tests := []struct {
		name                  string
		securityUpdatesOnJoin NodePoolSecurityUpdatesOnJoin
		wantInBody            bool
	}{
		{
			name:                  "unset is omitted from the request body",
			securityUpdatesOnJoin: "",
			wantInBody:            false,
		},
		{
			name:                  "off",
			securityUpdatesOnJoin: NodePoolSecurityUpdatesOnJoinOff,
			wantInBody:            true,
		},
		{
			name:                  "install",
			securityUpdatesOnJoin: NodePoolSecurityUpdatesOnJoinInstall,
			wantInBody:            true,
		},
		{
			name:                  "install-and-reboot",
			securityUpdatesOnJoin: NodePoolSecurityUpdatesOnJoinInstallAndReboot,
			wantInBody:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestBody map[string]any
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v1/orgs/org1/clusters/env1/cluster1/pools", func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Errorf("failed to read request body: %v", err)
				}
				if err := json.Unmarshal(body, &requestBody); err != nil {
					t.Errorf("failed to unmarshal request body: %v", err)
				}
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, `{"id":1,"name":"pool1","securityUpdatesOnJoin":%q}`, NodePoolSecurityUpdatesOnJoinInstall)
			})
			server := httptest.NewServer(mux)
			defer server.Close()

			c := &clientImpl{RestyClient: NewRestyClient(nil, ClientOpts{APIUrl: server.URL})}
			cluster := Cluster{CustomerSlug: "org1", EnvironmentSlug: "env1", Slug: "cluster1"}
			nodePool, err := c.CreateNodePool(context.Background(), cluster, CreateNodePool{
				Name:                  "pool1",
				SecurityUpdatesOnJoin: tt.securityUpdatesOnJoin,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got, inBody := requestBody["securityUpdatesOnJoin"]
			if inBody != tt.wantInBody {
				t.Fatalf("securityUpdatesOnJoin in request body = %v, want %v", inBody, tt.wantInBody)
			}
			if tt.wantInBody && got != string(tt.securityUpdatesOnJoin) {
				t.Fatalf("securityUpdatesOnJoin in request body = %v, want %v", got, tt.securityUpdatesOnJoin)
			}
			if nodePool.SecurityUpdatesOnJoin != NodePoolSecurityUpdatesOnJoinInstall {
				t.Fatalf("unexpected securityUpdatesOnJoin in response: %v", nodePool.SecurityUpdatesOnJoin)
			}
		})
	}
}

func TestUpdateNodePoolSecurityUpdatesOnJoinOmittedWhenUnset(t *testing.T) {
	var rawBody []byte
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/orgs/org1/clusters/env1/cluster1/pools/1", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read request body: %v", err)
		}
		rawBody = body
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":1,"name":"pool1"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	c := &clientImpl{RestyClient: NewRestyClient(nil, ClientOpts{APIUrl: server.URL})}
	cluster := Cluster{CustomerSlug: "org1", EnvironmentSlug: "env1", Slug: "cluster1"}
	nodePool, err := c.UpdateNodePool(context.Background(), cluster, 1, CreateNodePool{Name: "pool1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(string(rawBody), "securityUpdatesOnJoin") {
		t.Fatalf("expected securityUpdatesOnJoin to be omitted from request body, got: %s", rawBody)
	}
	if nodePool.SecurityUpdatesOnJoin != "" {
		t.Fatalf("unexpected securityUpdatesOnJoin in response: %v", nodePool.SecurityUpdatesOnJoin)
	}
}

func TestParseNodePoolSecurityUpdatesOnJoin(t *testing.T) {
	for _, securityUpdatesOnJoin := range AllNodePoolSecurityUpdatesOnJoin {
		t.Run(string(securityUpdatesOnJoin), func(t *testing.T) {
			got, err := ParseNodePoolSecurityUpdatesOnJoin(string(securityUpdatesOnJoin))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != securityUpdatesOnJoin {
				t.Fatalf("ParseNodePoolSecurityUpdatesOnJoin() = %v, want %v", got, securityUpdatesOnJoin)
			}
		})
	}

	t.Run("invalid", func(t *testing.T) {
		if _, err := ParseNodePoolSecurityUpdatesOnJoin("invalid"); err == nil {
			t.Fatal("expected error")
		}
	})
}
