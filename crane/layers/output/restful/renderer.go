package restful

import (
	"io"

	"clickyab.com/crane/crane/entity"
)

type render struct {
	data []byte
}

// need to register render{} somewhere
func (r *render) Render(w io.Writer, imp entity.Impression, cp entity.ClickProvider) error {
	err := parse(r, imp)
	if err != nil {
		return err
	}
	_, err = w.Write(r.data)
	return err

}
