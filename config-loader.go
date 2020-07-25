package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

var (
	serverConfig map[string]string
)

func lgServerConfigListLoad(configURL string, svLoc string) {
	configLogger := log.New(os.Stdout, "Config Loader: ", log.LstdFlags)
	serverConfig = make(map[string]string)
	serverConfig["whois"] = "disabled"
	serverConfig["nslookup"] = "disabled"
	serverConfig["ping"] = "disabled"
	serverConfig["icmp"] = "disabled"
	serverConfig["tracert"] = "disabled"
	serverConfig["webcontrol"] = "disabled"
	serverConfig["tcp"] = "disabled"
	serverConfig["IPv4"] = "disabled"
	serverConfig["IPv6"] = "disabled"
	serverConfig["curl"] = "disabled"
	serverConfig["mtr"] = "enabled"
	serverConfig["referrers"] = ""

	// Control the Config URL is defined
	if configURL == "" {
		configLogger.Println("Config url is not defined. please define with --config-url")
	}

	// Control the svLoc is defined
	if svLoc == "" {
		configLogger.Println("server Location is not defined. Please define with --svloc")
	}

	// If svLoc or configURL is empty, Enable every settings.
	if svLoc == "" || configURL == "" {
		configLogger.Println("\033[1;31mAll settings Enabled.\033[0m")
	} else {
		// If configuration url and server loc is given, Start to load conf.

		// Get config from remote.
		resp, err := http.Get(configURL)
		if err != nil {
			// If anything going bad, print err and exit.
			configLogger.Fatalln(err)
		}

		defer resp.Body.Close()

		// Check the configuration url has a file.
		if resp.StatusCode != http.StatusOK {
			configLogger.Fatalln("Config server response is ok. Response : " + resp.Status)
		}
		// Read response
		RemoteConfigJSON, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			configLogger.Fatalln(err)
		}

		// parse current server as a array to start nested loop.
		currentServerName := strings.Split(svLoc, ".")

		// Start minus 1 because firstly checks "ServerConfig" item.
		for i := -1; i < len(currentServerName); i++ {
			var currentLoc string
			for t := 0; t <= i; t++ {
				currentLoc = currentLoc + ".Servers." + currentServerName[t]
			}
			if i == len(currentServerName)-1 {
				currentLocWithRemovedFirstDot := strings.Replace(currentLoc, ".", "", 1)
				serverURLexpexted := gjson.Get(string(RemoteConfigJSON), currentLocWithRemovedFirstDot+".Url").String()
				serverListUnExpexted := gjson.Get(string(RemoteConfigJSON), currentLocWithRemovedFirstDot+".Servers").String()
				if serverURLexpexted != "" && serverListUnExpexted == "" {
					serverConfig["ThisServerURL"] = serverURLexpexted
				} else {
					configLogger.Fatalln("Given server location is not a server.")
				}

			}
			currentLoc = currentLoc + ".ServerConfig"
			currentLoc = strings.Replace(currentLoc, ".", "", 1)
			//configLogger.Println("\nCurrent loc " + currentLoc)
			if currentLoc != "" {
				for k := range serverConfig {
					if gjson.Get(string(RemoteConfigJSON), currentLoc+"."+k).String() != "" {
						serverConfig[k] = gjson.Get(string(RemoteConfigJSON), currentLoc+"."+k).String()
						//fmt.Printf("%s %s\n", k, serverConfig[k])
					}
				}
			}

		}
		// Get frontend server address from config url

		if serverConfig["referrers"] == "" {
			parsedConfigURL, err := url.Parse(configURL)
			if err != nil {
				panic(err)
			}
			configLogger.Println("Referer Domain is not given, you can set with --referrers. System is only allows incoming requests from ", parsedConfigURL.Host)
			allowedreferrers = []string{parsedConfigURL.Host} //parsedConfigURL.Host
		} else {
			allowedreferrers = strings.Split(serverConfig["referrers"], ",")
		}

		// Print everthing is successfully
		configLogger.Println("Setting loaded for " + svLoc)
		for k, v := range serverConfig {
			configLogger.Println(k+":", v)
		}
	}

}
