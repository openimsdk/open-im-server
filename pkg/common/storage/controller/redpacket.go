package controller

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

// RedPacketDatabase is a façade aggregating all redpacket-related collections.
// It mirrors the legacy Repository interface so the rpc service layer stays
// unaware of the underlying storage.
type RedPacketDatabase interface {
	CreateRedPacket(ctx context.Context, rp *model.RedPacket) error
	GetRedPacketByBizID(ctx context.Context, bizID string) (*model.RedPacket, error)
	GetRedPacketByPacketID(ctx context.Context, packetID string) (*model.RedPacket, error)
	UpdateRedPacketCreated(ctx context.Context, rp *model.RedPacket) error
	UpdateRedPacketStatus(ctx context.Context, packetID, status string) error
	UpdateRedPacketClaimProgress(ctx context.Context, packetID, claimedAmount, status string) error

	CreateClaimAuth(ctx context.Context, auth *model.RedPacketClaimAuth) error
	GetClaimAuth(ctx context.Context, packetID, claimer string) (*model.RedPacketClaimAuth, error)
	MarkClaimAuthUsed(ctx context.Context, authNonce string) error

	SaveClaim(ctx context.Context, claim *model.RedPacketClaim) error
	GetClaimByPacketIDAndClaimer(ctx context.Context, packetID, claimer string) (*model.RedPacketClaim, error)
	GetClaimByPacketIDAndUserID(ctx context.Context, packetID, userID string) (*model.RedPacketClaim, error)
	GetClaimsByPacketID(ctx context.Context, packetID string) ([]*model.RedPacketClaim, error)

	SaveRefund(ctx context.Context, refund *model.RedPacketRefund) error

	CreateWalletBindingChallenge(ctx context.Context, challenge *model.WalletBindingChallenge) error
	GetWalletBindingChallenge(ctx context.Context, challengeID string) (*model.WalletBindingChallenge, error)
	UpdateWalletBindingChallenge(ctx context.Context, challenge *model.WalletBindingChallenge) error

	UpsertWalletBinding(ctx context.Context, binding *model.WalletBinding) error
	GetActiveWalletBinding(ctx context.Context, userID, chainType, walletAddress string) (*model.WalletBinding, error)
}

type redPacketDatabase struct {
	rp        database.RedPacket
	claim     database.RedPacketClaim
	claimAuth database.RedPacketClaimAuth
	refund    database.RedPacketRefund
	challenge database.WalletBindingChallenge
	binding   database.WalletBinding
}

func NewRedPacketDatabase(
	rp database.RedPacket,
	claim database.RedPacketClaim,
	claimAuth database.RedPacketClaimAuth,
	refund database.RedPacketRefund,
	challenge database.WalletBindingChallenge,
	binding database.WalletBinding,
) RedPacketDatabase {
	return &redPacketDatabase{
		rp:        rp,
		claim:     claim,
		claimAuth: claimAuth,
		refund:    refund,
		challenge: challenge,
		binding:   binding,
	}
}

func (d *redPacketDatabase) CreateRedPacket(ctx context.Context, rp *model.RedPacket) error {
	return d.rp.Create(ctx, rp)
}

func (d *redPacketDatabase) GetRedPacketByBizID(ctx context.Context, bizID string) (*model.RedPacket, error) {
	return d.rp.GetByBizID(ctx, bizID)
}

func (d *redPacketDatabase) GetRedPacketByPacketID(ctx context.Context, packetID string) (*model.RedPacket, error) {
	return d.rp.GetByPacketID(ctx, packetID)
}

func (d *redPacketDatabase) UpdateRedPacketCreated(ctx context.Context, rp *model.RedPacket) error {
	return d.rp.UpdateCreated(ctx, rp)
}

func (d *redPacketDatabase) UpdateRedPacketStatus(ctx context.Context, packetID, status string) error {
	return d.rp.UpdateStatus(ctx, packetID, status)
}

func (d *redPacketDatabase) UpdateRedPacketClaimProgress(ctx context.Context, packetID, claimedAmount, status string) error {
	return d.rp.UpdateClaimProgress(ctx, packetID, claimedAmount, status)
}

func (d *redPacketDatabase) CreateClaimAuth(ctx context.Context, auth *model.RedPacketClaimAuth) error {
	return d.claimAuth.Create(ctx, auth)
}

func (d *redPacketDatabase) GetClaimAuth(ctx context.Context, packetID, claimer string) (*model.RedPacketClaimAuth, error) {
	return d.claimAuth.Get(ctx, packetID, claimer)
}

func (d *redPacketDatabase) MarkClaimAuthUsed(ctx context.Context, authNonce string) error {
	return d.claimAuth.MarkUsed(ctx, authNonce)
}

func (d *redPacketDatabase) SaveClaim(ctx context.Context, claim *model.RedPacketClaim) error {
	return d.claim.Save(ctx, claim)
}

func (d *redPacketDatabase) GetClaimByPacketIDAndClaimer(ctx context.Context, packetID, claimer string) (*model.RedPacketClaim, error) {
	return d.claim.GetByPacketIDAndClaimer(ctx, packetID, claimer)
}

func (d *redPacketDatabase) GetClaimByPacketIDAndUserID(ctx context.Context, packetID, userID string) (*model.RedPacketClaim, error) {
	return d.claim.GetByPacketIDAndUserID(ctx, packetID, userID)
}

func (d *redPacketDatabase) GetClaimsByPacketID(ctx context.Context, packetID string) ([]*model.RedPacketClaim, error) {
	return d.claim.ListByPacketID(ctx, packetID)
}

func (d *redPacketDatabase) SaveRefund(ctx context.Context, refund *model.RedPacketRefund) error {
	return d.refund.Save(ctx, refund)
}

func (d *redPacketDatabase) CreateWalletBindingChallenge(ctx context.Context, challenge *model.WalletBindingChallenge) error {
	return d.challenge.Create(ctx, challenge)
}

func (d *redPacketDatabase) GetWalletBindingChallenge(ctx context.Context, challengeID string) (*model.WalletBindingChallenge, error) {
	return d.challenge.Get(ctx, challengeID)
}

func (d *redPacketDatabase) UpdateWalletBindingChallenge(ctx context.Context, challenge *model.WalletBindingChallenge) error {
	return d.challenge.Update(ctx, challenge)
}

func (d *redPacketDatabase) UpsertWalletBinding(ctx context.Context, binding *model.WalletBinding) error {
	return d.binding.Upsert(ctx, binding)
}

func (d *redPacketDatabase) GetActiveWalletBinding(ctx context.Context, userID, chainType, walletAddress string) (*model.WalletBinding, error) {
	return d.binding.GetActive(ctx, userID, chainType, walletAddress)
}
