package api

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
	"github.com/msalopek/animus/engine"
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

func (api *AnimusAPI) GetUserKeys(c *gin.Context) {
	uid := c.GetInt("userID")
	ctx := repo.QueryCtxFromGin(c)
	keys, err := api.repo.GetUserApiKeys(ctx, uid)
	if err != nil {
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, keys)
}

// TODO: restrict number of keys
func (api *AnimusAPI) CreateUserKey(c *gin.Context) {
	uid := c.GetInt("userID")

	createdAt := time.Now()
	secret := uniuri.NewLen(32)
	key := model.Key{
		UserID:       int64(uid),
		ClientKey:    uniuri.NewLen(32),
		ClientSecret: secret,
		Rights:       model.ClientAccessRead,
		CreatedAt:    createdAt,
		ValidFrom:    createdAt,
	}
	res := api.repo.Create(&key)
	if res.Error != nil {
		// TODO: don't leak DB errors
		abortWithError(c, http.StatusInternalServerError, res.Error.Error())
		return
	}

	c.JSON(http.StatusCreated, engine.CreateKeyResponse{
		ID:        key.ID,
		UserID:    key.UserID,
		ClientKey: key.ClientKey,
		// expose the secret to the user after creation
		// the secret cannot be accessed after it is created
		ClientSecret: key.ClientSecret,
		Rights:       key.Rights,
		Disabled:     key.Disabled,
		CreatedAt:    key.CreatedAt,
		DeletedAt:    key.DeletedAt,
		ValidFrom:    key.ValidFrom,
		ValidTo:      key.ValidTo,
	})
}

func (api *AnimusAPI) UpdateUserKey(c *gin.Context) {
	uid := c.GetInt("userID")
	id, ok := c.Params.Get("id")
	if !ok {
		abortWithError(c, http.StatusBadRequest, engine.ErrInvalidQueryParam)
		return
	}

	keyId, err := strconv.Atoi(id)
	if !ok {
		abortWithError(c, http.StatusBadRequest, engine.ErrInvalidQueryParam)
		return
	}

	var req engine.UpdateKeyRequest
	if err := c.BindJSON(&req); err != nil {
		abortWithError(c, http.StatusBadRequest, engine.ErrInvalidRequestBody)
		return
	}

	updated, err := api.repo.UpdateUserApiKey(uid, int(keyId), &req)
	if err != nil {
		// TODO: don't leak DB errors
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, updated)

}

func (api *AnimusAPI) DeleteUserKey(c *gin.Context) {
	uid := c.GetInt("userID")
	id, ok := c.Params.Get("id")
	if !ok {
		abortWithError(c, http.StatusBadRequest, engine.ErrInvalidQueryParam)
		return
	}

	keyId, err := strconv.Atoi(id)
	if !ok {
		abortWithError(c, http.StatusBadRequest, engine.ErrInvalidQueryParam)
		return
	}

	err = api.repo.DeleteUserApiKey(uid, int(keyId))
	if err != nil {
		// TODO: don't leak DB errors
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (api *AnimusAPI) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
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
