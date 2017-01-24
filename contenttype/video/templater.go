package video

import "html/template"

var cTypeTemplate = `
<div class="center aligned column">
	<video
		src="http://localhost:4001/download/{{.Hash}}"
		controls
	>
	Error: Your browser does not support the video element.
	</video>
</div>
`

func Templater(t *template.Template) (interface{}, error) {
	_, err := t.New("contentType").Parse(cTypeTemplate)
	return nil, err
}
