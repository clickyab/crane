package asset

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/clickyab/services/xlog"

	"clickyab.com/crane/supplier/lists"

	"clickyab.com/crane/demand/entity"

	"clickyab.com/crane/models/item"
	openrtb "clickyab.com/crane/openrtb/v2.5"
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

	var err error

	var msg = ""

	defer func() {

		if err != nil {

			w.Header().Set("error", msg)

			_, _ = w.Write([]byte(msg))

			w.WriteHeader(http.StatusInternalServerError)

		}

	}()

	pl := &item.Asset{}

	u, ok := ctx.Value(user.KEY).(*openrtb.User)

	if !ok {
		err = fmt.Errorf("user not found")
		msg = err.Error()
		return
	}
	pl.User = u

	ul, err := url.Parse(r.URL.Query().Get("url"))
	if err != nil {
		msg = "list does not exists"
		return
	}
	ti := r.URL.Query().Get("title")
	if ti == "" {
		msg = "title does not exists"
		return
	}
	pl.FTitle = ti
	pl.FURL = ul.String()
	l, err := item.CheckList(r.URL.Query().Get("list"))
	if err != nil {
		msg = "list doesn't exists"
		return
	}
	if ul.Host != l.Domain {
		msg = "domain doesn't match"
		return
	}
	if err != nil {
		msg = "list does not exists"
		return
	}

	bl, err := strconv.ParseBool(r.URL.Query().Get("isavailable"))
	if err != nil {
		msg = "availability is not defined"
		return
	}
	pl.IsAvailable = bl

	pl.FBrand = r.URL.Query().Get("brand")
	d, err := strconv.ParseInt(r.URL.Query().Get("discount"), 10, 64)
	if err != nil && r.URL.Query().Get("discount") != "" {
		msg = "discount is not defined"
		return
	}
	pl.FDiscount = d

	p, err := strconv.ParseInt(r.URL.Query().Get("price"), 10, 64)
	if err != nil && r.URL.Query().Get("price") != "" {
		msg = "price is not defined"
		return
	}
	pl.FPrice = p

	img, err := url.Parse(r.URL.Query().Get("img"))
	if err != nil {
		msg = "image url is not valid"
		return
	}
	pl.FImg = img.String()

	pl.FSKU = r.URL.Query().Get("sku")
	pl.FBrand = r.URL.Query().Get("brand")

	// go func() {
	err = lists.SetLists(ctx, u.Id, pl.Hash(), l.KEY)
	fmt.Println("ASSET", u.Id, pl.Hash(), l.KEY)
	if err != nil {
		xlog.GetWithError(ctx, err).Debug("set list id")
	}
	assetChan <- pl
	// }()
}

func init() {
	router.Register(&controller{})
	go assetHandler()
}

var assetChan = make(chan entity.Item)

func assetHandler() {
	d := time.Second * 5
	t := time.After(d)
	for {
		select {
		case <-t:
			flush()
			t = time.After(d)
		case a := <-assetChan:
			items[a.Hash()] = a
			if len(items) > 100 {
				flush()
				t = time.After(d)
			}
		}
	}
}

var items = make(map[string]entity.Item)

func flush() {
	ts := make([]entity.Item, 0)
	for _, v := range items {
		ts = append(ts, v)
	}
	if len(ts) == 0 {
		return
	}
	_ = item.AddAssets(context.Background(), ts)
	items = make(map[string]entity.Item)
}
