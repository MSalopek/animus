package repo

import (
	"github.com/msalopek/animus/model"
	"gorm.io/gorm/clause"
)

func (rpo *Repo) GetUserNFTs(ctx QueryCtx, userID int) ([]*model.NFT, error) {
	var nfts []*model.NFT
	q := rpo.Where("user_id = ?", userID).
		Limit(ctx.Limit).
		Offset(ctx.Offset)
	if ctx.OrderBy != "" {
		q.Order(clause.OrderByColumn{
			Column: clause.Column{Name: ctx.OrderBy},
			Desc:   !ctx.Asc},
		)

	}
	res := q.Find(&nfts)
	if res.Error != nil {
		return nil, res.Error
	}
	return nfts, nil
}

func (rpo *Repo) GetUserNFTByStorageID(userID, storageID int) (*model.NFT, error) {
	var nft model.NFT
	res := rpo.Where("user_id = ? AND storage_id = ?", userID, storageID).First(&nft)
	if res.Error != nil {
		return nil, res.Error
	}
	return &nft, nil
}
