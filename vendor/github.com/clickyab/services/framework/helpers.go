package framework

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/clickyab/services/config"
	"github.com/rs/xhandler"
	"github.com/rs/xmux"
)

const (
	headerXForwardedFor   = "X-Forwarded-For"
	headerXForwardedProto = "X-Forwarded-Proto"
	headerXRealIP         = "X-Real-IP"
	headerCFConnectingIP  = "CF-Connecting-IP"
	headerContentType     = "Content-Type"
	jsonMIME              = "application/json;charset=UTF-8"

	// HTTP scheme
	HTTP string = "http"
	// HTTPS scheme
	HTTPS string = "https"
)

// Mix try to mix all middleware with the calling route
func Mix(final Handler, all ...Middleware) Handler {
	res := final
	for i := len(all) - 1; i >= 0; i-- {
		res = all[i](res)
	}

	return res
}

// Any is a function to route all type of request to one handler
func Any(mux *xmux.Mux, route string, handler Handler) {
	mux.GET(route, xhandler.HandlerFuncC(handler))
	mux.POST(route, xhandler.HandlerFuncC(handler))
	mux.PUT(route, xhandler.HandlerFuncC(handler))
	mux.PATCH(route, xhandler.HandlerFuncC(handler))
	mux.DELETE(route, xhandler.HandlerFuncC(handler))
	mux.HEAD(route, xhandler.HandlerFuncC(handler))
	mux.OPTIONS(route, xhandler.HandlerFuncC(handler))
}

// JSON is a helper function to write an json in output
func JSON(w http.ResponseWriter, code int, i interface{}) {
	w.Header().Set(headerContentType, jsonMIME)
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.Encode(i)
}

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

// Scheme extract scheme from request
func Scheme(r *http.Request) string {
	proto := HTTP
	if r.TLS != nil {
		proto = HTTPS
	}
	if xh := strings.ToLower(r.Header.Get(headerXForwardedProto)); xh == HTTPS {
		proto = HTTPS
	}

	return proto
}

var (
	maxPerPage = config.RegisterInt("services.framework.max_per_page", 100, "http maximum item per page")
	minPerPage = config.RegisterInt("services.framework.min_per_page", 1, "http minimum item per page")
	perPage    = config.RegisterInt("services.framework.per_page", 10, "http default item per page")
)

// GetPageAndCount return the p and c variable from the request, if not available
// return the default value
func GetPageAndCount(r *http.Request, offset bool) (int, int) {
	p64, err := strconv.ParseInt(r.URL.Query().Get("p"), 10, 0)
	p := int(p64)
	if err != nil || p < 1 {
		p = 1
	}

	c64, err := strconv.ParseInt(r.URL.Query().Get("c"), 10, 0)
	c := int(c64)
	if err != nil || c > maxPerPage.Int() || c < minPerPage.Int() {
		c = perPage.Int()
	}

	if offset {
		// If i need to make it to offset model then do it here
		p = (p - 1) * c
	}

	return p, c
}
