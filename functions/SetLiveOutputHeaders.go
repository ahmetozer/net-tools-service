package functions

import "net/http"

func SetLiveOutputHeaders(w http.ResponseWriter) {
	w.Header().Set("content-type", "application/x-javascript")
	w.Header().Set("expires", "10s")
	w.Header().Set("Pragma", "public")
	w.Header().Set("Cache-Control", "public, maxage=10, proxy-revalidate")
	w.Header().Set("X-Accel-Buffering", "no")
}
