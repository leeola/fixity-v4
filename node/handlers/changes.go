package handlers

// func GetMetaExportHandler(w http.ResponseWriter, r *http.Request) {
// 	hash := chi.URLParam(r, "hash")
// 	log := nodeware.GetLog(r).New("hash", hash)
//
// 	s, ok := GetStoreWithError(w, r)
// 	if !ok {
// 		return
// 	}
//
// 	// TODO(leeola): Add an index helper func to fetch the meta hash from a given hash,
// 	// which could be anchor or meta. This will allow the hash to be an anchor or
// 	// a meta, as the /download endpoint should be adaptable as possible. Eg,
// 	// you should be able to get the meta from an anchor hash or meta hash.
//
// 	cType, mb, err := store.GetContentTypeWithBytes(s, hash)
// 	if err != nil {
// 		log.Error("failed to get previous content type", "err", err)
// 		jsonutil.Error(w, "contenttype lookup failed", http.StatusInternalServerError)
// 		return
// 	}
//
// 	cts, ok := GetContentStorersWithError(w, r)
// 	if !ok {
// 		return
// 	}
//
// 	ct, ok := cts[cType]
// 	if !ok {
// 		log.Info("requested contentType not found", "cType", cType)
// 		jsonutil.Error(w, "requested contentType not found", http.StatusBadRequest)
// 		return
// 	}
//
// 	c, err := ct.MetaToChanges(mb)
// 	if err != nil {
// 		log.Error("failed to get changes from meta", "err", err)
// 		jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
// 			http.StatusInternalServerError)
// 		return
// 	}
//
// 	// Strip out common metadata keys that shouldn't be uploaded/imported. Eg, anchor
// 	// is a unique hash every time, so we cannot link to the anchor in the future
// 	// upload, as it does not exist.
// 	//
// 	// TODO(leeola): MetaToChanges should likely be refactored to be ExportMeta so it
// 	// can be responsible for deleting keys that will fail during an upload, like
// 	// Anchor.
// 	delete(c, "anchor")
// 	delete(c, "previousMeta")
// 	delete(c, "multiHash")
// 	delete(c, "multiPart")
//
// 	_, err = jsonutil.MarshalToWriter(w, ChangesResponse{
// 		Changes: c,
// 	})
// 	if err != nil {
// 		log.Error("failed to marshal response", "err", err)
// 		jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
// 			http.StatusInternalServerError)
// 		return
// 	}
// }
