package input

import (
	"net/http"

	"clickyab.com/crane/crane/entity"
)

type ImpressionLayer interface {
	Transform(string,*http.Request) (entity.Impression, error)
}
