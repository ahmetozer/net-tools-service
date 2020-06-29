package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"
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
	domainRegex = `^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z]$`
	portRegex   = "^((6553[0-5])|(655[0-2][0-9])|(65[0-4][0-9]{2})|(6[0-4][0-9]{3})|([1-5][0-9]{4})|([1-9][0-9]{3})|([1-9][0-9]{2})|([1-9][0-9])|([1-9]))$"

	//iframeStyle = "<pre style='white-space: pre-line; text-shadow: 3px 3px 4px #000; font-size: 20px; font-family: Arial, Helvetica, sans-serif;  color: #000'>"
	iframeStyle = "<pre style='white-space: pre-line; font-size: 20px; font-family: Arial, Helvetica, sans-serif;  color: #000'>"
)

func isMayIPv4(host string) bool {
	match, err := regexp.MatchString(ipv4Regex+`|`+domainRegex, host)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return match
}

func isMayOnlyIPv4(host string) bool {
	match, err := regexp.MatchString(ipv4Regex, host)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return match
}

func isMayIPv6(host string) bool {
	match, err := regexp.MatchString(ipv6Regex+`|`+domainRegex, host)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return match
}
func isMayOnlyIPv6(host string) bool {
	match, err := regexp.MatchString(ipv6Regex, host)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return match
}

func isHTTPURLScheme(scheme string) bool {
	match, err := regexp.MatchString(`^http$|^https$`, scheme)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return match
}

func isPortValid(port string) bool {
	match, err := regexp.MatchString(portRegex, port)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return match
}

func isMayDomain(host string) bool {
	match, err := regexp.MatchString(domainRegex, host)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return match
}

