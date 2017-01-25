package handlers

import (
	"net/http"

	"github.com/leeola/kala/contenttype"
	"github.com/leeola/kala/node/nodeware"
	"github.com/leeola/kala/util/jsonutil"
)

func PostUploadMetaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log := nodeware.GetLog(r)
	changes := contenttype.NewChangesFromValues(r.URL.Query())

	queryer, ok := GetQueryWithError(w, r)
	if !ok {
		return
	}

	s, ok := GetStoreWithError(w, r)
	if !ok {
		return
	}

	v, err := contenttype.ReadVersionFromChanges(s, queryer, changes)
	if err != nil {
		log.Error("failed to read version from changes", "err", err)
		jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	if v.ContentType == "" {
		// if no type has been defined even after loading anchors/etc, set it to default.
		v.ContentType = "data"
	}
	log = log.New("contentType", v.ContentType)

	css, ok := GetContentStorersWithError(w, r)
	if !ok {
		return
	}

	cs, ok := css[v.ContentType]
	if !ok {
		log.Info("requested contentType not found")
		jsonutil.Error(w, "requested contentType not found", http.StatusBadRequest)
		return
	}

	hashes, err := cs.StoreMeta(v, changes)
	if err != nil {
		log.Error("Meta returned error", "err", err)
		jsonutil.Error(w, "meta failed", http.StatusInternalServerError)
		return
	}

	_, err = jsonutil.MarshalToWriter(w, HashesResponse{
		Hashes: hashes,
	})
	if err != nil {
		log.Error("failed to marshal response", "err", err)
		jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}
