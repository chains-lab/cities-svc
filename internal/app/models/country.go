package models

import (
	"time"

	"github.com/chains-lab/cities-dir-svc/internal/enum"
	"github.com/google/uuid"
)

type Country struct {
	ID        uuid.UUID
	Name      string
	Status    enum.CountryStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}
