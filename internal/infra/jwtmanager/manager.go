package jwtmanager

import (
	"github.com/chains-lab/cities-svc/internal"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Manager struct {
	sk []byte
}

type inviteClaims struct {
	CityID uuid.UUID `json:"city_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

func NewManager(cfg internal.Config) Manager {
	return Manager{
		sk: []byte(cfg.JWT.Invites.SecretKey),
	}
}
