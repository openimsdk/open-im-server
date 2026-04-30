package service

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"redpacket/internal/authctx"
	"redpacket/internal/chain"
	"redpacket/internal/model"
	"redpacket/internal/repository"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RedPacketService struct {
	repo        repository.Repository
	chainClient *chain.ChainClient
	tronClient  *chain.TronClient
	signerKey   *ecdsa.PrivateKey
}

type CreateOrderRequest struct {
	ChainType       string   `json:"chain_type" binding:"required"`
	ChainID         int64    `json:"chain_id"`
	ContractAddress string   `json:"contract_address"`
	CreatorUserID   string   `json:"creator_user_id"`
	CreatorWallet   string   `json:"creator_wallet" binding:"required"`
	GroupID         string   `json:"group_id"`
	ScopeType       string   `json:"scope_type"`
	ReceiverUserID  string   `json:"receiver_user_id"`
	ReceiverUserIDs []string `json:"receiver_user_ids"`
	PacketType      int32    `json:"packet_type" binding:"required"`
	Token           string   `json:"token"`
	TotalAmount     string   `json:"total_amount" binding:"required"`
	TotalShares     int32    `json:"total_shares" binding:"required"`
	ExpiryAt        int64    `json:"expiry_at"`
	Remark          string   `json:"remark"`
}

type CreatedCallbackRequest struct {
	BizID           string   `json:"biz_id" binding:"required"`
	TxHash          string   `json:"tx_hash" binding:"required"`
	PacketID        string   `json:"packet_id"`
	GroupID         string   `json:"group_id"`
	ScopeType       string   `json:"scope_type"`
	ReceiverUserID  string   `json:"receiver_user_id"`
	ReceiverUserIDs []string `json:"receiver_user_ids"`
}

type ClaimResultRequest struct {
	PacketID string `json:"packet_id" binding:"required"`
	Claimer  string `json:"claimer" binding:"required"`
	UserID   string `json:"-"`
	TxHash   string `json:"tx_hash" binding:"required"`
}

type WalletBindChallengeRequest struct {
	UserID        string `json:"user_id"`
	ChainType     string `json:"chain_type" binding:"required"`
	ChainID       int64  `json:"chain_id"`
	WalletAddress string `json:"wallet_address" binding:"required"`
	Domain        string `json:"domain"`
	URI           string `json:"uri"`
}

type WalletBindConfirmRequest struct {
	ChallengeID string `json:"challenge_id" binding:"required"`
	Signature   string `json:"signature" binding:"required"`
}

func NewRedPacketService(repo repository.Repository, chainClient *chain.ChainClient, tronClient *chain.TronClient, signerPrivateKey string) *RedPacketService {
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
		tronClient:  tronClient,
		signerKey:   signerKey,
	}
}

