package fixity

type Ref string

type Mutation struct {
	ID   string `json:"id"`
	Time string `json:"time"`
	Blob Ref    `json:"blob"`
}

type Blob struct {
	Chunks        []Ref `json:"chunks"`
	BlobContinued Ref   `json:"blobContinued"`

	// Size is the total bytes for the blob.
	Size int64 `json:"size,omitempty"`
}
