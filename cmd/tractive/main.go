package main

import (
	"fmt"
	"log"

	"github.com/insomniacslk/tractive"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	flagToken    = pflag.StringP("token", "t", "", "Token. If empty, username and password must be specified. Requires --user-id")
	flagUserID   = pflag.StringP("user-id", "i", "", "User ID. If empty, username and password must be set. Requires --user-id")
	flagUsername = pflag.StringP("username", "u", "", "Username (e-mail)")
	flagPassword = pflag.StringP("password", "p", "", "Password")
	flagDebug    = pflag.BoolP("debug", "D", false, "Enable debug logs (might print sensitive information)")
)

func main() {
	pflag.Parse()
	if *flagDebug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	var (
		t   *tractive.Tractive
		err error
	)
	if *flagToken == "" {
		if *flagUsername == "" {
			log.Fatalf("Empty username and no token specified")
		}
		if *flagPassword == "" {
			log.Fatalf("Empty password and no token specified")
		}
		t, err = tractive.Authenticate(*flagUsername, *flagPassword)
		if err != nil {
			log.Fatalf("Failed to authenticate: %v", err)
		}
	} else {
		if *flagUserID == "" {
			log.Fatalf("Empty user ID")
		}
		t = &tractive.Tractive{
			Token:    *flagToken,
			ClientID: tractive.ClientID,
			UserID:   *flagUserID,
		}
	}
	fmt.Printf("%+v\n", t)

	// Account Info
	info, err := t.GetAccountInfo()
	if err != nil {
		log.Fatalf("Failed to get account info: %v", err)
	}
	fmt.Printf("Account info: %+v\n", info)

	// Subscriptions
	subscriptions, err := t.GetAccountSubscriptions()
	if err != nil {
		log.Fatalf("Failed to get account subscriptions: %v", err)
	}
	for _, s := range *subscriptions {
		sub, err := t.GetAccountSubscription(s.ID)
		if err != nil {
			logrus.Warningf("Failed to get subscription %q: %v", s.ID, err)
		}
		fmt.Printf("Subscription: %+v\n", sub)
	}

	// Account Shares
	shares, err := t.GetAccountShares()
	if err != nil {
		log.Fatalf("Failed to get account shares: %v", err)
	}
	fmt.Printf("Share: %+v\n", shares)

	// Pets
	pets, err := t.GetPets()
	if err != nil {
		log.Fatalf("Failed to get pets: %v", err)
	}
	for _, p := range *pets {
		pet, err := t.GetPet(p.ID)
		if err != nil {
			logrus.Warningf("Failed to get pet %q: %v", p.ID, err)
		}
		fmt.Printf("Pet: %+v\n", pet)
	}
	// Trackers
	trackers, err := t.GetAllTrackers()
	if err != nil {
		log.Fatalf("Failed to get trackers: %v", err)
	}
	fmt.Printf("Trackers: %+v\n", trackers)
	for _, tr := range *trackers {
		tracker, err := t.GetTracker(tr.ID)
		if err != nil {
			logrus.Warningf("Failed to get tracker %q: %v", tr.ID, err)
		}
		fmt.Printf("Tracker: %+v\n", tracker)
	}
}
