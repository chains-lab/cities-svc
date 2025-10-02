package jwtmanager

import (
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (m Manager) DecryptInviteToken(
	tokenStr string,
) (models.InviteTokenData, error) {
	var out models.InviteTokenData

	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	var claims inviteClaims
	token, err := parser.ParseWithClaims(string(tokenStr), &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.sk, nil
	})
	if err != nil {
		return out, err
	}
	if !token.Valid {
		return out, fmt.Errorf("invalid token")
	}

	if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
		return out, fmt.Errorf("token expired")
	}

	inviteID, err := uuid.Parse(claims.ID)
	if err != nil {
		return out, fmt.Errorf("invalid jti format: %w", err)
	}

	out = models.InviteTokenData{
		InviteID:  inviteID,
		CityID:    claims.CityID,
		Role:      claims.Role,
		ExpiresAt: claims.ExpiresAt.Time,
	}
	return out, nil
}
