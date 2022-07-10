package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msalopek/animus/engine"
	"github.com/msalopek/animus/engine/repo"
	"github.com/msalopek/animus/model"
	"github.com/msalopek/animus/queue"
	"github.com/msalopek/animus/storage"
)

// UploadFile extracts a file from gin.Context (multipart form),
// uploads it to default storage bucket and publishes a PinRequest message.
func (api *UserAPI) UploadFile(c *gin.Context) {
	ctxUID := c.GetInt("userID")

	file, err := c.FormFile("file")
	if err != nil {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrUnprocessableFormFile)
		return
	}

	objPath := fmt.Sprintf("%d/%s", ctxUID, file.Filename)
	info, err := storage.UploadFile(api.storage, file, api.cfg.Bucket, objPath)
	if err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrFileSaveFailed)
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
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrFileSaveFailed)
		return
	}

	if err := api.publishPinRequest(storage); err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}
	c.JSON(http.StatusCreated, storage)
}

// UploadDir extracts files from gin.Context (multipart form),
// uploads them to default storage bucket and publishes a PinRequest message.
func (api *UserAPI) UploadDir(c *gin.Context) {
	ctxUID := c.GetInt("userID")

	dirname := c.PostForm("name")
	if dirname == "" {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrMissingFormDirName)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrUnprocessableMultipartForm)
		return
	}

	meta := storage.Uploads{}
	files := form.File["files"]
	for _, f := range files {
		objPath := fmt.Sprintf("%d/%s/%s", ctxUID, dirname, f.Filename)
		info, err := storage.UploadFile(api.storage, f, api.cfg.Bucket, objPath)
		if err != nil {
			engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
			return
		}
		meta.Uploads = append(meta.Uploads, info)
	}

	if len(meta.Uploads) < 1 {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	// key is faked because dirs are not objects in s3 storage
	key := fmt.Sprintf("%d/%s", ctxUID, dirname)
	stage := model.UploadStageStorage
	storage := &model.Storage{
		UserID:        int64(ctxUID),
		Name:          dirname,
		Dir:           true,
		StorageBucket: &api.cfg.Bucket,
		StorageKey:    &key,
		UploadStage:   &stage,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if buf, err := json.Marshal(meta); err == nil {
		storage.Metadata = buf
	}

	if res := api.repo.Save(storage); res.Error != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrDirSaveFailed)
		return
	}

	if err := api.publishPinRequest(storage); err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	c.JSON(http.StatusCreated, storage)
}

func (api *UserAPI) GetUserUploads(c *gin.Context) {
	uid := c.GetInt("userID")

	ctx := repo.QueryCtxFromGin(c)
	storage, err := api.repo.GetUserUploads(ctx, uid)
	if err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	c.JSON(http.StatusOK, storage)
}

func (api *UserAPI) publishPinRequest(m *model.Storage) error {
	pr := queue.PinRequest{
		StorageID: int(m.ID),
		Key:       *m.StorageKey,
		Dir:       m.Dir,
		Source:    queue.SourceUser,
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
