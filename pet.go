package tractive

import (
	"encoding/json"
	"fmt"

	"github.com/insomniacslk/xjson"
)

type PetsResponse []struct {
	Envelope
}

type PetResponse struct {
	Envelope
	LeaderboardOptOut bool           `json:"leaderboard_opt_out"`
	DeviceID          string         `json:"device_id"`
	ReadOnly          bool           `json:"read_only"`
	CreatedAt         xjson.TimeUnix `json:"created_at"`
	Details           struct {
		Envelope
		Name                string         `json:"name"`
		PetType             string         `json:"pet_type"`
		BreedIDs            []string       `json:"breed_ids"`
		Gender              string         `json:"gender"`
		Birthday            xjson.TimeUnix `json:"birthday"`
		ProfilePictureFrame interface{}    `json:"profile_picture_frame"`
		Height              float64        `json:"height"`
		Length              interface{}    `json:"length"`
		Weight              int            `json:"weight"`
		ChipID              string         `json:"chip_id"`
		Neutered            bool           `json:"neutered"`
		Personality         []interface{}  `json:"personality"`
		LostOrDead          interface{}    `json:"lost_or_dead"`
		Lim                 interface{}    `json:"lim"`
		RibCage             interface{}    `json:"ribcage"`
		WeightIsDefault     interface{}    `json:"weight_is_default"`
		HeightIsDefault     interface{}    `json:"height_is_default"`
		BirthdayIsDefault   interface{}    `json:"birthday_is_default"`
		BreedIsDefault      interface{}    `json:"breed_is_default"`
		InstagramUsername   string         `json:"instagram_user_name"`
		ProfilePictureID    string         `json:"profile_picture_id"`
		CoverPictureID      string         `json:"cover_picture_id"`
		CharacteristicIDs   []interface{}  `json:"characteristic_ids"`
		GalleryPictureIDs   []interface{}  `json:"gallery_picture_ids"`
		ReadOnly            bool           `json:"read_only"`
		ActivitySettings    struct {
			Envelope
			DailyGoal                          int         `json:"daily_goal"`
			DailyDistanceGoal                  int         `json:"daily_distance_goal"`
			DailyActiveMinutesGoal             int         `json:"daily_active_minutes_goal"`
			ActivityCategoryThresholdsOverride interface{} `json:"activity_category_thresholds_override"`
		} `json:"activity_settings"`
	} `json:"details"`
}

func (t *Tractive) GetPets() (*PetsResponse, error) {
	u := getTractiveURL()
	u.Path = "/4/user/" + t.UserID + "/trackable_objects"
	body, err := tractiveRequest("GET", u, t.Token)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	var resp PetsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json response: %w", err)
	}
	return &resp, nil

}

func (t *Tractive) GetPet(petID string) (*PetResponse, error) {
	u := getTractiveURL()
	u.Path = "/4/trackable_object/" + petID
	body, err := tractiveRequest("GET", u, t.Token)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	var resp PetResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json response: %w", err)
	}
	return &resp, nil

}
