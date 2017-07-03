package restful

import (
	"bytes"
	"html/template"
	"strings"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/layers/output/internal"
)

type singleAd struct {
	Link   string
	Width  string
	Height string
	Src    string
	Tiny   bool
}

// makeSingleAdData returns web template single ad data
func makeSingleAdData(advertise entity.Advertise, impression entity.Impression, slot entity.Slot, cp entity.ClickProvider) (string, error) {
	scheme := impression.Protocol()

	var src string = advertise.Media()
	if scheme == "https" {
		src = strings.Replace(src, "http://", "https://", -1)
	}
	sa := singleAd{
		Link:   cp.ClickURL(slot, impression),
		Height: string(slot.Height()),
		Width:  string(slot.Width()),
		Src:    src,
		Tiny:   true,
	}
	var singleAdTemplate = template.Must(template.New("single_ad").Parse(internal.SingleAdTemp))
	buf := &bytes.Buffer{}
	if err := singleAdTemplate.Execute(buf, sa); err != nil {
		return "", err
	}

	return buf.String(), nil
}
