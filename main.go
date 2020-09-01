package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	//"log"
	//"os/exec"
	//"strconv"
	"context"
	"log"
	"net/http"
	"os/exec"
	"os/signal"
	"os/user"
	"time"

	"github.com/ahmetozer/net-tools-service/functions"
)

// Variables
var (
	listenAddr string
)

func main() {

	// Check is this program run in linux
	if runtime.GOOS != "linux" {
		fmt.Println("This program is only works on linux devices")
		return
	}

	// Get Current Username
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	// For security reason this program have to run in www-data user.
	if user.Name != "www-data" {
		fmt.Printf("To secure your system, please run this program as www-data."+
			"Current user %v \n", user.Name)
		os.Exit(3)
	}

	// Check the required programs
	requiredPrograms := []string{"ping", "traceroute", "whois", "nslookup", "curl", "mtr"}
	//var requiredapps [5]bool
	for i, s := range requiredPrograms {
		// Get path of required programs
		_, err := exec.LookPath(s)
		if err == nil {
			//fmt.Printf("Required program %v %v found at %v\n", i+1, s, path)
			//requiredapps[i] = true //save to array
		} else {
			fmt.Printf("Required program %v : %v cannot found.\n", i+1, s)
			//requiredapps[i] = false
			if i < len(requiredPrograms) { //sh and df is must required. If is not found in software than exit.
				fmt.Printf("Please install %v and run this program again\n", s)
				os.Exit(3)
			}
		}
	}

	flag.Bool("help", false, "")
	flag.Bool("h", false, "")
	flag.Usage = func() {}

	// Check the Web Server Certificates. If its not available create self cert.
	certDir := functions.CertCheck()

	// Parse the flags
	flag.Parse()

	listenAddr, ok := os.LookupEnv("listenaddr")
	if ok {
		if listenAddr == "" {
			log.Fatalln("\033[1;31mERR. listenaddr is defined to empty.\033[0m")
		} else {
			fmt.Println("HTTPS is listen at " + listenAddr)
		}

	} else {
		fmt.Println("HTTPS port 443 is used")
		listenAddr = ":443"
	}

	// Get the configs from web server.
	lgServerConfigListLoad()

	// Create a logger for https server
	logger := log.New(os.Stdout, "https: ", log.LstdFlags)
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	server := webServer(logger, listenAddr)
	go gracefullShutdown(server, logger, quit, done)
	logger.Println("Server is ready to handle requests at", listenAddr)
	if err := server.ListenAndServeTLS(certDir+"/cert.pem", certDir+"/key.pem"); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("\033[1;31m Could not listen on %s: %v\n\033[0m", listenAddr, err)
	}
	<-done
	logger.Println("Server stopped")

	fmt.Printf("Program Closed.")
}

var help = `
Please visit: https://github.com/ahmetozer/net-tools-service
`

//Grace Full Shutdown
func gracefullShutdown(server *http.Server, logger *log.Logger, quit <-chan os.Signal, done chan<- bool) {
	<-quit
	logger.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("\033[1;31mCould not gracefully shutdown the server: %v\n\033[0m", err)
	}
	close(done)
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
