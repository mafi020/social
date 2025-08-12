package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/mafi020/social/internal/dto"
	"github.com/mafi020/social/internal/env"
	"github.com/mafi020/social/internal/errs"
	"github.com/mafi020/social/internal/templates"
	"github.com/mafi020/social/internal/utils"
)

type createInvitationPayload struct {
	Email string `json:"email" validate:"required,email"`
}

func (app *application) createInvitationHandler(w http.ResponseWriter, r *http.Request) {
	var payload createInvitationPayload

	if err := utils.ReadJSON(r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := utils.ValidateStruct(&payload); err != nil {
		app.failedValidationError(w, r, err)
		return
	}

	ctx := r.Context()

	existingInv, err := app.getExistingInvitation(ctx, payload.Email)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if existingInv != nil {
		if app.handleExistingInvitation(ctx, w, r, existingInv) {
			return // handled — stop here
		}
	}

	inviterID := int64(1) // TODO: use auth user ID

	inv, err := app.createNewInvitation(ctx, inviterID, payload.Email)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.sendInvitationEmail(inv)

	now := time.Now()
	inv.EmailSentAt = &now

	if err := app.store.Invitations.UpdateEmailStatus(ctx, inv.ID, &now); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := utils.JSONResponse(w, http.StatusCreated, inv); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) createNewInvitation(ctx context.Context, inviterID int64, email string) (*dto.Invitation, error) {
	token, err := utils.GenerateToken(32)
	if err != nil {
		return nil, err
	}

	inv := &dto.Invitation{
		InviterID:   inviterID,
		Email:       email,
		Token:       token,
		ExpiresAt:   time.Now().Add(48 * time.Hour),
		Status:      "pending",
		EmailSentAt: nil,
	}

	if err := app.store.Invitations.Create(ctx, inv); err != nil {
		return nil, err
	}

	return inv, err
}

func (app *application) refreshAndResendInvitation(ctx context.Context, w http.ResponseWriter, r *http.Request, inv *dto.Invitation) error {

	token, err := utils.GenerateToken(32)
	if err != nil {
		return err
	}

	inv.Token = token
	inv.ExpiresAt = time.Now().Add(48 * time.Hour)
	inv.EmailSentAt = nil
	inv.Status = "pending"

	if err := app.store.Invitations.Update(ctx, inv); err != nil {
		app.internalServerError(w, r, err)
		return nil
	}

	app.sendInvitationEmail(inv)

	if err := utils.JSONResponse(w, http.StatusCreated, inv); err != nil {
		app.internalServerError(w, r, err)
	}
	return nil
}

func (app *application) getExistingInvitation(ctx context.Context, email string) (*dto.Invitation, error) {
	inv, err := app.store.Invitations.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return nil, err
	}
	return inv, nil
}

func (app *application) handleExistingInvitation(ctx context.Context, w http.ResponseWriter, r *http.Request, inv *dto.Invitation) bool {
	switch inv.Status {
	case "accepted":
		app.failedValidationError(w, r, map[string]string{
			"message": "This user has already accepted an invitation.",
		})
		return true

	case "pending":
		if time.Now().Before(inv.ExpiresAt) {
			app.failedValidationError(w, r, map[string]string{
				"message": "An active invitation is pending.",
			})
			return true
		}
		// Expired — refresh and resend
		if err := app.refreshAndResendInvitation(ctx, w, r, inv); err != nil {
			app.internalServerError(w, r, err)
		}
		return true
	}
	return false
}

func (app *application) sendInvitationEmail(inv *dto.Invitation) {
	plainTextContent, htmlContent := templates.EmailInvitation(inv.Token)
	go func() {
		if err := utils.SendEmail(
			"Social Golang Company",
			env.GetEnvOrPanic("COMPANY_EMAIL"), // must match SendGrid verified sender
			"",
			inv.Email,
			"Invitation to Join Social 🎉",
			plainTextContent,
			htmlContent,
		); err != nil {
			app.logger.Error("failed to send invitation email", "error", err)
		}
	}()
}

func (app *application) acceptInvitationHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		app.badRequestError(w, r, errors.New("token is required"))
		return
	}

	inv, err := app.store.Invitations.GetByToken(r.Context(), token)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if inv.Status != "pending" || time.Now().After(inv.ExpiresAt) {
		app.failedValidationError(w, r, map[string]string{"token": "Invitation expired or already used"})
		return
	}

	// At this point, you could prompt the invited user to set a password & register
	// Once registered, update status:
	if err := app.store.Invitations.UpdateStatus(r.Context(), inv.ID, "accepted"); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "Invitation accepted"}); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// type updateInvitationPayload struct {
// 	Status    *string    `json:"status" validate:"omitempty,oneof=pending accepted expired"`
// 	ExpiresAt *time.Time `json:"expires_at" validate:"omitempty"`
// }

// func (app *application) updateInvitationHandler(w http.ResponseWriter, r *http.Request) {
// 	idStr := chi.URLParam(r, "id")
// 	id, err := strconv.ParseInt(idStr, 10, 64)
// 	if err != nil {
// 		app.badRequestError(w, r, errors.New("invalid invitation id"))
// 		return
// 	}

// 	var payload updateInvitationPayload
// 	if err := utils.ReadJSON(r, &payload); err != nil {
// 		app.badRequestError(w, r, err)
// 		return
// 	}
// 	if err := utils.ValidateStruct(&payload); err != nil {
// 		app.failedValidationError(w, r, err)
// 		return
// 	}

// 	ctx := r.Context()

// 	// 1. Get the existing record
// 	inv, err := app.store.Invitations.GetByID(ctx, id)
// 	if err != nil {
// 		if errors.Is(err, errs.ErrNotFound) {
// 			app.notFoundError(w, r, err)
// 		} else {
// 			app.internalServerError(w, r, err)
// 		}
// 		return
// 	}

// 	// 2. Update only provided fields
// 	if payload.Status != nil {
// 		inv.Status = *payload.Status
// 	}
// 	if payload.ExpiresAt != nil {
// 		inv.ExpiresAt = *payload.ExpiresAt
// 	}

// 	// 3. Save changes
// 	if err := app.store.Invitations.Update(ctx, inv); err != nil {
// 		switch {
// 		case errors.Is(err, errs.ErrNotFound):
// 			app.notFoundError(w, r, err)
// 		default:
// 			app.internalServerError(w, r, err)
// 		}
// 		return
// 	}

// 	if err := utils.JSONResponse(w, http.StatusOK, inv); err != nil {
// 		app.internalServerError(w, r, err)
// 		return
// 	}
// }
