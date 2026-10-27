package querybuilder

import "errors"

var (
	errInvalidPointer = errors.New("INVALID_POINTER")

	// ErrEmptyFilter is returned by Query.GetDeleteFilter when a delete is
	// attempted with a query that carries no filter. An unfiltered delete
	// would match every document in the collection, so it is rejected
	// instead of being sent to the server.
	ErrEmptyFilter = errors.New("EMPTY_FILTER")
)
