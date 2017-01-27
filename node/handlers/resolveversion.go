package handlers

import (
	"net/http"

	"github.com/leeola/kala/index"
	"github.com/leeola/kala/node/nodeware"
	"github.com/leeola/kala/store"
	"github.com/leeola/kala/util/jsonutil"
	"github.com/pressly/chi"
)

func GetResolveVersionHandler(s store.Store, q index.Queryer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		resolve := chi.URLParam(r, "resolve")
		log := nodeware.GetLog(r).New("resolve", resolve)

		hash, err := index.ResolveHashOrAnchor(s, q, resolve)
		if err != nil {
			log.Error("failed to get hash from resolve", "err", err)
			jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		log = log.New("hash", hash)

		v, err := store.ReadVersion(s, hash)
		if err != nil {
			log.Error("failed to read blob hash from store", "err", err)
			jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		res := ResolveVersionResponse{
			Hash:    hash,
			Version: v,
		}

		if _, err = jsonutil.MarshalToWriter(w, res); err != nil {
			log.Error("failed to marshal response", "err", err)
			jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
	}
}
