package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mafi020/social/internal/dto"
	"github.com/mafi020/social/internal/errs"
	"github.com/mafi020/social/internal/utils"
)

type createCommentPayload struct {
	PostID  int64  `json:"post_id" validate:"required"`
	UserID  int64  `json:"user_id" validate:"required"`
	Content string `json:"content" validate:"required"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	payload := createCommentPayload{}
	if err := utils.ReadJSON(r, &payload); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := utils.ValidateStruct(&payload); err != nil {
		app.failedValidationError(w, r, err)
		return
	}

	ctx := r.Context()
	comment := &dto.Comment{
		PostID:  payload.PostID,
		UserID:  payload.UserID,
		Content: payload.Content,
	}

	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	userData, err := app.store.Users.GetById(ctx, payload.UserID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	comment.User = dto.CommentUser{
		ID:       userData.ID,
		UserName: userData.UserName,
	}

	if err := utils.JSONResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getCommentHandler(w http.ResponseWriter, r *http.Request) {
	paramCommentID := chi.URLParam(r, "commentID")
	commentId, err := strconv.ParseInt(paramCommentID, 10, 64)
	if err != nil {
		app.badRequestError(w, r, errors.New("invalid comment ID"))
		return
	}

	ctx := r.Context()

	comment, err := app.store.Comments.GetByID(ctx, commentId)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := utils.JSONResponse(w, http.StatusOK, comment); err != nil {
		app.internalServerError(w, r, err)
	}
}

type updateCommentPayload struct {
	Content *string `json:"content" validate:"omitempty,min=1"`
}

func (app *application) updateCommentHandler(w http.ResponseWriter, r *http.Request) {
	paramCommentID := chi.URLParam(r, "commentID")
	commentID, err := strconv.ParseInt(paramCommentID, 10, 64)
	if err != nil {
		app.badRequestError(w, r, errors.New("invalid comment ID"))
		return
	}
	var payload updateCommentPayload
	if err := utils.ReadJSON(r, &payload); err != nil {
		app.badRequestError(w, r, err)
	}

	validationErros := utils.ValidateStruct(&payload)
	if validationErros != nil {
		app.failedValidationError(w, r, validationErros)
	}
	ctx := r.Context()

	comment, err := app.store.Comments.GetByID(ctx, commentID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			app.notFoundError(w, r, err)
		} else {
			app.internalServerError(w, r, err)
		}
		return
	}

	if payload.Content != nil {
		comment.Content = *payload.Content
	}

	if err := app.store.Comments.Update(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := utils.JSONResponse(w, http.StatusOK, comment); err != nil {
		app.internalServerError(w, r, err)
	}

}
func (app *application) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	paramCommentID := chi.URLParam(r, "commentID")
	commentID, err := strconv.ParseInt(paramCommentID, 10, 64)
	if err != nil {
		app.badRequestError(w, r, errors.New("invalid comment ID"))
	}
	ctx := r.Context()
	if err := app.store.Comments.Delete(ctx, commentID); err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "Comment deleted successfully"}); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
