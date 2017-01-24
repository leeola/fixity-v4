package handlers

import (
	"fmt"
	"net/http"

	"github.com/leeola/kala/node/nodeware"
	"github.com/leeola/kala/util/jsonutil"
)

func GetNodeId(w http.ResponseWriter, r *http.Request) {
	log := nodeware.GetLog(r)

	db, ok := nodeware.GetDatabase(r)
	if !ok {
		log.Error("db middleware instance missing")
		jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	id, err := db.GetNodeId()
	if err != nil {
		log.Error("database GetNodeId failed", "err", err)
		jsonutil.Error(w, "database returned an error", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, id)
}
