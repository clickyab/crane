package routes

import (
	"context"
	"errors"
	"net/http"

	"fmt"

	"time"

	"clickyab.com/exchange/octopus/console/user/aaa"
	"github.com/Sirupsen/logrus"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/eav"
	"github.com/clickyab/services/random"
	"golang.org/x/crypto/bcrypt"
)

var loginTime = config.RegisterDuration("console.user.logged", 24*time.Hour, "user logged in duration")

type responseLoginOK struct {
	ID       int64        `json:"id"`
	Email    string       `json:"email"`
	Token    string       `json:"token"`
	UserType aaa.UserType `json:"user_type"`
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
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pl.Password)) != nil {
		c.ForbiddenResponse(w, errors.New("wrong password"))
		return
	}

	token := getNewToken(user)

	c.OKResponse(w, responseLoginOK{
		ID:       user.ID,
		Email:    user.Email,
		Token:    token,
		UserType: user.UserType,
	})
}

func getNewToken(user *aaa.User) string {
	t := fmt.Sprintf("%d:%s", user.ID, <-random.ID)
	logrus.Warn(t)
	generated := eav.NewEavStore(t).SetSubKey("token", user.Token)
	assert.Nil(generated.Save(loginTime.Duration()))
	return t
}
