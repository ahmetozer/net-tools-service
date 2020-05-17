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

func newWebserver(logger *log.Logger) *http.Server {

	// Crearte New HTTP Router
	router := http.NewServeMux()

	// Index Handler
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello  %s", r.URL.Path)
	} )


	// Get System Disk informations with linux util lsblk . // JSON
	router.HandleFunc("/lsblk.json", func(w http.ResponseWriter, r *http.Request) {


				// return success
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
