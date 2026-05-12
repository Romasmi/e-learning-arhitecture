package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id,omitzero"`
	Username  string    `json:"username,omitzero"`
	FirstName string    `json:"firstName,omitzero"`
	LastName  string    `json:"lastName,omitzero"`
	Email     string    `json:"email,omitzero"`
	Phone     string    `json:"phone,omitzero"`
	Password  string    `json:"password,omitzero"`
	CreatedAt time.Time `json:"createdAt,omitzero"`
}
