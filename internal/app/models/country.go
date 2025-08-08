package models

import (
	"time"

	"github.com/google/uuid"
)

type Country struct {
	ID        uuid.UUID
	Name      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
