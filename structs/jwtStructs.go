package structs

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTUser struct {
	Email           string      `json:"email"`
	ID              uuid.UUID   `json:"id"`
	UpdatedSecurity interface{} `json:"updated_security"`
	jwt.RegisteredClaims
}
