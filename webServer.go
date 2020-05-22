package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"time"
	//	"os"
	//"math/rand"
	//"encoding/base64"
	"strings"
	//"bytes"
	//"encoding/json"
)

// pass CMD output to HTTP
func writeCmdOutput(res http.ResponseWriter, pipeReader *io.PipeReader) {
	var BUFLEN = 4096 // for
	buffer := make([]byte, BUFLEN)
	for {
		n, err := pipeReader.Read(buffer)
		if err != nil {
			pipeReader.Close()
			break
		}

		data := buffer[0:n]
		res.Write(data)
		if f, ok := res.(http.Flusher); ok {
			f.Flush()
			fmt.Printf(fmt.Sprint(ok))

		}
		//reset buffer
		for i := 0; i < n; i++ {
			buffer[i] = 0
		}
	}
}

var (
	/*
		Regexs for checking input
	*/
	ipv6Regex   = `^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`
	ipv4Regex   = `^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`
	domainRegex = `^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$`
)

func isMayIPv4(host string) bool {
	match, _ := regexp.MatchString(ipv4Regex+`|`+domainRegex, host)
	return match
}

func isMayIPv6(host string) bool {
	match, _ := regexp.MatchString(ipv6Regex+`|`+domainRegex, host)
	return match
}

func webServer(logger *log.Logger) *http.Server {

	// Crearte New HTTP Router
	router := http.NewServeMux()

	// Index Handler
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `<title="Looking Glass Server"/><h2>Error page `+fmt.Sprint(r.URL)+
			` not found.</h2></br><p>Looking Glass Server. For more information, visit <a href="https://ahmetozer.org/">ahmetozer.org</a></p>`)
	})

	// Get System Disk informations with linux util lsblk . // JSON
	router.HandleFunc("/system-status.php", func(w http.ResponseWriter, r *http.Request) {
		/*
			Get arguments from url
		*/
		functype := r.URL.Query().Get("type")
		host := r.URL.Query().Get("host")
		nameserver := r.URL.Query().Get("nameserver")
		port := r.URL.Query().Get("port")
		id := r.URL.Query().Get("id")
		mobile := r.URL.Query().Get("mobile")
		term := r.URL.Query().Get("term")
		setLiveOutputHeaders(w)

		/*
			Check host input
		*/
		match, _ := regexp.MatchString(ipv4Regex+`|`+ipv6Regex+`|`+domainRegex, host)
		if !match {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "ERR: Host is not IPv4, IPv6 or domain")
			return
		}

		/*
			Server functions
		*/
		switch functype {
		case "icmp4", "icmp6", "icmp":

			//cmd.Dir = "/empty/"
			args := []string{"-l 3", "-c 5", "-i 0.3", "-s 64", "-t 64", "-W 1", "-q"}
			if functype == "icmp4" {
				if !isMayIPv4(host) {
					fmt.Fprintf(w, "{\"code\":\"FuncTypeMissMatch\"}")
					return
				}
				args = append(args, "-4")
			}
			if functype == "icmp6" {
				if !isMayIPv6(host) {
					fmt.Fprintf(w, "{\"code\":\"FuncTypeMissMatch\"}")
					return
				}
				args = append(args, "-6")
			}
			args = append(args, host)
			cmd := exec.Command("ping", args...)
			//err := cmd.Run()
			out, err := cmd.CombinedOutput()

			if err != nil { // If error occur on ping command.
				// If given input type wich is IPv4 or IPv6 and run type is not match this error will be occur
				if fmt.Sprint(err) == "exit status 2" {
					fmt.Fprintf(w, "{\"code\":\"FuncTypeMissMatchExecuted\"}")
					return
				}
				//	When the ping command cannot access the server, this error will be occur
				if fmt.Sprint(err) == "exit status 1" {
					fmt.Fprintf(w, "{\"code\":\"RemoteHostDown\"}")
					return
				}
				// If Any un expected occur, this will be shown
				fmt.Fprintf(w, "{\"code\":\"UnknownExit\",\"exitCode\":\""+fmt.Sprint(err)+"\",\"execOut:\""+string(out)+"\"}")

			} else {

				//  Execute output to convert string
				outString := string(out)

				// Get only rtt status
				mdevLoc := strings.Index(outString, "/mdev =")
				rttOut := outString[mdevLoc+8 : mdevLoc+strings.Index(outString[mdevLoc+1:], "ms")]
				// parse rtt status
				rttOutParsed := strings.Split(rttOut, "/") // [0] rtt min , [1] avg, [2] max, [3] mdev
				// Get other data from program output.
				transmittedPacketCount := outString[strings.Index(outString, "ping statistics ---")+20 : strings.Index(outString, " packets transmitted,")]
				recivedPacketCount := outString[strings.Index(outString, " packets transmitted,")+22 : strings.Index(outString, " received,")]
				packetLoss := outString[strings.Index(outString, " received,")+11 : strings.Index(outString, " packet loss,")]

				//fmt.Fprint(w)
				//fmt.Fprint(w, rttOutParsed[0]+"\n"+transmittedPacketCount+"\n"+recivedPacketCount+"\n"+packetLoss)

				fmt.Fprint(w, `{"code":"OK", "rttmin":"`+rttOutParsed[0]+`", "rttavg":"`+rttOutParsed[1]+`", "rttmax":"`+rttOutParsed[2]+`", "mdev":"`+
					rttOutParsed[3]+`", "packetloss":"`+packetLoss+`", "recivedPacketCount": "`+recivedPacketCount+`", "transmittedPacketCount":"`+transmittedPacketCount+`","functype":"`+functype+`"}`)

				// To debug output
				//fmt.Fprint(w, "\n=====================================================\n\n\n"+outString)
			}
		case "tcp":
		case "webkontrol":
			resp, err := http.Get(host)
			if err != nil {
				log.Fatal(err)
			}

			// Print the HTTP Status Code and Status Name
			fmt.Fprintf(w, `{ "code":"`+http.StatusText(resp.StatusCode)+`" }`)

		case "time":
			date := time.Now()
			fmt.Fprintf(w, `{ "time":"`+date.Format("15:04:05")+`","date":"`+date.Format("01/02/2006")+`" }`)
		case "whois":
			args := []string{host}
			if term == "" {
				w.Header().Set("content-type", "text/html; charset=utf-8")
				fmt.Fprintf(w, "<style>body {background-color: #2d3436}</style><pre style='white-space: pre-line;  font-size: 16px; font-family: Arial, Helvetica, sans-serif;  color: #dfe6e9'>")
			}
			cmd := exec.Command("whois", args...)
			// Organize pipelines
			out, err := cmd.CombinedOutput()
			if err != nil {
				// If error occur on ping command.
				if fmt.Sprint(err) == "exit status 1" {
					fmt.Fprintf(w, string(out))
					return
				}
				fmt.Fprintf(w, "{\"code\":\"UnknownExit\",\"exitCode\":\""+fmt.Sprint(err)+"\",\"execOut:\""+string(out)+"\"}")

			} else {
				//  Execute output to convert string
				fmt.Fprintf(w, string(out))
			}
			return
		case "nslookup":
			args := []string{host}
			if nameserver != "" {
				match, _ := regexp.MatchString(ipv4Regex+`|`+ipv6Regex+`|`+domainRegex, nameserver)
				if !match {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "ERR: Host is not IPv4, IPv6 or domain")
					return
				}
				args = append(args, nameserver)
			}
			if term == "" {
				w.Header().Set("content-type", "text/html; charset=utf-8")
				fmt.Fprintf(w, "<style>body {background-color: #2d3436}</style><pre style='white-space: pre-line;  font-size: 16px; font-family: Arial, Helvetica, sans-serif;  color: #dfe6e9'>")
			}
			cmd := exec.Command("nslookup", args...)
			out, err := cmd.CombinedOutput()
			if err != nil {
				// If error occur on command.
				if fmt.Sprint(err) == "exit status 1" {
					fmt.Fprintf(w, string(out))
					return
				}
				fmt.Fprintf(w, "{\"code\":\"UnknownExit\",\"exitCode\":\""+fmt.Sprint(err)+"\",\"execOut:\""+string(out)+"\"}")

			} else {
				//  Execute output to convert string
				fmt.Fprintf(w, string(out))
			}
			/*************************
			Live output for ping frame
			**************************/
		case "ping", "ping4", "ping6":

			args := []string{"-c 10", "-i 0.2"}

			if mobile == "1" { // No resolve domain names to reduce widht of ping output to shown in mobile in better
				args = append(args, "-n")
			}
			if functype == "ping4" {
				if !isMayIPv4(host) {
					fmt.Fprintf(w, "{\"code\":\"FuncTypeMissMatch\"}")
					return
				}
				args = append(args, "-4")
			}
			if functype == "ping6" {
				if !isMayIPv6(host) {
					fmt.Fprintf(w, "{\"code\":\"FuncTypeMissMatch\"}")
					return
				}
				args = append(args, "-6")
			}

			// If requests comes from iframe (not term), add style to iframe
			if term == "" {
				w.Header().Set("content-type", "text/html; charset=utf-8")
				fmt.Fprintf(w, "<style>body {background-color: #2d3436}</style><pre style='white-space: pre-line;  font-size: 16px; font-family: Arial, Helvetica, sans-serif;  color: #dfe6e9'>")
			}

			args = append(args, host) // add host to arguments

			cmd := exec.Command("ping", args...)
			// Organize pipelines
			pipeIn, pipeWriter := io.Pipe()
			cmd.Stdout = pipeWriter
			cmd.Stderr = pipeWriter
			// Pass to web output
			go writeCmdOutput(w, pipeIn)

			// Run commands
			cmd.Run()
			pipeWriter.Close()
			/****************************
			Live output for tracert frame
			*****************************/
		case "tracert", "tracert4", "tracert6":
			args := []string{}

			if mobile == "1" {
				args = append(args, "-n -q 1")
			} else {
				args = append(args, "-q 2")
			}

			if functype == "tracert4" {
				if !isMayIPv4(host) {
					fmt.Fprintf(w, "{\"code\":\"FuncTypeMissMatch\"}")
					return
				}
				args = append(args, "-4")
			} else if functype == "tracert6" {
				if !isMayIPv6(host) {
					fmt.Fprintf(w, "{\"code\":\"FuncTypeMissMatch\"}")
					return
				}
				args = append(args, "-6")
			}

			if term == "" {
				w.Header().Set("content-type", "text/html; charset=utf-8")
				fmt.Fprintf(w, "<style>body {background-color: #2d3436}</style><pre style='white-space: pre-line;  font-size: 16px; font-family: Arial, Helvetica, sans-serif;  color: #dfe6e9'>")
			}
			args = append(args, host) // add host to arguments

			cmd := exec.Command("traceroute", args...)
			// Organize pipelines
			pipeIn, pipeWriter := io.Pipe()
			cmd.Stdout = pipeWriter
			cmd.Stderr = pipeWriter
			// Pass to web output
			go writeCmdOutput(w, pipeIn) // live output
			// Run command
			cmd.Run()
			pipeWriter.Close()
			return
		case "curlkontrol":
		case "curltamkontrol":
		case "curldurum":
		default:
			// if any unknown function name given.
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "404 \n"+functype+host+nameserver+port+id+mobile+term)
			return
		}

	})

	// return as a webServer
	return &http.Server{

		Addr:     *listenAddr,
		Handler:  router,
		ErrorLog: logger,
		/* Close sockets */
		//ReadTimeout:  5 * time.Second,
		//WriteTimeout: 10 * time.Second,
		//IdleTimeout:  15 * time.Second,
	}
}

/*
	Set HTTP headers to show live output on browser
*/
func setLiveOutputHeaders(w http.ResponseWriter) {
	w.Header().Set("content-type", "application/x-javascript")
	w.Header().Set("expires", "10s")
	w.Header().Set("Pragma", "public")
	w.Header().Set("Cache-Control", "public, maxage=10, proxy-revalidate")
	w.Header().Set("X-Accel-Buffering", "no")
}
