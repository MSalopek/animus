package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msalopek/animus/engine"
	"github.com/msalopek/animus/model"
	"github.com/msalopek/animus/queue"
	log "github.com/sirupsen/logrus"
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
	// if c.Firstname == "" {
	// 	return errors.New("firstname not provided")
	// }
	// if c.Lastname == "" {
	// 	return errors.New("lastname not provided")
	// }
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
		Password:  hash,
		MaxKeys:   model.DefaultMaxKeys,
		CreatedAt: time.Now(),
	}

	if creds.Firstname != "" {
		user.Firstname = &creds.Firstname
	}

	if creds.Lastname != "" {
		user.Lastname = &creds.Lastname
	}

	if creds.Username == "" {
		user.Username = creds.Email
	}

	res := api.repo.Create(&user)
	if res.Error != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrCouldNotRegister)
		return
	}

	tkRes, err := api.repo.CreateRegisterToken(int(user.ID))
	if err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrCouldNotRegister)
		return
	}

	err = api.publishRegisterEmail(&user, tkRes.Token)
	if err != nil {
		api.logger.WithFields(
			log.Fields{
				"message":   "failed to publish register email event",
				"userID":    user.ID,
				"userEmail": user.Email,
			}).Error(err)
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
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (api *UserAPI) ActivateUser(c *gin.Context) {
	emailP, ok := c.Params.Get("email")
	if !ok {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrBadRequest)
		return
	}

	tokenP, ok := c.GetQuery("token")
	if !ok {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrBadRequest)
		return
	}

	token, err := api.repo.GetUserToken(emailP, tokenP, model.TokenTypeRegisterEmail)
	if err != nil {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrBadRequest)
		return
	}

	if time.Now().After(token.ValidTo.Time) {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrTokenExpired)
		return
	}

	user, err := api.repo.GetUserById(int(token.UserID))
	if err != nil {
		engine.AbortErr(c, http.StatusNotFound, engine.ErrNotFound)
		return
	}

	timestamp := time.Now()
	if res := api.repo.Model(token).Updates(
		model.Token{ValidTo: model.PgTime{Time: timestamp}, IsUsed: true}); res.Error != nil {

		// NOTE: just log errors, tokens expire anyway
		api.logger.WithFields(
			log.Fields{
				"message": "failed to deactivate token after use",
				"userID":  token.UserID,
				"token":   token.Token,
				"type":    token.Type,
			}).Error(res.Error)
	}

	user.Active = true
	user.UpdatedAt = timestamp
	if res := api.repo.Model(user).Updates(
		model.User{Active: true, UpdatedAt: timestamp}); res.Error != nil {
		api.logger.WithFields(
			log.Fields{
				"message": "failed to activate user",
				"userID":  token.UserID,
				"token":   token.Token,
				"type":    token.Type,
			}).Error(res.Error)
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	c.JSON(http.StatusOK, user)
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

func (api *UserAPI) publishRegisterEmail(m *model.User, token string) error {
	e := queue.RegisterEmail{
		Email:     m.Email,
		Token:     token,
		Username:  &m.Username,
		Firstname: m.Firstname,
		Lastname:  m.Lastname,
	}
	body, err := json.Marshal(e)
	if err != nil {
		return err
	}

	err = api.publisher.Publish(api.cfg.NsqEmailRegisterTopic, body)
	if err != nil {
		return err
	}
	return nil
}
