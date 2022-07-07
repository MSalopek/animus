package api

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msalopek/animus/engine/repo"
	"github.com/msalopek/animus/model"
	"github.com/msalopek/animus/queue"
	"github.com/msalopek/animus/storage"
)

// UploadFile extracts a file from gin.Context (multipart form),
// uploads it to default storage bucket and publishes a PinRequest message.
func (api *AnimusAPI) UploadFile(c *gin.Context) {
	ctxUID := c.GetInt("userID")

	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}

	objPath := fmt.Sprintf("%d/%s", ctxUID, file.Filename)
	info, err := api.uploadFile(file, objPath)
	if err != nil {
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	stage := model.UploadStageStorage
	storage := &model.Storage{
		UserID:        int64(ctxUID),
		Name:          file.Filename,
		StorageBucket: &info.Bucket,
		StorageKey:    &info.Key,
		UploadStage:   &stage,
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

	if err := api.publishPinRequest(storage); err != nil {
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, storage)
}

// UploadDir extracts a files from gin.Context (multipart form),
// uploads them it to default storage bucket and publishes a PinRequest message.
func (api *AnimusAPI) UploadDir(c *gin.Context) {
	ctxUID := c.GetInt("userID")

	dirname := c.PostForm("name")
	if dirname == "" {
		abortWithError(c, http.StatusBadRequest, "missing directory name")
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	meta := storage.Uploads{}
	files := form.File["files"]
	for _, f := range files {
		objPath := fmt.Sprintf("%d/%s/%s", ctxUID, dirname, f.Filename)
		info, err := api.uploadFile(f, objPath)
		if err != nil {
			abortWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		meta.Uploads = append(meta.Uploads, info)
	}

	if len(meta.Uploads) < 1 {
		abortWithError(c, http.StatusInternalServerError, "something went wrong")
		return
	}

	// key is faked because dirs are not objects in s3 storage
	key := fmt.Sprintf("%d/%s", ctxUID, dirname)
	stage := model.UploadStageStorage
	storage := &model.Storage{
		UserID:        int64(ctxUID),
		Name:          dirname,
		Dir:           true,
		StorageBucket: &meta.Uploads[0].Bucket,
		StorageKey:    &key,
		UploadStage:   &stage,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if buf, err := json.Marshal(meta); err == nil {
		storage.Metadata = buf
	}

	if res := api.repo.Save(storage); res.Error != nil {
		abortWithError(c, http.StatusInternalServerError, res.Error.Error())
		return
	}

	if err := api.publishPinRequest(storage); err != nil {
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, storage)
}

func (api *AnimusAPI) GetUserUploads(c *gin.Context) {
	uid := c.GetInt("userID")

	ctx := repo.QueryCtxFromGin(c)
	storage, err := api.repo.GetUserUploads(ctx, uid)
	if err != nil {
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, storage)
}

func (api *AnimusAPI) uploadFile(file *multipart.FileHeader, objName string) (*storage.Upload, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// TODO: add deadline
	ctx := context.Background()
	return api.storage.StreamFile(ctx, api.cfg.Bucket, objName, src, file.Size, storage.Opts{})
}

func (api *AnimusAPI) publishPinRequest(m *model.Storage) error {
	pr := queue.PinRequest{
		StorageID: int(m.ID),
		Key:       *m.StorageKey,
		Dir:       m.Dir,
	}
	body, err := json.Marshal(pr)
	if err != nil {
		return err
	}

	err = api.publisher.Publish(api.cfg.NsqTopic, body)
	if err != nil {
		return err
	}
	return nil
}
