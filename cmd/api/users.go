package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mafi020/social/internal/errs"
	"github.com/mafi020/social/internal/models"
	"github.com/mafi020/social/internal/utils"
)

type createUserPayload struct {
	UserName string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload createUserPayload

	if err := utils.ReadJSON(r, &payload); err != nil {
		app.badRequestError(w, r, errors.New(err.Error()))
		return
	}

	if err := utils.ValidateStruct(&payload); err != nil {
		app.failedValidationError(w, r, err)
		return
	}

	ctx := r.Context()

	validationErrors, err := app.store.Users.IsUserUnique(ctx, payload.Email, payload.UserName)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if validationErrors != nil {
		app.failedValidationError(w, r, validationErrors)
		return
	}

	user := &models.User{
		UserName: payload.UserName,
		Email:    payload.Email,
		Password: payload.Password,
	}

	if err := app.store.Users.Create(ctx, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := utils.JSONResponse(w, http.StatusCreated, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	paramUserId := chi.URLParam(r, "userID")
	userId, err := strconv.ParseInt(paramUserId, 10, 64)
	if err != nil {
		app.badRequestError(w, r, errors.New("invalid user ID"))
		return
	}
	ctx := r.Context()
	user, err := app.store.Users.GetById(ctx, userId)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	if err := utils.JSONResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	paramUserID := chi.URLParam(r, "userID")

	userID, err := strconv.ParseInt(paramUserID, 10, 64)
	if err != nil {
		app.badRequestError(w, r, errors.New("invalid user ID"))
		return
	}

	ctx := r.Context()
	if err := app.store.Users.Delete(ctx, userID); err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "User deleted successfully"}); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
