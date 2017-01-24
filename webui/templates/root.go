package templates

import "html/template"

var (
	Templates *template.Template

	header = `
<!doctype html>

<html lang="en">
<head>
  <meta charset="utf-8">

	<title>Kala</title>

	{{template "styles" .}}
</head>
<body>
`

	styles = `
<link rel="stylesheet" href="/public/semantic.min.css">
<link rel="stylesheet" href="/public/main.css">
`

	menu = `
<div class="ui menu">
  <a class="header item" href="/">
    Kala
  </a>
  <div class="item">
    Admin
  </div>
  <div class="item">
    Upload
  </div>
  <div class="item">
		Recent
  </div>
  <a class="item" href="/search">
    Search
  </a>
</div>
`

	footer = `
</body>
</html>
`

	root = `
{{template "header" .}}

<div id="root" class="ui one column middle aligned center aligned grid">
  <div class="eight wide column">
  	<div class="ui one column grid">

			<div id="search-header" class="column">
				<h2 class="ui column teal header">
					<div class="content">
						Search Your Kala Store
					</div>
				</h2>
			</div>

			<form class="ui column large form" action="/search">
				<div class="ui raised segment">
					<div class="field">
						<div class="ui left input">
							<input type="text" name="query" placeholder="Search..">
						</div>
					</div>
				</div>
			</form>

			<div class="column">
				<div class="ui center aligned grid">
					<div class="ui six wide column message">
						Need help? <a href="#">FAQ</a>
					</div>
				</div>
			</div>

		</div>
	</div>
</div>

{{template "footer" .}}
`
)

func init() {
	Templates = template.Must(template.New("root").Parse(root))
	template.Must(Templates.New("header").Parse(header))
	template.Must(Templates.New("menu").Parse(menu))
	template.Must(Templates.New("footer").Parse(footer))
	template.Must(Templates.New("styles").Parse(styles))
	template.Must(Templates.New("content").Parse(content))
	template.Must(Templates.New("search").Parse(search))
}
