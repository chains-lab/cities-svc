package models

import (
	"time"

	"github.com/chains-lab/cities-dir-svc/internal/enum"
	"github.com/google/uuid"
)

type CityAdmin struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	CityID    uuid.UUID
	Role      enum.CityAdminRole
	UpdatedAt time.Time
	CreatedAt time.Time
}
