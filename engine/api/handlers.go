package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/msalopek/animus/engine"
	"github.com/msalopek/animus/engine/repo"
	"gorm.io/gorm"
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
	uid := c.GetInt("userID")
	ctx := repo.QueryCtxFromGin(c)
	storage, err := api.repo.GetUserUploads(ctx, uid)
	if err != nil {
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, storage)
}

func (api *HttpAPI) ListsDir(c *gin.Context) {
	uid := c.GetInt("userID")
	cid := c.Param("cid")
	if cid == "" {
		abortWithError(c, http.StatusBadRequest, engine.ErrInvalidQueryParam)
	}
	storage, err := api.repo.GetUserUploadByCid(uid, cid)
	if err != nil && err == gorm.ErrRecordNotFound {
		abortWithError(c, http.StatusNotFound, engine.ErrNotFound)
		return
	}
	if err != nil {
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if !storage.Dir {
		abortWithError(c, http.StatusBadRequest, engine.ErrNotADirectory)
		return
	}

	data, err := api.ipfs.List(*storage.Cid)
	if err != nil {
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, ListResp{data})
}

func (api *HttpAPI) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

type ListResp struct {
	Objects []*shell.LsLink `json:"Objects"`
}
