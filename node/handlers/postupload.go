package handlers

import (
	"net/http"

	"github.com/leeola/kala/contenttype"
	"github.com/leeola/kala/node/nodeware"
	"github.com/leeola/kala/store"
	"github.com/leeola/kala/util/jsonutil"
	"github.com/leeola/kala/util/urlutil"
)

func PostUploadHandler(w http.ResponseWriter, r *http.Request) {
	log := nodeware.GetLog(r)
	changes := contenttype.NewChangesFromValues(r.URL.Query())

	// Temporarily disabled
	// anchor := urlutil.GetQueryString(r, "anchor")
	// previousVersion := urlutil.GetQueryString(r, "previousVersion")
	// // If there is no previous version to base this mutation off of, then query the
	// // indexer for the most recent hash for this anchor.
	// if previousVersion == "" && anchorHash != "" {
	// 	q := index.Query{
	// 		Metadata: index.Metadata{
	// 			// NOTE: Putting the hash in quotes because the querystring in bleve
	// 			// has issues with a hyphenated hashstring. This is annoying, and
	// 			// should be fixed somehow...
	// 			"anchor": `"` + anchorHash + `"`,
	// 		},
	// 	}
	// 	s := index.SortBy{
	// 		Field:      "uploadedAt",
	// 		Descending: true,
	// 	}

	queryer, ok := GetQueryWithError(w, r)
	if !ok {
		return
	}

	// 	result, err := queryer.QueryOne(q, s)
	// 	if err != nil {
	// 		log.Error("failed to query for previous meta hash", "err", err)
	// 		jsonutil.Error(w, "previous meta query failed", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	if result.Hash.Hash != "" {
	// 		previousMeta = result.Hash.Hash
	// 		changes.SetPreviousMeta(previousMeta)
	// 	}
	// }

	s, ok := GetStoreWithError(w, r)
	if !ok {
		return
	}

	v, err := contenttype.VersionFromChanges(s, queryer, changes)
	if v.ContentType == "" {
		// if no type has been defined even after loading anchors/etc, set it to default.
		v.ContentType = "data"
	}
	log = log.New("contentType", v.ContentType)

	// write a new anchor if specified
	if urlutil.GetQueryBool(r, "newAnchor") {
		a, err := store.NewAnchor()
		if err != nil {
			log.Error("failed to create new anchor", "err", err)
			jsonutil.Error(w, "newanchor failed", http.StatusInternalServerError)
			return
		}

		v.Anchor = a
	}

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

	hashes, err := cs.StoreContent(r.Body, v, changes)
	if err != nil {
		log.Error("StoreContent returned error", "err", err)
		jsonutil.Error(w, "upload failed", http.StatusInternalServerError)
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
