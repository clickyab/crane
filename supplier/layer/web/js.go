package web

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"clickyab.com/crane/demand/entity"
	"github.com/clickyab/services/framework/router"
)

var data = MustAsset("show-ad.es5.js")

// Serve js
func js(_ context.Context, w http.ResponseWriter, r *http.Request) {
	proto := entity.HTTP
	if r.TLS != nil {
		proto = entity.HTTPS
	}
	if xh := strings.ToLower(r.Header.Get("X-Forwarded-Proto")); xh == "https" {
		proto = entity.HTTPS
	}
	u := url.URL{
		Host:   r.Host,
		Scheme: proto.String(),
		Path:   router.MustPath("multi-ad", map[string]string{}),
	}
	// Exactly once!
	str := strings.Replace(string(data), "{{.URL}}", u.String(), 1)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(str))
}
