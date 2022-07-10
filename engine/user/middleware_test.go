package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/msalopek/animus/engine"
	"github.com/msalopek/animus/engine/user/auth"
	"github.com/msalopek/animus/model"
	"github.com/stretchr/testify/assert"
)

const (
	testsecret    = "testsecret"
	testauthority = "TestAuthAuthority"
	testemail     = "jtw-unit-test@example.com"
	testpassword  = "test-password!"
	// sigkey != testsecret; sigkey == "123456789";
	testinvalidtoken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJUZXN0QXV0aEF1dGhvcml0eSIsImlhdCI6MTY0OTY4ODg2MCwiZXhwIjpudWxsLCJhdWQiOiIiLCJzdWIiOiIiLCJlbWFpbCI6Imp0dy11bml0LXRlc3RAZXhhbXBsZS5jb20ifQ.QBH7M2WtLT8dbDOWRRXFmKeTNeMq4-P1p0A1yBfR_Fo"
)

func protectedHandler(c *gin.Context) {
	// auth middleware injects this
	email := c.GetString("email")
	if len(email) < 1 {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}
	c.JSON(200, gin.H{
		"email": testemail,
	})
}

func routerWithAuthMiddleware() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	auth := &auth.Auth{
		Secret:    testsecret,
		Authority: testauthority,
	}
	// router uses auth middleware
	router.Use(authorizeUserRequest(auth))

	return router
}

func getValidTestToken() string {
	a := auth.Auth{
		Secret:          testsecret,
		Authority:       testauthority,
		ExpirationHours: 1,
	}
	t, _ := a.GenerateToken(testemail)
	return t
}

func TestAuthorizeUserRequestMissingAuthHeader(t *testing.T) {
	router := routerWithAuthMiddleware()
	router.GET("/protected", protectedHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAuthInvalidTokenFormat(t *testing.T) {
	router := routerWithAuthMiddleware()
	router.GET("/protected", protectedHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Add("Authorization", "fail-unittest")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthorizeUserRequestInvalidToken(t *testing.T) {
	router := routerWithAuthMiddleware()
	router.GET("/protected", protectedHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	// the token was signed with a different key
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testinvalidtoken))

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestValidToken(t *testing.T) {
	// uses valid testdata and sets correct claims
	validToken := getValidTestToken()
	router := routerWithAuthMiddleware()
	router.GET("/protected", protectedHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", validToken))

	router.ServeHTTP(w, req)

	var resp model.User
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, testemail, resp.Email)
}
