package fixity

//go:generate stringer -type=BlobType -output=fixity_string.go

type Ref string

type Mutation struct {
	Schema
	ID        string `json:"id"`
	Namespace string `json:"namespace"`
	Signer    string `json:"signer"`
	Time      string `json:"time"`
	ValuesMap Ref    `json:"valuesMap,omitempty"`
	Data      Ref    `json:"data,omitempty"`
	Signature string `json:"signature"`
}

type Data struct {
	Parts

	// Size is the total bytes for the content.
	Size int64 `json:"size,omitempty"`
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
