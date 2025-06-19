package acloudapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRestyClient_CheckResponse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"message":"ok"}`)
		}))
		defer server.Close()

		c := NewRestyClient(nil, ClientOpts{APIUrl: server.URL})
		resp, err := c.R().Get("/")
		if err != nil {
			t.Fatalf("unexpected request error: %v", err)
		}
		if err := c.CheckResponse(resp, err); err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	})

	t.Run("nil response without trace", func(t *testing.T) {
		c := NewRestyClient(nil, ClientOpts{})
		err := c.CheckResponse(nil, nil)
		if err == nil || err.Error() != "no response received" {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("nil response with trace", func(t *testing.T) {
		c := NewRestyClient(nil, ClientOpts{Trace: true})
		err := c.CheckResponse(nil, nil)
		if err == nil || err.Error() != "no response received" {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
