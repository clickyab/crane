package user

import (
	"context"
	"net/http"
)

// logout is for the logout from the system
// @Route {
// 		url = /logout
//		method = get
//		middleware = Authenticate
//      200 = controller.NormalResponse
// }
func (c Controller) logout(ctx context.Context, w http.ResponseWriter, r *http.Request) {

}
