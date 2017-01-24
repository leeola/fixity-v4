package file

import (
	"html/template"

	"github.com/leeola/kala/webui/templates"
)

func Templater(t *template.Template) (interface{}, error) {
	_, err := t.New("contentType").Parse(templates.TmplErrCannotViewContent)
	return nil, err
}
