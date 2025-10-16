package models

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Pseudonym   *string   `json:"pseudonym"`
	Description *string   `json:"description"`
	Avatar      *string   `json:"avatar"`
	Official    bool      `json:"official"`
	Sex         *string   `json:"sex"`

	BirthDate *time.Time `json:"birth_date"`
	UpdatedAt time.Time  `json:"updated_at"`
	CreatedAt time.Time  `json:"created_at"`
}
