package native

import (
	"bytes"
	"text/template"

	"github.com/Sirupsen/logrus"
	"github.com/clickyab/services/assert"
)

type nativeContainer struct {
	Ads      []nativeAd
	Title    string
	Style    string
	FontSize string
	Position string
}

type nativeAd struct {
	Corners string
	Image   string
	Title   string
	// Tracking URL: E.g. native.clickyab.com/?id=3452345
	URL string
	// Destination URL: E.g. example.com
	Site string
}

func renderNative(imp nativeContainer) string {
	buf := &bytes.Buffer{}
	imp.Style = style
	e := native.Lookup("ads").Execute(buf, imp)
	assert.Nil(e)
	return string(buf.Bytes())
}

var addRenderer = func(ads []nativeAd) string {
	t, e := template.New("ad").Funcs(template.FuncMap{"isRound": func(s string) string {
		return "cyb-" + s
	}}).Parse(adTmpl)
	assert.Nil(e)

	b := &bytes.Buffer{}

	// remember to pack each two ad into one div like following example
	//         <div class="cyb-pack cyb-custom-pack">
	// 				<AD>
	//				<AD>
	// 			</div>
	// it's a simple hack to keep all ads (relatively) in same ratio
	p := 0
	for i, ad := range ads {
		if i != 0 && i == p {
			b.WriteString("</div>")
		}
		if i%2 == 0 {
			p += 2
			b.WriteString(`<div class="cyb-pack cyb-custom-pack">`)
		}
		e := t.Lookup("ad").Execute(b, ad)
		logrus.Warn(e)
		assert.Nil(e)

		if len(ads)-1 == i {
			b.WriteString("</div>")
		}

	}

	return b.String()
}

var native = template.New("native").Funcs(template.FuncMap{"renderAds": addRenderer})

func init() {
	native.Funcs(template.FuncMap{"isRound": func(s string) string {
		return "cyb-" + s
	}})
	_, e := native.Parse(nativeTmpl)
	assert.Nil(e)
	_, e = native.Parse(adTmpl)
	assert.Nil(e)
}

const nativeTmpl = `{{define "ads"}}
<div class="cyb-holder cyb-custom-holder" style="font-size: {{.FontSize}}">
    <style>
            {{.Style}}
    </style>
    <div class="cyb-title-holder cyb-custom-title-holder">
        <div class="cyb-title-before cyb-custom-title-before"></div>
        <div class="cyb-title cyb-custom-title">{{.Title}}</div>
        <div class="cyb-title-after cyb-custom-title-after"></div>
        <div class="cyb-logo">
            <a target="_blank" href="https://www.clickyab.com/?ref=icon" class="cyb-logo-container">
                <img class="cyb-logo-color"
                     src="data:image/png;base64,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">
            </a>
        </div>
    </div>
    <div class="cyb-suggests cyb-{{.Position}} cyb-custom-suggests">
        {{renderAds .Ads}}
    </div>
</div>
{{end}}
`

const adTmpl = `{{define "ad"}}
       <div class="cyb-suggest cyb-custom-suggest ">
                <div class="cyb-img-holder cyb-custom-img-holder">
                    <a target="_blank" href="{{.Site}}" onclick="cybOpen(event)" oncontextmenu="cybOpen(event)"
                       ondblclick="cybOpen(event)" data-href="{{.URL}}">
                        <img src="{{.Image}}" alt="{{.Title}}"
                             class="cyb-img {{isRound .Corners}} cyb-custom-img">
                    </a>
                </div>
                <div class="cyb-desc-holder cyb-custom-desc-holder">
                    <div class="cyb-desc cyb-custom-desc">
                        <a target="_blank" href="{{.Site}}" onclick="cybOpen(event)" oncontextmenu="cybOpen(event)"
                           ondblclick="cybOpen(event)" data-href="{{.URL}}">
                            {{.Title}}
                        </a>
                    </div>
                </div>
            </div>
            {{end}}
`
