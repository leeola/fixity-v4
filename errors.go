package kala

import "github.com/leeola/errors"

var (
	// HashNotFoundErr is to be returned by Store implementors when they cannot find
	// content for the given hash.
	ErrHashNotFound = errors.New("hash not found")

	// ErrHashNotMatchContent is to be returned by Store implementors if a given hash
	// does not match the expected content write.
	ErrHashNotMatchContent = errors.New("hash does not match content")

	// ErrNotVersion is returned when a hash's contents are being unmarshalled into
	// a Version, but the contents do not match a Version.
	ErrNotVersion = errors.New("hash is not a valid Version struct")
)
