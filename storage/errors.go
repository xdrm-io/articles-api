package storage

// Err type for this package
type Err string

func (e Err) Error() string {
	return string(e)
}

const (

	// ErrNotInitialized raised when using a non initialized database
	ErrNotInitialized = Err("db not initialized")

	// ErrInitFailed raised when db cannot be initialized
	ErrInitFailed = Err("db init failed")

	// ErrCreate raised on creation errors
	ErrCreate = Err("creation error")

	// ErrUpdate raised on update errors
	ErrUpdate = Err("update error")

	// ErrDelete raised on delete errors
	ErrDelete = Err("delete error")

	// ErrTransaction raised on transaction errors
	ErrTransaction = Err("transaction error")

	// ErrNotFound raised when a resource is not found
	ErrNotFound = Err("resource not found")

	// ErrUnexpected raised on unexpected errors
	ErrUnexpected = Err("unexpected error")
)
