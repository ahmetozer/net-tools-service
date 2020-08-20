package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	isFunctionEnabled map[string]bool
	///hostname string
	allowedReferrers []string
	//
	allowedDomain string = ""
)

func lgServerConfigListLoad() {
	configLogger := log.New(os.Stdout, "Config Loader: ", log.LstdFlags)

	availableFunctions := []string{"whois", "nslookup", "ping", "icmp", "tracert", "webcontrol", "tcp", "curl", "mtr"}
	availableIPVersion := []string{"IPv4", "IPv6", "DS", "disabled"}

	isFunctionEnabled = make(map[string]bool)
	// This is not required because it is false by default.
	for _, s := range availableFunctions {
		isFunctionEnabled[s] = false
	}

	// Look which functions is enabled
	enabledFunctionsENV, ok := os.LookupEnv("functions")
	if ok {
		enabledFunctions := strings.Split(enabledFunctionsENV, ",")
		for _, s := range enabledFunctions {
			if contains(availableFunctions, s) {
				isFunctionEnabled[s] = true
			} else {
				configLogger.Println(s + " is unknown function")
			}
		}
	} else {
		configLogger.Println("\033[1;31mEnvironment variable \"functions\" is not set. All functions is Enabled.\033[0m")
		for _, s := range availableFunctions {
			isFunctionEnabled[s] = true
		}
	}

	enabledIPVersionENV, ok := os.LookupEnv("ipver")
	if ok {
		enabledIPVersions := strings.Split(enabledIPVersionENV, ",")
		for _, s := range enabledIPVersions {
			if !contains(availableIPVersion, s) {
				configLogger.Println("IP version " + s + " is not known.")
			}
		}
		isFunctionEnabled["IPv4"] = false
		isFunctionEnabled["IPv6"] = false
		if !contains(enabledIPVersions, "disabled") {
			if contains(enabledIPVersions, "DS") {
				isFunctionEnabled["IPv4"] = true
				isFunctionEnabled["IPv6"] = true
			} else {
				if contains(enabledIPVersions, "IPv4") {
					isFunctionEnabled["IPv4"] = true
				}
				if contains(enabledIPVersions, "IPv6") {
					isFunctionEnabled["IPv6"] = true
				}
			}
		}

	} else {
		configLogger.Println("IP version is not defined, DualStack mode enabled.")
		isFunctionEnabled["IPv4"] = true
		isFunctionEnabled["IPv6"] = true
	}

	referrersENV, ok := os.LookupEnv("referrers")
	if ok {
		allowedReferrers := strings.Split(referrersENV, ",")
		configLogger.Print("Allowed sites to make request to this server : ")
		for _, s := range allowedReferrers {
			fmt.Print(s)
		}
		fmt.Println()
	} else {
		configLogger.Println("\033[1;31mEnvironment variable \"referrers\" is not set. All websites can make request to this server.\033[0m")
	}

	serverHostname, err := os.Hostname()
	if err != nil {
		configLogger.Fatalln(err)
	}

	hostnameCheck, ok := os.LookupEnv("hostname")
	if ok {
		if hostnameCheck == "" {
			allowedDomain = serverHostname
			configLogger.Println("hostname:", serverHostname)
		} else {
			configLogger.Println("hostname:", hostnameCheck)
			allowedDomain = hostnameCheck
		}
	} else {
		configLogger.Println("\033[1;31mRequested domain is not controlled. Any one can bind this ip with any domain.\033[0m")
	}

}
