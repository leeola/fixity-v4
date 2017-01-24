package image

import "html/template"

var cTypeTemplate = `
<div class="center aligned column">
	<h2>format: {{.Meta.Format}}</h2>
	<img src="http://localhost:4001/download/{{.Hash}}" />
</div>
`

func Templater(t *template.Template) (interface{}, error) {
	_, err := t.New("contentType").Parse(cTypeTemplate)
	return nil, err
}
