package client

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
	"gorm.io/datatypes"
)

// UploadFile extracts a file from gin.Context (multipart form),
// uploads it to default storage bucket and publishes a PinRequest message with Source == "client".
// If the multipart form contains "meta" field it will be parsed and persisted if valid.
// Max file size is equal to gin defaults -> should be 32MB.
func (api *ClientAPI) UploadFile(c *gin.Context) {
	ctxUID := c.GetInt("userID")

	file, err := c.FormFile("file")
	if err != nil {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrUnprocessableFormFile)
		return
	}

	meta := c.PostForm("meta")
	if meta != "" {
		// parse to check JSON validity
		check := make(map[string]interface{})
		if err := c.BindJSON(&check); err != nil {
			engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidMeta)
			return
		}
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
		// TODO: check if this actually works
		Metadata:  datatypes.JSON(meta),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
// If the multipart form contains "meta" field it will be parsed and persisted if valid.
func (api *ClientAPI) UploadDir(c *gin.Context) {
	ctxUID := c.GetInt("userID")

	dirname := c.PostForm("name")
	if dirname == "" {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrMissingFormDirName)
		return
	}

	m := c.PostForm("meta")
	meta := make(map[string]interface{})
	if m != "" {
		// parse to check JSON validity
		if err := c.BindJSON(&meta); err != nil {
			engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidMeta)
			return
		}
	}

	form, err := c.MultipartForm()
	if err != nil {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrUnprocessableMultipartForm)
		return
	}

	fup := storage.Uploads{}
	files := form.File["files"]
	for _, f := range files {
		objPath := fmt.Sprintf("%d/%s/%s", ctxUID, dirname, f.Filename)
		info, err := storage.UploadFile(api.storage, f, api.cfg.Bucket, objPath)
		if err != nil {
			engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
			return
		}
		fup.Uploads = append(fup.Uploads, info)
	}

	if len(fup.Uploads) < 1 {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}
	meta["uploads"] = fup.Uploads

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

	// NOTE: silently fail if meta is not processable
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

func (api *ClientAPI) GetUserUploads(c *gin.Context) {
	uid := c.GetInt("userID")

	ctx := repo.QueryCtxFromGin(c)
	storage, err := api.repo.GetUserUploads(ctx, uid)
	if err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	c.JSON(http.StatusOK, storage)
}

func (api *ClientAPI) publishPinRequest(m *model.Storage) error {
	pr := queue.PinRequest{
		StorageID: int(m.ID),
		Key:       *m.StorageKey,
		Dir:       m.Dir,
		Source:    queue.SourceClient,
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
