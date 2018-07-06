package fixity

type Ref string

type Mutation struct {
	ID        string `json:"id"`
	Namespace string `json:"namespace"`
	Signer    string `json:"signer"`
	Time      string `json:"time"`
	Content   Ref    `json:"content"`
	Signature string `json:"signature"`
}

type Content struct {
	Parts

	// Size is the total bytes for the content.
	Size int64 `json:"size,omitempty"`
}

type Parts struct {
	Parts     []Ref `json:"parts"`
	MoreParts *Ref  `json:"moreParts,omitempty"`
}
