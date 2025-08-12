package main

import (
	"errors"
	"net/http"

	"github.com/mafi020/social/internal/dto"
	"github.com/mafi020/social/internal/errs"
	"github.com/mafi020/social/internal/utils"
)

type registerUserPayload struct {
	UserName string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=25"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload registerUserPayload

	if err := utils.ReadJSON(r, &payload); err != nil {
		app.badRequestError(w, r, errors.New(err.Error()))
		return
	}

	if err := utils.ValidateStruct(&payload); err != nil {
		app.failedValidationError(w, r, err)
		return
	}

	ctx := r.Context()

	// ✅ Step 1: Check if they have an accepted invitation
	invitation, err := app.store.Invitations.GetByEmail(ctx, payload.Email)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			app.failedValidationError(w, r, map[string]string{
				"message": "No invitation found for this email",
			})
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if invitation.Status != "accepted" {
		app.failedValidationError(w, r, map[string]string{
			"invitation": "Invitation has not been accepted yet",
		})
		return
	}

	validationErrors, err := app.store.Users.IsUserUnique(ctx, payload.Email, payload.UserName)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if validationErrors != nil {
		app.failedValidationError(w, r, validationErrors)
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	user := &dto.User{
		UserName: payload.UserName,
		Email:    payload.Email,
		Password: hashedPassword,
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