var limiter = newIPRateLimiter(1, 1)

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
	router.HandleFunc("/svCheck", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
	})
	// Get System Disk informations with linux util lsblk . // JSON
	router.HandleFunc("/looking-glass-controller", func(w http.ResponseWriter, r *http.Request) {
		// Server conf checker
		var httpScheme string
		if r.TLS == nil {
			httpScheme = "http"
		} else {
			httpScheme = "https"
		}
		if httpScheme+"://"+r.Host != serverConfig["ThisServerURL"] {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprintf(w, `{"code":"NotAcceptable","err":Request Domain is different", "requestDomain": "%s"}`, httpScheme+"://"+r.Host)
			return
		}
		if serverConfig["IPv4"] != "enabled" && serverConfig["IPv6"] != "enabled" {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprintf(w, `{"code":"NotAcceptable","err":This server does not have a IPv4 and IPv6 connection, so this server is disabled."`)
			return
		}

		// All functions to be check connecting IP version except time,whois,nslookup
		switch r.URL.Query().Get("funcType") {
		case // IPversion control not required services
			"time",
			"whois",
			"nslookup":
		default:
			if serverConfig["IPv4"] == "enabled" && serverConfig["IPv6"] != "enabled" {
				if r.URL.Query().Get("IPVersion") != "IPv4" {
					w.WriteHeader(http.StatusNotAcceptable)
					fmt.Fprintf(w, `{"code":"NotAcceptable","err":"This server only allow IPv4 requests"}`)
					return
				}
			}
			if serverConfig["IPv6"] == "enabled" && serverConfig["IPv4"] != "enabled" {
				if r.URL.Query().Get("IPVersion") != "IPv6" {
					w.WriteHeader(http.StatusNotAcceptable)
					fmt.Fprintf(w, `{"code":"NotAcceptable","err":"This server only allow IPv6 requests"}`)
					return
				}
			}

		}
		if serverConfig[r.URL.Query().Get("funcType")] != "enabled" {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, `{"code":"Forbidden","err":"This function is disabled"`)
			return

		}

		// All functions require host variable

		host := r.URL.Query().Get("host")

		setLiveOutputHeaders(w)

		/*
			Check host input
		*/
		if host == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `"{"code":"BadRequest","err":"You have to define host."`)
			return
		}
		match, _ := regexp.MatchString(ipv4Regex+`|`+ipv6Regex+`|`+domainRegex, host)
		if !match {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `"{"code":"BadRequest","err":"Host is not IPv4, IPv6 or domain"`)
			return
		}

		/*
			Server functions
		*/
		switch r.URL.Query().Get("funcType") {
		case "icmp":

			//cmd.Dir = "/empty/"
			args := []string{"-l 3", "-c 5", "-i 0.3", "-s 64", "-t 64", "-W 1", "-q"}
			switch r.URL.Query().Get("IPVersion") {
			case "IPv4":
				if !isMayIPv4(host) {
					fmt.Fprintf(w, "{\"code\":\"IPVersionMissMatch\"}")
					return
				}
				args = append(args, "-4")
			case "IPv6":
				if !isMayIPv6(host) {
					fmt.Fprintf(w, "{\"code\":\"IPVersionMissMatch\"}")
					return
				}
				args = append(args, "-6")
			case "Default":
			default:
				fmt.Fprintf(w, "{\"code\":\"WrongIPVersion\"}")
				return
			}
			args = append(args, host)
			cmd := exec.Command("ping", args...)
			//err := cmd.Run()
			out, err := cmd.CombinedOutput()

			if err != nil { // If error occur on ping command.
				// If given input type wich is IPv4 or IPv6 and run type is not match this error will be occur
				if fmt.Sprint(err) == "exit status 2" {
					fmt.Fprintf(w, "{\"code\":\"funcTypeMissMatchExecuted\"}")
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
					rttOutParsed[3]+`", "packetloss":"`+packetLoss+`", "recivedPacketCount": "`+recivedPacketCount+`", "transmittedPacketCount":"`+transmittedPacketCount+`"}`)

				// To debug output
				//fmt.Fprint(w, "\n=====================================================\n\n\n"+outString)
			}
			return
		case "tcp":
			port := r.URL.Query().Get("port")
			if port == "" {
				port = "443"
			}

			//Check if it's a domain
			if isMayDomain(host) {
				// resolve domain, Pre resolvin important for net.Dialer. If its not pre resolved, Resolving time will be add to latency.
				ips, err := net.LookupIP(host)
				if err != nil {
					fmt.Fprintf(w, `{ "code":"DomainResolveErr", "err":"%s" }`, err)
					return
				}

				switch r.URL.Query().Get("IPVersion") {
				case "IPv4":
					for _, ip := range ips {
						host = ip.String()
						if isMayOnlyIPv4(host) {
							break
						}
					}
					if !isMayOnlyIPv4(host) {
						fmt.Fprintf(w, `{ "code":"DomainResolveErr", "err":"DomainDoesNotHaveAIPv4" }`)
						return
					}
				case "IPv6":
					for _, ip := range ips {
						host = ip.String()
						if isMayOnlyIPv6(host) {
							break
						}
					}
					if !isMayOnlyIPv6(host) {
						fmt.Fprintf(w, `{ "code":"DomainResolveErr", "err":"DomainDoesNotHaveAIPv6" }`)
						return
					}
				}
			}

			if isMayIPv6(host) { // Add brackets if IPv6
				host = "[" + host + "]"
			}

			if isPortValid(port) {
				host = host + ":" + port
			} else {
				fmt.Fprintf(w, `{ "code":"InvalidPort" }`)
				return
			}

			d := net.Dialer{Timeout: 5 * time.Second}
			dialStartTime := time.Now()
			conn, err := d.Dial("tcp", host)
			if err != nil {
				fmt.Fprintf(w, `{ "code"="Down","err":"`+fmt.Sprint(err)+`" }`)
				return
			}
			elapsedTime := time.Since(dialStartTime)
			fmt.Fprintf(w, `{ "code"="ok","latency":"%d ms" }`, elapsedTime.Milliseconds())
			defer conn.Close()

			return
		case "webcontrol":
			scheme := r.URL.Query().Get("scheme")
			if r.URL.Query().Get("scheme") == "" { // If scheme is not given, set to https
				scheme = "https"
			}
			port := r.URL.Query().Get("port")
			if port != "" {
				if isPortValid(port) {
					host = host + ":" + port
				} else {
					fmt.Fprintf(w, `{ "code":"InvalidPort" }`)
					return
				}
			}
			if isHTTPURLScheme(scheme) {
				resp, err := http.Get(scheme + "://" + host)
				if err != nil {
					fmt.Fprintf(w, `{ "code"="Down","err":"`+fmt.Sprint(err)+`" }`)
				} else { // Print the HTTP Status Code and Status Name
					fmt.Fprintf(w, `{ "code":"`+http.StatusText(resp.StatusCode)+`" }`)
				}
			} else {
				fmt.Fprintf(w, `{ "code":"SchemeDoesNotMatchHTTPorHTTPS" }`)
			}
			return
		case "time":
			date := time.Now()
			fmt.Fprintf(w, `{ "time":"`+date.Format("15:04:05")+`","date":"`+date.Format("01/02/2006")+`" }`)
			return
		case "whois":
			args := []string{host}
			if r.URL.Query().Get("term") == "" {
				w.Header().Set("content-type", "text/html; charset=utf-8")
				fmt.Fprintf(w, iframeStyle)
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
			if r.URL.Query().Get("nameserver") != "" {
				nameserver := r.URL.Query().Get("nameserver")
				match, _ := regexp.MatchString(ipv4Regex+`|`+ipv6Regex+`|`+domainRegex, nameserver)
				if !match {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "ERR: Host is not IPv4, IPv6 or domain")
					return
				}
				args = append(args, nameserver)
			}
			if r.URL.Query().Get("term") == "" {
				w.Header().Set("content-type", "text/html; charset=utf-8")
				fmt.Fprintf(w, iframeStyle)
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
				//  Execute output to convert string and send to web.
				fmt.Fprintf(w, string(out))
			}
			return
			/*************************
			Live output for ping frame
			**************************/
		case "ping":

			args := []string{"-c 10", "-i 0.2"}

			if r.URL.Query().Get("mobile") == "1" { // No resolve domain names to reduce widht of ping output to shown in mobile in better
				args = append(args, "-n")
			}
			switch r.URL.Query().Get("IPVersion") {
			case "IPv4":
				if !isMayIPv4(host) {
					fmt.Fprintf(w, "{\"code\":\"IPVersionMissMatch\"}")
					return
				}
				args = append(args, "-4")
			case "IPv6":
				if !isMayIPv6(host) {
					fmt.Fprintf(w, "{\"code\":\"IPVersionMissMatch\"}")
					return
				}
				args = append(args, "-6")
			case "Default":
			default:
				fmt.Fprintf(w, "{\"code\":\"WrongIPVersion\"}")
				return
			}

			// If requests comes from iframe (not term), add style to iframe
			if r.URL.Query().Get("term") == "" {
				w.Header().Set("content-type", "text/html; charset=utf-8")
				fmt.Fprintf(w, iframeStyle)
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
			return
			/****************************
			Live output for tracert frame
			*****************************/
		case "tracert":
			args := []string{}

			if r.URL.Query().Get("mobile") == "1" {
				args = append(args, "-n -q 1")
			} else {
				args = append(args, "-q 3")
				args = append(args, "-e")
			}

			switch r.URL.Query().Get("IPVersion") {
			case "IPv4":
				if !isMayIPv4(host) {
					fmt.Fprintf(w, "{\"code\":\"IPVersionMissMatch\"}")
					return
				}
				args = append(args, "-4")
			case "IPv6":
				if !isMayIPv6(host) {
					fmt.Fprintf(w, "{\"code\":\"IPVersionMissMatch\"}")
					return
				}
				args = append(args, "-6")
			case "Default":
			default:
				fmt.Fprintf(w, "{\"code\":\"WrongIPVersion\"}")
				return
			}

			if r.URL.Query().Get("term") == "" {
				w.Header().Set("content-type", "text/html; charset=utf-8")
				fmt.Fprintf(w, iframeStyle)
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
		default:
			// if any unknown function name given.
			w.WriteHeader(http.StatusNotFound)
			// requestDump, err := httputil.DumpRequest(r, true)
			// if err != nil {
			// 	fmt.Println(err)
			// }
			// fmt.Fprintf(w, string(requestDump))
			fmt.Fprintf(w, "Function "+r.URL.Query().Get("funcType")+" is not found")
			return
		}

	})

	// return as a webServer
	return &http.Server{

		Addr:     *listenAddr,
		Handler:  limitMiddleware(router),
		ErrorLog: logger,
		/* Close sockets */
		ReadTimeout:  5 * time.Second,  // Input Time Out
		WriteTimeout: 60 * time.Second, // Output Time Out
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-disposition", "filename='response.lg'")

}
