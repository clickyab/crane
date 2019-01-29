package asset

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	openrtb "clickyab.com/crane/openrtb/v2.5"
	"clickyab.com/crane/supplier/layers/entities"
	"clickyab.com/crane/supplier/middleware/user"

	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
)

type controller struct {
}

// Routes is for registering routes
func (controller) Routes(r framework.Mux) {
	r.GET("asset", "/api/asset", getAsset)
}

// https://t.clickyab.com/pixel/collecting?list=bamilo&url=https%3A%2F%2Fwww.bamilo.com%2Fproduct%2Factive-%25D9%25BE%25D9%2588%25D8%25AF%25D8%25B1-%25D9%2584%25D8%25A8%25D8%25A7%25D8%25B3%25D8%25B4%25D9%2588%25DB%258C%25DB%258C-%25D9%2585%25D8%25A7%25D8%25B4%25DB%258C%25D9%2586%25DB%258C-500-%25DA%25AF%25D8%25B1%25D9%2585%25DB%258C-9395631%2F&img=%2F%2Fmedia.bamilo.com%2Fp%2Factive-1843-1365939-1-zoom.jpg&title=%D9%BE%D9%88%D8%AF%D8%B1%20%D9%84%D8%A8%D8%A7%D8%B3%D8%B4%D9%88%DB%8C%DB%8C%20%D9%85%D8%A7%D8%B4%DB%8C%D9%86%DB%8C%20500%20%DA%AF%D8%B1%D9%85%DB%8C&price=5355&discount=10&sku=AC696OT084RNIALIYUN&isavailable=true&category=%D8%B3%D9%88%D9%BE%D8%B1%D9%85%D8%A7%D8%B1%DA%A9%D8%AA%2C%D8%A8%D9%87%D8%AF%D8%A7%D8%B4%D8%AA%20%D9%85%D9%86%D8%B2%D9%84%2C%D8%B4%D9%88%DB%8C%D9%86%D8%AF%D9%87%20%D9%84%D8%A8%D8%A7%D8%B3&brand=Active
func getAsset(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pl := &entities.Asset{}
	u, ok := ctx.Value(user.KEY).(*openrtb.User)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pl.User = u

	ls := strings.Split(r.URL.Query().Get("list"), ",")
	if len(ls) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	url, err := url.Parse(r.URL.Query().Get("url"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pl.URL = url.String()

	bl, err := strconv.ParseBool(r.URL.Query().Get("isavailable"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pl.IsAvailable = bl

	pl.Brand = r.URL.Query().Get("brand")
	d, err := strconv.ParseInt(r.URL.Query().Get("discount"), 10, 64)
	if err != nil && r.URL.Query().Get("discount") != "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pl.Discount = d

	p, err := strconv.ParseInt(r.URL.Query().Get("price"), 10, 64)
	if err != nil && r.URL.Query().Get("discount") != "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pl.Price = p

	err = json.NewDecoder(r.Body).Decode(pl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go func() {
		assetChan <- pl
	}()
}

func init() {
	router.Register(&controller{})
}

var assetChan = make(chan *entities.Asset)

// func assetHandler() {
// 	for {
// 		t := time.After(time.Second * 5)
// 		select {
// 		case <-t:
// 			flush()
// 		case a := <-assetChan:
//
// 		}
//
// 	}
// }

func init() {

}
