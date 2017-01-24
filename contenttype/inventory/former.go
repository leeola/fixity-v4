package inventory

import (
	"errors"
	"html/template"
)

var formTemplate = `
<div class="center aligned column segment">

<div class="ui labeled input">

	<form class="ui fluid action input">
		<div class="ui label">
			Name
		</div>
		<input name="name" type="text" placeholder="...">
		<div class="ui label">
			Description
		</div>
		<input name="description" type="text" placeholder="...">
		<div class="ui label">
			Container
		</div>
		<input name="container" type="text" placeholder="hash">
	</form>

</div>

</div>
`

func (tr *Templater) Form(h string, b []byte, t *template.Template) (interface{}, error) {
	return nil, errors.New("not implemented")
}
