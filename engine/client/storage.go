package client

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
	"gorm.io/datatypes"
	"gorm.io/gorm"
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
		Metadata:      datatypes.JSON(meta),
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
// If the multipart form contains "meta" field it will be parsed and persisted if valid.
func (api *ClientAPI) UploadDir(c *gin.Context) {
	ctxUID := c.GetInt("userID")

	dirname := c.PostForm("name")
	if dirname == "" {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrMissingFormDirName)
		return
	}

	m := c.PostForm("meta")
	var meta map[string]interface{}
	if m != "" {
		// parse to check JSON validity
		if err := json.Unmarshal([]byte(m), &meta); err != nil {
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

func (api *ClientAPI) GetStorageRecords(c *gin.Context) {
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

func (api *ClientAPI) GetStorageRecord(c *gin.Context) {
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

// DeleteStorageRecord will unpin record from IPFS and delete any storage objects.
func (api *ClientAPI) DeleteStorageRecord(c *gin.Context) {
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

func (api *ClientAPI) GetStorageRecordByCid(c *gin.Context) {
	uid := c.GetInt("userID")

	cid, ok := c.Params.Get("cid")
	if !ok {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidQueryParam)
		return
	}

	storage, err := api.repo.GetUserUploadByCid(uid, cid)
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
func (api *ClientAPI) RequestPin(c *gin.Context) {
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
func (api *ClientAPI) RequestUnpin(c *gin.Context) {
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

	err = api.publisher.Publish(api.cfg.NsqPinnerTopic, body)
	if err != nil {
		return err
	}
	return nil
}

func (api *ClientAPI) publishUnpinRequest(m *model.Storage) error {
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

	err = api.publisher.Publish(api.cfg.NsqPinnerTopic, body)
	if err != nil {
		return err
	}
	return nil
}
