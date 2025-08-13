package main

import (
	"net/http"
	"time"

	"github.com/mafi020/social/internal/utils"
)

func (app *application) refreshHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("refresh_token")
	if err != nil || c.Value == "" {
		http.Error(w, "missing refresh token", http.StatusUnauthorized)
		return
	}
	raw := c.Value
	hash := utils.HashToken(raw)

	ctx := r.Context()
	rt, err := app.store.RefreshTokens.GetByHash(ctx, hash)
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}
	if rt.RevokedAt != nil || time.Now().After(rt.ExpiresAt) {
		http.Error(w, "refresh token expired or revoked", http.StatusUnauthorized)
		return
	}

	// Rotate: revoke old
	if err := app.store.RefreshTokens.Revoke(ctx, hash); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Issue new refresh
	newRaw, err := utils.GenerateOpaqueToken(32)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	newHash := utils.HashToken(newRaw)
	newExp := utils.RefreshExpiry(7)

	ua := r.UserAgent()
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	if _, err := app.store.RefreshTokens.Create(ctx, rt.UserID, newHash, ua, ip, newExp); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// New access
	access, err := utils.GenerateAccessToken(rt.UserID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Set new cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRaw,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // true in prod
		Expires:  newExp,
	})

	if err := utils.JSONResponse(w, http.StatusOK, map[string]string{
		"access_token": access,
	}); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
