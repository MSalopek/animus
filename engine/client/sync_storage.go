package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/msalopek/animus/engine"
	"github.com/msalopek/animus/model"
	"github.com/msalopek/animus/storage"
	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
)

const (
	SYNC_FAILED  = "failed"
	SYNC_SUCCESS = "success"
)

func makeResponse(s *model.Storage, status, retryUrl string, err error) engine.SyncAddFileResponse {
	r := engine.SyncAddFileResponse{
		Object:   *s,
		Status:   status,
		RetryUrl: retryUrl,
	}

	if err != nil {
		r.Error = err.Error()
	}

	return r
}

// SyncUploadFile extracts a file from gin.Context (multipart form),
// synchronously uploads it to default storage bucket and pins it to IPFS node.
// If the multipart form contains "meta" field it will be parsed and persisted if valid.
// Max file size is equal to gin defaults -> should be 32MB.
func (api *ClientAPI) SyncUploadFile(c *gin.Context) {
	ctxUID := c.GetInt("userID")

	file, err := c.FormFile("file")
	if err != nil {
		api.logger.WithFields(log.Fields{
			"error":  err,
			"method": "SYNC UPLOAD FILE",
		}).Error("could not parse FormFile")
		engine.AbortErrWithStatusFailed(c, http.StatusBadRequest, engine.ErrUnprocessableFormFile)
		return
	}

	meta := c.PostForm("meta")
	if meta != "" {
		// parse to check JSON validity
		check := make(map[string]interface{})
		if err := c.BindJSON(&check); err != nil {
			engine.AbortErrWithStatusFailed(c, http.StatusBadRequest, engine.ErrInvalidMeta)
			return
		}
	}

	objPath := fmt.Sprintf("%d/%s", ctxUID, file.Filename)
	info, err := storage.UploadFile(api.storage, file, api.cfg.Bucket, objPath)
	if err != nil {
		api.logger.WithFields(log.Fields{
			"error":  err,
			"method": "SYNC UPLOAD FILE",
		}).Error("could not upload file")
		engine.AbortErrWithStatusFailed(c, http.StatusInternalServerError, engine.ErrFileSaveFailed)
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
		api.logger.WithFields(log.Fields{
			"error":     res.Error,
			"storageID": storage.ID,
			"bucket":    storage.StorageBucket,
			"key":       storage.StorageKey,
			"method":    "SYNC UPLOAD FILE",
		}).Error("could no repo save AFTER UPLOAD")
		engine.AbortErrWithStatusFailed(c, http.StatusInternalServerError, engine.ErrFileSaveFailed)
		return
	}

	cid, err := api.pinFile(storage)
	if err != nil {
		c.Abort()
		c.JSON(
			http.StatusInternalServerError,
			makeResponse(storage,
				SYNC_FAILED,
				fmt.Sprintf(`/auth/storage/pin/id/%d`, storage.ID),
				engine.ErrPinFailed),
		)
		return
	}

	storage.Cid = &cid
	storage.Pinned = true
	stage = model.UploadStageIPFS
	storage.UploadStage = &stage
	if res := api.repo.Save(storage); res.Error != nil {
		api.logger.WithFields(log.Fields{
			"error":     res.Error,
			"storageID": storage.ID,
			"bucket":    storage.StorageBucket,
			"key":       storage.StorageKey,
			"method":    "SYNC UPLOAD FILE",
		}).Error("could not repo save AFTER PIN")
		c.Abort()
		c.JSON(
			http.StatusInternalServerError,
			makeResponse(storage,
				SYNC_FAILED,
				fmt.Sprintf(`/auth/storage/pin/id/%d`, storage.ID),
				engine.ErrPinFailed),
		)
		return
	}

	c.JSON(http.StatusCreated, makeResponse(storage, "success", "", nil))
}

// returns CID and error
func (api *ClientAPI) pinFile(s *model.Storage) (string, error) {
	ctx := context.Background()
	var hash, tmp string
	var err error

	bucket := api.cfg.Bucket
	storageKey := s.StorageKey
	if storageKey == nil {
		return "", fmt.Errorf("missing storage key for record: %d", s.ID)
	}

	if s.Dir {
		tmp, err = api.storage.DownloadDir(ctx, bucket, *storageKey)
		if err != nil {
			return hash, err
		}

		hash, err := api.sh.AddDir(tmp)
		if err != nil {
			return hash, err
		}

		err = os.RemoveAll(tmp)
		if err != nil {
			api.logger.WithFields(log.Fields{
				"temp":    tmp,
				"message": "error removing temp IPFS dir",
			}).Error(err)
		}
		return hash, nil
	}

	obj, err := api.storage.GetObject(ctx, bucket, *storageKey, minio.GetObjectOptions{})
	if err != nil {
		api.logger.WithFields(log.Fields{
			"request": fmt.Sprintf("%#v", s),
			"message": "error streaming object",
		}).Error(err)
		return "", err
	}
	return api.sh.Add(obj)
}
