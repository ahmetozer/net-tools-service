package main

import (
	"io/ioutil"
	"log"
	"net/http"
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
	serverConfig["whois"] = "enabled"
	serverConfig["nslookup"] = "enabled"
	serverConfig["ping"] = "enabled"
	serverConfig["icmp"] = "enabled"
	serverConfig["tracert"] = "enabled"
	serverConfig["webcontrol"] = "enabled"
	serverConfig["tcp"] = "enabled"
	serverConfig["IPv4"] = "enabled"
	serverConfig["IPv6"] = "enabled"
	serverConfig["speedtest"] = "enabled"

	if configURL == "" {
		configLogger.Println("Config url is not defined. please define with --config-url")
	}
	if svLoc == "" {
		configLogger.Println("server Location is not defined. Please define with --svloc")
	}

	if svLoc == "" || configURL == "" {
		configLogger.Println("All settings Enabled.")
	} else {
		resp, err := http.Get(configURL)
		if err != nil {
			configLogger.Fatalln(err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			configLogger.Fatalln("Config server response is ok. Response : " + resp.Status)
		}
		RemoteConfigJSON, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			configLogger.Fatalln(err)
		}

		currentServerName := strings.Split(svLoc, ".")

		for i := -1; i < len(currentServerName); i++ {
			var temp string
			for t := 0; t <= i; t++ {
				temp = temp + ".Servers." + currentServerName[t]
			}
			if i == len(currentServerName)-1 {
				temp2 := strings.Replace(temp, ".", "", 1)
				serverURLexpexted := gjson.Get(string(RemoteConfigJSON), temp2+".Url").String()
				serverListUnExpexted := gjson.Get(string(RemoteConfigJSON), temp2+".Servers").String()
				if serverURLexpexted != "" && serverListUnExpexted == "" {
					serverConfig["ThisServerURL"] = serverURLexpexted
				} else {
					configLogger.Fatalln("Given server location is not a server.")
				}

			}
			temp = temp + ".ServerConfig"
			temp = strings.Replace(temp, ".", "", 1)
			//configLogger.Println("\nCurrent loc " + temp)
			if temp != "" {
				for k := range serverConfig {
					if gjson.Get(string(RemoteConfigJSON), temp+"."+k).String() != "" {
						serverConfig[k] = gjson.Get(string(RemoteConfigJSON), temp+"."+k).String()
						//fmt.Printf("%s %s\n", k, serverConfig[k])
					}
				}
			}

		}
		configLogger.Println("Setting loaded for " + svLoc)
		for k, v := range serverConfig {
			configLogger.Println(k, v)
		}
	}
}
