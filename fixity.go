package fixity

type Ref string

type Mutation struct {
	ID      string `json:"id"`
	Time    string `json:"time"`
	Content Ref    `json:"blob"`
}

type Content struct {
	Parts

	// Size is the total bytes for the blob.
	Size int64 `json:"size,omitempty"`
}

type Parts struct {
	Parts     []Ref `json:"parts"`
	MoreParts *Ref  `json:"moreParts,omitempty"`
}
