package jwtmanager

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/chains-lab/cities-svc/internal/config"
	"github.com/chains-lab/cities-svc/internal/constant"
)

// Manager подписывает/проверяет инвайт-JWT.
type Manager struct {
	iss string
	sk  string
}

type InviteData struct {
	JTI       string
	CityID    uuid.UUID
	Role      string
	InvitedBy uuid.UUID
	ExpiresAt time.Time
	Issuer    string
}

// наши claims внутри JWT
type inviteClaims struct {
	CityID    uuid.UUID `json:"city_id"`
	Role      string    `json:"role"`
	InvitedBy uuid.UUID `json:"invited_by,omitempty"`
	jwt.RegisteredClaims
}

func NewManager(cfg config.Config) Manager {
	return Manager{
		iss: constant.ServiceName,
		sk:  cfg.JWT.Invite.SecretKey,
	}
}

type InvitePayload struct {
	CityID    uuid.UUID
	Role      string
	InvitedBy uuid.UUID
	ExpiredAt time.Time
}

func (m Manager) CreateInviteToken(p InvitePayload) (string, uuid.UUID, error) {
	now := time.Now().UTC()
	id := uuid.New()

	claims := inviteClaims{
		CityID:    p.CityID,
		Role:      p.Role,
		InvitedBy: p.InvitedBy,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id.String(),
			Issuer:    m.iss,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(p.ExpiredAt),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString(m.sk)
	if err != nil {
		return "", uuid.Nil, err
	}
	return signed, id, nil
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
		return out, errors.New("invalid token")
	}

	if claims.Issuer != "" && claims.Issuer != m.iss {
		return out, errors.New("invalid issuer")
	}

	if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
		return out, errors.New("token expired")
	}

	out = InviteData{
		JTI:       claims.ID,
		CityID:    claims.CityID,
		Role:      claims.Role,
		InvitedBy: claims.InvitedBy,
		ExpiresAt: claims.ExpiresAt.Time,
		Issuer:    claims.Issuer,
	}
	return out, nil
}
