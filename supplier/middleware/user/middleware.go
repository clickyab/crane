package user

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/clickyab/services/kv"

	"github.com/clickyab/services/random"

	openrtb "clickyab.com/crane/openrtb/v2.5"

	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
)

// key of context value
type key int

const (
	uidKey     = "a8f5f167f44f4964e6c998dee827110c"
	lsKey      = "1ab482db173cf15b7fdbe7c38feb6748"
	userPrefix = "USR"
	// KEY of user in context
	KEY = key(1000)
)

type middleware struct {
}

func (middleware) PreRoute() bool {
	return true
}

func extractList(c string, u *openrtb.User) {

	res := &openrtb.UserData{
		Name:    "list",
		Id:      "1",
		Segment: []*openrtb.UserData_Segment{},
	}

	for _, v := range strings.Split(string(c), ",") {
		k := kv.NewEavStore(fmt.Sprintf("%s_%s_%s", userPrefix, u.GetId(), v))
		if len(k.AllKeys()) == 0 {
			continue
		}
		xl := make([]string, len(k.AllKeys()))
		for lk := range k.AllKeys() {
			xl = append(xl, lk)
		}

		res.Segment = append(res.Segment, &openrtb.UserData_Segment{
			Id:    v,
			Value: strings.Join(xl, ","),
		})
	}
}

func (middleware) Handler(next framework.Handler) framework.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		user := &openrtb.User{
			Data: make([]*openrtb.UserData, 0),
		}
		if uc, err := r.Cookie(uidKey); err != nil {
			user.Id = <-random.ID
			http.SetCookie(w,
				&http.Cookie{
					Domain:  ".clickyab.com",
					Expires: time.Now().AddDate(2, 0, 0),
					Value:   user.Id,
					Name:    uidKey,
					Path:    "/",
				})
		} else {
			user.Id = uc.Value
		}

		if ls, err := r.Cookie(lsKey); err != nil {
			extractList(ls.Value, user)
		}
		next(context.WithValue(ctx, KEY, user), w, r)
	}
}

func init() {
	router.RegisterGlobalMiddleware(&middleware{})
}
