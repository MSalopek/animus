package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/msalopek/animus/model"
	"github.com/stretchr/testify/assert"
)

const (
	testclientkey    = "test-client-key"
	testclientsecret = "my-supersecret-is-here"
)

type mockProvider struct{}

func (mp *mockProvider) GetApiClientByKey(key string) (*model.APIClient, error) {
	return &model.APIClient{
		UserID:       1,
		Email:        "email@example.com",
		ClientKey:    testclientkey,
		ClientSecret: testclientsecret}, nil
}

func clientApiRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(authorizeClientRequest(&mockProvider{}))

	return router
}

func TestMissingAPIHeaders(t *testing.T) {
	router := clientApiRouter()
	router.GET("/protected", protectedHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestInvalidSecret(t *testing.T) {
	router := clientApiRouter()
	router.GET("/protected", protectedHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)

	// message signed with a different key
	mac := hmac.New(sha256.New, []byte("wrong-key"))
	mac.Write([]byte("/protected"))
	msg := hex.EncodeToString(mac.Sum(nil))
	req.Header.Add("X-API-SIGN", msg)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestValidRequestHeaders(t *testing.T) {
	router := clientApiRouter()
	router.GET("/protected", protectedHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)

	// set request URI and headers
	req.Header.Add("X-API-KEY", testclientkey)
	req.RequestURI = "/protected"
	mac := hmac.New(sha256.New, []byte(testclientsecret))
	mac.Write([]byte("/protected"))
	sig := hex.EncodeToString(mac.Sum(nil))
	req.Header.Add("X-API-SIGN", sig)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
