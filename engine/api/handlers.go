package api

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msalopek/animus/engine"
	"github.com/msalopek/animus/model"
)

// TODO: remove
func WIPresponder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "method is under construction",
	})
}

func (api *HttpAPI) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
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
	email := c.GetString("email")
	if email == "" {
		abortWithError(c, http.StatusUnauthorized, engine.ErrUnauthorized)
		return
	}

	user, err := api.repo.GetUserByEmail(email)
	if err != nil {
		// return forbidden even if user does not exist
		// if execution reaches this point the request is authorized
		abortWithError(c, http.StatusForbidden, engine.ErrForbidden)
		return
	}

	meta := c.PostForm("meta")
	var parsed map[string]interface{}
	// TODO: currently this is just unmarshalled to validate JSON
	err = json.Unmarshal([]byte(meta), &parsed)
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
	storeRec := model.Storage{
		UserID:   user.ID,
		Name:     filename,
		Metadata: meta,
	}
	err = api.repo.CreateStorage(&storeRec)
	if err != nil {
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

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

	storeRec.Cid = &cid
	// TODO: make this configurable;
	// IPFS node does not need to autopin files
	storeRec.Pinned = true
	storeRec.UpdatedAt = time.Now()
	// TODO: check err
	api.repo.Save(&storeRec)

	c.JSON(http.StatusOK, gin.H{
		"name": filename,
		"CID":  cid,
		"meta": meta,
	})
}
