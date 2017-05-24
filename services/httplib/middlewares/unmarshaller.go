package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"clickyab.com/exchange/services/assert"
	"clickyab.com/exchange/services/trans"

	"clickyab.com/exchange/services/array"

	echo "gopkg.in/labstack/echo.v3"
)

const (
	// ContextBody is the context key for the body unmarshalled object
	ContextBody string = "_body"
)

type (
	// GroupError is the type of error is used by route validators
	GroupError map[string]trans.T9Error

	// Validator is used to validate the input
	Validator interface {
		// Validate return true, nil on no error, false ,error map in error string
		Validate(echo.Context) error
	}
)

// Error is to make this error object
func (ge GroupError) Error() string {
	tmp := ""

	for i := range ge {
		tmp = fmt.Sprintf("%s : %s\n", i, ge[i])
	}

	return tmp
}

// Translate is a helper function to translate the error to map error required by the
// middleware
func translate(err error) interface{} {

	if err == nil {
		return nil
	}
	switch e := err.(type) {
	case GroupError:
		return e
	default:
		return struct {
			Error error `json:"error"`
		}{
			Error: trans.EE(err),
		}
	}
}

// PayloadUnMarshallerGenerator create a middleware base on the pattern for the request body
func PayloadUnMarshallerGenerator(pattern interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Make sure the request is POST or PUT since DELETE and GET must not have payloads
			method := strings.ToUpper(c.Request().Method)
			assert.True(
				!array.StringInArray(method, "GET", "DELETE"),
				"[BUG] Get and Delete must not have request body",
			)
			// Create a copy
			cp := reflect.New(reflect.TypeOf(pattern)).Elem().Addr().Interface()
			decoder := json.NewDecoder(c.Request().Body)
			err := decoder.Decode(cp)
			if err != nil {
				c.Request().Header.Set("error", trans.T("invalid request body").String())
				e := struct {
					Error error `json:"error"`
				}{
					Error: trans.E("invalid request body"),
				}

				c.JSON(http.StatusBadRequest, e)
				return err
			}
			if valid, ok := cp.(Validator); ok {
				if errs := valid.Validate(c); errs == nil {
					c.Set(ContextBody, cp)
				} else {
					c.Request().Header.Set("error", trans.T("invalid request body").String())
					c.JSON(http.StatusBadRequest, translate(errs))

					return trans.E("invalid request body")
				}
			} else {
				// Just add it, no validation
				c.Set(ContextBody, cp)
			}
			return next(c)
		}
	}
}

// GetPayload from the request
func GetPayload(c echo.Context) (interface{}, bool) {
	t := c.Get(ContextBody)
	return t, t != nil
}
