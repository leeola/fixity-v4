package templates

import "html/template"

var (
	TmplErrTemplaterNotFound = `
<div class="ui center aligned three column grid">

<div class="ui negative message one column">
  <div class="header">No templater for content action</div>
	<p>{{.Meta.Action}} {{.Meta.Type}}</p>
</div>

</div>
`

	TmplErrCannotViewContent = `
<div class="center aligned column">
	Error: The requested content type cannot be viewed.
</div>
`
)

type NoContentTemplater struct {
	ContentType   string
	TemplaterType string
}

func (n NoContentTemplater) errorTmpler(t *template.Template, b []byte) (interface{}, error) {
	_, err := t.New("contentType").Parse(TmplErrTemplaterNotFound)
	tmplData := struct {
		Type   string
		Action string
	}{
		Type:   n.ContentType,
		Action: n.TemplaterType,
	}
	return tmplData, err
}

func (n NoContentTemplater) Display(h string, b []byte, t *template.Template) (interface{}, error) {
	return n.errorTmpler(t, b)
}

func (n NoContentTemplater) Form(h string, b []byte, t *template.Template) (interface{}, error) {
	return n.errorTmpler(t, b)
}
