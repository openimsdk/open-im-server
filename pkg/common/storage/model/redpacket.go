package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RedPacket struct {
	BizID           string    `bson:"biz_id"`
	ChainType       string    `bson:"chain_type"`
	PacketID        string    `bson:"packet_id"`
	ChainID         int64     `bson:"chain_id"`
	ContractAddress string    `bson:"contract_address"`
	CreatorUserID   string    `bson:"creator_user_id"`
	CreatorWallet   string    `bson:"creator_wallet"`
	GroupID         string    `bson:"group_id"`
	ScopeType       string    `bson:"scope_type"`
	ReceiverUserID  string    `bson:"receiver_user_id"`
	ReceiverUserIDs []string  `bson:"receiver_user_ids"`
	PacketType      int32     `bson:"packet_type"`
	Token           string    `bson:"token"`
	TotalAmount     string    `bson:"total_amount"`
	TotalShares     int32     `bson:"total_shares"`
	ClaimedAmount        string   `bson:"claimed_amount"`
	ClaimedShares        int32    `bson:"claimed_shares"`
	ProcessedClaimHashes []string `bson:"processed_claim_hashes"`
	ExpiryAt             int64    `bson:"expiry_at"`
	TxHash          string    `bson:"tx_hash"`
	Status          string    `bson:"status"`
	CreatedAt       time.Time `bson:"created_at"`
	UpdatedAt       time.Time `bson:"updated_at"`
}

type RedPacketClaim struct {
	PacketID      string    `bson:"packet_id"`
	UserID        string    `bson:"user_id"`
	ClaimerWallet string    `bson:"claimer_wallet"`
	AuthNonce     string    `bson:"auth_nonce"`
	ClaimTxHash   string    `bson:"claim_tx_hash"`
	ClaimedAmount string    `bson:"claimed_amount"`
	BlockNumber   uint64    `bson:"block_number"`
	Status        string    `bson:"status"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at"`
}

type RedPacketClaimAuth struct {
	PacketID   string    `bson:"packet_id"`
	Claimer    string    `bson:"claimer"`
	AuthNonce  string    `bson:"auth_nonce"`
	RandomSeed string    `bson:"random_seed"`
	Deadline   int64     `bson:"deadline"`
	Signature  string    `bson:"signature"`
	Used       bool      `bson:"used"`
	CreatedAt  time.Time `bson:"created_at"`
}

type RedPacketRefund struct {
	PacketID  string    `bson:"packet_id"`
	RefundTo  string    `bson:"refund_to"`
	TxHash    string    `bson:"tx_hash"`
	Amount    string    `bson:"amount"`
	CreatedAt time.Time `bson:"created_at"`
}

type WalletBindingChallenge struct {
	ChallengeID   string     `bson:"challenge_id"`
	UserID        string     `bson:"user_id"`
	ChainType     string     `bson:"chain_type"`
	ChainID       int64      `bson:"chain_id"`
	WalletAddress string     `bson:"wallet_address"`
	Nonce         string     `bson:"nonce"`
	Message       string     `bson:"message"`
	Protocol      string     `bson:"protocol"`
	SignMethod    string     `bson:"sign_method"`
	Status        string     `bson:"status"`
	Signature     string     `bson:"signature"`
	ExpiresAt     time.Time  `bson:"expires_at"`
	VerifiedAt    *time.Time `bson:"verified_at,omitempty"`
	CreatedAt     time.Time  `bson:"created_at"`
	UpdatedAt     time.Time  `bson:"updated_at"`
}

type WalletBinding struct {
	UserID        string     `bson:"user_id"`
	ChainType     string     `bson:"chain_type"`
	ChainID       int64      `bson:"chain_id"`
	WalletAddress string     `bson:"wallet_address"`
	Status        string     `bson:"status"`
	ChallengeID   string     `bson:"challenge_id"`
	VerifiedAt    time.Time  `bson:"verified_at"`
	RevokedAt     *time.Time `bson:"revoked_at,omitempty"`
	CreatedAt     time.Time  `bson:"created_at"`
	UpdatedAt     time.Time  `bson:"updated_at"`
}

// AdminAuditLog records each admin operation for accountability.
type AdminAuditLog struct {
	ID         primitive.ObjectID `bson:"_id"`
	OperatorID string             `bson:"operator_id"`
	Action     string             `bson:"action"`
	Params     string             `bson:"params"`   // JSON-encoded request
	Result     string             `bson:"result"`   // "success" | "failed"
	ErrMsg     string             `bson:"err_msg"`
	CreatedAt  time.Time          `bson:"created_at"`
}
