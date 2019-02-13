package user

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"clickyab.com/crane/supplier/lists"

	openrtb "clickyab.com/crane/openrtb/v2.5"

	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
)

// key of context value
type key int

const (
	uidKey = "a8f5f167f44f4964e6c998dee827110c"
	// KEY of user in context
	KEY = key(1000)
)

type middleware struct {
}

func (middleware) PreRoute() bool {
	return true
}

func (middleware) Handler(next framework.Handler) framework.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {

		fmt.Println("DOMAIN:", r.URL.Hostname())
		user := &openrtb.User{
			Data: make([]*openrtb.UserData, 0),
		}

		if uc, err := r.Cookie(uidKey); err != nil {

			user.Id = "cyb-" + fmt.Sprintf("%d", rand.Int31n(1e9)+1e9)
			http.SetCookie(w,
				&http.Cookie{
					Domain:  "clickyab.com",
					Expires: time.Now().AddDate(2, 0, 0),
					Value:   user.Id,
					Name:    uidKey,
					Path:    "/",
				})

			http.SetCookie(w,
				&http.Cookie{
					Domain:  "3rdad.com",
					Expires: time.Now().AddDate(2, 0, 0),
					Value:   user.Id,
					Name:    uidKey,
					Path:    "/",
				})

		} else {

			user.Id = uc.Value
			if ud, err := lists.GetLists(ctx, user.Id); err == nil {
				user.Data = append(user.Data, ud)
			}
		}
		x := context.WithValue(ctx, KEY, user)
		r = r.WithContext(x)
		next(x, w, r)
	}
}

func init() {
	router.RegisterGlobalMiddleware(&middleware{})
}