func (s *RedPacketService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (map[string]interface{}, error) {
	currentUserID, err := authctx.CurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	req.CreatorUserID = currentUserID

	bizID := uuid.NewString()
	chainType, err := normalizeChainType(req.ChainType)
	if err != nil {
		return nil, err
	}
	scopeType := normalizeScopeType(req.ScopeType)
	if err := validateCreateScope(scopeType, req.GroupID, req.ReceiverUserID, req.ReceiverUserIDs); err != nil {
		return nil, err
	}
	if err := s.validateCreateHook(ctx, req); err != nil {
		return nil, err
	}

	receiverUserIDs, err := encodeReceiverUserIDs(req.ReceiverUserIDs)
	if err != nil {
		return nil, fmt.Errorf("encode receiver_user_ids failed: %w", err)
	}

	chainID := req.ChainID
	contractAddress := strings.TrimSpace(req.ContractAddress)
	if chainType == "EVM" && s.chainClient != nil {
		if chainID == 0 {
			if chainValue := s.chainClient.ChainID(); chainValue != nil {
				chainID = chainValue.Int64()
			}
		}
		if contractAddress == "" {
			contractAddress = s.chainClient.ContractAddress().Hex()
		}
	}
	if chainType == "TRON" && s.tronClient != nil && contractAddress == "" {
		contractAddress = s.tronClient.ContractAddress()
	}

	rp := &model.RedPacket{
		BizID:           bizID,
		ChainType:       chainType,
		ChainID:         chainID,
		ContractAddress: contractAddress,
		CreatorUserID:   req.CreatorUserID,
		CreatorWallet:   req.CreatorWallet,
		GroupID:         req.GroupID,
		ScopeType:       scopeType,
		ReceiverUserID:  req.ReceiverUserID,
		ReceiverUserIDs: receiverUserIDs,
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
	rp, err := s.repo.GetRedPacketByBizID(ctx, req.BizID)
	if err != nil {
		return fmt.Errorf("biz record not found: %s", req.BizID)
	}

	groupID := firstNonEmpty(req.GroupID, rp.GroupID)
	scopeType := normalizeScopeType(firstNonEmpty(req.ScopeType, rp.ScopeType))
	receiverUserID := firstNonEmpty(req.ReceiverUserID, rp.ReceiverUserID)
	receiverUserIDs := rp.ReceiverUserIDs
	if len(req.ReceiverUserIDs) > 0 {
		receiverUserIDs, err = encodeReceiverUserIDs(req.ReceiverUserIDs)
		if err != nil {
			return fmt.Errorf("encode receiver_user_ids failed: %w", err)
		}
	}

	if err := validateCreateScope(scopeType, groupID, receiverUserID, decodeReceiverUserIDs(receiverUserIDs)); err != nil {
		return err
	}

	createdPacket, err := s.resolveCreatedPacket(ctx, rp, req.TxHash, req.PacketID)
	if err != nil {
		return err
	}

	return s.repo.UpdateRedPacketCreated(ctx, &model.RedPacket{
		BizID:           req.BizID,
		ChainType:       rp.ChainType,
		PacketID:        createdPacket.PacketID,
		ChainID:         createdPacket.ChainID,
		ContractAddress: createdPacket.ContractAddress,
		TxHash:          req.TxHash,
		GroupID:         groupID,
		ScopeType:       scopeType,
		ReceiverUserID:  receiverUserID,
		ReceiverUserIDs: receiverUserIDs,
		Status:          "ACTIVE",
	})
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
	rp, err := s.repo.GetRedPacketByPacketID(ctx, packetID)
	if err != nil {
		return fmt.Errorf("packet not found: %s", packetID)
	}

	if err := validateClaimBase(rp, userID, claimer); err != nil {
		return err
	}
	if err := s.ensureWalletBinding(ctx, userID, claimer, rp.ChainType); err != nil {
		return err
	}

	switch rp.PacketType {
	case 0:
		return s.validateFixedPacketClaim(ctx, rp, userID, claimer)
	case 1:
		return s.validateRandomPacketClaim(ctx, rp, userID, claimer)
	case 2:
		return s.validateTransferPacketClaim(ctx, rp, userID, claimer)
	default:
		return fmt.Errorf("unsupported packet_type: %d", rp.PacketType)
	}
}

// SignClaim generates signature for claim operation
func (s *RedPacketService) IssueClaimSign(ctx context.Context, packetID, claimer, randomSeed string) (map[string]interface{}, error) {
	userID, err := authctx.CurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.CanClaim(ctx, packetID, claimer, userID); err != nil {
		return nil, err
	}

	packetIDBig := new(big.Int)
	if _, ok := packetIDBig.SetString(packetID, 10); !ok {
		return nil, fmt.Errorf("invalid packet_id: %s", packetID)
	}

	claimerAddr := common.HexToAddress(claimer)

	// Generate nonce and deadline (5 minute expiry)
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())
	authNonceBig := new(big.Int)
	if _, ok := authNonceBig.SetString(nonce, 10); !ok {
		return nil, fmt.Errorf("invalid auth nonce")
	}
	deadline := time.Now().Add(5 * time.Minute).Unix()
	randomSeedBig := new(big.Int)
	if randomSeed != "" && randomSeed != "0" {
		if _, ok := randomSeedBig.SetString(randomSeed, 10); !ok {
			return nil, fmt.Errorf("invalid random_seed: %s", randomSeed)
		}
	} else {
		randomSeedBig.SetInt64(time.Now().UnixNano())
	}
	deadlineBig := big.NewInt(deadline)

	var digest [32]byte

	if s.chainClient != nil {
		// Use real contract call to getSignMessage
		digest, err = s.chainClient.GetSignMessage(ctx, packetIDBig, claimerAddr, authNonceBig, randomSeedBig, deadlineBig)
		if err != nil {
			return nil, fmt.Errorf("getSignMessage failed: %w", err)
		}
	} else {
		// Fallback for testing
		digest = crypto.Keccak256Hash([]byte(fmt.Sprintf("%s:%s:%s:%s:%d", packetID, claimer, nonce, randomSeedBig.String(), deadline)))
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
	userID, err := authctx.CurrentUserID(ctx)
	if err != nil {
		return err
	}
	req.UserID = userID

	rp, err := s.repo.GetRedPacketByPacketID(ctx, req.PacketID)
	if err != nil {
		return fmt.Errorf("packet not found: %s", req.PacketID)
	}

	if err := validateClaimBase(rp, req.UserID, req.Claimer); err != nil {
		return err
	}

	claim := &model.RedPacketClaim{
		PacketID:      req.PacketID,
		UserID:        req.UserID,
		ClaimerWallet: req.Claimer,
		ClaimTxHash:   req.TxHash,
		Status:        "PENDING",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.SaveClaim(ctx, claim); err != nil {
		return err
	}

	claimedEvent, err := s.resolveClaimedEvent(ctx, rp, req.TxHash)
	if err != nil {
		return nil
	}
	if claimedEvent == nil {
		return nil
	}
	if !strings.EqualFold(claimedEvent.ClaimerWallet, req.Claimer) {
		return fmt.Errorf("claim event claimer mismatch: got %s want %s", claimedEvent.ClaimerWallet, req.Claimer)
	}

	confirmed := &model.RedPacketClaim{
		PacketID:      req.PacketID,
		UserID:        req.UserID,
		ClaimerWallet: claimedEvent.ClaimerWallet,
		AuthNonce:     claimedEvent.AuthNonce,
		ClaimTxHash:   req.TxHash,
		ClaimedAmount: claimedEvent.Amount,
		BlockNumber:   claimedEvent.BlockNumber,
		Status:        "CONFIRMED",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := s.repo.SaveClaim(ctx, confirmed); err != nil {
		return err
	}

	if claimedEvent.AuthNonce != "" {
		if err := s.repo.MarkClaimAuthUsed(ctx, claimedEvent.AuthNonce); err != nil {
			return err
		}
	}

	nextStatus := derivePacketStatusAfterClaim(rp, claimedEvent.Amount)
	return s.repo.UpdateRedPacketClaimProgress(ctx, req.PacketID, claimedEvent.Amount, nextStatus)
}

func (s *RedPacketService) IssueWalletBindChallenge(ctx context.Context, req *WalletBindChallengeRequest) (map[string]interface{}, error) {
	currentUserID, err := authctx.CurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	req.UserID = currentUserID

	chainType, err := normalizeChainType(req.ChainType)
	if err != nil {
		return nil, err
	}

	walletAddress := strings.TrimSpace(req.WalletAddress)
	if walletAddress == "" {
		return nil, fmt.Errorf("wallet_address is required")
	}

	challengeID := uuid.NewString()
	nonce := uuid.NewString()
	issuedAt := time.Now().UTC()
	expiresAt := issuedAt.Add(10 * time.Minute)

	protocol := "siwe-eip4361"
	signMethod := "personal_sign"
	message := buildEVMBindMessage(req, challengeID, nonce, issuedAt, expiresAt)
	if chainType == "TRON" {
		protocol = "tron-signmessagev2"
		signMethod = "signMessageV2"
		message = buildTRONBindMessage(req, challengeID, nonce, issuedAt, expiresAt)
	}

	challenge := &model.WalletBindingChallenge{
		ChallengeID:   challengeID,
		UserID:        req.UserID,
		ChainType:     chainType,
		ChainID:       req.ChainID,
		WalletAddress: walletAddress,
		Nonce:         nonce,
		Message:       message,
		Protocol:      protocol,
		SignMethod:    signMethod,
		Status:        "PENDING",
		ExpiresAt:     expiresAt,
		CreatedAt:     issuedAt,
		UpdatedAt:     issuedAt,
	}
	if err := s.repo.CreateWalletBindingChallenge(ctx, challenge); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"challenge_id": challengeID,
		"user_id":      req.UserID,
		"chain_type":   chainType,
		"chain_id":     req.ChainID,
		"wallet":       walletAddress,
		"protocol":     protocol,
		"sign_method":  signMethod,
		"nonce":        nonce,
		"message":      message,
		"issued_at":    issuedAt.Format(time.RFC3339),
		"expires_at":   expiresAt.Format(time.RFC3339),
	}, nil
}

func (s *RedPacketService) ConfirmWalletBind(ctx context.Context, req *WalletBindConfirmRequest) (map[string]interface{}, error) {
	challenge, err := s.repo.GetWalletBindingChallenge(ctx, req.ChallengeID)
	if err != nil {
		return nil, fmt.Errorf("challenge not found: %s", req.ChallengeID)
	}
	if challenge.Status != "PENDING" {
		return nil, fmt.Errorf("challenge is not pending")
	}
	if time.Now().UTC().After(challenge.ExpiresAt) {
		challenge.Status = "EXPIRED"
		challenge.UpdatedAt = time.Now()
		_ = s.repo.UpdateWalletBindingChallenge(ctx, challenge)
		return nil, fmt.Errorf("challenge is expired")
	}

	switch challenge.ChainType {
	case "EVM":
		if err := verifyEVMBindSignature(challenge.Message, challenge.WalletAddress, req.Signature); err != nil {
			challenge.Status = "FAILED"
			challenge.Signature = req.Signature
			challenge.UpdatedAt = time.Now()
			_ = s.repo.UpdateWalletBindingChallenge(ctx, challenge)
			return nil, err
		}
	case "TRON":
		return nil, fmt.Errorf("TRON wallet binding verification is not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported chain_type: %s", challenge.ChainType)
	}

	now := time.Now().UTC()
	challenge.Status = "VERIFIED"
	challenge.Signature = req.Signature
	challenge.VerifiedAt = &now
	challenge.UpdatedAt = now
	if err := s.repo.UpdateWalletBindingChallenge(ctx, challenge); err != nil {
		return nil, err
	}

	binding := &model.WalletBinding{
		UserID:        challenge.UserID,
		ChainType:     challenge.ChainType,
		ChainID:       challenge.ChainID,
		WalletAddress: challenge.WalletAddress,
		Status:        "ACTIVE",
		ChallengeID:   challenge.ChallengeID,
		VerifiedAt:    now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := s.repo.UpsertWalletBinding(ctx, binding); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"user_id":        binding.UserID,
		"chain_type":     binding.ChainType,
		"chain_id":       binding.ChainID,
		"wallet_address": binding.WalletAddress,
		"status":         binding.Status,
		"verified_at":    binding.VerifiedAt.Format(time.RFC3339),
	}, nil
}

func (s *RedPacketService) GetWalletBinding(ctx context.Context, userID, chainType, walletAddress string) (map[string]interface{}, error) {
	currentUserID, err := authctx.CurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	userID = currentUserID

	normalizedChainType, err := normalizeChainType(chainType)
	if err != nil {
		return nil, err
	}
	binding, err := s.repo.GetActiveWalletBinding(ctx, userID, normalizedChainType, walletAddress)
	if err != nil {
		return nil, fmt.Errorf("active wallet binding not found")
	}
	return map[string]interface{}{
		"user_id":        binding.UserID,
		"chain_type":     binding.ChainType,
		"chain_id":       binding.ChainID,
		"wallet_address": binding.WalletAddress,
		"status":         binding.Status,
		"challenge_id":   binding.ChallengeID,
		"verified_at":    binding.VerifiedAt.Format(time.RFC3339),
	}, nil
}

type claimedEventSnapshot struct {
	ClaimerWallet string
	AuthNonce     string
	Amount        string
	BlockNumber   uint64
}

type createdPacketSnapshot struct {
	PacketID        string
	ChainID         int64
	ContractAddress string
	CreatorWallet   string
	PacketType      int32
	Token           string
	TotalAmount     string
	TotalShares     int32
	ExpiryAt        int64
}

func (s *RedPacketService) resolveCreatedPacket(ctx context.Context, rp *model.RedPacket, txHashHex, fallbackPacketID string) (*createdPacketSnapshot, error) {
	switch rp.ChainType {
	case "EVM":
		if s.chainClient == nil {
			if fallbackPacketID == "" {
				return nil, fmt.Errorf("packet_id is required when EVM client is unavailable")
			}
			return buildFallbackCreatedPacket(rp, fallbackPacketID), nil
		}

		events, err := s.chainClient.ParseTransactionReceipt(ctx, common.HexToHash(txHashHex))
		if err != nil {
			if fallbackPacketID == "" {
				return nil, fmt.Errorf("parse created tx failed: %w", err)
			}
			return buildFallbackCreatedPacket(rp, fallbackPacketID), nil
		}

		for _, event := range events {
			if event.Name != "PacketCreated" {
				continue
			}
			createdPacket := buildCreatedPacketSnapshot(rp, event)
			if chainValue := s.chainClient.ChainID(); chainValue != nil {
				createdPacket.ChainID = chainValue.Int64()
			}
			createdPacket.ContractAddress = s.chainClient.ContractAddress().Hex()
			if err := validateCreatedPacket(rp, createdPacket); err != nil {
				return nil, err
			}

			return createdPacket, nil
		}

		if fallbackPacketID == "" {
			return nil, fmt.Errorf("PacketCreated event not found in tx: %s", txHashHex)
		}
		return buildFallbackCreatedPacket(rp, fallbackPacketID), nil
	case "TRON":
		if s.tronClient == nil {
			if fallbackPacketID == "" {
				return nil, fmt.Errorf("packet_id is required when TRON client is unavailable")
			}
			return buildFallbackCreatedPacket(rp, fallbackPacketID), nil
		}

		events, err := s.tronClient.ParseTransactionReceipt(ctx, txHashHex)
		if err != nil {
			if fallbackPacketID == "" {
				return nil, fmt.Errorf("parse tron created tx failed: %w", err)
			}
			return buildFallbackCreatedPacket(rp, fallbackPacketID), nil
		}

		for _, event := range events {
			if event.Name != "PacketCreated" {
				continue
			}
			createdPacket := buildCreatedPacketSnapshot(rp, event)
			createdPacket.ContractAddress = firstNonEmpty(s.tronClient.ContractAddress(), rp.ContractAddress)
			if err := validateCreatedPacket(rp, createdPacket); err != nil {
				return nil, err
			}
			return createdPacket, nil
		}

		if fallbackPacketID == "" {
			return nil, fmt.Errorf("PacketCreated event not found in TRON tx: %s", txHashHex)
		}
		return buildFallbackCreatedPacket(rp, fallbackPacketID), nil
	default:
		return nil, fmt.Errorf("unsupported chain_type: %s", rp.ChainType)
	}
}

// validateCreateHook reserves a centralized validation extension point for
// create-order. Concrete centralized checks are deferred, but validation is
// already split by packet type so future rules can evolve independently.
func (s *RedPacketService) validateCreateHook(ctx context.Context, req *CreateOrderRequest) error {
	switch req.PacketType {
	case 0:
		return s.validateFixedPacketCreate(ctx, req)
	case 1:
		return s.validateRandomPacketCreate(ctx, req)
	case 2:
		return s.validateTransferPacketCreate(ctx, req)
	default:
		return fmt.Errorf("unsupported packet_type: %d", req.PacketType)
	}
}

// validateFixedPacketCreate reserves centralized checks for fixed red packets.
// todo: validate creator identity, group validity, and group membership.
func (s *RedPacketService) validateFixedPacketCreate(ctx context.Context, req *CreateOrderRequest) error {
	return nil
}

// validateRandomPacketCreate reserves centralized checks for random red packets.
// todo: validate creator identity, group validity, and group membership.
func (s *RedPacketService) validateRandomPacketCreate(ctx context.Context, req *CreateOrderRequest) error {
	return nil
}

// validateTransferPacketCreate reserves centralized checks for transfer packets.
// todo: validate creator identity and sender/receiver relationship.
func (s *RedPacketService) validateTransferPacketCreate(ctx context.Context, req *CreateOrderRequest) error {
	return nil
}

func buildFallbackCreatedPacket(rp *model.RedPacket, packetID string) *createdPacketSnapshot {
	return &createdPacketSnapshot{
		PacketID:        packetID,
		ChainID:         rp.ChainID,
		ContractAddress: rp.ContractAddress,
		CreatorWallet:   strings.ToLower(rp.CreatorWallet),
		PacketType:      rp.PacketType,
		Token:           normalizeTokenAddress(rp.Token),
		TotalAmount:     rp.TotalAmount,
		TotalShares:     rp.TotalShares,
		ExpiryAt:        rp.ExpiryAt,
	}
}

func buildCreatedPacketSnapshot(rp *model.RedPacket, event *chain.ParsedEvent) *createdPacketSnapshot {
	return &createdPacketSnapshot{
		PacketID:        chain.GetPacketIDFromEvent(event).String(),
		ChainID:         rp.ChainID,
		ContractAddress: rp.ContractAddress,
		CreatorWallet:   strings.ToLower(chain.GetAddressFromEvent(event, "creator").Hex()),
		PacketType:      int32(chain.GetUintFromEvent(event, "packetType").Int64()),
		Token:           strings.ToLower(chain.GetAddressFromEvent(event, "token").Hex()),
		TotalAmount:     chain.GetUintFromEvent(event, "totalAmount").String(),
		TotalShares:     int32(chain.GetUintFromEvent(event, "totalShares").Int64()),
		ExpiryAt:        chain.GetUintFromEvent(event, "expiryAt").Int64(),
	}
}

func validateCreatedPacket(rp *model.RedPacket, createdPacket *createdPacketSnapshot) error {
	if createdPacket == nil {
		return fmt.Errorf("created packet is nil")
	}

	if createdPacket.CreatorWallet != "" && strings.ToLower(rp.CreatorWallet) != createdPacket.CreatorWallet {
		return fmt.Errorf("creator mismatch: got %s want %s", createdPacket.CreatorWallet, rp.CreatorWallet)
	}
	if createdPacket.PacketType != rp.PacketType {
		return fmt.Errorf("packet type mismatch: got %d want %d", createdPacket.PacketType, rp.PacketType)
	}
	if createdPacket.TotalAmount != rp.TotalAmount {
		return fmt.Errorf("total amount mismatch: got %s want %s", createdPacket.TotalAmount, rp.TotalAmount)
	}
	if createdPacket.TotalShares != rp.TotalShares {
		return fmt.Errorf("total shares mismatch: got %d want %d", createdPacket.TotalShares, rp.TotalShares)
	}
	expectedToken := normalizeTokenAddress(rp.Token)
	if createdPacket.Token != expectedToken {
		return fmt.Errorf("token mismatch: got %s want %s", createdPacket.Token, expectedToken)
	}
	if rp.ExpiryAt > 0 && createdPacket.ExpiryAt != rp.ExpiryAt {
		return fmt.Errorf("expiry mismatch: got %d want %d", createdPacket.ExpiryAt, rp.ExpiryAt)
	}

	return nil
}

func validateClaimBase(rp *model.RedPacket, userID, claimer string) error {
	if rp == nil {
		return fmt.Errorf("packet not found")
	}
	if strings.TrimSpace(userID) == "" {
		return fmt.Errorf("user_id is required")
	}
	if strings.TrimSpace(claimer) == "" {
		return fmt.Errorf("claimer is required")
	}
	if rp.Status != "ACTIVE" {
		return fmt.Errorf("packet is not active, current status: %s", rp.Status)
	}
	if rp.ExpiryAt > 0 && rp.ExpiryAt <= time.Now().Unix() {
		return fmt.Errorf("packet is expired")
	}
	if rp.Status == "REFUNDED" {
		return fmt.Errorf("packet is refunded")
	}
	return nil
}

func (s *RedPacketService) validateFixedPacketClaim(ctx context.Context, rp *model.RedPacket, userID, claimer string) error {
	if strings.TrimSpace(rp.GroupID) == "" {
		return fmt.Errorf("group_id is required for fixed packet claim")
	}
	if err := s.ensureNotClaimed(ctx, rp.PacketID, userID, claimer); err != nil {
		return err
	}
	return s.ensureGroupEligibility(ctx, rp.GroupID, userID)
}

func (s *RedPacketService) validateRandomPacketClaim(ctx context.Context, rp *model.RedPacket, userID, claimer string) error {
	if strings.TrimSpace(rp.GroupID) == "" {
		return fmt.Errorf("group_id is required for random packet claim")
	}
	if err := s.ensureNotClaimed(ctx, rp.PacketID, userID, claimer); err != nil {
		return err
	}
	return s.ensureGroupEligibility(ctx, rp.GroupID, userID)
}

func (s *RedPacketService) validateTransferPacketClaim(ctx context.Context, rp *model.RedPacket, userID, claimer string) error {
	if err := s.ensureNotClaimed(ctx, rp.PacketID, userID, claimer); err != nil {
		return err
	}
	if strings.TrimSpace(rp.ReceiverUserID) == "" {
		return fmt.Errorf("receiver_user_id is required for transfer claim")
	}
	if rp.ReceiverUserID != userID {
		return fmt.Errorf("user is not the designated receiver")
	}
	return s.ensureFriendRelationship(ctx, rp.CreatorUserID, userID)
}

func (s *RedPacketService) ensureNotClaimed(ctx context.Context, packetID, userID, claimer string) error {
	if strings.TrimSpace(userID) != "" {
		claim, err := s.repo.GetClaimByPacketIDAndUserID(ctx, packetID, userID)
		if err == nil && claim != nil && claim.Status != "FAILED" {
			return fmt.Errorf("user already claimed")
		}
		if err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("failed to check user claim status: %w", err)
		}
	}

	claim, err := s.repo.GetClaimByPacketIDAndClaimer(ctx, packetID, claimer)
	if err == nil && claim != nil && claim.Status != "FAILED" {
		return fmt.Errorf("already claimed")
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check claim status: %w", err)
	}
	return nil
}

// ensureWalletBinding reserves the centralized identity check between Web2
// user identity and wallet address used for claiming.
func (s *RedPacketService) ensureWalletBinding(ctx context.Context, userID, claimer, chainType string) error {
	if _, err := s.repo.GetActiveWalletBinding(ctx, userID, chainType, claimer); err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("wallet is not bound to user")
		}
		return fmt.Errorf("check wallet binding failed: %w", err)
	}
	return nil
}

