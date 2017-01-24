package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/leeola/kala/contenttype"
	"github.com/leeola/kala/node/nodeware"
	"github.com/leeola/kala/store"
	"github.com/leeola/kala/webui/templates"
	"github.com/leeola/kala/webui/webware"
	"github.com/pressly/chi"
)

type GetHashContent struct {
	Hash string
	Meta interface{}
}

func GetHash(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	log := nodeware.GetLog(r).New("hash", hash)

	tmpl, err := templates.Templates.Clone()
	if err != nil {
		log.Error("failed to clone template", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	nodeClient, ok := webware.GetClient(r)
	if !ok {
		log.Error("node client missing from Context")
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	rc, err := nodeClient.GetDownloadBlob(hash)
	if err != nil {
		log.Error("failed to get blob content type", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	defer rc.Close()

	metaB, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Error("failed to read hash data", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	cType, err := store.UnmarshalContentType(metaB)
	if err != nil {
		log.Error("failed to get hash content type", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	var cTemplater contenttype.ContentDisplayer
	if cType != "" {
		t, _ := webware.GetContentTemplater(r, cType)
		cTemplater, _ = t.(contenttype.ContentDisplayer)
	}
	// If the templater still isn't set, set it to the default.
	if cTemplater == nil {
		cTemplater = templates.NoContentTemplater{
			ContentType:   cType,
			TemplaterType: "display",
		}
	}

	meta, err := cTemplater.Display(hash, metaB, tmpl)
	if err != nil {
		log.Error("content templater failed", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	tmplData := GetHashContent{
		Hash: hash,
		Meta: meta,
	}

	if err := tmpl.ExecuteTemplate(w, "content", tmplData); err != nil {
		log.Error("failed to execute template", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}
