package acloudapi

import (
	"context"
	"fmt"
)

func (c *clientImpl) GetMaintenanceSchedules(ctx context.Context, org string) ([]MaintenanceSchedule, error) {
	var result []MaintenanceSchedule
	response, err := c.R().
		SetContext(ctx).
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/organisations/%s/maintenance-schedule", org))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *clientImpl) GetMaintenanceSchedule(ctx context.Context, org, maintenanceScheduleID string) (*MaintenanceSchedule, error) {
	maintenanceSchedule := MaintenanceSchedule{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&maintenanceSchedule).
		Get(fmt.Sprintf("/api/v1/organisations/%s/maintenance-schedule/%s", org, maintenanceScheduleID))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return &maintenanceSchedule, nil
}

type CreateMaintenanceSchedule struct {
	Name string `json:"name"`
	// TimeZone is an IANA time zone name in which the windows are interpreted; empty: UTC.
	TimeZone                  string              `json:"timeZone,omitempty"`
	Windows                   []MaintenanceWindow `json:"windows"`
	DefaultForClusterUpgrades bool                `json:"defaultForClusterUpgrades,omitempty"`
	DefaultForReboots         bool                `json:"defaultForReboots,omitempty"`
}

func (c *clientImpl) CreateMaintenanceSchedule(ctx context.Context, org string, createMaintenanceSchedule CreateMaintenanceSchedule) (*MaintenanceSchedule, error) {
	maintenanceSchedule := MaintenanceSchedule{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&maintenanceSchedule).
		SetBody(&createMaintenanceSchedule).
		Post(fmt.Sprintf("/api/v1/organisations/%s/maintenance-schedule", org))

	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return &maintenanceSchedule, nil
}

func (c *clientImpl) DeleteMaintenanceSchedule(ctx context.Context, org, maintenanceScheduleID string) error {
	response, err := c.R().
		SetContext(ctx).
		Delete(fmt.Sprintf("/api/v1/organisations/%s/maintenance-schedule/%s", org, maintenanceScheduleID))
	if err := c.CheckResponse(response, err); err != nil {
		return err
	}
	return nil
}

// UpdateMaintenanceSchedule updates a schedule; empty or nil fields keep their current value.
type UpdateMaintenanceSchedule struct {
	Name     string              `json:"name,omitempty"`
	TimeZone *string             `json:"timeZone,omitempty"`
	Windows  []MaintenanceWindow `json:"windows,omitempty"`
	// DefaultForClusterUpgrades can only be set to true; the current default can't be unset.
	DefaultForClusterUpgrades *bool `json:"defaultForClusterUpgrades,omitempty"`
	// DefaultForReboots can only be set to true; the current default can't be unset.
	DefaultForReboots *bool `json:"defaultForReboots,omitempty"`
}

func (c *clientImpl) UpdateMaintenanceSchedule(ctx context.Context, org, maintenanceScheduleID string, updateMaintenanceSchedule UpdateMaintenanceSchedule) (*MaintenanceSchedule, error) {
	maintenanceSchedule := MaintenanceSchedule{}
	response, err := c.R().
		SetContext(ctx).
		SetResult(&maintenanceSchedule).
		SetBody(&updateMaintenanceSchedule).
		Patch(fmt.Sprintf("/api/v1/organisations/%s/maintenance-schedule/%s", org, maintenanceScheduleID))
	if err := c.CheckResponse(response, err); err != nil {
		return nil, err
	}
	return &maintenanceSchedule, nil
}
