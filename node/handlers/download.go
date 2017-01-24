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
	hash := chi.URLParam(r, "hash")
	log := nodeware.GetLog(r).New("hash", hash)

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

// func GetDownloadBlobHandler(s store.Store, q index.Queryer) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")
//
// 		hash := chi.URLParam(r, "hash")
// 		log := nodeware.GetLog(r).New("hash", hash)
//
// 		isAnchor, hashB, err := store.IsVersionWithBytes(s, hash)
// 		if err != nil {
// 			log.Error("failed to check if hash is anchor", "err", err)
// 			jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
// 				http.StatusInternalServerError)
// 			return
// 		}
//
// 		if isAnchor {
// 			result, err := q.QueryOne(index.Query{
// 				// NOTE: Putting the hash in quotes because the querystring in bleve
// 				// has issues with a hyphenated hashstring. This is annoying, and
// 				// should be fixed somehow...
// 				Metadata: index.Metadata{"anchor": `"` + hash + `"`},
// 			}, index.SortBy{Field: "uploadedAt", Descending: true})
// 			if err != nil {
// 				log.Error("failed to query metadata", "err", err)
// 				jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
// 					http.StatusInternalServerError)
// 				return
// 			}
//
// 			if result.Hash.Hash == "" {
// 				log.Error("no meta found with anchor", "anchor", hash)
// 				jsonutil.Error(w, "no meta found with anchor", http.StatusNotFound)
// 				return
// 			}
//
// 			rc, err := s.Read(result.Hash.Hash)
// 			if err != nil {
// 				log.Error("failed to read meta hash from store", "err", err)
// 				jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
// 					http.StatusInternalServerError)
// 				return
// 			}
// 			defer rc.Close()
//
// 			b, err := ioutil.ReadAll(rc)
// 			if err != nil {
// 				log.Error("failed to read meta hash from store reader", "err", err)
// 				jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
// 					http.StatusInternalServerError)
// 				return
// 			}
// 			hashB = b
// 		}
//
// 		if _, err := io.Copy(w, bytes.NewReader(hashB)); err != nil {
// 			log.Error("response copy failed", "err", err)
// 			jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
// 				http.StatusInternalServerError)
// 			return
// 		}
// 	}
// }
