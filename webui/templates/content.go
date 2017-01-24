package templates

var content = `
{{template "header" .}}
{{template "menu" .}}

<div class="ui one column grid container contenttype">
{{template "contentType" .}}
</div>

{{template "footer" .}}
`
