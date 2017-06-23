package routes

import (
	"context"
	"net/http"

	"errors"
	"strings"
	"time"

	"clickyab.com/exchange/octopus/console/user/aaa"
	"clickyab.com/exchange/octopus/console/user/routes"
	"clickyab.com/exchange/octopus/models"
	"github.com/clickyab/services/array"
	"github.com/clickyab/services/framework"
	"github.com/rs/xmux"
)

type supplierReportResponse struct {
	Data  []models.SupplierReporter `json:"data"`
	Count int64                     `json:"count"`
}

// supplier report in system
// @Route {
// 		url = /supplier/:from/:to
//		method = get
//		middleware = routes.Authenticate
//		400 = controller.ErrorResponseSimple
// }
func (c Controller) supplier(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	currentUser := routes.MustGetUser(ctx)
	var toTime time.Time
	var res supplierReportResponse
	p, count := framework.GetPageAndCount(r, false)
	from := xmux.Param(ctx, "from")
	if from == "" {
		c.BadResponse(w, errors.New("start date not valid"))
		return
	}
	to := xmux.Param(ctx, "to")
	fromTime, err := time.Parse("20060102", from)
	if err != nil {
		c.BadResponse(w, errors.New("start date not valid"))
		return
	}
	toTime, err = time.Parse("20060102", to)
	if err != nil {
		toTime = fromTime.AddDate(0, 0, 1)
	}
	fromTimeInt := models.FactTableID(fromTime)
	toTimeInt := models.FactTableID(toTime)
	s := r.URL.Query().Get("sort")
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		parts = append(parts, "asc")
	}
	sort := parts[0]
	if !array.StringInArray(sort, "id") {
		sort = ""
	}
	order := strings.ToUpper(parts[1])
	if !array.StringInArray(order, "ASC", "DESC") {
		order = aaa.DefaultOrder
	}
	result, num := models.NewManager().FillSupplierReport(p, count, sort, order, fromTimeInt, toTimeInt, currentUser)
	res.Data = result
	res.Count = num
	c.OKResponse(w, res)
}
