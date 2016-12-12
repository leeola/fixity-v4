package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
	"github.com/leeola/kala/util/jsonutil"
	"github.com/leeola/kala/util/urlutil"
	"github.com/pressly/chi"
)

func (n *Node) GetNodeId(w http.ResponseWriter, r *http.Request) {
	log := GetLog(r)

	id, err := n.db.GetNodeId()
	if err != nil {
		log.Error("database GetNodeId failed", "err", err)
		http.Error(w, "database returned an error", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, id)
}

func (n *Node) HeadBlobHandler(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	log := GetLog(r).New("hash", hash)

	exists, err := n.store.Exists(hash)
	if err != nil {
		log.Error("store.Exists failed", "err", err)
		http.Error(w, "store Exists failed", http.StatusInternalServerError)
		return
	}

	// If it does not exist, return 404.
	if !exists {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	// return 200 if it exists.
}

func (n *Node) GetBlobHandler(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	log := GetLog(r).New("hash", hash)

	rc, err := n.store.Read(hash)
	if err == store.HashNotFoundErr {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Error("store.Read failed", "err", err)
		http.Error(w, "store Read failed", http.StatusInternalServerError)
		return
	}
	defer rc.Close()

	if _, err := io.Copy(w, rc); err != nil {
		log.Error("response write failed", "err", err)
		http.Error(w, "response write failed", http.StatusInternalServerError)
		return
	}
}

func (n *Node) PutBlobHandler(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	log := GetLog(r).New("hash", hash)

	err := store.WriteHashReader(n.store, hash, r.Body)
	if err == store.HashNotMatchContentErr {
		log.Error("write of nonmatching content for hash attempted")
		http.Error(w, "content does not match hash", http.StatusForbidden)
		return
	}
	if err != nil {
		log.Error("store write failed", "err", err)
		http.Error(w, "store write failed", http.StatusInternalServerError)
		return
	}
}

func (n *Node) PostBlobHandler(w http.ResponseWriter, r *http.Request) {
	log := GetLog(r)
	h, err := store.WriteReader(n.store, r.Body)
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

func (n *Node) GetQueryHandler(w http.ResponseWriter, r *http.Request) {
	log := GetLog(r)

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
		case "fromEntry":
			i, err := strconv.Atoi(v[0])
			if err != nil {
				jsonutil.Error(w, "fromEntry must be integer", http.StatusBadRequest)
				return
			}
			q.FromEntry = i
		case "limit":
			i, err := strconv.Atoi(v[0])
			if err != nil {
				jsonutil.Error(w, "limit must be integer", http.StatusBadRequest)
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
			q.Metadata[k] = v[0]
		}
	}

	result, err := n.index.Query(q, sorts)
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

func (n *Node) GetIndexContentHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (n *Node) PostUploadHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "handler")
	log := GetLog(r).New("handler", key)

	metaChanges := store.NewMetaChangesFromValues(r.URL.Query())

	newAnchor := urlutil.GetQueryBool(r, "newAnchor")
	anchorHash := urlutil.GetQueryString(r, "anchor")
	if newAnchor && anchorHash != "" {
		log.Info("newAnchor and anchor query fields both defined, rejecting request")
		jsonutil.Error(w, "query params newAnchor and anchor must not be used together",
			http.StatusBadRequest)
	}

	uploadHandler, ok := n.upload[key]
	if !ok {
		log.Info("requested upload handler not found")
		http.Error(w, "requested upload handler not found",
			http.StatusBadRequest)
		return
	}

	hashes, err := uploadHandler.Upload(r.Body, metaChanges)
	if err != nil {
		log.Error("uplad handler returned error", "err", err)
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

func (n *Node) GetDownloadHandler(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	log := GetLog(r).New("hash", hash)

	reader, err := store.NewReader(store.ReaderConfig{
		Hash:  hash,
		Store: n.store,
	})
	if err != nil {
		log.Error("failed to marshal response", "err", err)
		jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(w, reader); err != nil {
		log.Error("response write failed", "err", err)
		jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}
