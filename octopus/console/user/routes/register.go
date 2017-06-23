package routes

import (
	"context"
	"net/http"

	"errors"

	"clickyab.com/exchange/octopus/console/user/aaa"
)

type registrationPayload struct {
	Email    string       `json:"email"`
	Password string       `json:"password"`
	UserType aaa.UserType `json:"user_type"`
}

// login user in system
// @Route {
// 		url = /register
//		method = post
//      	payload = registrationPayload
//		200 = responseLoginOK
//		400 = controller.ErrorResponseSimple
// }
func (c Controller) register(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pl := c.MustGetPayload(ctx).(*registrationPayload)
	if !pl.UserType.IsValid() {
		c.BadResponse(w, errors.New("wrong user type"))
		return
	}
	m := aaa.NewAaaManager()

	usr, err := m.RegisterUser(pl.Email, pl.Password, pl.UserType)
	if err != nil {
		c.BadResponse(w, errors.New("error registering user"))
		return
	}
	token := getNewToken(usr)
	c.OKResponse(w, responseLoginOK{
		ID:       usr.ID,
		Email:    usr.Email,
		Token:    token,
		UserType: usr.UserType,
	})
}
