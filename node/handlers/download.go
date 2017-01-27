package handlers

import (
	"io"
	"net/http"

	"github.com/leeola/kala/index/indexreader"
	"github.com/leeola/kala/node/nodeware"
	"github.com/leeola/kala/util/jsonutil"
	"github.com/pressly/chi"
)

func GetDownloadHandler(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "anchor")
	log := nodeware.GetLog(r).New("anchor", hash)

	s, ok := GetStoreWithError(w, r)
	if !ok {
		return
	}

	q, ok := GetQueryWithError(w, r)
	if !ok {
		return
	}

	reader, err := indexreader.New(indexreader.Config{
		HashOrAnchor: hash,
		Store:        s,
		Query:        q,
	})
	if err != nil {
		log.Error("failed to marshal response", "err", err)
		jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(w, reader); err != nil {
		log.Error("response copy failed", "err", err)
		jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}
