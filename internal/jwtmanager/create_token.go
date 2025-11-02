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
}

func (m Manager) CreateInviteToken(
	inviteID uuid.UUID,
	role string,
	cityID uuid.UUID,
	ExpiredAt time.Time,
) (string, error) {

	claims := inviteClaims{
		CityID: cityID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        inviteID.String(),
			ExpiresAt: jwt.NewNumericDate(ExpiredAt),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString(m.sk)
	if err != nil {
		return "", err
	}
	return signed, nil
}
