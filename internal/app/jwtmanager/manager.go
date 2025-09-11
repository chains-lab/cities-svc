package jwtmanager

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/chains-lab/cities-svc/internal/config"
	"github.com/chains-lab/enum"
)

// Manager подписывает/проверяет инвайт-JWT.
type Manager struct {
	iss string
	sk  string
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
		sk:  cfg.JWT.Invite.SecretKey,
	}
}

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
