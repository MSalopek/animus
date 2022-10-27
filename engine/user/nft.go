package user

import (
	"bytes"
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
	"gorm.io/gorm"
)

func (api *UserAPI) GetUserNFTs(c *gin.Context) {
	uid := c.GetInt("userID")
	ctx := repo.QueryCtxFromGin(c)
	nfts, err := api.repo.GetUserNFTs(ctx, uid)
	if err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	c.JSON(http.StatusOK, nfts)
}

func (api *UserAPI) CreateUserNFT(c *gin.Context) {
	uid := c.GetInt("userID")

	var req model.NFT
	if err := c.BindJSON(&req); err != nil {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidRequestBody)
		return
	}

	createdAt := time.Now()
	nft := model.NFT{
		UserID:      int64(uid),
		StorageID:   req.StorageID,
		TokenID:     req.TokenID,
		ItemType:    req.ItemType,
		ExternalURL: req.ExternalURL,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}
	res := api.repo.Create(&nft)
	if res.Error != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	c.JSON(http.StatusCreated, nft)
}

func (api *UserAPI) SendNFTWebhook(c *gin.Context) {
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

	token, ok := c.GetQuery("token")
	if !ok {
		engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidQueryParam)
		return
	}

	nft, err := api.repo.GetUserNFTByStorageID(uid, id)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		engine.AbortErr(c, http.StatusNotFound, engine.ErrNotFound)
		return
	} else if err != nil {
		fmt.Println("################ GET", err)
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	buf, err := json.Marshal(engine.AfterMintWebhookRequest{
		ItemId:      nft.ExternalItemId,
		Type:        nft.ItemType,
		ExternalURL: nft.ExternalURL,
		TokenID:     nft.TokenID,
	})
	if err != nil {
		fmt.Println("################ MARSHAL", err)
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}
	r, err := http.NewRequest("POST",
		fmt.Sprintf("%s/v2/animus/nft/webhook", api.cfg.WebhookBaseURL),
		bytes.NewBuffer(buf))
	if err != nil {
		fmt.Println("################ POST", err)
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Println("################# DO", api.cfg.WebhookBaseURL, "DPN", err)
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	c.Status(http.StatusOK)
}
