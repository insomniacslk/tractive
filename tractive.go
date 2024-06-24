package tractive

import (
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
