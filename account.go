package tractive

import (
	"encoding/json"
	"fmt"

	"github.com/insomniacslk/xjson"
)

type Envelope struct {
	ID      string `json:"_id"`
	Version string `json:"_version"`
	Type    string `json:"_type"`
}

type AccountInfoSubSettings struct {
	Envelope
	UserRegistered                bool `json:"user_registered"`
	GeofenceIn                    bool `json:"geofence_in"`
	GeofenceOut                   bool `json:"geofence_out"`
	BatteryFull                   bool `json:"battery_full"`
	BatteryLow                    bool `json:"battery_low"`
	BatteryCritical               bool `json:"battery_critical"`
	BatteryEmpty                  bool `json:"battery_empty"`
	ChargingReminder              bool `json:"charging_reminder"`
	TrackerClipLost               bool `json:"tracker_clip_lost"`
	TrackerTemperatureHighWarning bool `json:"tracker_temperature_high_warning"`
	TrackerTemperatureLowWarning  bool `json:"tracker_temperature_low_warning"`
	TrackerStartupShutdown        bool `json:"tracker_startup_shutdown"`
	Sharing                       bool `json:"sharing"`
	WellnessAndActivity           bool `json:"wellness_and_activity"`
	PointOfInterestWarnings       bool `json:"points_of_interest_warnings"`
}

type AccountInfoResponse struct {
	Envelope
	Email             string         `json:"email"`
	ActivatedAt       xjson.TimeUnix `json:"activated_at"`
	ProfilePictureID  *string        `json:"profile_picture_id"`
	MembershipType    string         `json:"membership_type"`
	ReferralBonusType string         `json:"referral_bonus_type"`
	GUID              string         `json:"guid"`
	Details           struct {
		Envelope
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
		Envelope
		Locale              string `json:"locale"`
		Language            string `json:"language"`
		Country             string `json:"country"`
		IsLanguageSetByUser bool   `json:"is_language_set_by_user"`
	} `json:"demographics"`
	Settings struct {
		Envelope
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
		Envelope
		FirstName    string      `json:"first_name"`
		LastName     string      `json:"last_name"`
		StreetName   interface{} `json:"street_name"`
		StreetNumber interface{} `json:"street_number"`
		City         interface{} `json:"city"`
		ZipCode      interface{} `json:"zip_code"`
		State        interface{} `json:"state"`
	} `json:"invoice_address"`
	Shelter                 interface{}    `json:"shelter"`
	ProfilePictures         []interface{}  `json:"profile_pictures"`
	Role                    []string       `json:"role"`
	ReferralLink            string         `json:"referral_link"`
	TermsAcceptedAt         xjson.TimeUnix `json:"terms_accepted_at"`
	PrivacyPolicyAcceptedAt xjson.TimeUnix `json:"privacy_policy_accepted_at"`
	UserFederatedLogin      bool           `json:"user_federated_login"`
}

type AccountSubscriptionsResponse []struct {
	Envelope
}

type AccountSubscriptionResponse struct {
	Envelope
	Status                   string         `json:"status"`
	ValidFrom                xjson.TimeUnix `json:"valid_from"`
	ValidTo                  xjson.TimeUnix `json:"valid_to"`
	TrialEndDate             interface{}    `json:"trial_end_date"`
	Recurring                bool           `json:"recurring"`
	PlanTypeUsed             string         `json:"plan_type_used"`
	ExtraMonths              []interface{}  `json:"extra_months"`
	RecurringPausedUntil     interface{}    `json:"recurring_paused_until"`
	ServicePausedUntil       interface{}    `json:"service_paused_until"`
	TrackerShipped           bool           `json:"trackers_shipped"`
	DunningSince             interface{}    `json:"dunning_since"`
	PaymentPlanID            string         `json:"payment_plan_id"`
	AdditionalServiceIDs     []string       `json:"additional_service_ids"`
	SuspendedAt              interface{}    `json:"suspended_at"`
	HomeCountry              string         `json:"home_country"`
	StateCode                interface{}    `json:"state_code"`
	ZipCode                  interface{}    `json:"zip_code"`
	ReadOnly                 bool           `json:"read_only"`
	TrackerID                string         `json:"tracker_id"`
	BillingInterval          string         `json:"billing_interval"`
	InsuranceActive          bool           `json:"insurance_active"`
	InsuranceClaimsAvailable bool           `json:"insurance_claims_available"`
	ServiceCanBePaused       bool           `json:"service_can_be_paused"`
	Transferable             bool           `json:"transferable"`
	Refundable               bool           `json:"refundable"`
	IsShelterTransferrable   interface{}    `json:"is_shelter_transferable"`
	Care                     struct {
		Status          string `json:"status"`
		AvailableClaims int    `json:"available_claims"`
	} `json:"care"`
}

type AccountSharesResponse []struct {
	Envelope
}

func (t *Tractive) GetAccountInfo() (*AccountInfoResponse, error) {
	u := getTractiveURL()
	u.Path = "/4/user/" + t.UserID
	body, err := tractiveRequest("GET", u, t.Token)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	var resp AccountInfoResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json response: %w", err)
	}
	return &resp, nil
}

func (t *Tractive) GetAccountSubscriptions() (*AccountSubscriptionsResponse, error) {
	u := getTractiveURL()
	u.Path = "/4/user/" + t.UserID + "/subscriptions"
	body, err := tractiveRequest("GET", u, t.Token)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	var resp AccountSubscriptionsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json response: %w", err)
	}
	return &resp, nil
}

func (t *Tractive) GetAccountSubscription(subscriptionID string) (*AccountSubscriptionResponse, error) {
	u := getTractiveURL()
	u.Path = "/4/subscription/" + subscriptionID
	body, err := tractiveRequest("GET", u, t.Token)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	var resp AccountSubscriptionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json response: %w", err)
	}
	return &resp, nil
}

func (t *Tractive) GetAccountShares() (*AccountSharesResponse, error) {
	u := getTractiveURL()
	u.Path = "/4/user/" + t.UserID + "/shares"
	body, err := tractiveRequest("GET", u, t.Token)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	var resp AccountSharesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json response: %w", err)
	}
	return &resp, nil
}
