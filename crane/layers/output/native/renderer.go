package native

import (
	"errors"
	"fmt"
	"io"

	"clickyab.com/crane/crane/entity"
)

type renderer struct {
}

// Render write native template to writer
func (n renderer) Render(w io.Writer, imp entity.Context, cl entity.ClickProvider) error {

	slots := make([]entity.Slot, 0)
	for _, s := range imp.Slots() {
		if a := s.WinnerAdvertise(); a != nil && a.Type() == entity.AdTypeNative {
			slots = append(slots, s)
		}
	}

	if len(slots) < 1 {
		return errors.New("there is no ad to render")
	}

	nc := &nativeContainer{}
	nc.Title = imp.Attributes()["title"]
	nc.Style = imp.Attributes()["style"]
	nc.FontSize = imp.Attributes()["fontSize"]
	nc.Position = imp.Attributes()["position"]

	nads := make([]nativeAd, 0)

	for _, a := range slots {
		v, ok := a.WinnerAdvertise().Attributes()["title"].(string)
		if !ok {
			return fmt.Errorf("ad with ID %s does't have title in attributes", a.WinnerAdvertise().ID())
		}
		nads = append(nads, nativeAd{
			Title: v,
			URL:   cl.ClickURL(a, imp),
			Site:  a.WinnerAdvertise().TargetURL(),
			Image: a.WinnerAdvertise().Media(),
		})
	}

	nc.Ads = nads
	_, e := w.Write([]byte(renderNative(*nc)))
	return e
}

// NewRenderer just fill Renderer interface
func NewRenderer() entity.Renderer {
	return renderer{}
}
