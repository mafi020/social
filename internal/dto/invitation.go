package dto

import "time"

type Invitation struct {
	ID          int64      `json:"id"`
	InviterID   int64      `json:"inviter_id"`
	Email       string     `json:"email"`
	EmailSentAt *time.Time `json:"email_sent_at,omitempty"` // pointer to handle null
	Token       string     `json:"token"`
	ExpiresAt   time.Time  `json:"expires_at"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
