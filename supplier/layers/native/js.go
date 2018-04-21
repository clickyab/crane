package native

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"clickyab.com/crane/supplier/layers/internal/js"
	"github.com/clickyab/services/framework/router"
)

var (
	nativejs = js.MustAsset("native.js")
)

func getNativeJS(_ context.Context, w http.ResponseWriter, r *http.Request) {
	//proto := framework.Scheme(r)
	u := url.URL{
		Host: r.Host,
		//	Scheme: proto,
		Path: router.MustPath("native", map[string]string{}),
	}
	// Exactly once!
	str := strings.Replace(string(nativejs), "{{.URL}}", u.String(), 1)
	w.Header().Set("Content-Type", "application/javascript")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(str))
}
