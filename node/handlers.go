package node

import (
	"fmt"
	"io"
	"net/http"

	"github.com/leeola/kala/store"
	"github.com/pressly/chi"
)

func (n *Node) HeadContentHandler(w http.ResponseWriter, r *http.Request) {
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

func (n *Node) GetContentHandler(w http.ResponseWriter, r *http.Request) {
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

func (n *Node) PutContentHandler(w http.ResponseWriter, r *http.Request) {
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

func (n *Node) PostContentHandler(w http.ResponseWriter, r *http.Request) {
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
