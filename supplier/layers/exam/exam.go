package native

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httputil"

	"github.com/clickyab/services/version"
)

func exam(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/plain")
	v, err := json.Marshal(version.GetVersion())
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	d, err := httputil.DumpRequest(r, true)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return

	}
	_, _ = w.Write(v)
	_, _ = w.Write(d)
}
