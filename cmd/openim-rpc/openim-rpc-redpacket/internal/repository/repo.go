package repository

import (
	"context"
	"math/big"

	"redpacket/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	CreateRedPacket(ctx context.Context, rp *model.RedPacket) error
	GetRedPacketByBizID(ctx context.Context, bizID string) (*model.RedPacket, error)
	GetRedPacketByPacketID(ctx context.Context, packetID string) (*model.RedPacket, error)
	UpdateRedPacketCreated(ctx context.Context, rp *model.RedPacket) error
	UpdateRedPacketStatus(ctx context.Context, packetID, status string) error
	UpdateRedPacketClaimProgress(ctx context.Context, packetID, claimedAmount, status string) error
	CreateClaimAuth(ctx context.Context, auth *model.RedPacketClaimAuth) error
	GetClaimAuth(ctx context.Context, packetID, claimer string) (*model.RedPacketClaimAuth, error)
	MarkClaimAuthUsed(ctx context.Context, authNonce string) error
	GetClaimByPacketIDAndClaimer(ctx context.Context, packetID, claimer string) (*model.RedPacketClaim, error)
	GetClaimByPacketIDAndUserID(ctx context.Context, packetID, userID string) (*model.RedPacketClaim, error)
	SaveClaim(ctx context.Context, claim *model.RedPacketClaim) error
	GetClaimsByPacketID(ctx context.Context, packetID string) ([]model.RedPacketClaim, error)
	SaveRefund(ctx context.Context, refund *model.RedPacketRefund) error
	CreateWalletBindingChallenge(ctx context.Context, challenge *model.WalletBindingChallenge) error
	GetWalletBindingChallenge(ctx context.Context, challengeID string) (*model.WalletBindingChallenge, error)
	UpdateWalletBindingChallenge(ctx context.Context, challenge *model.WalletBindingChallenge) error
	UpsertWalletBinding(ctx context.Context, binding *model.WalletBinding) error
	GetActiveWalletBinding(ctx context.Context, userID, chainType, walletAddress string) (*model.WalletBinding, error)
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

func (r *repository) UpdateRedPacketCreated(ctx context.Context, rp *model.RedPacket) error {
	return r.db.WithContext(ctx).Model(&model.RedPacket{}).
		Where("biz_id = ?", rp.BizID).
		Updates(map[string]interface{}{
			"chain_type":        rp.ChainType,
			"packet_id":         rp.PacketID,
			"tx_hash":           rp.TxHash,
			"chain_id":          rp.ChainID,
			"contract_address":  rp.ContractAddress,
			"group_id":          rp.GroupID,
			"scope_type":        rp.ScopeType,
			"receiver_user_id":  rp.ReceiverUserID,
			"receiver_user_ids": rp.ReceiverUserIDs,
			"status":            rp.Status,
		}).Error
}

func (r *repository) UpdateRedPacketStatus(ctx context.Context, packetID, status string) error {
	return r.db.WithContext(ctx).Model(&model.RedPacket{}).
		Where("packet_id = ?", packetID).
		Update("status", status).Error
}

func (r *repository) UpdateRedPacketClaimProgress(ctx context.Context, packetID, claimedAmount, status string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var rp model.RedPacket
		if err := tx.Where("packet_id = ?", packetID).First(&rp).Error; err != nil {
			return err
		}

		totalClaimed := addNumericStrings(rp.ClaimedAmount, claimedAmount)
		nextShares := rp.ClaimedShares + 1

		updates := map[string]interface{}{
			"claimed_amount": totalClaimed,
			"claimed_shares": nextShares,
			"updated_at":     gorm.Expr("CURRENT_TIMESTAMP"),
		}
		if status != "" {
			updates["status"] = status
		}

		return tx.Model(&model.RedPacket{}).
			Where("id = ?", rp.ID).
			Updates(updates).Error
	})
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

func (r *repository) GetClaimByPacketIDAndClaimer(ctx context.Context, packetID, claimer string) (*model.RedPacketClaim, error) {
	var claim model.RedPacketClaim
	err := r.db.WithContext(ctx).
		Where("packet_id = ? AND claimer_wallet = ?", packetID, claimer).
		Order("created_at desc").
		First(&claim).Error
	return &claim, err
}

