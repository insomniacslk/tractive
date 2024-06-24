package tractive

import (
	"encoding/json"
	"fmt"

	"github.com/insomniacslk/xjson"
)

type GetAllTrackersResponse []struct {
	Envelope
}

type GetTrackerResponse struct {
	Envelope
	HwID                      string         `json:"hw_id"`
	HwEdition                 string         `json:"hw_edition"`
	ModelNumber               string         `json:"model_number"`
	BluetoothMAC              interface{}    `json:"bluetooth_mac"`
	GeofenceSensitivity       string         `json:"geofence_sensitivity"`
	BatterySaveMode           interface{}    `json:"battery_save_mode"`
	ReadOnly                  bool           `json:"read_only"`
	Demo                      bool           `json:"demo"`
	SelfTestAvailable         bool           `json:"self_test_available"`
	Capabilities              []string       `json:"capabilities"`
	SupportedGeofenceTypes    []string       `json:"supported_geofence_types"`
	FwVersion                 string         `json:"fw_version"`
	State                     string         `json:"state"`
	StateReason               string         `json:"state_reason"`
	ChargingState             string         `json:"charging_state"`
	BatteryState              string         `json:"battery_state"`
	PowerSavingZoneID         string         `json:"power_saving_zone_id"`
	PrioritizedZoneID         string         `json:"prioritized_zone_id"`
	PrioritizedZoneType       string         `json:"prioritized_zone_type"`
	PrioritizedZoneLastSeenAt xjson.TimeUnix `json:"prioritized_zone_last_seen_at"`
	PrioritizedZoneEnteredAt  xjson.TimeUnix `json:"prioritized_zone_entered_at"`
}

func (t *Tractive) GetAllTrackers() (*GetAllTrackersResponse, error) {
	u := getTractiveURL()
	u.Path = "/4/user/" + t.UserID + "/trackers"
	body, err := tractiveRequest("GET", u, t.Token)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	var resp GetAllTrackersResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json response: %w", err)
	}
	return &resp, nil
}

func (t *Tractive) GetTracker(trackerID string) (*GetTrackerResponse, error) {
	u := getTractiveURL()
	u.Path = "/4/tracker/" + trackerID
	body, err := tractiveRequest("GET", u, t.Token)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	var resp GetTrackerResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json response: %w", err)
	}
	return &resp, nil
}
