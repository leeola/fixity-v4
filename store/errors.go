package store

import "github.com/leeola/errors"

var (
	// HashNotFoundErr is to be returned by Store implementors when they cannot find
	// content for the given hash.
	HashNotFoundErr = errors.New("hash not found")

	// HashNotMatchContent is to be returned by Store implementors if a given hash
	// does not match the expected content write.
	HashNotMatchContentErr = errors.New("hash does not match content")
)
