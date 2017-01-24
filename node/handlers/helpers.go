package handlers

import (
	"net/http"

	"github.com/leeola/kala/contenttype"
	"github.com/leeola/kala/database"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/node/nodeware"
	"github.com/leeola/kala/store"
	"github.com/leeola/kala/util/jsonutil"
)

// GetContentStorerWithError helps to write an error if the ContentStorers cannot be found.
func GetContentStorersWithError(w http.ResponseWriter, r *http.Request) (map[string]contenttype.ContentType, bool) {
	s, ok := nodeware.GetContentStorers(r)
	if !ok {
		nodeware.GetLog(r).Error("contentStorer middleware instance missing")
		jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	return s, ok
}

// GetDatabaseWithError is a helper to write an error if the db cannot be found.
func GetDatabaseWithError(w http.ResponseWriter, r *http.Request) (database.Database, bool) {
	db, ok := nodeware.GetDatabase(r)
	if !ok {
		nodeware.GetLog(r).Error("db middleware instance missing")
		jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	return db, ok
}

// GetIndexWithError is a helper to write an error if the indexer cannot be found.
func GetIndexWithError(w http.ResponseWriter, r *http.Request) (index.Indexer, bool) {
	s, ok := nodeware.GetIndex(r)
	if !ok {
		nodeware.GetLog(r).Error("index middleware instance missing")
		jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	return s, ok
}

// GetStoreWithError is a helper to write an error if the store cannot be found.
func GetStoreWithError(w http.ResponseWriter, r *http.Request) (store.Store, bool) {
	s, ok := nodeware.GetStore(r)
	if !ok {
		nodeware.GetLog(r).Error("store middleware instance missing")
		jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	return s, ok
}

// GetQueryWithError is a helper to write an error if the queryer cannot be found.
func GetQueryWithError(w http.ResponseWriter, r *http.Request) (index.Queryer, bool) {
	s, ok := nodeware.GetQuery(r)
	if !ok {
		nodeware.GetLog(r).Error("query middleware instance missing")
		jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	return s, ok
}
