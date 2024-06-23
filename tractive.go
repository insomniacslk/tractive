package tractive

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	TractiveScheme   = "https"
	TractiveHost     = "graph.tractive.com"
	TractiveClientID = "6536c228870a3c8857d452e8"
)

type Tractive struct {
	Username       string
	Password       string
	Token          string
	TokenExpiresAt time.Time
	UserID         string
	ClientID       string
}

type TokenRequestResponse struct {
	UserID      string `json:"user_id"`
	ClientID    string `json:"client_id"`
	ExpiresAt   int64  `json:"expires_at"`
	AccessToken string `json:"access_token"`
}

type SubSettings struct {
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
		ID                            string      `json:"_id"`
		Version                       string      `json:"_version"`
		Type                          string      `json:"_type"`
		Email                         string      `json:"email"`
		MetricSystem                  bool        `json:"metric_system"`
		PreferredMapTypeStreet        string      `json:"preferred_map_type_street"`
		PreferredMapTypeHybrid        string      `json:"preferred_map_type_hybrid"`
		PosRequestAllowed             bool        `json:"pos_request_allowed"`
		GetLivePositionFeatureEnabled bool        `json:"get_live_position_feature_enabled"`
		NoPetSurvey                   interface{} `json:"no_pet_survey"`
		DistanceUnit                  string      `json:"distance_unit"`
		WeightUnit                    string      `json:"weight_unit"`
		BadgeCelebrationsDisabled     interface{} `json:"badge_celebrations_disabled"`
		PushSoundSettings             interface{} `json:"push_sound_settings"`
		MailSettings                  SubSettings `json:"mail_settings"`
		PushSettings                  SubSettings `json:"push_settings"`
		WebPushSettings               SubSettings `json:"web_push_settings"`
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

func getTractiveURL() url.URL {
	return url.URL{
		Scheme: TractiveScheme,
		Host:   TractiveHost,
	}
}

func tractiveRequest(method string, u url.URL, token string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("X-Tractive-Client", TractiveClientID)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// only if debug requested
	if logrus.GetLevel() == logrus.DebugLevel {
		reqDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			logrus.Warningf("Failed to dump HTTP request: %v", err)
		} else {
			logrus.Debugf("HTTP REQUEST:\n%s\n", string(reqDump))
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute http request: %w", err)
	}

	// only if debug requested
	if logrus.GetLevel() == logrus.DebugLevel {
		respDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			logrus.Warningf("Failed to dump HTTP request: %v", err)
		} else {
			fmt.Printf("HTTP RESPONSE:\n%s\n", string(respDump))
		}
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http status is %s, expected 200 OK", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to get http body: %w", err)
	}
	return body, nil
}

func Authenticate(username, password string) (*Tractive, error) {
	u := getTractiveURL()
	u.Path = "/4/auth/token"
	v := url.Values{}
	v.Set("grant_type", "tractive")
	v.Set("platform_email", username)
	v.Set("platform_token", password)
	u.RawQuery = v.Encode()
	body, err := tractiveRequest("POST", u, "")
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	var trr TokenRequestResponse
	if err := json.Unmarshal(body, &trr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json response: %w", err)
	}
	return &Tractive{
		Username:       username,
		Password:       password,
		UserID:         trr.UserID,
		ClientID:       trr.ClientID,
		Token:          trr.AccessToken,
		TokenExpiresAt: time.Unix(trr.ExpiresAt, 0),
	}, nil
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
