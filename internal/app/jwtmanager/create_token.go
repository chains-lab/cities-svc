package jwtmanager

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type InvitePayload struct {
	ID        uuid.UUID
	CityID    uuid.UUID
	Role      string
	ExpiredAt time.Time
	CreatedAt time.Time
}

func (m Manager) CreateInviteToken(p InvitePayload) (string, error) {
	claims := inviteClaims{
		CityID: p.CityID,
		Role:   p.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        p.ID.String(),
			Issuer:    m.iss,
			IssuedAt:  jwt.NewNumericDate(p.CreatedAt),
			ExpiresAt: jwt.NewNumericDate(p.ExpiredAt),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString(m.sk)
	if err != nil {
		return "", err
	}
	return signed, nil
}
