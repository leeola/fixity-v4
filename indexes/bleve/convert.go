package bleve

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/fatih/structs"
	"github.com/leeola/errors"
	kq "github.com/leeola/kala/q"
	ops "github.com/leeola/kala/q/operators"
)

func ConvertQuery(kq *kq.Query) (*bleve.SearchRequest, error) {
	var bq query.Query
	if !structs.IsZero(kq.Constraint) {
		constQuery, err := ConvertConstraint(kq.Constraint)
		if err != nil {
			return nil, err
		}
		bq = constQuery
	} else {
		// If no constraint is specified, match all documents.
		bq = bleve.NewMatchAllQuery()
	}

	if kq.SkipBy != 0 {
		return nil, errors.New("SkipBy not implemented for bleve index")
	}

	search := bleve.NewSearchRequest(bq)
	search.Size = kq.LimitBy

	if len(kq.SortBy) > 0 {
		sortBy := make([]string, len(kq.SortBy))
		for i, sb := range kq.SortBy {
			if sb.Descending {
				sortBy[i] = "-" + sb.Field
			} else {
				sortBy[i] = sb.Field
			}
		}
		search.SortBy(sortBy)
	}

	return search, nil
}

func ConvertConstraint(c kq.Constraint) (query.Query, error) {
	switch c.Operator {
	case ops.Equal:
		// TODO(leeola): implement bool and int casts.
		if v, ok := c.Value.(string); ok {
			return bleve.NewTermQuery(v), nil
		} else {
			return nil, errors.Errorf(
				"%q operator not supported for value: %s", c.Operator, c.Value)
		}

	case ops.And, ops.Or, ops.Not:
		qs := make([]query.Query, len(c.Constraints))
		for i, subc := range c.Constraints {
			q, err := ConvertConstraint(subc)
			if err != nil {
				return nil, err
			}
			qs[i] = q
		}

		boolQuery := bleve.NewBooleanQuery()
		switch c.Operator {
		case ops.And:
			boolQuery.AddMust(qs...)
		case ops.Or:
			// Not sure how to implement Or with Bleve at the moment, failing here.
			//
			// I'm fairly sure what is needed here is a Disjunction query though.
			// Eg: https://github.com/blevesearch/bleve/blob/master/query.go#L69
			//
			// A list of ORs simply satisfy one or more of the queries.
			return nil, errors.New("Or operator not yet supported")
		case ops.Not:
			boolQuery.AddMustNot(qs...)
		default:
			// should never happen
			return nil, errors.Errorf(
				"unexpected operator inside And/Or/Not case: %s", c.Operator)
		}
		return boolQuery, nil

	default:
		return nil, errors.Errorf("unsupported operator: %s", c.Operator)
	}
}
