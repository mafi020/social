package main

import (
	"errors"
	"net/http"

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
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := utils.WriteJSON(w, http.StatusCreated, user); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
