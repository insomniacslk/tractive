package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/insomniacslk/tractive"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	flagTractiveToken     = pflag.StringP("tractive-token", "t", "", "Token. If empty, username and password must be specified. Requires --user-id")
	flagTractiveUserID    = pflag.StringP("tractive-user-id", "i", "", "Tractive user ID. If empty, username and password must be set. Requires --token")
	flagTractiveUsername  = pflag.StringP("tractive-username", "u", "", "Tractive username (e-mail)")
	flagTractivePassword  = pflag.StringP("tractive-password", "p", "", "Tractive password")
	flagOwntracksEndpoint = pflag.StringP("owntracks-endpoint", "E", "http://localhost:8083/pub", "OwnTracks endpoint URL to publish datapoints to")
	flagOwntracksUsername = pflag.StringP("owntracks-username", "U", "", "OwnTracks username")
	flagOwntracksPassword = pflag.StringP("owntracks-password", "P", "", "OwnTracks password")
	flagOwntracksDevice   = pflag.StringP("owntracks-device", "D", "", "OwnTracks device name")
	flagOwntracksTID      = pflag.StringP("owntracks-tid", "T", "", "OwnTracks tracker ID (two letters)")
	flagStartTime         = pflag.IntP("start-time", "s", -1, "Start time as UNIX timestamp (if not specified, default to now-1h)")
	flagEndTime           = pflag.IntP("end-time", "e", -1, "End time as UNIX timestamp (if not specified, default to now)")
	flagDebug             = pflag.BoolP("debug", "d", false, "Enable debug logs (might print sensitive information)")
)

type OwnTracksDatapoint struct {
	Type      string  `json:"_type"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Timestamp int64   `json:"tst"`
	Accuracy  int     `json:"acc"`
	Altitude  int     `json:"alt"`
	TID       string  `json:"tid"`
}

func main() {
	pflag.Parse()
	if *flagDebug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	var (
		t   *tractive.Tractive
		err error
	)
	if *flagTractiveToken == "" {
		if *flagTractiveUsername == "" {
			logrus.Fatalf("Empty username and no token specified")
		}
		if *flagTractivePassword == "" {
			logrus.Fatalf("Empty password and no token specified")
		}
		t, err = tractive.Authenticate(*flagTractiveUsername, *flagTractivePassword)
		if err != nil {
			logrus.Fatalf("Failed to authenticate: %v", err)
		}
	} else {
		if *flagTractiveUserID == "" {
			logrus.Fatalf("Empty user ID")
		}
		t = &tractive.Tractive{
			Token:    *flagTractiveToken,
			ClientID: tractive.ClientID,
			UserID:   *flagTractiveUserID,
		}
	}
	if *flagOwntracksEndpoint == "" {
		logrus.Fatalf("owntracks-endpoint is not set")
	}
	if *flagOwntracksTID == "" {
		logrus.Fatalf("owntracks-tid is not set")
	}

	// Pets
	trackersToPets := make(map[string]*tractive.PetResponse, 0)
	pets, err := t.GetPets()
	if err != nil {
		logrus.Fatalf("Failed to get pets: %v", err)
	}
	logrus.Infof("Found %d pets", len(*pets))
	for _, p := range *pets {
		pet, err := t.GetPet(p.ID)
		if err != nil {
			logrus.Warningf("Failed to get pet %q: %v", p.ID, err)
		}
		logrus.Infof("  Pet: %+v\n", pet)
		trackersToPets[pet.DeviceID] = pet
	}

	// Trackers
	trackers, err := t.GetAllTrackers()
	if err != nil {
		logrus.Fatalf("Failed to get trackers: %v", err)
	}
	logrus.Infof("Found %d trackers", len(*trackers))
	start := time.Now().Add(-time.Hour)
	end := time.Now()
	if *flagStartTime != -1 {
		start = time.Unix(int64(*flagStartTime), 0)
	}
	if *flagEndTime != -1 {
		end = time.Unix(int64(*flagEndTime), 0)
	}
	logrus.Infof("Querying time range: %s   -->   %s\n", start, end)
	owntracksEndpoint, err := url.Parse(*flagOwntracksEndpoint)
	if err != nil {
		logrus.Fatalf("Failed to parse OwnTracks endpoint URL: %v", err)
	}
	for _, tr := range *trackers {
		tracker, err := t.GetTracker(tr.ID)
		if err != nil {
			logrus.Warningf("Failed to get tracker %q: %v", tr.ID, err)
			continue
		}
		pet, ok := trackersToPets[tr.ID]
		if !ok {
			logrus.Fatalf("No pet found for tracker ID %q", tr.ID)
		}
		q := owntracksEndpoint.Query()
		q.Add("u", pet.Details.Name)
		q.Add("d", "tractive")
		owntracksEndpoint.RawQuery = q.Encode()
		logrus.Infof("Tracker: %+v\n", tracker)
		positions, err := t.GetTrackerPositions(tr.ID, start, end)
		if err != nil {
			logrus.Warningf("Failed to get tracker %q 's positions: %v", tr.ID, err)
			continue
		}
		logrus.Debugf("Tracker positions:\n")
		client := http.Client{}
		for _, pos := range (*positions)[0] {
			dp := OwnTracksDatapoint{
				Type:      "location",
				Latitude:  pos.LatLong[0],
				Longitude: pos.LatLong[1],
				Timestamp: pos.Time,
				Accuracy:  pos.PosUncertainty,
				Altitude:  pos.Alt,
				TID:       *flagOwntracksTID,
			}
			logrus.Debugf("  pos=%s\n", pos.String())
			logrus.Debugf("  dp=%+v\n", dp)
			postData, err := json.Marshal(dp)
			if err != nil {
				logrus.Fatalf("Failed to marshal owntracks JSON payload: %v", err)
			}
			logrus.Debugf("Sending request to %s", owntracksEndpoint)
			logrus.Debugf("POST payload: %s", postData)
			req, err := http.NewRequest(http.MethodPost, owntracksEndpoint.String(), bytes.NewBuffer(postData))
			req.Header.Set("Content-Type", "application/json")
			// TODO set owntracks username and password
			resp, err := client.Do(req)
			if err != nil {
				logrus.Fatalf("Failed to make POST request to OwnTracks: %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				logrus.Fatalf("POST request to OwnTracks failed with %s", resp.Status)
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logrus.Fatalf("Failed to read response body: %v", err)
			}
			logrus.Debugf("Response: %s\n", string(body))
		}
		logrus.Infof("Pushed %d positions for tracker %s", len((*positions)[0]), tr.ID)
	}
}
