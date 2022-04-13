package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msalopek/animus/engine"
)

// UploadFile extracts a file from gin.Context (multipart form)
// and synchronously uploads it to the attached IPFS node.
func (api *HttpAPI) UploadFile(c *gin.Context) {
	ctxUID := c.GetInt("userID")

	meta := c.PostForm("meta")
	var parsed map[string]interface{}
	// TODO: currently this is just unmarshalled to validate JSON
	if err := json.Unmarshal([]byte(meta), &parsed); err != nil {
		abortWithError(c, http.StatusBadRequest, engine.ErrInvalidMeta)
		return
	}

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}

	src, err := file.Open()
	if err != nil {
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer src.Close()

	cid, err := api.IPFSUpload(src)
	if err != nil {
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	upload, err := api.repo.GetUserUploadByCid(ctxUID, cid)
	if err != nil {
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// create only if not exists to avoid duplication
	if upload.ID == 0 {
		upload.UserID = int64(ctxUID)
		upload.Name = file.Filename
		upload.Metadata = meta
		upload.Cid = &cid
		upload.CreatedAt = time.Now()
		upload.UpdatedAt = upload.CreatedAt
		upload.Pinned = true // TODO: make configurable

		if res := api.repo.Save(upload); res.Error != nil {
			abortWithError(c, http.StatusInternalServerError, res.Error.Error())
			return
		}
	}

	c.JSON(http.StatusOK, upload)
}

func (api *HttpAPI) GetUserUploads(c *gin.Context) {
	ctxUID := c.GetInt("userID")
	// TODO: paginate with limit, offset
	storage, err := api.repo.GetUserUploads(ctxUID)
	if err != nil {
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, storage)
}

func (api *HttpAPI) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
