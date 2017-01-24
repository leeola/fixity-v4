package contenttype

import "html/template"

type ContentDisplayer interface {
	Display(string, []byte, *template.Template) (interface{}, error)
}

type ContentFormer interface {
	Form(string, []byte, *template.Template) (interface{}, error)
}
