package domain

import (
	"errors"
	"time"
)

type AccountStatus string

const (
	AccountStatusActive   AccountStatus = "ACTIVE"
	AccountStatusArchived AccountStatus = "ARCHIVED"
)

type Account struct {
	ID        string        `json:"id"`
	PortalID  string        `json:"portal_id"`
	Name      string        `json:"name"`
	Status    AccountStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
}

type Admin struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrAdminNotFound   = errors.New("admin not found")
)
