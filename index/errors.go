package index

import "github.com/leeola/errors"

var (
	// ErrNoQueryResults is to be returned if there are no results matching the
	// given query.
	ErrNoQueryResults = errors.New("no results match query")

	// ErrIndexVersionsDoNotMatch is to be returned if the expected IndexVersion does
	// not match the current Index version.
	ErrIndexVersionsDoNotMatch = errors.New(
		"the expected index version does not match current version")
)
