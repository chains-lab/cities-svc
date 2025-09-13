package jwtmanager

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/chains-lab/cities-svc/internal/config"
	"github.com/chains-lab/enum"
)

// Manager подписывает/проверяет инвайт-JWT.
type Manager struct {
	iss string
	sk  []byte
}

type InviteData struct {
	JTI       uuid.UUID
	CityID    uuid.UUID
	Role      string
	ExpiresAt time.Time
	Issuer    string
}

// наши claims внутри JWT
type inviteClaims struct {
	CityID uuid.UUID `json:"city_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

func NewManager(cfg config.Config) Manager {
	return Manager{
		iss: enum.CitiesSVC,
		sk:  []byte(cfg.JWT.Invites.SecretKey),
	}
}
