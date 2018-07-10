package bleve

import (
	"fmt"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/leeola/fixity"
	"github.com/leeola/fixity/index"
	"github.com/leeola/fixity/q"
	"github.com/leeola/fixity/q/operator"
)

const (
	fieldNameRef = index.FRefKey
	fieldNameID  = index.FIDKey
)

func (ix *Index) Query(qu q.Query) ([]fixity.Match, error) {
	var index bleve.Index
	if qu.IncludeVersions {
		index = ix.refIndex
	} else {
		index = ix.idIndex
	}

	return queryIndex(index, qu)
}

func queryIndex(ix bleve.Index, qu q.Query) ([]fixity.Match, error) {
	bq, err := fixQtoBleveQ(qu.Constraint)
	if err != nil {
		return nil, err // avoiding helper context to callers
	}

	search := bleve.NewSearchRequest(bq)
	search.Fields = []string{fieldNameID, fieldNameRef}

	searchResults, err := ix.Search(search)
	if err != nil {
		return nil, fmt.Errorf("search: %v", err)
	}

	matches := make([]fixity.Match, len(searchResults.Hits))

	for i, hit := range searchResults.Hits {
		refIfc, ok := hit.Fields[fieldNameRef]
		if !ok {
			return nil, fmt.Errorf("hit does not contain field: %s", fieldNameRef)
		}

		refStr, ok := refIfc.(string)
		if !ok {
			return nil, fmt.Errorf("hit field ref not valid string")
		}

		idIfc, ok := hit.Fields[fieldNameID]
		if !ok {
			return nil, fmt.Errorf("hit does not contain field: %s", fieldNameRef)
		}

		id, ok := idIfc.(string)
		if !ok {
			return nil, fmt.Errorf("hit field ref not valid string")
		}

		matches[i] = fixity.Match{
			ID:  id,
			Ref: fixity.Ref(refStr),
		}
	}

	return matches, nil
}

func fixQtoBleveQ(c q.Constraint) (query.Query, error) {
	switch c.Operator {
	case operator.Equal:
		if c.Value == nil {
			return nil, fmt.Errorf("field or value nil on equal op")
		}
		s, err := c.Value.ToString()
		if err != nil {
			return nil, fmt.Errorf("equal tostring: %v", err)
		}
		bq := bleve.NewTermQuery(s)
		// allow fieldless matches
		if c.Field != nil {
			bq.FieldVal = *c.Field
		}
		return bq, nil
	default:
		return nil, fmt.Errorf("unsupported constraint operator: %q", c.Operator)
	}
}
