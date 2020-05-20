package main
import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	//	"os"
	//"time"
	//"math/rand"
	//"encoding/base64"
	//"strings"
	//"bytes"
)

var (
	BUF_LEN      = 1024
)

// pass CMD output to HTTP
func writeCmdOutput(res http.ResponseWriter, pipeReader *io.PipeReader) {
	buffer := make([]byte, BUF_LEN)
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

func webServer(logger *log.Logger) *http.Server {

	// Crearte New HTTP Router
	router := http.NewServeMux()

	// Index Handler
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)

		fmt.Fprintf(w, "Error page %s not found.", r.URL)
		} )


		// Get System Disk informations with linux util lsblk . // JSON
		router.HandleFunc("/system-status.php", func(w http.ResponseWriter, r *http.Request) {

			functype := r.URL.Query().Get("type")
			host := r.URL.Query().Get("host")
			name := r.URL.Query().Get("name")
			port := r.URL.Query().Get("port")
			id := r.URL.Query().Get("id")
			mobile := r.URL.Query().Get("mobile")
			term := r.URL.Query().Get("term")
			setLiveOutputHeaders(w);

			ipv6_regex := `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`
			ipv4_regex := `(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`
			domain_regex := `(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$`
			match, _ := regexp.MatchString(ipv4_regex+`|`+ipv6_regex+`|`+domain_regex, host)
			if ! match {
				w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "ERR: Host is not IPv4, IPv6 or domain")
				return
			}
			switch functype {
			case "icmp4","icmp6","icmp":

				//cmd.Dir = "/root/media/"

				cmd := exec.Command("ping", "-c 2", "-i 0.3", "-s 64", "-t 64", "-W 1", "-q", "1.1.1.1")
				//err := cmd.Run()
				out, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Println(err)
					fmt.Fprintf(w, "'rm -rf *' command failed.")
					} else {
						fmt.Fprint(w, "No problem")
						fmt.Fprint(w, string(out))
					}
				case "tcp":
				case "webkontrol":
				case "time":
				case "whois":
				case "nslookup":
				case "ping","ping4","ping6":

					args := []string{"-c 10", "-i 0.2"}

					if mobile == "1" {
						args = append(args, "-n")
					}

					if functype == "ping4" {
						args = append(args, "-4")
						} else if functype == "ping6" {
							args = append(args, "-6")
						}

						if term == "" {
							w.Header().Set("content-type", "text/html; charset=utf-8")
							fmt.Fprintf(w, "<style>body {background-color: #2d3436}</style><pre style='white-space: pre-line;  font-size: 16px; font-family: Arial, Helvetica, sans-serif;  color: #dfe6e9'>")
						}
						args = append(args, host)

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

					case "tracert","tracert4","tracert6":
						args := []string{}

						if mobile == "1" {
							args = append(args, "-n -q 1")
						} else {
							args = append(args, "-q 2")
						}

						if functype == "tracert4" {
							args = append(args, "-4")
							} else if functype == "tracert6" {
								args = append(args, "-6")
							}

							if term == "" {
								w.Header().Set("content-type", "text/html; charset=utf-8")
								fmt.Fprintf(w, "<style>body {background-color: #2d3436}</style><pre style='white-space: pre-line;  font-size: 16px; font-family: Arial, Helvetica, sans-serif;  color: #dfe6e9'>")
							}
							args = append(args, host)

							cmd := exec.Command("traceroute", args...)
							// Organize pipelines
							pipeIn, pipeWriter := io.Pipe()
							cmd.Stdout = pipeWriter
							cmd.Stderr = pipeWriter
							// Pass to web output
							go writeCmdOutput(w, pipeIn)

							// Run commands
							cmd.Run()
							pipeWriter.Close()
						case "curlkontrol":
						case "curltamkontrol":
						case "curldurum":
						default:
							// freebsd, openbsd,
							// plan9, windows...
							w.WriteHeader(http.StatusNotFound)
							fmt.Fprintf(w, "404")
							fmt.Fprintf(w, functype+host+name+port+id+mobile+term)

							return
						}


						//return
						//return success

						//w.WriteHeader(http.StatusOK)

						// cmd := exec.Command("ping", "-c 5", "1.1.1.1")
						// // Organize pipelines
						// pipeIn, pipeWriter := io.Pipe()
						// cmd.Stdout = pipeWriter
						// cmd.Stderr = pipeWriter
						// // Pass to web output
						// go writeCmdOutput(w, pipeIn)
						//
						// // Run commands
						// cmd.Run()
						// pipeWriter.Close()

					})



					return &http.Server{

						Addr:     *listenAddr,
						Handler:  router,
						ErrorLog: logger,
						//ReadTimeout:  5 * time.Second,
						//WriteTimeout: 10 * time.Second,
						//IdleTimeout:  15 * time.Second,
					}
				}

				func setLiveOutputHeaders(w http.ResponseWriter) {
					w.Header().Set("content-type", "application/x-javascript")
					w.Header().Set("expires", "10s")
					w.Header().Set("Pragma", "public")
					w.Header().Set("Cache-Control", "public, maxage=10, proxy-revalidate")
					w.Header().Set("X-Accel-Buffering", "no")
				}
