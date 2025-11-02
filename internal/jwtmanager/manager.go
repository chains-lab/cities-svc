package jwtmanager

import (
	"crypto/sha256"

	"github.com/chains-lab/cities-svc/internal"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

const bcryptCost = bcrypt.DefaultCost

func (m Manager) HashInviteToken(tokenStr string) (string, error) {
	sum := sha256.Sum256([]byte(tokenStr)) // 32 байта
	b, err := bcrypt.GenerateFromPassword(sum[:], bcryptCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (m Manager) VerifyInviteToken(tokenStr, hashed string) error {
	sum := sha256.Sum256([]byte(tokenStr))
	return bcrypt.CompareHashAndPassword([]byte(hashed), sum[:])
}
