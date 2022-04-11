package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/msalopek/animus/models"
	"golang.org/x/crypto/bcrypt"
)

// TODO: remove
func WIPresponder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "method is under construction",
	})
}

type Credentials struct {
	Username  string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (api *HttpAPI) Register(c *gin.Context) {
	var creds Credentials

	if err := c.BindJSON(&creds); err != nil {
		abortWithError(c, http.StatusBadRequest, ErrCouldNotRegister)
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(creds.Password), 12)
	user := models.User{
		Username:  creds.Username,
		Email:     creds.Email,
		Password:  hash,
		CreatedAt: time.Now(),
	}
	// TODO: mv database operatios to Repo
	api.db.Create(&user)

	c.JSON(http.StatusCreated, user)
}

func (api *HttpAPI) Login(c *gin.Context) {
	var creds Credentials

	// TODO: log body for debugging
	if err := c.BindJSON(&creds); err != nil {
		abortWithError(c, http.StatusBadRequest, ErrInvalidCredentials)
		return
	}

	if len(creds.Email) < 1 || len(creds.Password) < 1 {
		abortWithError(c, http.StatusBadRequest, ErrInvalidCredentials)
		return
	}

	var user models.User
	// TODO: mv database operatios to Repo
	api.db.Where("email = ?", creds.Email).First(&user)

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(creds.Password)); err != nil {
		abortWithError(c, http.StatusBadRequest, ErrInvalidCredentials)
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
	})

	token, err := claims.SignedString([]byte(api.secret))

	if err != nil {
		abortWithError(c, http.StatusInternalServerError, ErrCouldNotLogin)
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func abortWithError(c *gin.Context, httpCode int, err string) {
	c.JSON(httpCode, gin.H{
		"error": err,
	})
}
