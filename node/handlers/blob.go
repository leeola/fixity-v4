package handlers

import (
	"io"
	"net/http"

	"github.com/leeola/kala/node/nodeware"
	"github.com/leeola/kala/store"
	"github.com/leeola/kala/util/jsonutil"
	"github.com/pressly/chi"
)

func GetBlobHandler(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	log := nodeware.GetLog(r).New("hash", hash)

	s, ok := GetStoreWithError(w, r)
	if !ok {
		return
	}

	rc, err := s.Read(hash)
	if err == store.HashNotFoundErr {
		jsonutil.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Error("store.Read failed", "err", err)
		jsonutil.Error(w, "store Read failed", http.StatusInternalServerError)
		return
	}
	defer rc.Close()

	if _, err := io.Copy(w, rc); err != nil {
		log.Error("response write failed", "err", err)
		jsonutil.Error(w, "response write failed", http.StatusInternalServerError)
		return
	}
}

func HeadBlobHandler(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	log := nodeware.GetLog(r).New("hash", hash)

	s, ok := GetStoreWithError(w, r)
	if !ok {
		return
	}

	exists, err := s.Exists(hash)
	if err != nil {
		log.Error("store.Exists failed", "err", err)
		jsonutil.Error(w, "store Exists failed", http.StatusInternalServerError)
		return
	}

	// If it does not exist, return 404.
	if !exists {
		jsonutil.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// return 200 if it exists.
}

func PutBlobHandler(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	log := nodeware.GetLog(r).New("hash", hash)

	s, ok := GetStoreWithError(w, r)
	if !ok {
		return
	}

	err := store.WriteHashReader(s, hash, r.Body)
	if err == store.HashNotMatchContentErr {
		log.Error("write of nonmatching content for hash attempted")
		jsonutil.Error(w, "content does not match hash", http.StatusForbidden)
		return
	}
	if err != nil {
		log.Error("store write failed", "err", err)
		jsonutil.Error(w, "store write failed", http.StatusInternalServerError)
		return
	}
}

func PostBlobHandler(w http.ResponseWriter, r *http.Request) {
	log := nodeware.GetLog(r)

	s, ok := GetStoreWithError(w, r)
	if !ok {
		return
	}

	h, err := store.WriteReader(s, r.Body)
	if err != nil {
		log.Error("store write failed", "err", err)
		jsonutil.Error(w, "store write failed", http.StatusInternalServerError)
		return
	}

	log.Debug("POSTed content", "hash", h)

	_, err = jsonutil.MarshalToWriter(w, HashResponse{
		Hash: h,
	})
	if err != nil {
		log.Error("store write failed", "err", err)
		jsonutil.Error(w, "store write failed", http.StatusInternalServerError)
		return
	}
}
