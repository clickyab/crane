package web

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
)

var (
	showV2 = MustAsset("show.v2.js")
	showV1 = jsV1(MustAsset("show.js"))
)

// Serve jsV2
func jsV2(_ context.Context, w http.ResponseWriter, r *http.Request) {
	proto := framework.Scheme(r)
	u := url.URL{
		Host:   r.Host,
		Scheme: proto,
		Path:   router.MustPath("multi-ad", map[string]string{}),
	}
	// Exactly once!
	str := strings.Replace(string(showV2), "{{.URL}}", u.String(), 1)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(str))
}

type jsV1 []byte

// Serve jsV2
func (data jsV1) ServeHTTPC(_ context.Context, w http.ResponseWriter, r *http.Request) {
	proto := framework.Scheme(r)
	u := url.URL{
		Host:   r.Host,
		Scheme: proto,
		Path:   router.MustPath("multi-js", map[string]string{}),
	}
	// Exactly once!
	str := strings.Replace(string(data), "{{.URL}}", u.String(), 1)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(str))
}
