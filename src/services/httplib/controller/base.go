package controller

import (
	"net/http"
	"services/assert"
	"services/httplib/middlewares"
	"services/trans"

	"gopkg.in/labstack/echo.v3"
)

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

// Base is the base controller for all other controllers
type Base struct {
}

// BadResponse is 400 request
func (c Base) BadResponse(ctx echo.Context, err error) error {
	ctx.JSON(http.StatusBadRequest, ErrorResponseSimple{Error: trans.EE(err)})
	return err
}

// ForbiddenResponse is 403 request
func (c Base) ForbiddenResponse(ctx echo.Context, err error) error {
	ctx.JSON(http.StatusForbidden, ErrorResponseSimple{Error: trans.EE(err)})
	return err
}

// NotFoundResponse is 404 request
func (c Base) NotFoundResponse(ctx echo.Context, err error) error {
	var res = ErrorResponseSimple{}
	if err != nil {
		res.Error = trans.EE(err)
	} else {
		res.Error = trans.E(http.StatusText(http.StatusNotFound))
	}
	ctx.Response().Header().Add("error", res.Error.Error())
	ctx.JSON(http.StatusNotFound, res)

	return res.Error
}

// OKResponse is 200 request
func (c Base) OKResponse(ctx echo.Context, res interface{}) error {
	if res == nil {
		res = NormalResponse{}
	}
	ctx.JSON(http.StatusOK, res)

	return nil
}

// MustGetPayload is for payload middleware
func (c Base) MustGetPayload(ctx echo.Context) interface{} {
	obj, ok := middlewares.GetPayload(ctx)
	assert.True(ok, "[BUG] payload un-marshaller failed")

	return obj
}
