package model

import (
	"time"

	"github.com/google/uuid"
)

// User keeps user data.
type User struct {
	ID        uuid.UUID
	Login     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
