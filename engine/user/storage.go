package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msalopek/animus/engine"
	"github.com/msalopek/animus/engine/repo"
	"github.com/msalopek/animus/model"
	"github.com/msalopek/animus/queue"
	"github.com/msalopek/animus/storage"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

	// NOTE: the file is on s3 and can be manually pinned on error
	if err := api.publishPinRequest(storage); err != nil {
		api.logger.WithFields(log.Fields{
			"error":     err.Error(),
			"storageID": storage.ID,
			"bucket":    storage.StorageBucket,
			"key":       storage.StorageKey,
			"dir":       storage.Dir,
		}).Error("could not publish pin request")
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

	// NOTE: the file is on s3 and can be manually pinned on error
	if err := api.publishPinRequest(storage); err != nil {
		api.logger.WithFields(log.Fields{
			"error":     err.Error(),
			"storageID": storage.ID,
			"bucket":    storage.StorageBucket,
			"key":       storage.StorageKey,
			"dir":       storage.Dir,
		}).Error("could not publish pin request")
	}

	c.JSON(http.StatusCreated, storage)
}

func (api *UserAPI) GetUserUploads(c *gin.Context) {
	uid := c.GetInt("userID")

	ctx := repo.QueryCtxFromGin(c)
	rows, err := api.repo.GetCountedUserUploads(ctx, uid)
	if err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	resp := engine.GetStorageResponse{
		Returned: len(rows),
		Rows:     rows,
	}
	// total record count is the same for each returned record
	if len(rows) > 0 {
		resp.Total = rows[0].TotalRows
	}

	c.JSON(http.StatusOK, resp)
}

func (api *UserAPI) GetStorageRecord(c *gin.Context) {
	uid := c.GetInt("userID")

	idParam, ok := c.Params.Get("id")
	if !ok {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidQueryParam)
		return
	}
	id, err := strconv.Atoi(idParam)
	if !ok {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidQueryParam)
		return
	}

	storage, err := api.repo.GetUserUploadByID(uid, id)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		engine.AbortErr(c, http.StatusNotFound, engine.ErrNotFound)
		return
	} else if err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	c.JSON(http.StatusOK, storage)
}

// Send a pin request if the record is not already pinned.
// Pinning may be attempted again by adding "forced" query param.
func (api *UserAPI) RequestPin(c *gin.Context) {
	uid := c.GetInt("userID")

	idParam, ok := c.Params.Get("id")
	if !ok {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidQueryParam)
		return
	}
	id, err := strconv.Atoi(idParam)
	if !ok {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidQueryParam)
		return
	}

	// force on any truthy value
	_, force := c.GetQuery("force")

	storage, err := api.repo.GetUserUploadByID(uid, id)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		engine.AbortErr(c, http.StatusNotFound, engine.ErrNotFound)
		return
	} else if err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	if !force && storage.Pinned {
		c.Status(http.StatusOK)
		return
	}

	if err := api.publishPinRequest(storage); err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	c.Status(http.StatusOK)
}

// Send an unpin request if the record is pinned.
// Unpinning may be attempted again by adding "forced" query param.
// Attempting to unpin a recored without CID will returns "CID is missing" error.
func (api *UserAPI) RequestUnpin(c *gin.Context) {
	uid := c.GetInt("userID")

	idParam, ok := c.Params.Get("id")
	if !ok {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidQueryParam)
		return
	}
	id, err := strconv.Atoi(idParam)
	if !ok {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidQueryParam)
		return
	}

	// force on any truthy value
	_, force := c.GetQuery("force")

	storage, err := api.repo.GetUserUploadByID(uid, id)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		engine.AbortErr(c, http.StatusNotFound, engine.ErrNotFound)
		return
	} else if err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	if storage.Cid == nil {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrNoCID)
		return
	}

	if !force && !storage.Pinned {
		c.Status(http.StatusOK)
		return
	}

	if err := api.publishUnpinRequest(storage); err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	c.Status(http.StatusOK)
}

// DeleteStorageRecord will unpin record from IPFS and delete any storage objects.
func (api *UserAPI) DeleteStorageRecord(c *gin.Context) {
	uid := c.GetInt("userID")

	idParam, ok := c.Params.Get("id")
	if !ok {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidQueryParam)
		return
	}
	id, err := strconv.Atoi(idParam)
	if !ok {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidQueryParam)
		return
	}

	storage, err := api.repo.GetUserUploadByID(uid, id)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		engine.AbortErr(c, http.StatusNotFound, engine.ErrNotFound)
		return
	} else if err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	if storage.Pinned {
		if err := api.publishUnpinRequest(storage); err != nil {
			engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
			return
		}
	}

	path := fmt.Sprintf("%d/%s", uid, storage.Name)
	if err := api.storage.RemoveDirObjects(context.Background(), api.cfg.Bucket, path); err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	if err := api.repo.DeleteUserUploadById(uid, id); err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	c.Status(http.StatusNoContent)
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

func (api *UserAPI) publishUnpinRequest(m *model.Storage) error {
	if m.Cid == nil {
		return engine.ErrNoCID
	}

	pr := queue.PinRequest{
		StorageID: int(m.ID),
		CID:       *m.Cid,
		Unpin:     true,
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
