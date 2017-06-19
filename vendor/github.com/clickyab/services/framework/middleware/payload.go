package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/clickyab/services/array"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/trans"
)

type contextKey string

const (
	// ContextBody is the context key for the body unmarshalled object
	ContextBody contextKey = "_body"
)

type (
	// GroupError is the type of error is used by route validators
	GroupError map[string]trans.T9Error

	// Validator is used to validate the input
	Validator interface {
		// Validate return true, nil on no error, false ,error map in error string
		Validate(context.Context, http.ResponseWriter, *http.Request) error
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
func PayloadUnMarshallerGenerator(pattern interface{}) framework.Middleware {
	return func(next framework.Handler) framework.Handler {
		return func(c context.Context, w http.ResponseWriter, r *http.Request) {
			// Make sure the request is POST or PUT since DELETE and GET must not have payloads
			method := strings.ToUpper(r.Method)
			assert.True(
				!array.StringInArray(method, "GET", "DELETE"),
				"[BUG] Get and Delete must not have request body",
			)
			// Create a copy
			cp := reflect.New(reflect.TypeOf(pattern)).Elem().Addr().Interface()
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(cp)
			if err != nil {
				w.Header().Set("error", trans.T("invalid request body").String())
				e := struct {
					Error error `json:"error"`
				}{
					Error: trans.E("invalid request body"),
				}

				framework.JSON(w, http.StatusBadRequest, e)
				return
			}
			if valid, ok := cp.(Validator); ok {
				if errs := valid.Validate(c, w, r); errs == nil {
					c = context.WithValue(c, ContextBody, cp)
				} else {
					w.Header().Set("error", trans.T("invalid request body").String())
					framework.JSON(w, http.StatusBadRequest, translate(errs))
					return
				}
			} else {
				// Just add it, no validation
				c = context.WithValue(c, ContextBody, cp)
			}
			next(c, w, r)
		}
	}
}

// GetPayload from the request
func GetPayload(c context.Context) (interface{}, bool) {
	t := c.Value(ContextBody)
	return t, t != nil
}
