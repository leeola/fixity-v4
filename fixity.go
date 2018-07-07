package fixity

type Ref string

type Mutation struct {
	ID        string `json:"id"`
	Namespace string `json:"namespace"`
	Signer    string `json:"signer"`
	Time      string `json:"time"`
	Values    Ref    `json:"values,omitempty"`
	Data      Ref    `json:"data,omitempty"`
	Signature string `json:"signature"`
}

type Data struct {
	Parts

	// Size is the total bytes for the content.
	Size int64 `json:"size,omitempty"`
}

type Parts struct {
	Parts     []Ref `json:"parts"`
	MoreParts *Ref  `json:"moreParts,omitempty"`
}

type Values struct {
	ValueMap ValueMap `json:"valueMap"`
}
