package repository

import (
	"context"

	"redpacket/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	CreateRedPacket(ctx context.Context, rp *model.RedPacket) error
	GetRedPacketByBizID(ctx context.Context, bizID string) (*model.RedPacket, error)
	GetRedPacketByPacketID(ctx context.Context, packetID string) (*model.RedPacket, error)
	UpdateRedPacketTxHash(ctx context.Context, bizID, txHash, packetID string) error
	CreateClaimAuth(ctx context.Context, auth *model.RedPacketClaimAuth) error
	GetClaimAuth(ctx context.Context, packetID, claimer string) (*model.RedPacketClaimAuth, error)
	MarkClaimAuthUsed(ctx context.Context, authNonce string) error
	CreateClaim(ctx context.Context, claim *model.RedPacketClaim) error
	GetClaimsByPacketID(ctx context.Context, packetID string) ([]model.RedPacketClaim, error)
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateRedPacket(ctx context.Context, rp *model.RedPacket) error {
	return r.db.WithContext(ctx).Create(rp).Error
}

func (r *repository) GetRedPacketByBizID(ctx context.Context, bizID string) (*model.RedPacket, error) {
	var rp model.RedPacket
	err := r.db.WithContext(ctx).Where("biz_id = ?", bizID).First(&rp).Error
	return &rp, err
}

func (r *repository) GetRedPacketByPacketID(ctx context.Context, packetID string) (*model.RedPacket, error) {
	var rp model.RedPacket
	err := r.db.WithContext(ctx).Where("packet_id = ?", packetID).First(&rp).Error
	return &rp, err
}

func (r *repository) UpdateRedPacketTxHash(ctx context.Context, bizID, txHash, packetID string) error {
	return r.db.WithContext(ctx).Model(&model.RedPacket{}).
		Where("biz_id = ?", bizID).
		Updates(map[string]interface{}{
			"tx_hash":   txHash,
			"packet_id": packetID,
			"status":    "ACTIVE",
		}).Error
}

func (r *repository) CreateClaimAuth(ctx context.Context, auth *model.RedPacketClaimAuth) error {
	return r.db.WithContext(ctx).Create(auth).Error
}

func (r *repository) GetClaimAuth(ctx context.Context, packetID, claimer string) (*model.RedPacketClaimAuth, error) {
	var auth model.RedPacketClaimAuth
	err := r.db.WithContext(ctx).Where("packet_id = ? AND claimer = ? AND used = false", packetID, claimer).First(&auth).Error
	return &auth, err
}

func (r *repository) MarkClaimAuthUsed(ctx context.Context, authNonce string) error {
	return r.db.WithContext(ctx).Model(&model.RedPacketClaimAuth{}).
		Where("auth_nonce = ?", authNonce).
		Update("used", true).Error
}

func (r *repository) CreateClaim(ctx context.Context, claim *model.RedPacketClaim) error {
	return r.db.WithContext(ctx).Create(claim).Error
}

func (r *repository) GetClaimsByPacketID(ctx context.Context, packetID string) ([]model.RedPacketClaim, error) {
	var claims []model.RedPacketClaim
	err := r.db.WithContext(ctx).Where("packet_id = ?", packetID).Order("created_at desc").Find(&claims).Error
	return claims, err
}
