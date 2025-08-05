package main

import (
	"net/http"

	"github.com/mafi020/social/internal/models"
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
	comment := &models.Comment{
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

	comment.User = models.CommentUser{
		ID:       userData.ID,
		UserName: userData.UserName,
	}

	if err := utils.WriteJSON(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
