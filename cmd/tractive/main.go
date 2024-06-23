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
	for _, sub := range *subscriptions {
		fmt.Printf("%+v\n", sub)
	}
}
