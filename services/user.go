package services

import (
	"context"

	"github.com/xdrm-io/articles-api/model"
)

type byID struct {
	ID uint
}
type userList struct {
	Users []model.User
}

type createUserRequest struct {
	Username  string
	Firstname string
	Lastname  string
}

type updateUserRequest struct {
	ID        uint
	Username  *string
	Firstname *string
	Lastname  *string
}

func (h *Handler) getAllUsers(ctx context.Context) (*userList, error) {
	users, err := h.db.GetAllUsers()
	if err != nil {
		return nil, storageError(err)
	}
	return &userList{Users: users}, nil
}

func (h *Handler) getUserByID(ctx context.Context, req byID) (*model.User, error) {
	user, err := h.db.GetUserByID(req.ID)
	if err != nil {
		return nil, storageError(err)
	}
	return user, nil
}

func (h *Handler) createUser(ctx context.Context, req createUserRequest) (*model.User, error) {
	user, err := h.db.CreateUser(req.Username, req.Firstname, req.Lastname)
	if err != nil {
		return nil, storageError(err)
	}
	return user, nil
}

func (h *Handler) updateUser(ctx context.Context, req updateUserRequest) (*model.User, error) {
	// nothing to update, ignore
	if req.Username == nil && req.Firstname == nil && req.Lastname == nil {
		return h.getUserByID(ctx, byID{req.ID})
	}

	user, err := h.db.UpdateUser(req.ID, req.Username, req.Firstname, req.Lastname)
	if err != nil {
		return nil, storageError(err)
	}

	return user, nil
}

func (h *Handler) deleteUser(ctx context.Context, req byID) error {
	err := h.db.DeleteUser(req.ID)
	if err != nil {
		return storageError(err)
	}
	return nil
}
