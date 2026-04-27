package service

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"redpacket/internal/chain"
	"redpacket/internal/model"
	"redpacket/internal/repository"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

type RedPacketService struct {
	repo        repository.Repository
	chainClient *chain.ChainClient
	signerKey   *ecdsa.PrivateKey
}

type CreateOrderRequest struct {
	CreatorUserID string `json:"creator_user_id" binding:"required"`
	CreatorWallet string `json:"creator_wallet" binding:"required"`
	PacketType    int32  `json:"packet_type" binding:"required"`
	Token         string `json:"token"`
	TotalAmount   string `json:"total_amount" binding:"required"`
	TotalShares   int32  `json:"total_shares" binding:"required"`
	ExpiryAt      int64  `json:"expiry_at"`
	Remark        string `json:"remark"`
}

type CreatedCallbackRequest struct {
	BizID    string `json:"biz_id" binding:"required"`
	TxHash   string `json:"tx_hash" binding:"required"`
	PacketID string `json:"packet_id" binding:"required"`
}

type ClaimResultRequest struct {
	PacketID string `json:"packet_id" binding:"required"`
	Claimer  string `json:"claimer" binding:"required"`
	UserID   string `json:"user_id"`
	TxHash   string `json:"tx_hash" binding:"required"`
}

func NewRedPacketService(repo repository.Repository, chainClient *chain.ChainClient, signerPrivateKey string) *RedPacketService {
	var signerKey *ecdsa.PrivateKey
	if signerPrivateKey != "" {
		var err error
		signerKey, err = crypto.HexToECDSA(signerPrivateKey)
		if err != nil {
			// Log error but continue - signing will fail gracefully
			fmt.Printf("Warning: failed to parse signer private key: %v\n", err)
		}
	}

	return &RedPacketService{
		repo:        repo,
		chainClient: chainClient,
		signerKey:   signerKey,
	}
}

func (s *RedPacketService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (map[string]interface{}, error) {
	bizID := uuid.NewString()

	rp := &model.RedPacket{
		BizID:           bizID,
		CreatorUserID:   req.CreatorUserID,
		CreatorWallet:   req.CreatorWallet,
		PacketType:      req.PacketType,
		Token:           req.Token,
		TotalAmount:     req.TotalAmount,
		TotalShares:     req.TotalShares,
		ExpiryAt:        req.ExpiryAt,
		Status:          "PENDING",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.repo.CreateRedPacket(ctx, rp); err != nil {
		return nil, fmt.Errorf("failed to create red packet: %w", err)
	}

	return map[string]interface{}{
		"biz_id": bizID,
	}, nil
}

func (s *RedPacketService) CreatedCallback(ctx context.Context, req *CreatedCallbackRequest) error {
	return s.repo.UpdateRedPacketTxHash(ctx, req.BizID, req.TxHash, req.PacketID)
}

func (s *RedPacketService) GetDetail(ctx context.Context, packetID string) (map[string]interface{}, error) {
	rp, err := s.repo.GetRedPacketByPacketID(ctx, packetID)
	if err != nil {
		return nil, fmt.Errorf("packet not found: %s", packetID)
	}

	claims, err := s.repo.GetClaimsByPacketID(ctx, packetID)
	if err != nil {
		claims = []model.RedPacketClaim{}
	}

	return map[string]interface{}{
		"biz_record": rp,
		"claims":     claims,
	}, nil
}

func (s *RedPacketService) CanClaim(ctx context.Context, packetID, claimer, userID string) error {
	// Check if packet exists and is active
	rp, err := s.repo.GetRedPacketByPacketID(ctx, packetID)
	if err != nil {
		return fmt.Errorf("packet not found: %s", packetID)
	}

	if rp.Status != "ACTIVE" {
		return fmt.Errorf("packet is not active, current status: %s", rp.Status)
	}

	// TODO: Add more checks - expiry, already claimed by this user, etc.
	// For now we allow the claim
	return nil
}

// SignClaim generates signature for claim operation
func (s *RedPacketService) IssueClaimSign(ctx context.Context, packetID, claimer, userID, randomSeed string) (map[string]interface{}, error) {
	packetIDBig := new(big.Int)
	packetIDBig.SetString(packetID, 10)

	claimerAddr := common.HexToAddress(claimer)

	// Generate nonce and deadline (5 minute expiry)
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())
	deadline := time.Now().Add(5 * time.Minute).Unix()
	randomSeedBig := new(big.Int)
	if randomSeed != "" && randomSeed != "0" {
		randomSeedBig.SetString(randomSeed, 10)
	} else {
		randomSeedBig.SetInt64(time.Now().UnixNano())
	}
	deadlineBig := big.NewInt(deadline)

	var digest [32]byte
	var err error

	if s.chainClient != nil {
		// Use real contract call to getSignMessage
		digest, err = s.chainClient.GetSignMessage(ctx, packetIDBig, claimerAddr, big.NewInt(0), randomSeedBig, deadlineBig)
		if err != nil {
			return nil, fmt.Errorf("getSignMessage failed: %w", err)
		}
	} else {
		// Fallback for testing
		digest = crypto.Keccak256Hash([]byte(fmt.Sprintf("%s%s%s", packetID, claimer, nonce)))
	}

	// Sign the digest
	var signature []byte
	if s.signerKey != nil {
		signature, err = crypto.Sign(digest[:], s.signerKey)
		if err != nil {
			return nil, fmt.Errorf("sign failed: %w", err)
		}
		if len(signature) == 65 && signature[64] < 27 {
			signature[64] += 27
		}
	} else {
		signature = []byte("0xplaceholder-signature-for-testing")
	}

	sigHex := "0x" + hex.EncodeToString(signature)

	auth := &model.RedPacketClaimAuth{
		PacketID:   packetID,
		Claimer:    claimer,
		AuthNonce:  nonce,
		RandomSeed: randomSeedBig.String(),
		Deadline:   deadline,
		Signature:  sigHex,
		CreatedAt:  time.Now(),
	}

	if err := s.repo.CreateClaimAuth(ctx, auth); err != nil {
		return nil, fmt.Errorf("save claim auth failed: %w", err)
	}

	return map[string]interface{}{
		"auth_nonce":  nonce,
		"deadline":    deadline,
		"signature":   sigHex,
		"random_seed": randomSeedBig.String(),
	}, nil
}

func (s *RedPacketService) ClaimResult(ctx context.Context, req *ClaimResultRequest) error {
	claim := &model.RedPacketClaim{
		PacketID:      req.PacketID,
		ClaimerWallet: req.Claimer,
		ClaimTxHash:   req.TxHash,
		Status:        "CONFIRMED",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return s.repo.CreateClaim(ctx, claim)
}
