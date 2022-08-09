package user

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/msalopek/animus/engine"
	"github.com/msalopek/animus/model"
)

func (api *UserAPI) UpdateUser(c *gin.Context) {
	uid := c.GetInt("userID")

	var req engine.UpdateUserRequest
	if err := c.BindJSON(&req); err != nil {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidRequestBody)
		return
	}

	u := model.User{
		ID: int64(uid),
	}

	if req.Email != nil && *req.Email != "" {
		u.Email = *req.Email
	}
	if req.Firstname != nil && *req.Firstname != "" {
		u.Firstname = req.Firstname
	}
	if req.Lastname != nil && *req.Lastname != "" {
		u.Lastname = req.Lastname
	}
	if req.Username != nil && *req.Username != "" {
		u.Username = *req.Username
	}

	if req.WebhooksURL != nil {
		u.WebhooksURL = req.WebhooksURL

		// can be set to empty string but not to an invalid URL
		if *req.WebhooksURL != "" && !isValidUrl(*req.WebhooksURL) {
			engine.AbortErr(c, http.StatusBadRequest, engine.ErrBadRequest)
			return
		}
	}

	if req.WebhooksActive != nil {
		u.WebhooksActive = req.WebhooksActive
	}

	updated, err := api.repo.UpdateUser(uid, &u)
	if err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	c.JSON(http.StatusOK, updated)

}

// check if string is a valid url.
func isValidUrl(s string) bool {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}

	u, err := url.Parse(s)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}
