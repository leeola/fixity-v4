package handlers

import (
	"net/http"

	"github.com/leeola/kala/node/nodeware"
	"github.com/leeola/kala/webui/templates"
)

func GetRoot(w http.ResponseWriter, r *http.Request) {
	log := nodeware.GetLog(r)

	tmpl, err := templates.Templates.Clone()
	if err != nil {
		log.Error("failed to clone template", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "root", nil); err != nil {
		log.Error("failed to execute template", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}
