package services

import (
	"context"

	"github.com/xdrm-io/aicra/api"
	"github.com/xdrm-io/articles-api/model"
)

type articleList struct {
	Articles []model.Article
}

type createArticleRequest struct {
	Author uint
	Title  string
	Body   string
}

type voteRequest struct {
	User    uint
	Article uint
}

func (h *Handler) getArticlesByAuthor(ctx context.Context, req byID) (*articleList, api.Err) {
	articles, err := h.db.GetArticlesByAuthor(req.ID)
	if err != nil {
		return nil, storageError(err)
	}
	return &articleList{Articles: articles}, api.ErrSuccess
}

func (h *Handler) getAllArticles(ctx context.Context) (*articleList, api.Err) {
	articles, err := h.db.GetAllArticles()
	if err != nil {
		return nil, storageError(err)
	}
	return &articleList{Articles: articles}, api.ErrSuccess
}

func (h *Handler) getArticleByID(ctx context.Context, req byID) (*model.Article, api.Err) {
	article, err := h.db.GetArticleByID(req.ID)
	if err != nil {
		return nil, storageError(err)
	}
	return article, api.ErrSuccess
}

func (h *Handler) createArticle(ctx context.Context, param createArticleRequest) (*model.Article, api.Err) {
	article, err := h.db.CreateArticle(param.Title, param.Body, param.Author)
	if err != nil {
		return nil, storageError(err)
	}
	return article, api.ErrSuccess
}

func (h *Handler) deleteArticle(ctx context.Context, req byID) api.Err {
	err := h.db.DeleteArticle(req.ID)
	if err != nil {
		return storageError(err)
	}
	return api.ErrSuccess
}

func (h *Handler) upVote(ctx context.Context, req voteRequest) (*model.Article, api.Err) {
	_, err := h.db.UpVote(req.User, req.Article)
	if err != nil {
		return nil, storageError(err)
	}
	// get article
	return h.getArticleByID(ctx, byID{ID: req.Article})
}

func (h *Handler) downVote(ctx context.Context, req voteRequest) (*model.Article, api.Err) {
	_, err := h.db.DownVote(req.User, req.Article)
	if err != nil {
		return nil, storageError(err)
	}
	// get article
	return h.getArticleByID(ctx, byID{ID: req.Article})
}
