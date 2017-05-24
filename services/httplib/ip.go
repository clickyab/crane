package httplib

import (
	"net"
	"net/http"
)

const (
	headerXForwardedFor  = "X-Forwarded-For"
	headerXRealIP        = "X-Real-IP"
	headerCFConnectingIP = "CF-Connecting-IP"
)

// RealIP extract ip from request
func RealIP(r *http.Request) string {
	ra := r.RemoteAddr
	if ip := r.Header.Get(headerCFConnectingIP); ip != "" {
		ra = ip
	} else if ip := r.Header.Get(headerXForwardedFor); ip != "" {
		ra = ip
	} else if ip := r.Header.Get(headerXRealIP); ip != "" {
		ra = ip
	} else {
		ra, _, _ = net.SplitHostPort(ra)
	}
	return ra
}