// ensureGroupEligibility reserves centralized group validation, including
// whether the group exists and whether the user is currently a member.
func (s *RedPacketService) ensureGroupEligibility(ctx context.Context, groupID, userID string) error {
	return nil
}

// ensureFriendRelationship reserves centralized relation validation for
// transfer packets.
func (s *RedPacketService) ensureFriendRelationship(ctx context.Context, creatorUserID, receiverUserID string) error {
	return nil
}

func (s *RedPacketService) resolveClaimedEvent(ctx context.Context, rp *model.RedPacket, txHash string) (*claimedEventSnapshot, error) {
	var (
		events []*chain.ParsedEvent
		err    error
	)

	switch rp.ChainType {
	case "EVM":
		if s.chainClient == nil {
			return nil, nil
		}
		events, err = s.chainClient.ParseTransactionReceipt(ctx, common.HexToHash(txHash))
	case "TRON":
		if s.tronClient == nil {
			return nil, nil
		}
		events, err = s.tronClient.ParseTransactionReceipt(ctx, txHash)
	default:
		return nil, fmt.Errorf("unsupported chain_type: %s", rp.ChainType)
	}
	if err != nil {
		return nil, err
	}

	for _, event := range events {
		if event.Name != "PacketClaimed" {
			continue
		}
		packetID := chain.GetPacketIDFromEvent(event).String()
		claimerWallet := strings.ToLower(chain.GetAddressFromEvent(event, "claimer").Hex())
		if packetID != rp.PacketID {
			return nil, fmt.Errorf("claim event packet mismatch: got %s want %s", packetID, rp.PacketID)
		}
		return &claimedEventSnapshot{
			ClaimerWallet: claimerWallet,
			AuthNonce:     chain.GetUintFromEvent(event, "authNonce").String(),
			Amount:        chain.GetAmountFromEvent(event).String(),
			BlockNumber:   event.BlockNumber,
		}, nil
	}

	return nil, nil
}

