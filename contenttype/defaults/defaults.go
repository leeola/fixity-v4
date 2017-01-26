package defaults

import (
	ct "github.com/leeola/kala/contenttype"
	"github.com/leeola/kala/contenttype/data"
	"github.com/leeola/kala/contenttype/file"
	"github.com/leeola/kala/contenttype/image"
	"github.com/leeola/kala/contenttype/inventory"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
)

func DefaultTypes(s store.Store, i index.Indexer) (map[string]ct.ContentType, error) {
	m := map[string]ct.ContentType{}

	var cs ct.ContentType
	cs, err := data.New(data.Config{Store: s, Index: i})
	if err != nil {
		return nil, err
	}
	m[inventory.TypeKey] = cs

	cs, err = inventory.New(inventory.Config{Store: s, Index: i})
	if err != nil {
		return nil, err
	}
	m[inventory.TypeKey] = cs

	cs, err = file.New(file.Config{Store: s, Index: i})
	if err != nil {
		return nil, err
	}
	m[file.TypeKey] = cs

	// cs, err = video.New(video.Config{Store: s, Index: i})
	// if err != nil {
	// 	return nil, err
	// }
	// m["video"] = cs

	cs, err = image.New(image.Config{Store: s, Index: i})
	if err != nil {
		return nil, err
	}
	m[image.TypeKey] = cs

	// cs, err = folder.New(folder.Config{Store: s, Index: i})
	// if err != nil {
	// 	return nil, err
	// }
	// m["folder"] = cs

	return m, nil
}
