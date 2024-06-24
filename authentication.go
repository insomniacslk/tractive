package tractive

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/insomniacslk/xjson"
)

type AuthenticationResponse struct {
	UserID      string         `json:"user_id"`
	ClientID    string         `json:"client_id"`
	ExpiresAt   xjson.TimeUnix `json:"expires_at"`
	AccessToken string         `json:"access_token"`
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
	var ar AuthenticationResponse
	if err := json.Unmarshal(body, &ar); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json response: %w", err)
	}
	return &Tractive{
		Username:       username,
		Password:       password,
		UserID:         ar.UserID,
		ClientID:       ar.ClientID,
		Token:          ar.AccessToken,
		TokenExpiresAt: time.Unix(time.Time(ar.ExpiresAt).Unix(), 0),
	}, nil
}
