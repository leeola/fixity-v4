package blev

import (
	"fmt"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
	"github.com/blevesearch/bleve/search/query"
	"github.com/leeola/errors"
	"github.com/leeola/kala/index"
)

func (b *Bleve) QueryOne(q index.Query, sb ...index.SortBy) (index.Result, error) {
	q.Limit = 1
	results, err := b.Query(q, sb...)
	if err != nil {
		return index.Result{}, err
	}

	var h index.Hash
	// technically Query() should have returned ErrNoQueryResults and been
	// returned above, so there should always be at least one hash. Nevertheless,
	// prevent a slice bounds panic.
	if len(results.Hashes) > 0 {
		h = results.Hashes[0]
	}

	return index.Result{
		IndexVersion: results.IndexVersion,
		Hash:         h,
	}, nil
}

func (b *Bleve) Query(q index.Query, ss ...index.SortBy) (index.Results, error) {
	var indexToUse bleve.Index
	if q.SearchVersions {
		indexToUse = b.entryIndex
	} else {
		indexToUse = b.anchorIndex
	}

	matches, err := queryIndex(indexToUse, q, ss...)
	if err != nil {
		return index.Results{}, errors.Stack(err)
	}

	hashes := make([]index.Hash, len(matches))

	for i, documentMatch := range matches {
		entryInterface, ok := documentMatch.Fields["indexEntry"]
		if !ok {
			return index.Results{}, errors.New("indexEntry value not found")
		}

		entry, ok := entryInterface.(float64)
		if !ok {
			return index.Results{}, errors.New("indexEntry value not int")
		}

		var h string
		if q.SearchVersions {
			h = documentMatch.ID
		} else {
			// if we're not matcing all, we need to get the stored hash value since the
			// document id will not equal the metadata, it will equal the anchor.
			f, ok := documentMatch.Fields["_metaHash"]
			if !ok {
				return index.Results{}, errors.New("_metaHash value not found")
			}

			v, ok := f.(string)
			if !ok {
				return index.Results{}, errors.New("_metaHash value not string")
			}
			h = v
		}

		hashes[i] = index.Hash{
			Entry: int(entry),
			Hash:  h,
		}
	}

	return index.Results{
		IndexVersion: b.indexVersion,
		Hashes:       hashes,
	}, nil
}

func queryIndex(i bleve.Index, q index.Query, ss ...index.SortBy) ([]*search.DocumentMatch, error) {
	var queries []query.Query

	if q.FromEntry != 0 {
		queries = append(queries, bleve.NewQueryStringQuery(fmt.Sprintf(
			"indexEntry:>=%d", q.FromEntry,
		)))
	}

	if q.Metadata != nil {
		for k, v := range q.Metadata {
			queries = append(queries, bleve.NewQueryStringQuery(fmt.Sprintf(
				"%s:%s", k, v,
			)))
		}
	}

	if len(queries) == 0 {
		return nil, nil
	}

	conjQuery := bleve.NewConjunctionQuery(queries...)
	search := bleve.NewSearchRequest(conjQuery)
	search.Size = q.Limit
	search.Fields = []string{"indexEntry", "_metaHash"}

	if len(ss) > 0 {
		sortBy := make([]string, len(ss))
		for i, s := range ss {
			if s.Descending {
				sortBy[i] = "-" + s.Field
			} else {
				sortBy[i] = s.Field
			}
		}
		search.SortBy(sortBy)
	}

	searchResults, err := i.Search(search)
	if err != nil {
		return nil, errors.Stack(err)
	}

	return searchResults.Hits, nil
}
