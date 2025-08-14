package main

import (
	"errors"
	"net"
	"net/http"

	"github.com/mafi020/social/internal/dto"
	"github.com/mafi020/social/internal/env"
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

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := utils.ReadJSON(r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	user, err := app.store.Users.GetByEmail(ctx, payload.Email)
	if err != nil {
		app.failedValidationError(w, r, map[string]string{"email": "Invalid email"})
		return
	}

	if !utils.CheckPassword(user.Password, payload.Password) {
		app.failedValidationError(w, r, map[string]string{"credentials": "invalid email or password"})
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Refresh token (opaque)
	refreshRaw, err := utils.GenerateOpaqueToken(32)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	refreshHash := utils.HashToken(refreshRaw)
	expiresAt := utils.RefreshExpiry(7) // 7 days

	ua := r.UserAgent()
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	if _, err := app.store.RefreshTokens.Create(ctx, user.ID, refreshHash, ua, ip, expiresAt); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	env := env.GetEnvOrPanic("ENVIRONMENT")

	// httpOnly cookie for refresh token
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshRaw,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   env == "production",
		Expires:  expiresAt,
	})

	if err := utils.JSONResponse(w, http.StatusOK, map[string]string{"access_token": accessToken}); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("refresh_token")
	if err == nil && c.Value != "" {
		hash := utils.HashToken(c.Value)
		_ = app.store.RefreshTokens.Revoke(r.Context(), hash)
	}
	// Clear cookie
	env := env.GetEnvOrPanic("ENVIRONMENT")
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   env == "production",
		MaxAge:   -1,
	})
	w.WriteHeader(http.StatusNoContent)
}
