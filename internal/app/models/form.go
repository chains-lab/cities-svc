package models

import (
	"time"

	"github.com/google/uuid"
)

type Form struct {
	ID           uuid.UUID `json:"id"`
	Status       string    `json:"status"`
	CityName     string    `json:"city_name"`
	CountryID    uuid.UUID `json:"country_id"`
	InitiatorID  uuid.UUID `json:"initiator_id"`
	ContactEmail string    `json:"contact_email"`
	ContactPhone string    `json:"contact_phone"`
	Text         string    `json:"text"`
	UserRevID    uuid.UUID `json:"user_reviewed_id"`
	CreatedAt    time.Time `json:"created_at"`
}
