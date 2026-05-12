package database

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

type RedPacket interface {
	Create(ctx context.Context, rp *model.RedPacket) error
	GetByBizID(ctx context.Context, bizID string) (*model.RedPacket, error)
	GetByPacketID(ctx context.Context, packetID string) (*model.RedPacket, error)
	UpdateCreated(ctx context.Context, rp *model.RedPacket) error
	UpdateStatus(ctx context.Context, packetID, status string) error
	// UpdateClaimProgress atomically increments the claim counter for packetID.
	// claimTxHash is used as an idempotency key so that re-processing the same
	// on-chain transaction never double-counts. When status is empty the method
	// auto-derives the correct status (COMPLETED or ACTIVE).
	UpdateClaimProgress(ctx context.Context, packetID, claimedAmount, status, claimTxHash string) error
	// GetExpiredPending returns ACTIVE packets whose expiry_at < now (unix seconds).
	GetExpiredPending(ctx context.Context, now int64) ([]*model.RedPacket, error)
}

type RedPacketClaim interface {
	Save(ctx context.Context, claim *model.RedPacketClaim) error
	GetByPacketIDAndClaimer(ctx context.Context, packetID, claimer string) (*model.RedPacketClaim, error)
	GetByPacketIDAndUserID(ctx context.Context, packetID, userID string) (*model.RedPacketClaim, error)
	ListByPacketID(ctx context.Context, packetID string) ([]*model.RedPacketClaim, error)
}

type RedPacketClaimAuth interface {
	Create(ctx context.Context, auth *model.RedPacketClaimAuth) error
	Get(ctx context.Context, packetID, claimer string) (*model.RedPacketClaimAuth, error)
	MarkUsed(ctx context.Context, authNonce string) error
}

type RedPacketRefund interface {
	Save(ctx context.Context, refund *model.RedPacketRefund) error
	GetByPacketID(ctx context.Context, packetID string) (*model.RedPacketRefund, error)
}

type AdminAuditLog interface {
	Create(ctx context.Context, log *model.AdminAuditLog) error
}

type WalletBindingChallenge interface {
	Create(ctx context.Context, challenge *model.WalletBindingChallenge) error
	Get(ctx context.Context, challengeID string) (*model.WalletBindingChallenge, error)
	Update(ctx context.Context, challenge *model.WalletBindingChallenge) error
}

type WalletBinding interface {
	Upsert(ctx context.Context, binding *model.WalletBinding) error
	GetActive(ctx context.Context, userID, chainType, walletAddress string) (*model.WalletBinding, error)
}
