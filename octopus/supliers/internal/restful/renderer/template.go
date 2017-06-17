package renderer

import (
	"bytes"
	"html/template"

	"github.com/clickyab/services/assert"
)

type templateContext struct {
	Landing       string
	IsFilled      bool
	URL           string
	Pixel         string
	Width, Height int
}

const tpl = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .Landing }}</title>
    <style>
    #adhere iframe {max-width:100%;margin: 0 auto;}
    .pixel {position: absolute; top : -1000px; left : -1000px}
    </style>
</head>
<body style="margin: 0; padding: 0;">
    {{ if .IsFilled }}<img class="pixel" src="{{ .Pixel }}" alt="">{{ end }}
    <div id="adhere"><iframe id="thirdad_frame" width="{{ .Width }}" height="{{ .Height }}" frameborder=0 marginwidth="0" marginheight="0" vspace="0" hspace="0" allowtransparency="true" scrolling="no" src="{{ .URL }}" class="thirdad{{ if .IsFilled }} thrdadok{{ else }} thirdadempty{{ end }}"></iframe></div>
</body>
</html>`

var (
	restTemplate = template.Must(template.New("rest_tpl").Parse(tpl))
)

func renderTemplate(ctx templateContext) string {
	buf := &bytes.Buffer{}
	assert.Nil(restTemplate.Execute(buf, ctx))
	return buf.String()
}
