package models

import "github.com/google/uuid"

type User struct {
	Email string    `json:"email"`
	ID    uuid.UUID `json:"id"`
}