func derivePacketStatusAfterClaim(rp *model.RedPacket, claimedAmount string) string {
	if rp == nil {
		return ""
	}
	if rp.PacketType == 2 {
		return "COMPLETED"
	}

	nextShares := rp.ClaimedShares + 1
	if rp.TotalShares > 0 && nextShares >= rp.TotalShares {
		return "COMPLETED"
	}

	totalClaimed := addNumericStrings(rp.ClaimedAmount, claimedAmount)
	if rp.TotalAmount != "" && totalClaimed == rp.TotalAmount {
		return "COMPLETED"
	}

	return "ACTIVE"
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

func buildEVMBindMessage(req *WalletBindChallengeRequest, challengeID, nonce string, issuedAt, expiresAt time.Time) string {
	domain := strings.TrimSpace(req.Domain)
	if domain == "" {
		domain = "redpacket"
	}
	uri := strings.TrimSpace(req.URI)
	if uri == "" {
		uri = "https://redpacket.local/wallet-bind"
	}

	var b strings.Builder
	fmt.Fprintf(&b, "%s wants you to sign in with your Ethereum account:\n", domain)
	b.WriteString(strings.TrimSpace(req.WalletAddress))
	b.WriteString("\n\n")
	fmt.Fprintf(&b, "Bind wallet %s to user %s.\n", strings.TrimSpace(req.WalletAddress), strings.TrimSpace(req.UserID))
	fmt.Fprintf(&b, "URI: %s\n", uri)
	fmt.Fprintf(&b, "Version: 1\n")
	fmt.Fprintf(&b, "Chain ID: %d\n", req.ChainID)
	fmt.Fprintf(&b, "Nonce: %s\n", nonce)
	fmt.Fprintf(&b, "Issued At: %s\n", issuedAt.Format(time.RFC3339))
	fmt.Fprintf(&b, "Expiration Time: %s\n", expiresAt.Format(time.RFC3339))
	fmt.Fprintf(&b, "Request ID: %s", challengeID)
	return b.String()
}

func buildTRONBindMessage(req *WalletBindChallengeRequest, challengeID, nonce string, issuedAt, expiresAt time.Time) string {
	return fmt.Sprintf(
		"Bind TRON wallet %s to user %s\nchallenge_id: %s\nnonce: %s\nchain_id: %d\nissued_at: %s\nexpires_at: %s",
		strings.TrimSpace(req.WalletAddress),
		strings.TrimSpace(req.UserID),
		challengeID,
		nonce,
		req.ChainID,
		issuedAt.Format(time.RFC3339),
		expiresAt.Format(time.RFC3339),
	)
}

func verifyEVMBindSignature(message, walletAddress, signature string) error {
	if strings.TrimSpace(message) == "" {
		return fmt.Errorf("bind message is empty")
	}
	if !common.IsHexAddress(walletAddress) {
		return fmt.Errorf("invalid evm wallet address")
	}

	sig, err := hex.DecodeString(strings.TrimPrefix(signature, "0x"))
	if err != nil {
		return fmt.Errorf("decode signature failed: %w", err)
	}
	if len(sig) != 65 {
		return fmt.Errorf("invalid signature length: %d", len(sig))
	}
	if sig[64] >= 27 {
		sig[64] -= 27
	}
	if sig[64] > 1 {
		return fmt.Errorf("invalid signature recovery id")
	}

	hash := crypto.Keccak256Hash([]byte(personalSignMessage(message)))
	pubKey, err := crypto.SigToPub(hash.Bytes(), sig)
	if err != nil {
		return fmt.Errorf("recover signer failed: %w", err)
	}

	recovered := crypto.PubkeyToAddress(*pubKey)
	if !strings.EqualFold(recovered.Hex(), walletAddress) {
		return fmt.Errorf("signature does not match wallet address")
	}
	return nil
}

func personalSignMessage(message string) string {
	return fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
}

func normalizeScopeType(scopeType string) string {
	switch strings.ToUpper(strings.TrimSpace(scopeType)) {
	case "GROUP", "DIRECT", "PUBLIC":
		return strings.ToUpper(strings.TrimSpace(scopeType))
	default:
		return "PUBLIC"
	}
}

func normalizeChainType(chainType string) (string, error) {
	switch strings.ToUpper(strings.TrimSpace(chainType)) {
	case "EVM":
		return "EVM", nil
	case "TRON":
		return "TRON", nil
	default:
		return "", fmt.Errorf("unsupported chain_type: %s", chainType)
	}
}

func validateCreateScope(scopeType, groupID, receiverUserID string, receiverUserIDs []string) error {
	switch scopeType {
	case "GROUP":
		if strings.TrimSpace(groupID) == "" {
			return fmt.Errorf("group_id is required when scope_type=GROUP")
		}
	case "DIRECT":
		if strings.TrimSpace(receiverUserID) == "" && len(receiverUserIDs) == 0 {
			return fmt.Errorf("receiver_user_id or receiver_user_ids is required when scope_type=DIRECT")
		}
	}
	return nil
}

func encodeReceiverUserIDs(userIDs []string) (string, error) {
	if len(userIDs) == 0 {
		return "", nil
	}
	encoded, err := json.Marshal(userIDs)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

func decodeReceiverUserIDs(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var userIDs []string
	if err := json.Unmarshal([]byte(raw), &userIDs); err != nil {
		return nil
	}
	return userIDs
}

func normalizeTokenAddress(token string) string {
	if strings.TrimSpace(token) == "" {
		return strings.ToLower(common.Address{}.Hex())
	}
	return strings.ToLower(common.HexToAddress(token).Hex())
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
