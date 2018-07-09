package bleve

import (
	"fmt"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/leeola/fixity"
	"github.com/leeola/fixity/q"
	"github.com/leeola/fixity/value/operator"
)

const (
	fieldNameRef = "fixityRef"
	fieldNameID  = "fixityID"
)

func (ix *Index) Query(qu q.Query) ([]fixity.Ref, error) {
	// for prototype, use only id index
	return queryIndex(ix.idIndex, qu)
}

func queryIndex(ix bleve.Index, qu q.Query) ([]fixity.Ref, error) {

	bq, err := fixQtoBleveQ(qu.Constraint)
	if err != nil {
		return nil, err // avoiding helper context to callers
	}

	search := bleve.NewSearchRequest(bq)
	search.Fields = []string{fieldNameRef}

	searchResults, err := ix.Search(search)
	if err != nil {
		return nil, fmt.Errorf("search: %v", err)
	}

	refs := make([]fixity.Ref, len(searchResults.Hits))

	for i, hit := range searchResults.Hits {
		fv, ok := hit.Fields[fieldNameRef]
		if !ok {
			return nil, fmt.Errorf("hit does not contain field: %s", fieldNameRef)
		}

		s, ok := fv.(string)
		if !ok {
			return nil, fmt.Errorf("hit field ref not valid string")
		}

		refs[i] = fixity.Ref(s)
	}

	return refs, nil
}

func fixQtoBleveQ(c q.Constraint) (query.Query, error) {
	switch c.Operator {
	case operator.Equal:
		if c.Field == nil || c.Value == nil {
			return nil, fmt.Errorf("field or value nil on equal op")
		}
		s, err := c.Value.ToString()
		if err != nil {
			return nil, fmt.Errorf("equal tostring: %v", err)
		}
		bq := bleve.NewTermQuery(s)
		bq.FieldVal = *c.Field
		return bq, nil
	default:
		return nil, fmt.Errorf("unsupported constraint operator: %q", c.Operator)
	}
}
