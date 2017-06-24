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
	URL     string
	Site    string
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

// DO NOT EDIT THIS STYLE
// this can regenerate from nativeStyle.less
// please update "less" file if it's necessary
const style = `
.cyb-holder{direction:rtl;box-sizing:border-box;width:100%;height:auto;font-size:12pt;line-height:1.4em;
text-rendering:optimizeLegibility}.cyb-holder .cyb-zero{height:0!important;margin:0!important;
padding:0!important}.cyb-holder .cyb-none{display:none;visibility:hidden}.cyb-holder .cyb-title-holder{
display:flex;margin-bottom:.3em}.cyb-holder .cyb-title-holder .cyb-title-after{flex:1;background:#eee;margin:.7em 0}
.cyb-holder .cyb-title-holder .cyb-title{color:inherit;padding:10px;font-size:1.1em;font-weight:500}.cyb-holder
.cyb-suggests{display:flex;flex-wrap:wrap-reverse}.cyb-holder .cyb-suggests .cyb-pack{flex-basis:500px;display:flex;
flex:1;flex-wrap:wrap-reverse;min-width:250px}.cyb-holder .cyb-suggests.cyb-left .cyb-suggest{font-size:.9em;
flex-basis:250px;flex-direction:row}.cyb-holder .cyb-suggests.cyb-left .cyb-suggest .cyb-desc-holder,.cyb-holder
.cyb-suggests.cyb-left .cyb-suggest .cyb-img-holder{flex:1}.cyb-holder .cyb-suggests.cyb-left .cyb-suggest
.cyb-img-holder{flex-grow:2;padding-top:6px}.cyb-holder .cyb-suggests.cyb-left .cyb-suggest .cyb-desc-holder{
flex-grow:3;padding:0 5px}.cyb-holder .cyb-suggests.cyb-right .cyb-suggest{font-size:.9em;flex-basis:250px;
flex-direction:row-reverse}.cyb-holder .cyb-suggests.cyb-right .cyb-suggest .cyb-desc-holder,.cyb-holder
.cyb-suggests.cyb-right .cyb-suggest .cyb-img-holder{flex:1}.cyb-holder .cyb-suggests.cyb-right .cyb-suggest
.cyb-img-holder{flex-grow:2;padding-top:6px}.cyb-holder .cyb-suggests.cyb-right .cyb-suggest .cyb-desc-holder{
flex-grow:3;padding:0 5px}.cyb-holder .cyb-suggests.cyb-top .cyb-suggest{flex-basis:180px;flex-direction:column}
.cyb-holder .cyb-suggests.cyb-top .cyb-suggest .cyb-desc-holder,.cyb-holder .cyb-suggests.cyb-top .cyb-suggest
.cyb-img-holder{flex:1;margin:4px}.cyb-holder .cyb-suggests.cyb-top .cyb-suggest .cyb-desc-holder{margin-bottom:10px;
display:flex;font-size:1em}.cyb-holder .cyb-suggests.cyb-bottom .cyb-suggest{flex-basis:180px;flex-direction:column-reverse}
.cyb-holder .cyb-suggests.cyb-bottom .cyb-suggest .cyb-desc-holder,.cyb-holder .cyb-suggests.cyb-bottom .cyb-suggest
.cyb-img-holder{flex:1;margin:4px}.cyb-holder .cyb-suggests.cyb-bottom .cyb-suggest .cyb-desc-holder{margin-bottom:10px}
.cyb-holder .cyb-suggests.cyb-bottom .cyb-suggest .cyb-img-holder{display:flex}.cyb-holder .cyb-suggests .cyb-suggest{
display:flex;flex:1;margin-bottom:15px;position:relative;flex-basis:150px}.cyb-holder .cyb-suggests .cyb-suggest
.cyb-img-holder .cyb-img{width:100%;height:auto}.cyb-holder .cyb-suggests .cyb-suggest .cyb-img-holder .cyb-img.cyb-round{
border-radius:5px}.cyb-holder .cyb-suggests .cyb-more{border-radius:1em;border:.1em solid;border-color:rgba(128,128,128,.4);
padding:6px 13px;display:inline;font-size:.7em}.cyb-holder a{color:inherit;text-decoration:none;display:block}.cyb-vertical{
flex-basis:180px}.cyb-vertical .cyb-desc-holder,.cyb-vertical .cyb-img-holder{flex:1;margin:4px}.cyb-vertical .cyb-desc-holder{
margin-bottom:10px}.cyb-horizontal{font-size:.9em;flex-basis:250px}.cyb-horizontal .cyb-desc-holder,.cyb-horizontal
.cyb-img-holder{flex:1}.cyb-horizontal .cyb-img-holder{flex-grow:2;padding-top:6px}.cyb-horizontal .cyb-desc-holder{
flex-grow:3;padding:0 5px}
`
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
