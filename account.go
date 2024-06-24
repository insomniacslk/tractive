package tractive

import (
	"encoding/json"
	"fmt"
)

type AccountInfoSubSettings struct {
	ID                            string `json:"_id"`
	Version                       string `json:"_version"`
	Type                          string `json:"_type"`
	UserRegistered                bool   `json:"user_registered"`
	GeofenceIn                    bool   `json:"geofence_in"`
	GeofenceOut                   bool   `json:"geofence_out"`
	BatteryFull                   bool   `json:"battery_full"`
	BatteryLow                    bool   `json:"battery_low"`
	BatteryCritical               bool   `json:"battery_critical"`
	BatteryEmpty                  bool   `json:"battery_empty"`
	ChargingReminder              bool   `json:"charging_reminder"`
	TrackerClipLost               bool   `json:"tracker_clip_lost"`
	TrackerTemperatureHighWarning bool   `json:"tracker_temperature_high_warning"`
	TrackerTemperatureLowWarning  bool   `json:"tracker_temperature_low_warning"`
	TrackerStartupShutdown        bool   `json:"tracker_startup_shutdown"`
	Sharing                       bool   `json:"sharing"`
	WellnessAndActivity           bool   `json:"wellness_and_activity"`
	PointOfInterestWarnings       bool   `json:"points_of_interest_warnings"`
}

type AccountInfoResponse struct {
	ID                string  `json:"_id"`
	Version           string  `json:"_version"`
	Type              string  `json:"_type"`
	Email             string  `json:"email"`
	ActivatedAt       int64   `json:"activated_at"`
	ProfilePictureID  *string `json:"profile_picture_id"`
	MembershipType    string  `json:"membership_type"`
	ReferralBonusType string  `json:"referral_bonus_type"`
	GUID              string  `json:"guid"`
	Details           struct {
		ID              string  `json:"_id"`
		Version         string  `json:"_version"`
		Type            string  `json:"_type"`
		FirstName       string  `json:"first_name"`
		LastName        string  `json:"last_name"`
		PhoneNumber     *string `json:"phone_number"`
		Gender          *string `json:"gender"`
		Birthday        *string `json:"birthday"`
		UnitDistance    string  `json:"unit_distance"`
		UnitWeight      string  `json:"unit_weight"`
		UnitTemperature string  `json:"unit_temperature"`
	} `json:"details"`
	Demographics struct {
		ID                  string `json:"_id"`
		Version             string `json:"_version"`
		Type                string `json:"_type"`
		Locale              string `json:"locale"`
		Language            string `json:"language"`
		Country             string `json:"country"`
		IsLanguageSetByUser bool   `json:"is_language_set_by_user"`
	} `json:"demographics"`
	Settings struct {
		ID                            string                 `json:"_id"`
		Version                       string                 `json:"_version"`
		Type                          string                 `json:"_type"`
		Email                         string                 `json:"email"`
		MetricSystem                  bool                   `json:"metric_system"`
		PreferredMapTypeStreet        string                 `json:"preferred_map_type_street"`
		PreferredMapTypeHybrid        string                 `json:"preferred_map_type_hybrid"`
		PosRequestAllowed             bool                   `json:"pos_request_allowed"`
		GetLivePositionFeatureEnabled bool                   `json:"get_live_position_feature_enabled"`
		NoPetSurvey                   interface{}            `json:"no_pet_survey"`
		DistanceUnit                  string                 `json:"distance_unit"`
		WeightUnit                    string                 `json:"weight_unit"`
		BadgeCelebrationsDisabled     interface{}            `json:"badge_celebrations_disabled"`
		PushSoundSettings             interface{}            `json:"push_sound_settings"`
		MailSettings                  AccountInfoSubSettings `json:"mail_settings"`
		PushSettings                  AccountInfoSubSettings `json:"push_settings"`
		WebPushSettings               AccountInfoSubSettings `json:"web_push_settings"`
	} `json:"settings"`
	InvoiceAddress struct {
		ID           string      `json:"_id"`
		Version      string      `json:"_version"`
		Type         string      `json:"_type"`
		FirstName    string      `json:"first_name"`
		LastName     string      `json:"last_name"`
		StreetName   interface{} `json:"street_name"`
		StreetNumber interface{} `json:"street_number"`
		City         interface{} `json:"city"`
		ZipCode      interface{} `json:"zip_code"`
		State        interface{} `json:"state"`
	} `json:"invoice_address"`
	Shelter                 interface{}   `json:"shelter"`
	ProfilePictures         []interface{} `json:"profile_pictures"`
	Role                    []string      `json:"role"`
	ReferralLink            string        `json:"referral_link"`
	TermsAcceptedAt         int64         `json:"terms_accepted_at"`
	PrivacyPolicyAcceptedAt int64         `json:"privacy_policy_accepted_at"`
	UserFederatedLogin      bool          `json:"user_federated_login"`
}

type AccountSubscriptionsResponse []struct {
	ID      string `json:"_id"`
	Version string `json:"_version"`
	Type    string `json:"_type"`
}

func (t *Tractive) GetAccountInfo() (*AccountInfoResponse, error) {
	u := getTractiveURL()
	u.Path = "/4/user/" + t.UserID
	body, err := tractiveRequest("GET", u, t.Token)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	var air AccountInfoResponse
	if err := json.Unmarshal(body, &air); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json response: %w", err)
	}
	return &air, nil
}

func (t *Tractive) GetAccountSubscriptions() (*AccountSubscriptionsResponse, error) {
	u := getTractiveURL()
	u.Path = "/4/user/" + t.UserID + "/subscriptions"
	body, err := tractiveRequest("GET", u, t.Token)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	var asr AccountSubscriptionsResponse
	if err := json.Unmarshal(body, &asr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json response: %w", err)
	}
	return &asr, nil
}