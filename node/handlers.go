package node

// func (n *Node) PostUploadMetaHandler(w http.ResponseWriter, r *http.Request) {
// 	log := nodeware.GetLog(r)
// 	metaChanges := contenttype.NewChangesFromValues(r.URL.Query())
//
// 	anchorHash := urlutil.GetQueryString(r, "anchor")
// 	previousMeta := urlutil.GetQueryString(r, "previousMeta")
// 	// If there is no previous meta to base this mutation off of, then query the
// 	// indexer for the most recent hash for this anchor.
// 	if previousMeta == "" && anchorHash != "" {
// 		q := index.Query{
// 			Metadata: index.Metadata{
// 				// NOTE: Putting the hash in quotes because the querystring in bleve
// 				// has issues with a hyphenated hashstring. This is annoying, and
// 				// should be fixed somehow...
// 				"anchor": `"` + anchorHash + `"`,
// 			},
// 		}
// 		s := index.SortBy{
// 			Field:      "uploadedAt",
// 			Descending: true,
// 		}
//
// 		result, err := n.index.QueryOne(q, s)
// 		if err != nil {
// 			log.Error("failed to query for previous meta hash", "err", err)
// 			jsonutil.Error(w, "previous meta query failed", http.StatusInternalServerError)
// 			return
// 		}
//
// 		if result.Hash.Hash != "" {
// 			previousMeta = result.Hash.Hash
// 			metaChanges.SetPreviousMeta(previousMeta)
// 		}
// 	}
//
// 	var metaBytes []byte
// 	cType, ok := metaChanges.GetContentType()
// 	if !ok {
// 		// The caller did not specify the content type, so look it up from the
// 		// previousMeta
// 		if previousMeta != "" {
// 			ct, mb, err := store.GetContentTypeWithBytes(n.store, previousMeta)
// 			if err != nil {
// 				log.Error("failed to get previous content type", "err", err)
// 				jsonutil.Error(w, "contenttype lookup failed", http.StatusInternalServerError)
// 				return
// 			}
// 			cType = ct
// 			metaBytes = mb
// 		}
//
// 		// if even after loading the meta and checking for content type we *still*
// 		// don't have the contentType, set it to the default.
// 		if cType == "" {
// 			cType = "data"
// 			metaChanges.SetContentType(cType)
// 		}
// 	}
// 	log = log.New("contentType", cType)
//
// 	// write a new anchor if specified
// 	if urlutil.GetQueryBool(r, "newAnchor") {
// 		h, err := store.NewAnchor(n.store)
// 		if err != nil {
// 			log.Error("failed to create new anchor", "err", err)
// 			jsonutil.Error(w, "newanchor failed", http.StatusInternalServerError)
// 			return
// 		}
//
// 		if err := n.index.Entry(h); err != nil {
// 			log.Error("failed to index new anchor", "err", err)
// 			jsonutil.Error(w, "newanchor failed", http.StatusInternalServerError)
// 			return
// 		}
//
// 		metaChanges.SetAnchor(h)
// 	}
//
// 	cs, ok := n.contentStorers[cType]
// 	if !ok {
// 		log.Info("requested contentType not found")
// 		jsonutil.Error(w, "requested contentType not found", http.StatusBadRequest)
// 		return
// 	}
//
// 	hashes, err := cs.StoreMeta(metaBytes, metaChanges)
// 	if err != nil {
// 		log.Error("Meta returned error", "err", err)
// 		jsonutil.Error(w, "meta failed", http.StatusInternalServerError)
// 		return
// 	}
//
// 	_, err = jsonutil.MarshalToWriter(w, handlers.HashesResponse{
// 		Hashes: hashes,
// 	})
// 	if err != nil {
// 		log.Error("failed to marshal response", "err", err)
// 		jsonutil.Error(w, http.StatusText(http.StatusInternalServerError),
// 			http.StatusInternalServerError)
// 		return
// 	}
// }
