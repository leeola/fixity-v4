package fixity

type Ref string

type Mutation struct {
	ID   string `json:"id"`
	Time string `json:"time"`
	Blob Ref    `json:"blob"`
}

type Content struct {
	Part

	// Size is the total bytes for the blob.
	Size int64 `json:"size,omitempty"`
}

type Part struct {
	Chunks   []Ref `json:"chunks"`
	NextPart *Ref  `json:"nextPart,omitempty"`
}
