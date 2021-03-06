package fixity

type Schema struct {
	SchemaType BlobType `json:"_fixitySchema"`
}

type DataSchema struct {
	PartsSchema

	// Size is the total bytes for the content.
	Size int64 `json:"size,omitempty"`

	// Checksum of the bytes referenced in all the parts, not including any
	// schema information.
	//
	// Ie, just the raw user uploaded data.
	//
	// Hex encoded for user convenience, but using the same hashing algorithm
	// that the content address of this dataschema. That is to say, if the
	// content address of this dataschema is a Blake2b multihash, this
	// checksum will be a plain Blake2b hash, not a multihash.
	//
	// IMPORTANT: For ease of comparison, this hash string *does not*
	// include multihash identification prefixes.
	Checksum string `json:"checksum"`
}

type PartsSchema struct {
	Schema
	Parts     []Ref `json:"parts"`
	MoreParts *Ref  `json:"moreParts,omitempty"`
}

type ValuesSchema struct {
	Schema
	Values Values `json:"values"`
}
