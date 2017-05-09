package renderer

import (
	"bytes"
	"html/template"

	"services/assert"
)

type templateContext struct {
	Landing  string
	IsFilled bool
	URL      string
	Pixel    string
}

const tpl = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{ .Landing }</title>
</head>
<body>
<iframe>
    {{ if .IsFilled }}<img src="{{ .Pixel }}" alt="">{{ end }}
    <iframe id="thirdad_frame" src="{{ .URL }}" class="thirdad{{ if .IsFilled }}thrdadok{{ else }}thirdadempty{{ end }}"></iframe>
</iframe>
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
