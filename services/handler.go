package services

import (
	"github.com/xdrm-io/aicra/api"
	"github.com/xdrm-io/articles-api/storage"
)

// Handler for aicra services
type Handler struct {
	db *storage.DB
}

// NewHandler injects dependencies to create a proper services' handler
func NewHandler(db *storage.DB) *Handler {
	return &Handler{
		db: db,
	}
}

// storageError transforms a storage error into an api error
func storageError(err error) api.Err {
	switch err {
	case storage.ErrTransaction:
		return api.ErrTransaction
	case storage.ErrCreate:
		return api.ErrCreation
	case storage.ErrUpdate:
		return api.ErrModification
	case storage.ErrDelete:
		return api.ErrDeletion
	case storage.ErrNotFound:
		return api.ErrNoMatchFound
	default:
		return api.ErrUnknown
	}
}
