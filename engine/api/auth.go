package api

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type Auth struct {
	Secret          string
	Authority       string
	ExpirationHours time.Duration
}

// AuthClaim adds email claim to the standard token claims
type AuthClaim struct {
	Email string
	jwt.StandardClaims
}

// GenerateToken generates a JWT token
func (a *Auth) GenerateToken(email string) (string, error) {
	claims := &AuthClaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * a.ExpirationHours).Unix(),
			Issuer:    a.Authority,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(a.Secret))
	if err != nil {
		return "", err
	}

	return signed, nil
}

// ValidateToken validates a JWT token
func (a *Auth) ValidateToken(signedToken string) (*AuthClaim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&AuthClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(a.Secret), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AuthClaim)
	if !ok {
		return nil, errors.New(ErrJWTClaimUnprocessable)
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New(ErrJWTExpired)
	}

	return claims, nil

}
