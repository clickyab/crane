package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clickyab/services/config"
)

var monitoredApps = config.RegisterString(
	"crane.supplier.app.monitored",
	`{"status":1,"apps":[{"name":"snapp","packaage":"cab.snapp.passenger"},{"name":"tap30","packaage":"taxi.tap30.passenger"},{"name":"ajancy","packaage":"com.mammutgroup.ajancy.passenger"},{"name":"digikala","packaage":"com.digikala"},{"name":"bamilo","packaage":"com.bamilo.android"},{"name":"pintapin","packaage":"com.pintapin.pintapin"},{"name":"alibaba","packaage":"ir.alibaba"}]}`,
	"a json file contain the monitored apps, must be valid json, pass it as-is to the client, be careful!",
)

func getInappJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, monitoredApps.String())
}
