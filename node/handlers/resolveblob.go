package handlers

import (
	"io"
	"net/http"

	"github.com/leeola/kala/index"
	"github.com/leeola/kala/node/nodeware"
	"github.com/leeola/kala/store"
	"github.com/leeola/kala/util/jsonutil"
	"github.com/pressly/chi"
)

func GetResolveBlobHandler(s store.Store, q index.Queryer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		anchor := chi.URLParam(r, "anchor")
		log := nodeware.GetLog(r).New("anchor", anchor)

		hash, err := index.ResolveHashOrAnchor(s, q, anchor)
		if err != nil {
			log.Error("failed to resolve hash from anchor", "err", err)
			jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		log = log.New("hash", hash)

		rc, err := s.Read(hash)
		if err != nil {
			log.Error("failed to read meta hash from store", "err", err)
			jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		defer rc.Close()

		if _, err := io.Copy(w, rc); err != nil {
			log.Error("response copy failed", "err", err)
			jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
	}
}
