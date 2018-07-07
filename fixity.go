package fixity

import "time"

//go:generate stringer -type=BlobType,ValueType -output=fixity_string.go

type Ref string

type Mutation struct {
	Schema
	ID        string    `json:"id"`
	Namespace string    `json:"namespace"`
	Signer    string    `json:"signer"`
	Time      time.Time `json:"time"`
	ValuesMap Ref       `json:"valuesMap,omitempty"`
	Data      Ref       `json:"data,omitempty"`
	Signature string    `json:"signature"`
}

type Data struct {
	Parts

	// Size is the total bytes for the content.
	Size int64 `json:"size,omitempty"`

	// Checksum of the bytes referenced in all the parts, not including any
	// schema information.
	//
	// Ie, just the raw user uploaded data.
	//
	// Hex encoded for user convenience, but using the same hashing algorithm
	// that the content address was decided as.
	//
	// NOTE: For ease of comparison, this hash string *does not* include
	// multihash identification prefixes.
	Checksum string `json:"checksum"`
}

type Parts struct {
	Schema
	Parts     []Ref `json:"parts"`
	MoreParts *Ref  `json:"moreParts,omitempty"`
}

type ValuesMap struct {
	Schema
	Values Values `json:"value"`
}
