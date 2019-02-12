package helper

import (
	"context"
	"fmt"
	"net/http"
	"time"

	openrtb "clickyab.com/crane/openrtb/v2.5"
	"clickyab.com/crane/supplier/lists"
	"github.com/clickyab/services/random"
)

const uidKey = "a8f5f167f44f4964e6c998dee827110c"

func GetUser(ctx context.Context, w http.ResponseWriter, r *http.Request) *openrtb.User {
	fmt.Println("KKKSKSKSKSKSKSKKSKSKS", len(r.Cookies()))

	for _, v := range r.Cookies() {
		fmt.Println("COOOO", v.Name, v.Value)
	}
	user := &openrtb.User{
		Data: make([]*openrtb.UserData, 0),
	}
	if uc, err := r.Cookie(uidKey); err != nil {
		fmt.Println("NOT NIL!!!!+", err)

		user.Id = <-random.ID
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
		fmt.Println("sssssssssssssssss+", uc.Value)

		user.Id = uc.Value
		if ud, err := lists.GetLists(ctx, user.Id); err == nil {
			user.Data = append(user.Data, ud)
		}
	}
	fmt.Println("sssssssssssssssss", user.Id)

	return user
}
