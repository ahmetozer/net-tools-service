package main
import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
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

		fmt.Fprintf(w, "Error page %s not found.", r.URL.Query().Get("param1"))
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

			switch functype {
			case "icmp4":
				fmt.Fprintf(w, "Error page %s not found.", r.URL.Query().Get("param1"))
			case "icmp6":
				fmt.Println("Linux.")
			case "tcp":
			case "webkontrol":
			case "time":
			case "whois":
			case "nslookup":
			case "ping":
			case "ping6":
			case "tracert":
			case "tracert6":
			case "curlkontrol":
			case "curltamkontrol":
			case "curldurum":
			default:
				// freebsd, openbsd,
				// plan9, windows...
				fmt.Printf(functype+host+name+port+id+mobile+term)
			}



			// return success
			setLiveOutputHeaders(w);
			w.WriteHeader(http.StatusOK)

			cmd := exec.Command("ping", "-c 5", "1.1.1.1")
			// Organize pipelines
			pipeIn, pipeWriter := io.Pipe()
			cmd.Stdout = pipeWriter
			cmd.Stderr = pipeWriter
			// Pass to web output
			go writeCmdOutput(w, pipeIn)

			// Run commands
			cmd.Run()
			pipeWriter.Close()

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
		//w.Header().Set("X-Accel-Buffering", "no")
	}
