package vast

import (
	"io"

	"clickyab.com/crane/crane/entity"
)

type data struct {
	data []byte
}

// Render implementation
func (v *data) Render(writer io.Writer, imp entity.Context, cp entity.ClickProvider) error {
	v.parse(imp, cp)
	_, err := writer.Write(v.data)
	return err
}

/*func init() {

}*/
