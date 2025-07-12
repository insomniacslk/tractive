package tractive

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

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

type GetTrackerPositionsResponse [][]TrackerPosition

type TrackerPosition struct {
	Time           int64      `json:"time"`
	LatLong        [2]float64 `json:"latlong"`
	Alt            int        `json:"alt"`
	Speed          float64    `json:"speed"`
	Course         int        `json:"course"`
	PosUncertainty int        `json:"pos_uncertainty"`
	SensorUsed     string     `json:"sensor_used"`
}

func (p *TrackerPosition) String() string {
	return fmt.Sprintf("[%s] latitude=%.3f longitude=%.3f altitude=%d speed=%.3f course=%d pos_uncertainty=%d sensor_used=%s", time.Unix(p.Time, 0), p.LatLong[0], p.LatLong[1], p.Alt, p.Speed, p.Course, p.PosUncertainty, p.SensorUsed)
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

func (t *Tractive) GetTrackerPositions(trackerID string, start, end time.Time) (*GetTrackerPositionsResponse, error) {
	u := getTractiveURL()
	// FIXME: using API version 3 because I couldn't find the equivalent method for API v4
	u.Path = "/3/tracker/" + trackerID + "/positions"
	q := u.Query()
	q.Add("time_from", strconv.FormatInt(start.Unix(), 10))
	q.Add("time_to", strconv.FormatInt(end.Unix(), 10))
	q.Add("format", "json_segments")
	u.RawQuery = q.Encode()
	body, err := tractiveRequest("GET", u, t.Token)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	var resp GetTrackerPositionsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json response: %w", err)
	}
	return &resp, nil
}
