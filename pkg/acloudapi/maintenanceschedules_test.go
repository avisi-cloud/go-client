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
)

const maintenanceScheduleJSON = `{
	"identity": "schedule1", "name": "Weekdays night", "clusterCount": 3,
	"timeZone": "Europe/Amsterdam",
	"windows": [
		{"day": "MONDAY", "startTime": "01:00", "duration": 240},
		{"day": "FRIDAY", "startTime": "22:00", "duration": 180}
	],
	"defaultForClusterUpgrades": true, "defaultForReboots": false
}`

func TestGetClusterMaintenanceSchedules(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/orgs/org1/clusters/env1/cluster1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{
			"slug": "cluster1",
			"maintenanceSchedule": %s,
			"rebootMaintenanceSchedule": {
				"identity": "schedule2", "name": "Sunday", "clusterCount": 1, "timeZone": "UTC",
				"windows": [{"day": "SUNDAY", "startTime": "02:00", "duration": 60}],
				"defaultForClusterUpgrades": false, "defaultForReboots": true
			}
		}`, maintenanceScheduleJSON)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	c := &clientImpl{RestyClient: NewRestyClient(nil, ClientOpts{APIUrl: server.URL})}
	cluster, err := c.GetCluster(context.Background(), "org1", "env1", "cluster1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantSchedule := &MaintenanceSchedule{
		Identity:     "schedule1",
		Name:         "Weekdays night",
		ClusterCount: 3,
		TimeZone:     "Europe/Amsterdam",
		MaintenanceWindows: []MaintenanceWindow{
			{Day: "MONDAY", StartTime: "01:00", Duration: 240},
			{Day: "FRIDAY", StartTime: "22:00", Duration: 180},
		},
		DefaultForClusterUpgrades: true,
	}
	if !reflect.DeepEqual(cluster.MaintenanceSchedule, wantSchedule) {
		t.Fatalf("unexpected maintenanceSchedule in response: %+v", cluster.MaintenanceSchedule)
	}
	wantRebootSchedule := &MaintenanceSchedule{
		Identity:           "schedule2",
		Name:               "Sunday",
		ClusterCount:       1,
		TimeZone:           "UTC",
		MaintenanceWindows: []MaintenanceWindow{{Day: "SUNDAY", StartTime: "02:00", Duration: 60}},
		DefaultForReboots:  true,
	}
	if !reflect.DeepEqual(cluster.RebootMaintenanceSchedule, wantRebootSchedule) {
		t.Fatalf("unexpected rebootMaintenanceSchedule in response: %+v", cluster.RebootMaintenanceSchedule)
	}
}

func TestCreateMaintenanceScheduleRequestBody(t *testing.T) {
	tests := []struct {
		name   string
		create CreateMaintenanceSchedule
		want   map[string]any
	}{
		{
			name: "time zone and default flags are omitted when unset",
			create: CreateMaintenanceSchedule{
				Name:    "Weekdays night",
				Windows: []MaintenanceWindow{{Day: "MONDAY", StartTime: "01:00", Duration: 240}},
			},
			want: map[string]any{
				"name":    "Weekdays night",
				"windows": []any{map[string]any{"day": "MONDAY", "startTime": "01:00", "duration": float64(240)}},
			},
		},
		{
			name: "all fields",
			create: CreateMaintenanceSchedule{
				Name:     "Weekdays night",
				TimeZone: "Europe/Amsterdam",
				Windows: []MaintenanceWindow{
					{Day: "MONDAY", StartTime: "01:00", Duration: 240},
					{Day: "FRIDAY", StartTime: "23:00", Duration: 180},
				},
				DefaultForClusterUpgrades: true,
				DefaultForReboots:         true,
			},
			want: map[string]any{
				"name":     "Weekdays night",
				"timeZone": "Europe/Amsterdam",
				"windows": []any{
					map[string]any{"day": "MONDAY", "startTime": "01:00", "duration": float64(240)},
					map[string]any{"day": "FRIDAY", "startTime": "23:00", "duration": float64(180)},
				},
				"defaultForClusterUpgrades": true,
				"defaultForReboots":         true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestBody map[string]any
			server := maintenanceScheduleServer(t, http.MethodPost, "/api/v1/organisations/org1/maintenance-schedule", &requestBody)
			defer server.Close()

			c := &clientImpl{RestyClient: NewRestyClient(nil, ClientOpts{APIUrl: server.URL})}
			schedule, err := c.CreateMaintenanceSchedule(context.Background(), "org1", tt.create)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(requestBody, tt.want) {
				t.Fatalf("request body = %v, want %v", requestBody, tt.want)
			}
			if schedule.TimeZone != "Europe/Amsterdam" {
				t.Fatalf("unexpected time zone in response: %q", schedule.TimeZone)
			}
		})
	}
}

func TestUpdateMaintenanceScheduleRequestBody(t *testing.T) {
	timeZone := "Europe/Amsterdam"
	defaultForReboots := true
	tests := []struct {
		name   string
		update UpdateMaintenanceSchedule
		want   map[string]any
	}{
		{
			name:   "nothing is sent when nothing is set",
			update: UpdateMaintenanceSchedule{},
			want:   map[string]any{},
		},
		{
			name:   "switching a default sends only the flag",
			update: UpdateMaintenanceSchedule{DefaultForReboots: &defaultForReboots},
			want:   map[string]any{"defaultForReboots": true},
		},
		{
			name: "all fields",
			update: UpdateMaintenanceSchedule{
				Name:              "Sunday",
				TimeZone:          &timeZone,
				Windows:           []MaintenanceWindow{{Day: "SUNDAY", StartTime: "02:30", Duration: 60}},
				DefaultForReboots: &defaultForReboots,
			},
			want: map[string]any{
				"name":              "Sunday",
				"timeZone":          "Europe/Amsterdam",
				"windows":           []any{map[string]any{"day": "SUNDAY", "startTime": "02:30", "duration": float64(60)}},
				"defaultForReboots": true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestBody map[string]any
			server := maintenanceScheduleServer(t, http.MethodPatch, "/api/v1/organisations/org1/maintenance-schedule/schedule1", &requestBody)
			defer server.Close()

			c := &clientImpl{RestyClient: NewRestyClient(nil, ClientOpts{APIUrl: server.URL})}
			if _, err := c.UpdateMaintenanceSchedule(context.Background(), "org1", "schedule1", tt.update); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(requestBody, tt.want) {
				t.Fatalf("request body = %v, want %v", requestBody, tt.want)
			}
		})
	}
}

// maintenanceScheduleServer serves maintenanceScheduleJSON on method and path, and decodes the request body into requestBody.
func maintenanceScheduleServer(t *testing.T, method, path string, requestBody *map[string]any) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Errorf("method = %s, want %s", r.Method, method)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read request body: %v", err)
		}
		if err := json.Unmarshal(body, requestBody); err != nil {
			t.Errorf("failed to unmarshal request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, maintenanceScheduleJSON)
	})
	return httptest.NewServer(mux)
}
