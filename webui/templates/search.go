package templates

var search = `
{{template "header" .}}
{{template "menu" .}}

<div class="ui one column grid container">
	<div class="sixteen wide column">
		<form class="ui fluid action input">
			<input type="text" name="query" placeholder="Search.." value="{{.Query}}">
			<button class="ui button">Search</button>
		</form>
	</div>
	<div class="column">
		<table class="ui celled striped table">
			<thead>
				<tr><th colspan="5">
					Search Results
				</th>
			</tr></thead>
			<tbody>

				{{range $result := .Results}}
					<tr>
						<td>
							<a href="/hash/{{$result.Hash.Hash}}">{{$result.Name}}</a>
						</td>
						<td>
							{{range $index, $tag := $result.Tags}}
								{{if $index}},{{end}}
								<a href="/search?query=tags:&quot;{{$tag}}&quot;">{{$tag}}</a>
							{{end}}
						</td>
						<td class="right aligned collapsing">
							{{if $result.Anchor}}
								<a href="/search?query=anchor%3A{{$result.Anchor}}+searchVersions%3Atrue">
									{{$result.ShortAnchor}}
								</a>
							{{end}}
						</td>
						<td class="right aligned collapsing">
							{{$result.ContentType}}
						</td>
						<td class="right aligned collapsing">
							{{$result.HumanTime}}
						</td>
					</tr>
				{{end}}

			</tbody>
		</table>
	</div>
</div>

{{template "footer" .}}
`
