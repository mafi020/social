package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mafi020/social/internal/dto"
	"github.com/mafi020/social/internal/errs"
	"github.com/mafi020/social/internal/utils"
)

type createPostPayload struct {
	Title   string   `json:"title" validate:"required"`
	Content string   `json:"content" validate:"required"`
	UserID  int64    `json:"user_id" validate:"required"`
	Tags    []string `json:"tags" validate:"required"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload createPostPayload

	if err := utils.ReadJSON(r, &payload); err != nil {
		app.badRequestError(w, r, errors.New(err.Error()))
		return
	}

	if err := utils.ValidateStruct(&payload); err != nil {
		app.failedValidationError(w, r, err)
		return
	}

	ctx := r.Context()

	post := dto.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  payload.UserID,
	}

	if err := app.store.Posts.Create(ctx, &post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = []dto.Comment{}

	if err := utils.JSONResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {

	postIDParam := chi.URLParam(r, "postID")

	postID, err := strconv.ParseInt(postIDParam, 10, 64)
	if err != nil {
		app.badRequestError(w, r, errors.New("invalid post id"))
		return
	}

	ctx := r.Context()

	post, err := app.store.Posts.GetByID(ctx, postID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	comments, err := app.store.Comments.GetCommentsByPostID(ctx, postID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	post.Comments = comments

	if err := utils.JSONResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	paramPostID := chi.URLParam(r, "postID")
	postId, err := strconv.ParseInt(paramPostID, 10, 64)
	if err != nil {
		app.badRequestError(w, r, errors.New("invalid Post ID"))
		return
	}

	ctx := r.Context()
	if err := app.store.Posts.Delete(ctx, postId); err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "Post Deleted Successfully"}); err != nil {
		app.internalServerError(w, r, err)
	}
}

// Important to use pointer to distinguish between intentional empty fields and struct generated nil fields if not provided
type updatePostPayload struct {
	// Ttile type as pointer string means, it'a an optional field
	Title   *string   `json:"title" validate:"omitempty,min=1"`
	Content *string   `json:"content" validate:"omitempty,min=1"`
	Tags    *[]string `json:"tags" validate:"omitempty,min=1"`
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	postIDParam := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(postIDParam, 10, 64)
	if err != nil {
		app.badRequestError(w, r, errors.New("invalid post id"))
		return
	}

	var payload updatePostPayload
	if err := utils.ReadJSON(r, &payload); err != nil {
		app.badRequestError(w, r, errors.New(err.Error()))
		return
	}

	log.Printf("After reading the data %+v", payload)

	if validationErrors := utils.ValidateStruct(&payload); validationErrors != nil {
		app.failedValidationError(w, r, validationErrors)
		return
	}

	ctx := r.Context()

	post, err := app.store.Posts.GetByID(ctx, postID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			app.notFoundError(w, r, err)
		} else {
			app.internalServerError(w, r, err)
		}
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Tags != nil {
		post.Tags = *payload.Tags
	}

	if err := app.store.Posts.Update(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = []dto.Comment{}

	if err := utils.JSONResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
