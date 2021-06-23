package services

import (
	"net/http"

	"github.com/xdrm-io/aicra"
)

// route represents a service route with an associated handler
type route struct {
	uri, method string
	fn          interface{}
}

// Wire all methods to their associated routes in the service definition
func (h *Handler) Wire(b *aicra.Builder) error {
	routes := []route{
		{uri: "/users", method: http.MethodGet, fn: h.getAllUsers},
		{uri: "/user/{id}", method: http.MethodGet, fn: h.getUserByID},
		{uri: "/user", method: http.MethodPost, fn: h.createUser},
		{uri: "/user/{id}", method: http.MethodPut, fn: h.updateUser},
		{uri: "/user/{id}", method: http.MethodDelete, fn: h.deleteUser},
		{uri: "/user/{id}/articles", method: http.MethodGet, fn: h.getArticlesByAuthor},

		{uri: "/articles", method: http.MethodGet, fn: h.getAllArticles},
		{uri: "/article/{id}", method: http.MethodGet, fn: h.getArticleByID},
		{uri: "/article/{author}", method: http.MethodPost, fn: h.createArticle},
		{uri: "/article/{id}", method: http.MethodDelete, fn: h.deleteArticle},

		{uri: "/article/{id}/up", method: http.MethodPost, fn: h.upVote},
		{uri: "/article/{id}/down", method: http.MethodPost, fn: h.downVote},
	}

	for _, r := range routes {
		err := b.Bind(r.method, r.uri, r.fn)
		if err != nil {
			return err
		}
	}
	return nil
}
