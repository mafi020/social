package main

import (
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/mafi020/social/internal/utils"
)

func (app *application) refreshHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("refresh_token")
	if err != nil || c.Value == "" {
		app.unAuthorizedError(w, r, errors.New("missing refresh token"))
		return
	}
	raw := c.Value
	hash := utils.HashToken(raw)

	ctx := r.Context()

	// Fetch the token *even if revoked* — critical for reuse detection
	rt, err := app.store.RefreshTokens.GetByHash(ctx, hash)
	if err != nil {
		// Unknown/garbage token
		app.unAuthorizedError(w, r, errors.New("invalid refresh token"))
		return
	}

	// ===== Reuse detection =====
	// If a token that was already rotated (revoked) shows up again, it’s almost certainly stolen.
	if rt.RevokedAt != nil {
		// Defensive response: revoke ALL refresh tokens for this user
		_ = app.store.RefreshTokens.RevokeAllForUser(ctx, rt.UserID)

		// Optional: log a security event, notify the user, etc.
		app.logger.Warnw("refresh token reuse detected",
			"user_id", rt.UserID,
			"ip", rt.IPAddress,
			"ua", rt.UserAgent,
			"revoked_at", rt.RevokedAt,
		)

		app.unAuthorizedError(w, r, errors.New("token reuse detected; all sessions revoked"))
		return
	}

	// Normal validation
	if time.Now().After(rt.ExpiresAt) {
		app.unAuthorizedError(w, r, errors.New("refresh token expired"))
		return
	}

	// Rotate: revoke the current token
	if err := app.store.RefreshTokens.Revoke(ctx, hash); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Mint a new refresh token (opaque)
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
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	if _, err := app.store.RefreshTokens.Create(ctx, rt.UserID, newHash, ua, ip, newExp); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// New access token
	access, err := utils.GenerateAccessToken(rt.UserID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Set new refresh cookie
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
