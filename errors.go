package fixity

import "github.com/leeola/errors"

var (
	// ErrEmptyBlockchain is to be returned when no blockchain is available.
	//
	// Such as when a Blockchain has not yet been written to.
	ErrEmptyBlockchain = errors.New("empty blockchain")

	// HashNotFoundErr is to be returned by Store implementors when they cannot find
	// content for the given hash.
	ErrHashNotFound = errors.New("hash not found")

	// ErrHashNotMatchBytes is to be returned by Store implementors if a given hash
	// does not match the expected content write.
	ErrHashNotMatchBytes = errors.New("hash does not match bytes")

	// ErrNotContent is returned if the json of a hash is not a Content.
	ErrNotContent = errors.New("hash is not a valid Content struct")

	// ErrFieldNotFound is returned when a FieldUnmarshaller cannot unmarshal the field.
	ErrFieldNotFound = errors.New("field ummarshaller cannot find field")
)
