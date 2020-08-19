package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
)

// This is SECURITY Control Middleware
func middlewareHTTPHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Client IP address for Rate Limitter.
		ip, _, err := net.SplitHostPort(r.RemoteAddr) // _ is port but not required.

		// If system cannot parse the addr. (may be socket in next)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"code":"%s","err":"%s"}`, http.StatusText(http.StatusInternalServerError), err.Error())
			return
		}

		// If every basic security checks is ok, Let's control the ratio of request from client to preventing over usage and also any proxy.

		limiter := limiter.GetLimiter(ip)
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		// Define request method.
		httpScheme := "https"
		if r.TLS == nil {
			httpScheme = "http"
		}

		// This may be requests comes from different domain, port or http scheme.
		if httpScheme+"://"+r.Host != serverConfigOldToBeReplaced["ThisServerURL"] {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprintf(w, `{"code":"NotAcceptable","err":"Request Domain is different", "requestDomain": "%s"}`, httpScheme+"://"+r.Host)
			return
		}

		// Referer Control for forbiding embedding this service to unknown websites.
		middleReferer := r.Header.Get("referer")

		// Parse referer url to get host.
		u, err := url.Parse(middleReferer)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"code":"%s","err":"%s"}`, http.StatusText(http.StatusInternalServerError), err.Error())
			return
		}
		// Check is referer in list.
		if !contains(allowedReferrers, u.Host) {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprintf(w, `{"code":"NotAcceptable","err":"Referrer is not allowed", "referrer": "%s"}`, u.Host)
			return
		}

		w.Header().Set("CONTENT-SECURITY-POLICY", "default-src 'none'; style-src 'unsafe-inline';base-uri 'self';")
		w.Header().Set("Access-Control-Allow-Origin", u.Host)

		next.ServeHTTP(w, r)
	})
}
