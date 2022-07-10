package user

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
	"github.com/msalopek/animus/engine"
	"github.com/msalopek/animus/engine/repo"
	"github.com/msalopek/animus/model"
)

func (api *UserAPI) GetUserKeys(c *gin.Context) {
	uid := c.GetInt("userID")
	ctx := repo.QueryCtxFromGin(c)
	keys, err := api.repo.GetUserApiKeys(ctx, uid)
	if err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	c.JSON(http.StatusOK, keys)
}

// CreateUserKey creates a new API key.
// If the user has reached their maximum key allocation the
// creation will fail. The key allocation is tracked on the DB using triggers.
func (api *UserAPI) CreateUserKey(c *gin.Context) {
	uid := c.GetInt("userID")

	createdAt := time.Now()
	secret := uniuri.NewLen(32)
	key := model.Key{
		UserID:       int64(uid),
		ClientKey:    uniuri.NewLen(32),
		ClientSecret: secret,
		Rights:       model.ClientAccessRead,
		CreatedAt:    createdAt,
		ValidFrom:    createdAt,
	}
	res := api.repo.Create(&key)
	if res.Error != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	c.JSON(http.StatusCreated, engine.CreateKeyResponse{
		ID:        key.ID,
		UserID:    key.UserID,
		ClientKey: key.ClientKey,
		// expose the secret to the user after creation
		// the secret cannot be accessed after it is created
		ClientSecret: key.ClientSecret,
		Rights:       key.Rights,
		Disabled:     key.Disabled,
		CreatedAt:    key.CreatedAt,
		DeletedAt:    key.DeletedAt,
		ValidFrom:    key.ValidFrom,
	})
}

func (api *UserAPI) UpdateUserKey(c *gin.Context) {
	uid := c.GetInt("userID")
	id, ok := c.Params.Get("id")
	if !ok {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidQueryParam)
		return
	}

	keyId, err := strconv.Atoi(id)
	if !ok {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidQueryParam)
		return
	}

	var req engine.UpdateKeyRequest
	if err := c.BindJSON(&req); err != nil {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidRequestBody)
		return
	}

	updated, err := api.repo.UpdateUserApiKey(uid, int(keyId), &req)
	if err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	c.JSON(http.StatusOK, updated)

}

func (api *UserAPI) DeleteUserKey(c *gin.Context) {
	uid := c.GetInt("userID")
	id, ok := c.Params.Get("id")
	if !ok {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidQueryParam)
		return
	}

	keyId, err := strconv.Atoi(id)
	if !ok {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidQueryParam)
		return
	}

	err = api.repo.DeleteUserApiKey(uid, int(keyId))
	if err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	c.Status(http.StatusNoContent)
}
