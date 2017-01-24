package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/leeola/kala/index"
	"github.com/leeola/kala/node/nodeware"
	"github.com/leeola/kala/util/jsonutil"
)

func GetQueryHandler(w http.ResponseWriter, r *http.Request) {
	log := nodeware.GetLog(r)

	q := index.Query{
		// default limit of 5
		Limit: 5,
	}
	sorts := []index.SortBy{}
	for k, v := range r.URL.Query() {
		if !strings.HasPrefix(k, "sort") && len(v) != 1 {
			jsonutil.Error(w, "duplicate query values not supported",
				http.StatusBadRequest)
			return
		}
		switch k {
		case "searchVersions":
			v, err := strconv.ParseBool(v[0])
			if err != nil {
				jsonutil.Error(w, "searchVersions must be a bool", http.StatusBadRequest)
				return
			}
			q.SearchVersions = v
		case "fromEntry":
			i, err := strconv.Atoi(v[0])
			if err != nil {
				jsonutil.Error(w, "fromEntry must be an integer", http.StatusBadRequest)
				return
			}
			q.FromEntry = i
		case "limit":
			i, err := strconv.Atoi(v[0])
			if err != nil {
				jsonutil.Error(w, "limit must be an integer", http.StatusBadRequest)
				return
			}
			q.Limit = i
		case "indexVersion":
			q.IndexVersion = v[0]
		case "sortAscending":
			for _, sort := range v {
				sorts = append(sorts, index.SortBy{Field: sort})
			}
		case "sortDescending":
			for _, field := range v {
				sorts = append(sorts, index.SortBy{
					Field:      field,
					Descending: true,
				})
			}
		default:
			if q.Metadata == nil {
				q.Metadata = map[string]interface{}{}
			}
			if v[0] != "" {
				q.Metadata[k] = v[0]
			}
		}
	}

	queryer, ok := GetQueryWithError(w, r)
	if !ok {
		return
	}

	result, err := queryer.Query(q, sorts...)
	switch err {
	case index.ErrIndexVersionsDoNotMatch:
		jsonutil.Error(w, "index Versions do not match", http.StatusBadRequest)
		return
	case nil:
		// do nothing here. we use this so that default: doesn't catch nil err
	default:
		log.Error("index.Query failed", "err", err)
		jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(result)
	if err != nil {
		log.Error("result marshalling failed", "err", err)
		jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	io.Copy(w, bytes.NewReader(b))
}
