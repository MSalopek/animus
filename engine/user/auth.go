package user

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msalopek/animus/engine"
	"github.com/msalopek/animus/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Credentials struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (c *Credentials) Validate() error {
	if c.Username == "" {
		return errors.New("username not provided")
	}
	if c.Firstname == "" {
		return errors.New("firstname not provided")
	}
	if c.Lastname == "" {
		return errors.New("lastname not provided")
	}
	if c.Email == "" {
		return errors.New("email not provided")
	}
	if len(c.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}
func (api *UserAPI) Register(c *gin.Context) {
	var creds Credentials

	if err := c.BindJSON(&creds); err != nil {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrCouldNotRegister)
		return
	}

	if err := creds.Validate(); err != nil {
		engine.AbortErr(c, http.StatusBadRequest, err)
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
	res := api.repo.Create(&user)
	if res.Error != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrCouldNotRegister)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (api *UserAPI) Login(c *gin.Context) {
	var creds Credentials

	// TODO: log body for debugging
	if err := c.BindJSON(&creds); err != nil {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidCredentials)
		return
	}

	if len(creds.Email) < 1 || len(creds.Password) < 8 {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidCredentials)
		return
	}

	user, err := api.repo.GetUserByEmail(creds.Email)
	if err == gorm.ErrRecordNotFound {
		engine.AbortErr(c, http.StatusNotFound, engine.ErrNotFound)
		return
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(creds.Password)); err != nil {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidCredentials)
		return
	}

	token, err := api.auth.GenerateToken(user.Email)
	if err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrCouldNotLogin)
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (api *UserAPI) WhoAmI(c *gin.Context) {
	email := c.GetString("email")
	if len(email) < 1 {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	user, err := api.repo.GetUserByEmail(email)
	if err == gorm.ErrRecordNotFound {
		engine.AbortErr(c, http.StatusNotFound, engine.ErrNotFound)
		return
	}

	if err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	c.JSON(200, user)
}
