package node

type HashResponse struct {
	Hash  string `json:",omitempty"`
	Error string `json:",omitempty"`
}

type HashesResponse struct {
	Hashes []string `json:",omitempty"`
	Error  string   `json:",omitempty"`
}
