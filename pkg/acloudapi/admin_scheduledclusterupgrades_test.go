package acloudapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUpdateScheduledClusterUpgradeProgress(t *testing.T) {
	var body map[string]any
	mux := http.NewServeMux()
	mux.HandleFunc("/admin/v1/scheduled-cluster-upgrades/upgrade1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method = %s, want PATCH", r.Method)
		}
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"identity": "upgrade1",
			"status": "STALLED",
			"startedAt": "2026-11-05T08:00:00Z",
			"lastProgressAt": "2026-11-05T09:30:00Z",
			"progress": "control plane 1.31.4; 18/30 nodes on target version (18 ready)"
		}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	lastProgressAt := time.Date(2026, 11, 5, 9, 30, 0, 0, time.UTC)
	c := &adminClientImpl{RestyClient: NewRestyClient(nil, ClientOpts{APIUrl: server.URL})}
	upgrade, err := c.UpdateScheduledClusterUpgrade(context.Background(), UpdateScheduledClusterUpgradeRequest{
		Identity:       "upgrade1",
		LastProgressAt: &lastProgressAt,
		Progress:       "control plane 1.31.4; 18/30 nodes on target version (18 ready)",
	})
	if err != nil {
		t.Fatalf("UpdateScheduledClusterUpgrade() error = %v", err)
	}

	if _, ok := body["status"]; ok {
		t.Errorf("request body contains status, want it omitted: %v", body)
	}
	if got := body["lastProgressAt"]; got != "2026-11-05T09:30:00Z" {
		t.Errorf("request lastProgressAt = %v", got)
	}
	if got := body["progress"]; got != "control plane 1.31.4; 18/30 nodes on target version (18 ready)" {
		t.Errorf("request progress = %v", got)
	}

	if upgrade.Status != Stalled {
		t.Errorf("Status = %q, want %q", upgrade.Status, Stalled)
	}
	if upgrade.StartedAt == nil || !upgrade.StartedAt.Equal(time.Date(2026, 11, 5, 8, 0, 0, 0, time.UTC)) {
		t.Errorf("StartedAt = %v", upgrade.StartedAt)
	}
	if upgrade.LastProgressAt == nil || !upgrade.LastProgressAt.Equal(lastProgressAt) {
		t.Errorf("LastProgressAt = %v", upgrade.LastProgressAt)
	}
	if upgrade.Progress != "control plane 1.31.4; 18/30 nodes on target version (18 ready)" {
		t.Errorf("Progress = %q", upgrade.Progress)
	}
	if upgrade.SucceededClusterVersion != "" {
		t.Errorf("SucceededClusterVersion = %q, want it empty for a STALLED upgrade", upgrade.SucceededClusterVersion)
	}
}

func TestScheduledClusterUpgradeDecodesSucceededClusterVersion(t *testing.T) {
	var upgrade ScheduledClusterUpgrade
	err := json.Unmarshal([]byte(`{
		"identity": "upgrade1",
		"status": "SUCCEEDED",
		"toClusterVersion": "v1.31.4-ame.0",
		"succeededClusterVersion": "v1.31.5-ame.0"
	}`), &upgrade)
	if err != nil {
		t.Fatal(err)
	}
	if upgrade.ToClusterVersion != "v1.31.4-ame.0" {
		t.Errorf("ToClusterVersion = %q", upgrade.ToClusterVersion)
	}
	if upgrade.SucceededClusterVersion != "v1.31.5-ame.0" {
		t.Errorf("SucceededClusterVersion = %q, want %q", upgrade.SucceededClusterVersion, "v1.31.5-ame.0")
	}
}

func TestUpdateScheduledClusterUpgradeRequestOmitsEmptyProgress(t *testing.T) {
	raw, err := json.Marshal(UpdateScheduledClusterUpgradeRequest{Identity: "upgrade1", Status: Succeeded})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(raw), `{"identity":"upgrade1","status":"SUCCEEDED"}`; got != want {
		t.Errorf("json = %s, want %s", got, want)
	}
}
