package inventory

import (
	"encoding/json"
	"html/template"

	"github.com/leeola/errors"
	"github.com/leeola/kala/client"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/util/jsonutil"
)

type templateMeta struct {
	Meta

	// The primary container that this inventory item is in.
	Container Meta

	// A slice of containers that this item is under. Eg, [box12, closet, home]
	Containers []Meta

	// Any children items that this inventory item contains.
	Children []Meta

	// If the item has no children, we display it's siblings, if any.
	Siblings []Meta
}

var displayTemplate = `
<div class="center aligned column segment">

<div class="ui raised very padded text container segment">
  <h2 class="ui header">{{.Meta.Name}}</h2>

	{{ if .Meta.Containers }}
		<div class="ui mini breadcrumb">
		{{ range $i, $element := .Meta.Containers }}
			{{ if $i }} <div class="divider"> / </div> {{ end }}
			<a href="/hash/{{$element.Anchor}}">{{$element.Name}}</a>
		{{ end }}
		</div>
	{{ end }}

	{{ if .Meta.Description }}
		<p>{{.Meta.Description}}</p>
	{{ end  }}

	{{ if .Meta.Children }}
		<table class="ui very compact table">
			<thead>
				<tr>
					<th colspan="2">Children</th>
				</tr>
			</thead>
			<tbody>
				{{ range $elm := .Meta.Children }}
					<tr>
						<td><a href="/hash/{{$elm.Anchor}}">{{$elm.Name}}</a></td>
						<td>{{$elm.Description}}</td>
					</tr>
				{{ end }}
			</tbody>
		</table>
	{{ end }}

	{{ if .Meta.Siblings }}
		<table class="ui very compact table">
			<thead>
				<tr>
					<th colspan="2">Siblings</th>
				</tr>
			</thead>
			<tbody>
				{{ range $elm := .Meta.Siblings }}
					<tr>
						<td><a href="/hash/{{$elm.Anchor}}">{{$elm.Name}}</a></td>
						<td>{{$elm.Description}}</td>
					</tr>
				{{ end }}
			</tbody>
		</table>
	{{ end }}
</div>

</div>
`

type Templater struct {
	client *client.Client
}

func NewTemplater(c *client.Client) *Templater {
	return &Templater{
		client: c,
	}
}

func (tr *Templater) Display(h string, b []byte, t *template.Template) (interface{}, error) {
	_, err := t.New("contentType").Parse(displayTemplate)
	if err != nil {
		return nil, errors.Stack(err)
	}

	var meta templateMeta
	if err := json.Unmarshal(b, &meta); err != nil {
		return nil, errors.Stack(err)
	}

	containers, err := tr.FetchContainers(meta.Meta.Container)
	if err != nil {
		return nil, err
	}
	reverseMetas(containers)
	meta.Containers = containers

	if len(meta.Containers) > 0 {
		meta.Container = containers[0]
	}

	children, err := tr.FetchChildren(h)
	if err != nil {
		return nil, err
	}
	meta.Children = children

	if len(children) == 0 && meta.Meta.Container != "" {
		siblings, err := tr.FetchSiblings(h, meta.Meta.Container)
		if err != nil {
			return nil, err
		}
		meta.Siblings = siblings
	}

	return meta, nil
}

func (tr *Templater) FetchContainers(startingContainer string) ([]Meta, error) {
	var (
		// containerHashes is used to check each container against and ensure we
		// are not ever fetching a container we already fetched - if that were to
		// happen, an endless loop would occur.
		containerHashes = map[string]struct{}{}

		containers    []Meta
		nextContainer = startingContainer
	)

	for nextContainer != "" {
		// always exit if hash exists, see containerHashes docstring.
		if _, ok := containerHashes[nextContainer]; ok {
			break
		}
		containerHashes[nextContainer] = struct{}{}

		var meta Meta
		if err := tr.getAndUnmarshal(nextContainer, &meta); err != nil {
			return nil, err
		}

		nextContainer = meta.Container

		containers = append(containers, meta)
	}

	return containers, nil
}

func (tr *Templater) FetchChildren(parent string) ([]Meta, error) {
	q := index.Query{
		Metadata: index.Metadata{
			"container": `"` + parent + `"`,
		},
	}
	results, err := tr.client.Query(q)
	if err != nil {
		return nil, err
	}

	children := make([]Meta, len(results.Hashes))
	for i, h := range results.Hashes {
		var m Meta
		if err := tr.getAndUnmarshal(h.Hash, &m); err != nil {
			return nil, err
		}
		children[i] = m
	}

	return children, nil
}

func (tr *Templater) FetchSiblings(currentChild, parent string) ([]Meta, error) {
	// Disabled until it can be updated to the new anchor spec/etc.
	// q := index.Query{
	// 	Metadata: index.Metadata{
	// 		"container": `"` + parent + `"`,
	// 	},
	// }
	// results, err := tr.client.Query(q)
	// if err != nil {
	// 	return nil, err
	// }

	var children []Meta
	// Disabled until it can be updated to the new anchor spec/etc.
	// for _, h := range results.Hashes {
	// 	var m Meta
	// 	if err := tr.getAndUnmarshal(h.Hash, &m); err != nil {
	// 		return nil, err
	// 	}
	// 	if m.Anchor == currentChild {
	// 		continue
	// 	}
	// 	children = append(children, m)
	// }

	return children, nil
}

func (tr *Templater) getAndUnmarshal(h string, m *Meta) error {
	rc, err := tr.client.GetDownloadBlob(h)
	if err != nil {
		return err
	}
	defer rc.Close()

	if err := jsonutil.UnmarshalReader(rc, m); err != nil {
		return err
	}

	return nil
}

func reverseMetas(m []Meta) {
	for i := 0; i < len(m)/2; i++ {
		j := len(m) - i - 1
		m[i], m[j] = m[j], m[i]
	}
}
