package api

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGenerateToken tests token generation
// and sets the generated token as env var jwtTestToken
// to be used by other tests.
func TestGenerateToken(t *testing.T) {
	jwtWrapper := Auth{
		Secret:          "123456789",
		Authority:       "TestAuthAuthority",
		ExpirationHours: 1,
	}

	generatedToken, err := jwtWrapper.GenerateToken("test.jwt@example.com")
	assert.NoError(t, err)

	os.Setenv("jwtTestToken", generatedToken)
}

// TestValidateToken validates generated token.
// The token to validate is read from env var jwtTestToken.
func TestValidateToken(t *testing.T) {
	encodedToken := os.Getenv("jwtTestToken")

	jwtWrapper := Auth{
		Secret:    "123456789",
		Authority: "TestAuthAuthority",
	}

	claims, err := jwtWrapper.ValidateToken(encodedToken)
	assert.NoError(t, err)

	assert.Equal(t, "test.jwt@example.com", claims.Email)
	assert.Equal(t, "TestAuthAuthority", claims.Issuer)
}

// TestValidateTokenWrongAuthority validates generated token.
// The token to validate is read from env var jwtTestToken.
func TestValidateTokenWrongAuthority(t *testing.T) {
	encodedToken := os.Getenv("jwtTestToken")

	jwtWrapper := Auth{
		Secret:    "123456789",
		Authority: "WrongTestAuthAuthority",
	}

	claims, err := jwtWrapper.ValidateToken(encodedToken)
	assert.NoError(t, err)

	assert.Equal(t, "test.jwt@example.com", claims.Email)
	assert.NotEqual(t, "WrongTestAuthAuthority", claims.Issuer)
}

// TestValidateTokenWrongSecret validates generated token.
// The token to validate is read from env var jwtTestToken.
func TestValidateTokenWrongSecret(t *testing.T) {
	encodedToken := os.Getenv("jwtTestToken")

	jwtWrapper := Auth{
		Secret:    "xxxxxxxxx",
		Authority: "TestAuthAuthority",
	}

	claims, err := jwtWrapper.ValidateToken(encodedToken)
	assert.Error(t, err)
	assert.EqualError(t, err, "signature is invalid")
	assert.Nil(t, claims)
}
