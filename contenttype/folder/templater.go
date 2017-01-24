package folder

import (
	"html/template"

	"github.com/leeola/kala/webui/templates"
)

func Templater(t *template.Template) error {
	_, err := t.New("contentType").Parse(templates.TmplErrCannotViewContent)
	return err
}
