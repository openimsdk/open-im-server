package model

import (
	"time"
)

type RedPacket struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	BizID           string    `gorm:"uniqueIndex;size:64" json:"biz_id"`
	ChainType       string    `gorm:"index;size:16" json:"chain_type"` // EVM, TRON
	PacketID        string    `gorm:"index;size:32" json:"packet_id"`
	ChainID         int64     `json:"chain_id"`
	ContractAddress string    `json:"contract_address"`
	CreatorUserID   string    `gorm:"size:64" json:"creator_user_id"`
	CreatorWallet   string    `gorm:"size:66" json:"creator_wallet"`
	GroupID         string    `gorm:"index;size:64" json:"group_id"`
	ScopeType       string    `gorm:"size:20" json:"scope_type"` // GROUP, DIRECT, PUBLIC
	ReceiverUserID  string    `gorm:"size:64" json:"receiver_user_id"`
	ReceiverUserIDs string    `gorm:"type:text" json:"receiver_user_ids"`
	PacketType      int32     `json:"packet_type"` // 0=fixed, 1=random, 2=transfer
	Token           string    `gorm:"size:66" json:"token"`
	TotalAmount     string    `gorm:"size:50" json:"total_amount"`
	TotalShares     int32     `json:"total_shares"`
	ClaimedAmount   string    `gorm:"size:50" json:"claimed_amount"`
	ClaimedShares   int32     `json:"claimed_shares"`
	ExpiryAt        int64     `json:"expiry_at"`
	TxHash          string    `gorm:"size:66" json:"tx_hash"`
	Status          string    `gorm:"size:20" json:"status"` // PENDING, ACTIVE, EXPIRED, COMPLETED, REFUNDED
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type RedPacketClaim struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	PacketID      string    `gorm:"index;index:idx_packet_user;size:32" json:"packet_id"`
	UserID        string    `gorm:"index;index:idx_packet_user;size:64" json:"user_id"`
	ClaimerWallet string    `gorm:"size:66" json:"claimer_wallet"`
	AuthNonce     string    `gorm:"size:32" json:"auth_nonce"`
	ClaimTxHash   string    `gorm:"uniqueIndex;size:66" json:"claim_tx_hash"`
	ClaimedAmount string    `gorm:"size:50" json:"claimed_amount"`
	BlockNumber   uint64    `json:"block_number"`
	Status        string    `gorm:"size:20" json:"status"` // PENDING, CONFIRMED, FAILED
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type RedPacketClaimAuth struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	PacketID   string    `gorm:"index;size:32" json:"packet_id"`
	Claimer    string    `gorm:"size:66" json:"claimer"`
	AuthNonce  string    `gorm:"uniqueIndex;size:32" json:"auth_nonce"`
	RandomSeed string    `gorm:"size:32" json:"random_seed"`
	Deadline   int64     `json:"deadline"`
	Signature  string    `gorm:"size:132" json:"signature"`
	Used       bool      `json:"used"`
	CreatedAt  time.Time `json:"created_at"`
}

type RedPacketRefund struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	PacketID  string    `gorm:"index;size:32" json:"packet_id"`
	RefundTo  string    `gorm:"size:66" json:"refund_to"`
	TxHash    string    `gorm:"uniqueIndex;size:66" json:"tx_hash"`
	Amount    string    `gorm:"size:50" json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type WalletBindingChallenge struct {
	ID            uint       `gorm:"primarykey" json:"id"`
	ChallengeID   string     `gorm:"uniqueIndex;size:64" json:"challenge_id"`
	UserID        string     `gorm:"index;size:64" json:"user_id"`
	ChainType     string     `gorm:"index;size:16" json:"chain_type"`
	ChainID       int64      `json:"chain_id"`
	WalletAddress string     `gorm:"index;size:128" json:"wallet_address"`
	Nonce         string     `gorm:"size:64" json:"nonce"`
	Message       string     `gorm:"type:text" json:"message"`
	Protocol      string     `gorm:"size:32" json:"protocol"`
	SignMethod    string     `gorm:"size:32" json:"sign_method"`
	Status        string     `gorm:"size:20" json:"status"` // PENDING, VERIFIED, EXPIRED, FAILED
	Signature     string     `gorm:"type:text" json:"signature"`
	ExpiresAt     time.Time  `json:"expires_at"`
	VerifiedAt    *time.Time `json:"verified_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type WalletBinding struct {
	ID            uint       `gorm:"primarykey" json:"id"`
	UserID        string     `gorm:"index:idx_user_chain_wallet,unique;size:64" json:"user_id"`
	ChainType     string     `gorm:"index:idx_user_chain_wallet,unique;size:16" json:"chain_type"`
	ChainID       int64      `json:"chain_id"`
	WalletAddress string     `gorm:"index:idx_user_chain_wallet,unique;size:128" json:"wallet_address"`
	Status        string     `gorm:"size:20" json:"status"` // ACTIVE, REVOKED
	ChallengeID   string     `gorm:"size:64" json:"challenge_id"`
	VerifiedAt    time.Time  `json:"verified_at"`
	RevokedAt     *time.Time `json:"revoked_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
