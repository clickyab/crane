package pixel

import (
	"context"
	"encoding/base64"
	"net/http"

	"clickyab.com/crane/demand/entity"
	"github.com/clickyab/services/assert"
)

const emptyPixel = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABAQMAAAAl21bKAAAAA1BMVEUAAACnej3aAAAAAXRSTlMAQObYZgAAAApJREFUCNdjYAAAAAIAAeIhvDMAAAAASUVORK5CYII="

// Render an empty pixel
func Render(_ context.Context, w http.ResponseWriter, _ entity.Context) error {
	data, err := base64.StdEncoding.DecodeString(emptyPixel)
	assert.Nil(err)
	w.Header().Set("Content-Type", "image/png")
	_, _ = w.Write(data)
	return nil
}
