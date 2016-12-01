package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
	"github.com/pressly/chi"
)

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
		http.Error(w, "store write failed", http.StatusInternalServerError)
		return
	}

	log.Debug("POSTed content", "hash", h)
	fmt.Fprint(w, h)
}

func (n *Node) GetQueryHandler(w http.ResponseWriter, r *http.Request) {
	log := GetLog(r)

	q := index.Query{
		// default limit of 5
		Limit: 5,
	}
	for k, v := range r.URL.Query() {
		switch k {
		case "fromEntry":
			i, err := strconv.Atoi(v[0])
			if err != nil {
				http.Error(w, "fromEntry must be integer", http.StatusInternalServerError)
				return
			}
			q.FromEntry = i
		case "limit":
			i, err := strconv.Atoi(v[0])
			if err != nil {
				http.Error(w, "limit must be integer", http.StatusInternalServerError)
				return
			}
			q.Limit = i
		case "indexVersion":
			q.IndexVersion = v[0]
		default:
			log.Error("unhandled query param", "key", k, "value", v)
			http.Error(w,
				fmt.Sprintf("invalid query param %q", k),
				http.StatusInternalServerError)
			return
		}
	}

	result, err := n.index.Query(q)
	switch err {
	case index.ErrNoQueryResults:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	case index.ErrIndexVersionsDoNotMatch:
		http.Error(w, "index Versions do not match", http.StatusInternalServerError)
		return
	case nil:
		// do nothing here. we use this so that default: doesn't catch nil err
	default:
		log.Error("index.Query failed", "err", err)
		http.Error(w, "index Query failed", http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(result)
	if err != nil {
		log.Error("index.Query failed", "err", err)
		http.Error(w, "index Query failed", http.StatusInternalServerError)
		return
	}

	io.Copy(w, bytes.NewReader(b))
}

func (n *Node) GetIndexContentHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
