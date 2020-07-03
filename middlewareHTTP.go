package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
)

func middlewareHTTPHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr) // _ is port but not required.
		if err != nil {
			fmt.Fprintf(w, "userip: %q is not IP:port", r.RemoteAddr)
		}

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
		middleReferer := r.Header.Get("referer")

		u, err := url.Parse(middleReferer)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+err.Error(), http.StatusInternalServerError)
			return
		}

		if u.Scheme+"://"+u.Host != serverConfig["expectedRefererURL"] {
			http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
			return
		}
		limiter := limiter.GetLimiter(ip)
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
