package routes

import (
	"context"
	"errors"
	"net/http"

	"fmt"

	"clickyab.com/exchange/octopus/console/internal/aaa"
	"github.com/clickyab/services/eav"
	"github.com/clickyab/services/random"
	"golang.org/x/crypto/bcrypt"
)

type responseLoginOK struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// login user in system
// @Route {
// 		url = /login
//		method = post
//      	payload = loginPayload
//		200 = responseLoginOK
//		400 = controller.ErrorResponseSimple
// }
func (c Controller) login(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pl := c.MustGetPayload(ctx).(*loginPayload)

	user, err := aaa.NewAaaManager().FindUserByEmail(pl.Email)

	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pl.Password)) == nil {
		c.ForbiddenResponse(w, errors.New("wrong password"))
	}

	token := getNewToken(user)

	c.OKResponse(w, responseLoginOK{
		ID:    user.ID,
		Email: user.Email,
		Token: token,
	})
}

func getNewToken(user *aaa.User) string {
	t := fmt.Sprintf("%d:%s", user.ID, <-random.ID)
	eav.NewEavStore(t).SetSubKey("token", user.Token)

	return t
}
