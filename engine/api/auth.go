package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msalopek/animus/engine"
	"github.com/msalopek/animus/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (api *HttpAPI) Register(c *gin.Context) {
	var creds Credentials

	if err := c.BindJSON(&creds); err != nil {
		abortWithError(c, http.StatusBadRequest, engine.ErrCouldNotRegister)
		return
	}

	if err := creds.Validate(); err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(creds.Password), 12)
	user := model.User{
		Username:  creds.Username,
		Email:     creds.Email,
		Firstname: &creds.Firstname,
		Lastname:  &creds.Lastname,
		Password:  hash,
		CreatedAt: time.Now(),
	}
	// TODO: mv database operatios to Repo
	res := api.repo.Create(&user)
	if res.Error != nil {
		// TODO: don't leak DB errors
		abortWithError(c, http.StatusInternalServerError, res.Error.Error())
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (api *HttpAPI) Login(c *gin.Context) {
	var creds Credentials

	// TODO: log body for debugging
	if err := c.BindJSON(&creds); err != nil {
		abortWithError(c, http.StatusBadRequest, engine.ErrInvalidCredentials)
		return
	}

	if len(creds.Email) < 1 || len(creds.Password) < 1 {
		abortWithError(c, http.StatusBadRequest, engine.ErrInvalidCredentials)
		return
	}

	user, err := api.repo.GetUserByEmail(creds.Email)
	if err == gorm.ErrRecordNotFound {
		abortWithError(c, http.StatusNotFound, engine.ErrNotFound)
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(creds.Password)); err != nil {
		abortWithError(c, http.StatusBadRequest, engine.ErrInvalidCredentials)
		return
	}

	token, err := api.auth.GenerateToken(user.Email)
	if err != nil {
		abortWithError(c, http.StatusInternalServerError, engine.ErrCouldNotLogin)
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (api *HttpAPI) WhoAmI(c *gin.Context) {
	// auth middleware injects this
	email := c.GetString("email")
	if len(email) < 1 {
		abortWithError(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	user, err := api.repo.GetUserByEmail(email)
	if err == gorm.ErrRecordNotFound {
		abortWithError(c, http.StatusNotFound, engine.ErrNotFound)
		return
	}

	if err != nil {
		// abortWithError(c, http.StatusInternalServerError, engine.ErrInternalError)
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(200, user)
}
