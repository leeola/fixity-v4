package inventory

import (
	"html/template"

	"github.com/leeola/errors"
	"github.com/leeola/kala/client"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
)

type verMeta struct {
	Meta
	VersionHash string
	Version     store.Version
}

type templateMeta struct {
	Meta

	// The primary container that this inventory item is in.
	Container verMeta

	// A slice of containers that this item is under. Eg, [box12, closet, home]
	Containers []verMeta

	// Any children items that this inventory item contains.
	Children []verMeta

	// If the item has no children, we display it's siblings, if any.
	Siblings []verMeta
}

var displayTemplate = `
<div class="center aligned column segment">

<div class="ui raised very padded text container segment">
  <h2 class="ui header">{{.Meta.Name}}</h2>

	{{ if .Meta.Containers }}
		<div class="ui mini breadcrumb">
		{{ range $i, $element := .Meta.Containers }}
			{{ if $i }} <div class="divider"> / </div> {{ end }}
			<a href="/hash/{{$element.VersionHash}}">{{$element.Name}}</a>
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
						<td><a href="/hash/{{$elm.VersionHash}}">{{$elm.Name}}</a></td>
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
						<td><a href="/hash/{{$elm.VersionHash}}">{{$elm.Name}}</a></td>
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

func (tr *Templater) Display(h string, v store.Version, t *template.Template) (interface{}, error) {
	_, err := t.New("contentType").Parse(displayTemplate)
	if err != nil {
		return nil, errors.Stack(err)
	}

	var meta templateMeta
	if err := tr.client.GetBlobAndUnmarshal(v.Meta, &meta); err != nil {
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

	children, err := tr.FetchChildren(v.Anchor)
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

func (tr *Templater) FetchContainers(startingContainer string) ([]verMeta, error) {
	var (
		// containerHashes is used to check each container against and ensure we
		// are not ever fetching a container we already fetched - if that were to
		// happen, an endless loop would occur.
		containerHashes = map[string]struct{}{}

		containers    []verMeta
		nextContainer = startingContainer
	)

	for nextContainer != "" {
		// always exit if hash exists, see containerHashes docstring.
		if _, ok := containerHashes[nextContainer]; ok {
			break
		}
		containerHashes[nextContainer] = struct{}{}

		var vm verMeta
		verHash, err := tr.client.GetResolveBlobAndUnmarshal(nextContainer, &vm.Version)
		if err != nil {
			return nil, err
		}
		vm.VersionHash = verHash

		if err := tr.client.GetBlobAndUnmarshal(vm.Version.Meta, &vm.Meta); err != nil {
			return nil, err
		}

		nextContainer = vm.Container

		containers = append(containers, vm)
	}

	return containers, nil
}

func (tr *Templater) FetchChildren(parent string) ([]verMeta, error) {
	q := index.Query{
		Metadata: index.Metadata{
			"container": `"` + parent + `"`,
		},
	}
	results, err := tr.client.Query(q)
	if err != nil {
		return nil, err
	}

	children := make([]verMeta, len(results.Hashes))
	for i, h := range results.Hashes {
		vm := verMeta{VersionHash: h.Hash}
		if err := tr.client.GetBlobAndUnmarshal(h.Hash, &vm.Version); err != nil {
			return nil, err
		}
		if err := tr.client.GetBlobAndUnmarshal(vm.Version.Meta, &vm.Meta); err != nil {
			return nil, err
		}
		children[i] = vm
	}

	return children, nil
}

func (tr *Templater) FetchSiblings(currentChild, parent string) ([]verMeta, error) {
	q := index.Query{
		Metadata: index.Metadata{
			"container": `"` + parent + `"`,
		},
	}
	results, err := tr.client.Query(q)
	if err != nil {
		return nil, err
	}

	var children []verMeta
	for _, h := range results.Hashes {
		if h.Hash == currentChild {
			continue
		}

		vm := verMeta{VersionHash: h.Hash}
		if err := tr.client.GetBlobAndUnmarshal(h.Hash, &vm.Version); err != nil {
			return nil, err
		}

		if err := tr.client.GetBlobAndUnmarshal(vm.Version.Meta, &vm.Meta); err != nil {
			return nil, err
		}

		children = append(children, vm)
	}

	return children, nil
}

func reverseMetas(m []verMeta) {
	for i := 0; i < len(m)/2; i++ {
		j := len(m) - i - 1
		m[i], m[j] = m[j], m[i]
	}
}
