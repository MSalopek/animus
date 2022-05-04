package api

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msalopek/animus/engine"
	"github.com/msalopek/animus/engine/repo"
	"github.com/msalopek/animus/model"
	"github.com/msalopek/animus/storage"
)

// UploadFile extracts a file from gin.Context (multipart form)
// and synchronously uploads it to default storage bucket.
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

	objPath := fmt.Sprintf("%s/%s", ctxUID, file.Filename)
	info, err := api.uploadFile(file, objPath)
	if err != nil {
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	storage := &model.Storage{
		UserID:        int64(ctxUID),
		Name:          file.Filename,
		StorageBucket: &info.Bucket,
		StorageKey:    &info.Key,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if buf, err := json.Marshal(info); err == nil {
		storage.Metadata = buf
	}

	if res := api.repo.Save(storage); res.Error != nil {
		abortWithError(c, http.StatusInternalServerError, res.Error.Error())
		return
	}

	c.JSON(http.StatusCreated, storage)
}

// UploadDir extracts a files from gin.Context (multipart form)
// and synchronously uploads them it to default storage bucket.
func (api *HttpAPI) UploadDir(c *gin.Context) {
	ctxUID := c.GetInt("userID")

	dirname := c.PostForm("name")
	if dirname == "" {
		abortWithError(c, http.StatusBadRequest, "missing directory name")
	}

	form, err := c.MultipartForm()
	if err != nil {
		abortWithError(c, http.StatusBadRequest, res.Error.Error())
		return
	}

	meta := storage.UploadInfo{}
	files := form.File["upload[]"]
	for _, f := range files {
		objPath := fmt.Sprintf("%s/%s/%s", ctxUID, dirname, f.Filename)
		info, err := api.uploadFile(f, objPath)
		if err != nil {
			abortWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		meta = append(meta, info)
	}

	storage := &model.Storage{
		UserID:        int64(ctxUID),
		Name:          dirname,
		Dir:           true,
		StorageBucket: &meta.Uploads[0].Bucket,
		// key is faked because dirs don't have keys in s3 compatible storage
		StorageKey: fmt.Sprintf("%s/%s", ctxUID, dirname),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if buf, err := json.Marshal(info); err == nil {
		storage.Metadata = buf
	}

	if res := api.repo.Save(storage); res.Error != nil {
		abortWithError(c, http.StatusInternalServerError, res.Error.Error())
		return
	}

	c.JSON(http.StatusCreated, storage)
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

func (api *HttpAPI) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func (api *HttpAPI) uploadFile(file *multipart.FileHeader, objName string) (storage.UploadInfo, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// TODO: add deadline
	ctx := context.Background()
	return api.storage.StreamFile(ctx, api.cfg.Bucket, objName, src, file.Size, storage.Opts{})
}
