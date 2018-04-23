package ortb

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"clickyab.com/crane/demand/layers/internal/js"
)

var (
	tracker             = js.MustAsset("tracker.js")
	conversionTrackerJS = "/conversion/clickyab-tracking.js"
)

// Serve jsV2
func trackerJS(_ context.Context, w http.ResponseWriter, r *http.Request) {
	//proto := framework.Scheme(r)
	u := url.URL{
		Host: r.Host,
		//	Scheme: proto,

		// Path: router.MustPath("conversion", map[string]string{}),
		Path: "/conversion", // TODO: we re-write the path for BC, must switch to real path
	}
	// Exactly once!
	str := strings.Replace(string(tracker), "{{.URL}}", u.String(), 1)
	w.Header().Set("Content-Type", "application/javascript")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(str))
}
