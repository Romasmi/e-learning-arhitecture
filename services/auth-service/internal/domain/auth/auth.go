package auth

import (
	"time"

	"github.com/google/uuid"
)

type Auth struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"-"`
	PortalID     string    `json:"portal_id"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AuthLog struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Login     string    `json:"login"`
	Action    string    `json:"action"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}
