package contenttype

import (
	"html/template"

	"github.com/leeola/kala/store"
)

type ContentDisplayer interface {
	Display(string, store.Version, *template.Template) (interface{}, error)
}

type ContentFormer interface {
	Form(string, store.Version, *template.Template) (interface{}, error)
}