func (r *repository) GetClaimByPacketIDAndUserID(ctx context.Context, packetID, userID string) (*model.RedPacketClaim, error) {
	var claim model.RedPacketClaim
	err := r.db.WithContext(ctx).
		Where("packet_id = ? AND user_id = ?", packetID, userID).
		Order("created_at desc").
		First(&claim).Error
	return &claim, err
}

func (r *repository) SaveClaim(ctx context.Context, claim *model.RedPacketClaim) error {
	if claim.UserID != "" {
		var existing model.RedPacketClaim
		err := r.db.WithContext(ctx).
			Where("packet_id = ? AND user_id = ?", claim.PacketID, claim.UserID).
			First(&existing).Error
		if err == nil {
			claim.ID = existing.ID
			return r.db.WithContext(ctx).Model(&model.RedPacketClaim{}).
				Where("id = ?", existing.ID).
				Updates(map[string]interface{}{
					"claimer_wallet": existing.ClaimerWallet,
					"auth_nonce":     claim.AuthNonce,
					"claim_tx_hash":  claim.ClaimTxHash,
					"claimed_amount": claim.ClaimedAmount,
					"block_number":   claim.BlockNumber,
					"status":         claim.Status,
					"updated_at":     claim.UpdatedAt,
				}).Error
		}
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
	}

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "claim_tx_hash"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"user_id",
			"packet_id",
			"claimer_wallet",
			"auth_nonce",
			"claimed_amount",
			"block_number",
			"status",
			"updated_at",
		}),
	}).Create(claim).Error
}

func (r *repository) GetClaimsByPacketID(ctx context.Context, packetID string) ([]model.RedPacketClaim, error) {
	var claims []model.RedPacketClaim
	err := r.db.WithContext(ctx).Where("packet_id = ?", packetID).Order("created_at desc").Find(&claims).Error
	return claims, err
}

func (r *repository) SaveRefund(ctx context.Context, refund *model.RedPacketRefund) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "tx_hash"}},
		DoNothing: true,
	}).Create(refund).Error
}

func (r *repository) CreateWalletBindingChallenge(ctx context.Context, challenge *model.WalletBindingChallenge) error {
	return r.db.WithContext(ctx).Create(challenge).Error
}

func (r *repository) GetWalletBindingChallenge(ctx context.Context, challengeID string) (*model.WalletBindingChallenge, error) {
	var challenge model.WalletBindingChallenge
	err := r.db.WithContext(ctx).Where("challenge_id = ?", challengeID).First(&challenge).Error
	return &challenge, err
}

func (r *repository) UpdateWalletBindingChallenge(ctx context.Context, challenge *model.WalletBindingChallenge) error {
	return r.db.WithContext(ctx).Model(&model.WalletBindingChallenge{}).
		Where("challenge_id = ?", challenge.ChallengeID).
		Updates(map[string]interface{}{
			"status":      challenge.Status,
			"signature":   challenge.Signature,
			"verified_at": challenge.VerifiedAt,
			"updated_at":  challenge.UpdatedAt,
		}).Error
}

func (r *repository) UpsertWalletBinding(ctx context.Context, binding *model.WalletBinding) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "chain_type"},
			{Name: "wallet_address"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"chain_id",
			"status",
			"challenge_id",
			"verified_at",
			"revoked_at",
			"updated_at",
		}),
	}).Create(binding).Error
}

func (r *repository) GetActiveWalletBinding(ctx context.Context, userID, chainType, walletAddress string) (*model.WalletBinding, error) {
	var binding model.WalletBinding
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND chain_type = ? AND wallet_address = ? AND status = ?", userID, chainType, walletAddress, "ACTIVE").
		First(&binding).Error
	return &binding, err
}

func addNumericStrings(current, delta string) string {
	left := new(big.Int)
	if current != "" {
		left.SetString(current, 10)
	}
	right := new(big.Int)
	if delta != "" {
		right.SetString(delta, 10)
	}
	return new(big.Int).Add(left, right).String()
}
