package video

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"clickyab.com/crane/supplier/layer/internal/js"
	"github.com/clickyab/services/framework/router"
)

var (
	videojs  = js.MustAsset("videojs.js")
	jwplayer = js.MustAsset("jwplayer.js")
)

func getVideojs(_ context.Context, w http.ResponseWriter, r *http.Request) {
	//proto := framework.Scheme(r)
	u := url.URL{
		Host: r.Host,
		//	Scheme: proto,
		Path: router.MustPath("vast", map[string]string{}),
	}
	// Exactly once!
	str := strings.Replace(string(videojs), "{{.URL}}", u.String(), 1)
	w.Header().Set("Content-Type", "application/javascript")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(str))
}

func getJwplayer(_ context.Context, w http.ResponseWriter, r *http.Request) {
	//proto := framework.Scheme(r)
	u := url.URL{
		Host: r.Host,
		//	Scheme: proto,
		Path: router.MustPath("vast", map[string]string{}),
	}
	// Exactly once!
	str := strings.Replace(string(jwplayer), "{{.URL}}", u.String(), 1)
	w.Header().Set("Content-Type", "application/javascript")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(str))
}
