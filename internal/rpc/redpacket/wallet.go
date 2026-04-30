package redpacket

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	pbredpacket "github.com/openimsdk/protocol/redpacket"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/mcontext"
)

func (s *redPacketServer) IssueWalletBindChallenge(ctx context.Context, req *pbredpacket.IssueWalletBindChallengeReq) (*pbredpacket.IssueWalletBindChallengeResp, error) {
	currentUserID := mcontext.GetOpUserID(ctx)
	if currentUserID == "" {
		return nil, servererrs.ErrNoPermission.WrapMsg("op user id is empty")
	}

	chainType, err := normalizeChainType(req.ChainType)
	if err != nil {
		return nil, err
	}

	walletAddress := strings.TrimSpace(req.WalletAddress)
	if walletAddress == "" {
		return nil, errs.ErrArgs.WrapMsg("wallet_address is required")
	}

	challengeID := uuid.NewString()
	nonce := uuid.NewString()
	issuedAt := time.Now().UTC()
	expiresAt := issuedAt.Add(10 * time.Minute)

	protocol := "siwe-eip4361"
	signMethod := "personal_sign"
	message := buildEVMBindMessage(currentUserID, walletAddress, req.Domain, req.Uri, req.ChainID, challengeID, nonce, issuedAt, expiresAt)
	if chainType == "TRON" {
		protocol = "tron-signmessagev2"
		signMethod = "signMessageV2"
		message = buildTRONBindMessage(currentUserID, walletAddress, req.ChainID, challengeID, nonce, issuedAt, expiresAt)
	}

	challenge := &model.WalletBindingChallenge{
		ChallengeID:   challengeID,
		UserID:        currentUserID,
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
	if err := s.db.CreateWalletBindingChallenge(ctx, challenge); err != nil {
		return nil, err
	}

	return &pbredpacket.IssueWalletBindChallengeResp{
		ChallengeID: challengeID,
		UserID:      currentUserID,
		ChainType:   chainType,
		ChainID:     req.ChainID,
		Wallet:      walletAddress,
		Protocol:    protocol,
		SignMethod:  signMethod,
		Nonce:       nonce,
		Message:     message,
		IssuedAt:    issuedAt.Format(time.RFC3339),
		ExpiresAt:   expiresAt.Format(time.RFC3339),
	}, nil
}

func (s *redPacketServer) ConfirmWalletBind(ctx context.Context, req *pbredpacket.ConfirmWalletBindReq) (*pbredpacket.ConfirmWalletBindResp, error) {
	if strings.TrimSpace(req.ChallengeID) == "" || strings.TrimSpace(req.Signature) == "" {
		return nil, errs.ErrArgs.WrapMsg("challenge_id and signature are required")
	}
	challenge, err := s.db.GetWalletBindingChallenge(ctx, req.ChallengeID)
	if err != nil {
		return nil, err
	}
	if challenge.Status != "PENDING" {
		return nil, errs.ErrArgs.WrapMsg("challenge is not pending")
	}
	if time.Now().UTC().After(challenge.ExpiresAt) {
		challenge.Status = "EXPIRED"
		challenge.UpdatedAt = time.Now()
		_ = s.db.UpdateWalletBindingChallenge(ctx, challenge)
		return nil, errs.ErrArgs.WrapMsg("challenge is expired")
	}

	var verifyErr error
	switch challenge.ChainType {
	case "EVM":
		verifyErr = verifyEVMBindSignature(challenge.Message, challenge.WalletAddress, req.Signature)
	case "TRON":
		verifyErr = verifyTRONBindSignature(challenge.Message, challenge.WalletAddress, req.Signature)
	default:
		return nil, errs.ErrArgs.WrapMsg("unsupported chain_type: " + challenge.ChainType)
	}
	if verifyErr != nil {
		challenge.Status = "FAILED"
		challenge.Signature = req.Signature
		challenge.UpdatedAt = time.Now()
		_ = s.db.UpdateWalletBindingChallenge(ctx, challenge)
		return nil, verifyErr
	}

	now := time.Now().UTC()
	challenge.Status = "VERIFIED"
	challenge.Signature = req.Signature
	challenge.VerifiedAt = &now
	challenge.UpdatedAt = now
	if err := s.db.UpdateWalletBindingChallenge(ctx, challenge); err != nil {
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
	if err := s.db.UpsertWalletBinding(ctx, binding); err != nil {
		return nil, err
	}

	return &pbredpacket.ConfirmWalletBindResp{
		UserID:        binding.UserID,
		ChainType:     binding.ChainType,
		ChainID:       binding.ChainID,
		WalletAddress: binding.WalletAddress,
		Status:        binding.Status,
		VerifiedAt:    binding.VerifiedAt.Format(time.RFC3339),
	}, nil
}

func (s *redPacketServer) GetWalletBinding(ctx context.Context, req *pbredpacket.GetWalletBindingReq) (*pbredpacket.GetWalletBindingResp, error) {
	currentUserID := mcontext.GetOpUserID(ctx)
	if currentUserID == "" {
		return nil, servererrs.ErrNoPermission.WrapMsg("op user id is empty")
	}

	normalizedChainType, err := normalizeChainType(req.ChainType)
	if err != nil {
		return nil, err
	}
	binding, err := s.db.GetActiveWalletBinding(ctx, currentUserID, normalizedChainType, req.WalletAddress)
	if err != nil {
		return nil, err
	}
	return &pbredpacket.GetWalletBindingResp{
		UserID:        binding.UserID,
		ChainType:     binding.ChainType,
		ChainID:       binding.ChainID,
		WalletAddress: binding.WalletAddress,
		Status:        binding.Status,
		ChallengeID:   binding.ChallengeID,
		VerifiedAt:    binding.VerifiedAt.Format(time.RFC3339),
	}, nil
}

func buildEVMBindMessage(userID, walletAddress, domainIn, uriIn string, chainID int64, challengeID, nonce string, issuedAt, expiresAt time.Time) string {
	domain := strings.TrimSpace(domainIn)
	if domain == "" {
		domain = "redpacket"
	}
	uri := strings.TrimSpace(uriIn)
	if uri == "" {
		uri = "https://redpacket.local/wallet-bind"
	}

	var b strings.Builder
	fmt.Fprintf(&b, "%s wants you to sign in with your Ethereum account:\n", domain)
	b.WriteString(strings.TrimSpace(walletAddress))
	b.WriteString("\n\n")
	fmt.Fprintf(&b, "Bind wallet %s to user %s.\n", strings.TrimSpace(walletAddress), strings.TrimSpace(userID))
	fmt.Fprintf(&b, "URI: %s\n", uri)
	fmt.Fprintf(&b, "Version: 1\n")
	fmt.Fprintf(&b, "Chain ID: %d\n", chainID)
	fmt.Fprintf(&b, "Nonce: %s\n", nonce)
	fmt.Fprintf(&b, "Issued At: %s\n", issuedAt.Format(time.RFC3339))
	fmt.Fprintf(&b, "Expiration Time: %s\n", expiresAt.Format(time.RFC3339))
	fmt.Fprintf(&b, "Request ID: %s", challengeID)
	return b.String()
}

func buildTRONBindMessage(userID, walletAddress string, chainID int64, challengeID, nonce string, issuedAt, expiresAt time.Time) string {
	return fmt.Sprintf(
		"Bind TRON wallet %s to user %s\nchallenge_id: %s\nnonce: %s\nchain_id: %d\nissued_at: %s\nexpires_at: %s",
		strings.TrimSpace(walletAddress),
		strings.TrimSpace(userID),
		challengeID,
		nonce,
		chainID,
		issuedAt.Format(time.RFC3339),
		expiresAt.Format(time.RFC3339),
	)
}

func verifyEVMBindSignature(message, walletAddress, signature string) error {
	if strings.TrimSpace(message) == "" {
		return errs.ErrArgs.WrapMsg("bind message is empty")
	}
	if !common.IsHexAddress(walletAddress) {
		return errs.ErrArgs.WrapMsg("invalid evm wallet address")
	}

	sig, err := hex.DecodeString(strings.TrimPrefix(signature, "0x"))
	if err != nil {
		return errs.ErrArgs.WrapMsg("decode signature failed: " + err.Error())
	}
	if len(sig) != 65 {
		return errs.ErrArgs.WrapMsg(fmt.Sprintf("invalid signature length: %d", len(sig)))
	}
	if sig[64] >= 27 {
		sig[64] -= 27
	}
	if sig[64] > 1 {
		return errs.ErrArgs.WrapMsg("invalid signature recovery id")
	}

	hash := crypto.Keccak256Hash([]byte(personalSignMessage(message)))
	pubKey, err := crypto.SigToPub(hash.Bytes(), sig)
	if err != nil {
		return errs.ErrInternalServer.WrapMsg("recover signer failed: " + err.Error())
	}

	recovered := crypto.PubkeyToAddress(*pubKey)
	if !strings.EqualFold(recovered.Hex(), walletAddress) {
		return errs.ErrNoPermission.WrapMsg("signature does not match wallet address")
	}
	return nil
}

func personalSignMessage(message string) string {
	return fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
}

// verifyTRONBindSignature verifies a TRON signMessageV2 (TronLink) signature.
// TRON uses the same secp256k1 curve as Ethereum; the only differences are:
//   - message prefix: "\x19TRON Signed Message:\n<decimal_len>"
//   - wallet address: base58check-encoded with a leading 0x41 byte
func verifyTRONBindSignature(message, walletAddress, signature string) error {
	if strings.TrimSpace(message) == "" {
		return errs.ErrArgs.WrapMsg("bind message is empty")
	}

	sig, err := hex.DecodeString(strings.TrimPrefix(signature, "0x"))
	if err != nil {
		return errs.ErrArgs.WrapMsg("decode tron signature failed: " + err.Error())
	}
	if len(sig) != 65 {
		return errs.ErrArgs.WrapMsg(fmt.Sprintf("invalid tron signature length: %d", len(sig)))
	}
	// Some TRON wallets encode v as 27/28; normalise to 0/1.
	if sig[64] >= 27 {
		sig[64] -= 27
	}

	prefix := fmt.Sprintf("\x19TRON Signed Message:\n%d", len(message))
	hash := crypto.Keccak256Hash([]byte(prefix + message))

	pubKey, err := crypto.SigToPub(hash.Bytes(), sig)
	if err != nil {
		return errs.ErrInternalServer.WrapMsg("recover tron signer failed: " + err.Error())
	}

	// Derive the raw 20-byte address (identical derivation to Ethereum).
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	// Decode the TRON base58check address to its 20 raw bytes.
	addrBytes, err := decodeTRONAddress(walletAddress)
	if err != nil {
		return errs.ErrArgs.WrapMsg("invalid tron address: " + err.Error())
	}

	if !bytes.Equal(recoveredAddr.Bytes(), addrBytes) {
		return errs.ErrNoPermission.WrapMsg("tron signature does not match wallet address")
	}
	return nil
}

// decodeTRONAddress decodes a TRON base58check address and returns the 20
// raw address bytes (i.e., without the leading 0x41 network prefix byte).
func decodeTRONAddress(addr string) ([]byte, error) {
	decoded := tronBase58Decode(addr)
	if len(decoded) != 25 {
		return nil, fmt.Errorf("invalid length %d", len(decoded))
	}

	payload := decoded[:21]
	checksum := decoded[21:25]
	h1 := sha256.Sum256(payload)
	h2 := sha256.Sum256(h1[:])
	if !bytes.Equal(h2[:4], checksum) {
		return nil, fmt.Errorf("invalid base58check checksum")
	}
	if payload[0] != 0x41 {
		return nil, fmt.Errorf("invalid tron address prefix byte: 0x%02x", payload[0])
	}
	return payload[1:], nil
}

const tronBase58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func tronBase58Decode(s string) []byte {
	n := new(big.Int)
	base := big.NewInt(58)
	for _, c := range s {
		idx := strings.IndexRune(tronBase58Alphabet, c)
		if idx < 0 {
			return nil
		}
		n.Mul(n, base)
		n.Add(n, big.NewInt(int64(idx)))
	}

	decoded := n.Bytes()
	leadingOnes := 0
	for _, c := range s {
		if c == '1' {
			leadingOnes++
		} else {
			break
		}
	}
	out := make([]byte, leadingOnes+len(decoded))
	copy(out[leadingOnes:], decoded)
	return out
}
