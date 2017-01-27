package handlers

import (
	"net/http"

	"github.com/leeola/kala/index"
	"github.com/leeola/kala/node/nodeware"
	"github.com/leeola/kala/store"
	"github.com/leeola/kala/util/strutil"
	"github.com/leeola/kala/webui/templates"
	"github.com/leeola/kala/webui/webware"
)

type SearchPage struct {
	Query   string
	Results []MetaResult
}

type MetaResult struct {
	index.Hash
	store.Version
	store.Meta

	ShortAnchor string
	HumanTime   string
}

func GetSearch(w http.ResponseWriter, r *http.Request) {
	log := nodeware.GetLog(r)

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

	q := index.Query{
		Limit:    5,
		Metadata: index.Metadata{},
	}

	qString := r.URL.Query().Get("query")
	if qString != "" {
		for _, field := range strutil.QuotedFields(qString) {
			k, v := strutil.SplitQueryField(field)
			if k != "" {
				// if the key was specified, set the requested metadata key with the val
				q.Metadata[k] = v
			} else if _, ok := q.Metadata["name"]; !ok {
				// if the key was not specified, default the value to filename
				// only if filename was not already set.
				//
				// Since this is implicit logic, we should not overwrite the users explicit
				// filename value.
				q.Metadata["name"] = v
			}
		}
	}

	results, err := nodeClient.Query(q)
	if err != nil {
		log.Error("node query failed", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	metaResults := make([]MetaResult, len(results.Hashes))
	for i, hash := range results.Hashes {
		mr := MetaResult{Hash: hash}
		if err := nodeClient.GetBlobAndUnmarshal(hash.Hash, &mr.Version); err != nil {
			log.Error("failed to get blob of hash", "err", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		if err := nodeClient.GetBlobAndUnmarshal(mr.Version.Meta, &mr.Meta); err != nil {
			log.Error("failed to get blob of hash", "err", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		if !mr.UploadedAt.IsZero() {
			mr.HumanTime = mr.UploadedAt.Format("Mon Jan _2 03:04PM")
		}
		mr.ShortAnchor = strutil.ShortHash(mr.Anchor, 8)
		metaResults[i] = mr
	}

	page := SearchPage{
		Results: metaResults,
		Query:   qString,
	}
	if err := tmpl.ExecuteTemplate(w, "search", page); err != nil {
		log.Error("failed to execute template", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}
