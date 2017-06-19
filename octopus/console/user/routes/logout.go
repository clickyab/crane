package routes

import (
	"context"
	"net/http"

	"github.com/clickyab/services/assert"

	"github.com/clickyab/services/eav"
)

// logout is for the logout from the system
// @Route {
// 		url = /logout
//		method = get
//		middleware = Authenticate
//      	200 = controller.NormalResponse
// }
func (c Controller) logout(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	d := MustGetToken(ctx)
	if d == "" {
		c.OKResponse(w, struct {
			Status string `json:"status"`
		}{
			Status: "Already logged out",
		})
	}
	err := eav.NewEavStore(d).Drop()
	assert.Nil(err)

	c.OKResponse(w, struct {
		Status string `json:"status"`
	}{
		Status: "logged out successfully",
	})
}
