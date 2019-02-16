package asset

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
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
		xlog.GetWithError(ctx, err).Debug()
		msg = "invalid url"
		return
	}
	ti := r.URL.Query().Get("title")
	if ti == "" {
		msg = "title does not exists"
		xlog.GetWithError(ctx, fmt.Errorf(msg)).Debug()
		return
	}
	pl.FTitle = ti
	pl.FURL = ul.String()
	l, err := item.CheckList(r.URL.Query().Get("list"))
	if err != nil {
		xlog.GetWithError(ctx, err).Debug()
		msg = "list doesn't exists"
		return
	}
	if !strings.HasSuffix(ul.Host, l.Domain) {
		err = fmt.Errorf("domain doesn't match")
		xlog.GetWithError(ctx, err).Debug()
		msg = err.Error()
		return
	}

	bl, err := strconv.ParseBool(r.URL.Query().Get("isavailable"))
	if err != nil {
		xlog.GetWithError(ctx, err).Debug()
		msg = "availability is not defined"
		return
	}
	pl.IsAvailable = bl

	pl.FBrand = r.URL.Query().Get("brand")
	d, err := strconv.ParseInt(r.URL.Query().Get("discount"), 10, 64)
	if err != nil && r.URL.Query().Get("discount") != "" {
		xlog.GetWithError(ctx, err).Debug()
		msg = "discount is not defined"
		return
	}
	pl.FDiscount = d

	p, err := strconv.ParseInt(r.URL.Query().Get("price"), 10, 64)
	if err != nil && r.URL.Query().Get("price") != "" {
		xlog.GetWithError(ctx, err).Debug()
		msg = "price is not defined"
		return
	}
	pl.FPrice = p

	img, err := url.Parse(r.URL.Query().Get("img"))
	if err != nil {
		xlog.GetWithError(ctx, err).Debug()
		msg = "image url is not valid"
		return
	}
	pl.FImg = img.String()
	pl.FSKU = r.URL.Query().Get("sku")
	pl.FBrand = r.URL.Query().Get("brand")

	go func() {
		err = lists.SetLists(ctx, u.Id, pl.Hash(), l.KEY)
		if err != nil {
			xlog.GetWithError(ctx, err).Debug("set list id")
		}
		xlog.GetWithError(ctx, err).Debug("add list to channel")
		assetChan <- pl
	}()
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
			lock.Lock()
			items[a.Hash()] = a
			lock.Unlock()
			if len(items) > 100 {
				flush()
				t = time.After(d)
			}
		}
	}
}

var items = make(map[string]entity.Item)
var lock = sync.Mutex{}

func flush() {
	lock.Lock()
	defer lock.Unlock()
	ts := make([]entity.Item, 0)
	for _, v := range items {
		ts = append(ts, v)
	}
	if len(ts) == 0 {
		return
	}

	err := item.AddAssets(context.Background(), ts)
	if err != nil {
		fmt.Println("ASSET FLUSH: ", err.Error())
	}

	items = make(map[string]entity.Item)
}
