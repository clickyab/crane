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
	native.Parse(nativeTmpl)
	native.Parse(adTmpl)
}

const nativeTmpl = `{{define "ads"}}<div class="cyb-holder cyb-custom-holder" style="font-size: {{.FontSize}}">
	<style>
	{{.Style}}
	</style>
    <div class="cyb-title-holder cyb-custom-title-holder">
        <div class="cyb-title-before cyb-custom-title-before"></div>
        <div class="cyb-title cyb-custom-title">{{.Title}}</div>
        <div class="cyb-title-after cyb-custom-title-after"></div>
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
