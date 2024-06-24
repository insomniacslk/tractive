package main

import (
	"fmt"
	"log"

	"github.com/insomniacslk/tractive"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	flagUsername = pflag.StringP("username", "u", "", "Username (e-mail)")
	flagPassword = pflag.StringP("password", "p", "", "Password")
	flagDebug    = pflag.BoolP("debug", "D", false, "Enable debug logs (might print sensitive information)")
)

func main() {
	pflag.Parse()
	if *flagUsername == "" {
		log.Fatalf("Empty username")
	}
	if *flagPassword == "" {
		log.Fatalf("Empty password")
	}
	if *flagDebug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	t, err := tractive.Authenticate(*flagUsername, *flagPassword)
	if err != nil {
		log.Fatalf("Failed to authenticate: %v", err)
	}
	fmt.Printf("%+v\n", t)
	info, err := t.GetAccountInfo()
	if err != nil {
		log.Fatalf("Failed to get account info: %v", err)
	}
	fmt.Printf("%+v\n", info)
	subscriptions, err := t.GetAccountSubscriptions()
	if err != nil {
		log.Fatalf("Failed to get account subscriptions: %v", err)
	}
	for _, s := range *subscriptions {
		sub, err := t.GetAccountSubscription(s.ID)
		if err != nil {
			logrus.Warningf("Failed to get subscription %q: %v", s.ID, err)
		}
		fmt.Printf("%+v\n", sub)
	}

	shares, err := t.GetAccountShares()
	if err != nil {
		log.Fatalf("Failed to get account shares: %v", err)
	}
	fmt.Printf("%+v\n", shares)

	pets, err := t.GetPets()
	if err != nil {
		log.Fatalf("Failed to get pets: %v", err)
	}
	for _, p := range *pets {
		pet, err := t.GetPet(p.ID)
		if err != nil {
			logrus.Warningf("Failed to get pet %q: %v", p.ID, err)
		}
		fmt.Printf("%+v\n", pet)
	}
}
