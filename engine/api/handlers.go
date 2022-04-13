package api

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/msalopek/animus/engine"
	"github.com/msalopek/animus/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

func (api *HttpAPI) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func (api *HttpAPI) Register(c *gin.Context) {
	var creds Credentials

	if err := c.BindJSON(&creds); err != nil {
		abortWithError(c, http.StatusBadRequest, engine.ErrCouldNotRegister)
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
		abortWithError(c, http.StatusBadRequest, engine.ErrInvalidCredentials)
		return
	}

	if len(creds.Email) < 1 || len(creds.Password) < 1 {
		abortWithError(c, http.StatusBadRequest, engine.ErrInvalidCredentials)
		return
	}

	var user models.User
	// TODO: mv database operatios to Repo
	api.db.Where("email = ?", creds.Email).First(&user)

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(creds.Password)); err != nil {
		abortWithError(c, http.StatusBadRequest, engine.ErrInvalidCredentials)
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
	})

	token, err := claims.SignedString([]byte(api.secret))

	if err != nil {
		abortWithError(c, http.StatusInternalServerError, engine.ErrCouldNotLogin)
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (api *HttpAPI) WhoAmI(c *gin.Context) {
	var user models.User

	// auth middleware injects this
	email := c.GetString("email")
	if len(email) < 1 {
		abortWithError(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	result := api.db.Where("email = ?", email).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		abortWithError(c, http.StatusNotFound, engine.ErrNotFound)
		return
	}

	if result.Error != nil {
		abortWithError(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	c.JSON(200, user)
}

// UploadFile extracts a file from gin.Context (multipart form)
// and synchronously uploads it to the attached IPFS node.
//
// TODO:
// - abstract DB operations
//     - extract user operations (create Repo interface)
//     - extract storage operations (create Repo interface)
// - make file upload async
func (api *HttpAPI) UploadFile(c *gin.Context) {
	email, ok := c.Get("email")
	if !ok {
		abortWithError(c, http.StatusUnauthorized, engine.ErrUnauthorized)
		return
	}

	var user models.User
	res := api.db.Where("email = ?", email).First(&user)
	if res.Error == gorm.ErrRecordNotFound {
		abortWithError(c, http.StatusNotFound, engine.ErrUserNotFound)
		return
	}

	if res.Error != nil {
		// TODO: don't leak this error
		abortWithError(c, http.StatusInternalServerError, res.Error.Error())
		return
	}

	meta := c.PostForm("meta")
	var parsed map[string]interface{}
	// TODO: currently this is just unmarshalled to validate JSON
	err := json.Unmarshal([]byte(meta), &parsed)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, engine.ErrInvalidMeta)
		return
	}

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}

	filename := filepath.Base(file.Filename)
	src, err := file.Open()
	if err != nil {
		// TODO: don't expose this error
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer src.Close()
	cid, err := api.IPFSUpload(src)
	if err != nil {
		// TODO: don't expose this error
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	storeRec := models.Storage{
		Cid:      &cid,
		UserID:   user.ID,
		Name:     filename,
		Metadata: meta,
		// TODO: make this configurable;
		// IPFS node does not need to autopin files
		Pinned: true,
	}

	// TODO: check errs
	api.db.Create(storeRec)

	c.JSON(http.StatusOK, gin.H{
		"name": filename,
		"CID":  cid,
		"meta": meta,
	})
}
