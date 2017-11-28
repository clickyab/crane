package banner

import (
	"net/http"

	"strings"

	"html/template"

	"fmt"

	"clickyab.com/crane/crane/builder"
	"clickyab.com/crane/crane/entity"
)

const bannerTemplateText = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta name="robots" content="nofollow">
	<meta content="text/html; charset=utf-8" http-equiv="Content-Type">
	<style>
	html, body, div, span, applet, object, iframe,h1, h2, h3, h4, h5, h6, p, blockquote, pre,a, abbr, acronym, address, big, cite, code,del, dfn, em, font, img, ins, kbd, q, s, samp,small, strike, strong, sub, sup, tt, var,dl, dt, dd, ol, ul, li,fieldset, form, label, legend,table, caption, tbody, tfoot, thead, tr, th, td { margin: 0; padding: 0; border: 0; outline: 0; font-weight: inherit; font-style: inherit; font-size: 100%; font-family: inherit; vertical-align: baseline;}
	:focus {outline: 0;}
	body {line-height: 1;color: black;background: white;}
	ol, ul {list-style: none;}
	table {border-collapse: separate;border-spacing: 0;}
	caption, th, td {text-align: left;font-weight: normal;}
	blockquote:before, blockquote:after,q:before, q:after {content: "";}
	blockquote, q {quotes: "" "";}
	body{ margin: 0; padding: 0; text-align: center; }
	.o{ position:absolute; top:0; left:0; border:0; height:250px; width:300px; z-index: 99; }
	#showb{ position:absolute; top:0; left:0; border:0; line-height: 250px; height:250px; width:300px; z-index: 100; background: rgba(0, 0, 0, 0.60); text-align: center; }
	{{ if .Tiny }}
	.tiny2{ height: 18px; width: 19px; position: absolute; bottom: 0px; right: 0; z-index: 100; background: url("//static.clickyab.com/img/clickyab-tiny.png") right top no-repeat; border-top-left-radius:4px; -moz-border-radius-topleft:4px  }
	.tiny2:hover{ width: 66px;  }
	.tiny{ height: 18px; width: 19px; position: absolute; top: 0px; right: 0; z-index: 100; background: url("//static.clickyab.com/img/clickyab-tiny.png") right top no-repeat; border-bottom-left-radius:4px; -moz-border-radius-bottomleft:4px  }
	.tiny:hover{ width: 66px;  }
	.tiny3{ position: absolute; top: 0px; right: 0; z-index: 100; }
	{{ end }}
	.butl {background: #4474CB;color: #FFF;padding: 10px;text-decoration: none;border: 2px solid #FFFFFF;font-family: tahoma;font-size: 13px;}
	img.adhere {max-width:100%;height:auto;}
	video {background: #232323 none repeat scroll 0 0;}
	</style>
</head>
<body>
    {{ if .Tiny }}<a class="tiny" href="http://clickyab.com/?ref=icon" target="_blank"></a>{{ end }}
    <a href="{{ .Link }}" target="_blank"><img  src="{{ .Src }}" border="0" height="{{ .Height }}" width="{{ .Width }}"/></a>
    <br style="clear: both;"/>
    {{ if .ShowT }}<br style="clear: both;"/><iframe src="//t.clickyab.com/" width="1" height="1" frameborder="0"></iframe>{{ end }}
</body>
</html>`

var bannerTemplate = template.Must(template.New("banner_template").Parse(bannerTemplateText))

// bannerData is the single ad id
type bannerData struct {
	Link   string
	Width  string
	Height string
	Src    string
	Tiny   bool
	ShowT  bool
}

func renderWebBanner(w http.ResponseWriter, ctx *builder.Context, slot entity.Slot, ad entity.Advertise) error {
	src := ad.Media()
	if ctx.GetCommon().Scheme == "https" {
		src = strings.Replace(src, "http://", "https://", -1)
	}
	sa := bannerData{
		Link:   slot.ClickURL(),
		Height: fmt.Sprint(slot.Height()),
		Width:  fmt.Sprint(slot.Width()),
		Src:    src,
		Tiny:   !ctx.GetCommon().NoTiny,
		ShowT:  ctx.ShowT(),
	}

	return bannerTemplate.Execute(w, sa)
}
