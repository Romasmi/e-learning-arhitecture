package account

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrAccountNotFound = errors.New("account not found")

type Account struct {
	UserID    uuid.UUID `json:"user_id"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
