package controller

import (
	"net/http"

	"context"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/middleware"
	"github.com/clickyab/services/trans"
)

// Base is the controller base for all controllers
type Base struct {
}

// NormalResponse is for 2X responses
type NormalResponse struct {
}

// ComplexResponse for the result, when the result type in not in the structure
type ComplexResponse map[string]trans.T9Error

// ErrorResponseMap is the map for the response with detail error mapping
type ErrorResponseMap map[string]trans.T9Error

// ErrorResponseSimple is the type for response when the error is simply a string
type ErrorResponseSimple struct {
	Error trans.T9Error `json:"error"`
}

// BadResponse is 400 request
func (b Base) BadResponse(w http.ResponseWriter, err error) {
	b.JSON(w, http.StatusBadRequest, ErrorResponseSimple{Error: trans.EE(err)})
}

// ForbiddenResponse is 403 request
func (b Base) ForbiddenResponse(w http.ResponseWriter, err error) {
	b.JSON(w, http.StatusForbidden, ErrorResponseSimple{Error: trans.EE(err)})
}

// NotFoundResponse is 404 request
func (b Base) NotFoundResponse(w http.ResponseWriter, err error) {
	var res = ErrorResponseSimple{}
	if err != nil {
		res.Error = trans.EE(err)
	} else {
		res.Error = trans.E(http.StatusText(http.StatusNotFound))
	}
	w.Header().Add("error", res.Error.Error())
	b.JSON(w, http.StatusNotFound, res)
}

// OKResponse is 200 request
func (b Base) OKResponse(w http.ResponseWriter, res interface{}) {
	if res == nil {
		res = NormalResponse{}
	}
	b.JSON(w, http.StatusOK, res)
}

// MustGetPayload is for payload middleware
func (b Base) MustGetPayload(ctx context.Context) interface{} {
	obj, ok := middleware.GetPayload(ctx)
	assert.True(ok, "[BUG] payload un-marshaller failed")

	return obj
}

// JSON is a helper function to write an json in output
func (b Base) JSON(w http.ResponseWriter, code int, data interface{}) {
	framework.JSON(w, code, data)
}
