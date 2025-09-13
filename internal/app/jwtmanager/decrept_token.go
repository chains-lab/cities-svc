package jwtmanager

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (m Manager) DecryptInviteToken(tokenStr string) (InviteData, error) {
	var out InviteData

	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	var claims inviteClaims
	token, err := parser.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
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

	if claims.Issuer != "" && claims.Issuer != m.iss {
		return out, fmt.Errorf("invalid issuer")
	}

	if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
		return out, fmt.Errorf("token expired")
	}

	JTI, err := uuid.Parse(claims.ID)
	if err != nil {
		return out, fmt.Errorf("invalid jti format: %w", err)
	}

	out = InviteData{
		JTI:       JTI,
		CityID:    claims.CityID,
		Role:      claims.Role,
		ExpiresAt: claims.ExpiresAt.Time,
		Issuer:    claims.Issuer,
	}
	return out, nil
}
