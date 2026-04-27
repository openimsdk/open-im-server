package service

import (
	"context"
	"fmt"
	"math/big"

	"redpacket/internal/chain"

	"github.com/ethereum/go-ethereum/common"
)

// AdminService handles administrative operations on the RedPacket contract
type AdminService struct {
	ethClient  *chain.ChainClient
	tronClient *chain.TronClient
}

func NewAdminService(ethClient *chain.ChainClient, tronClient *chain.TronClient) *AdminService {
	return &AdminService{
		ethClient:  ethClient,
		tronClient: tronClient,
	}
}

func (s *AdminService) SetSigner(ctx context.Context, signerAddress string) error {
	if s.ethClient != nil {
		// For ETH: call setSigner through contract
		// In real implementation this would use admin key to send transaction
		fmt.Printf("ETH: Setting signer to %s (mock)\n", signerAddress)
		return nil
	}

	if s.tronClient != nil {
		_, err := s.tronClient.SendAdminTransaction(ctx, "setSigner", signerAddress)
		return err
	}

	return fmt.Errorf("no blockchain client configured")
}

func (s *AdminService) SetToken(ctx context.Context, tokenAddress string, allowed bool, minAmount string) error {
	minAmountBig := new(big.Int)
	if minAmount != "" {
		minAmountBig.SetString(minAmount, 10)
	} else {
		minAmountBig.SetInt64(0)
	}

	if s.ethClient != nil {
		fmt.Printf("ETH: Setting token %s allowed=%v minAmount=%s (mock)\n", tokenAddress, allowed, minAmount)
		return nil
	}

	if s.tronClient != nil {
		_, err := s.tronClient.SendAdminTransaction(ctx, "setAllowedToken", tokenAddress, allowed, minAmountBig)
		return err
	}

	return fmt.Errorf("no blockchain client configured")
}

func (s *AdminService) SetExpiry(ctx context.Context, expirySeconds int64) error {
	if s.ethClient != nil {
		fmt.Printf("ETH: Setting default expiry to %d seconds (mock)\n", expirySeconds)
		return nil
	}

	if s.tronClient != nil {
		_, err := s.tronClient.SendAdminTransaction(ctx, "setDefaultExpiryDuration", expirySeconds)
		return err
	}

	return fmt.Errorf("no blockchain client configured")
}

func (s *AdminService) SetAllowAllTokens(ctx context.Context, allowAll bool) error {
	if s.ethClient != nil {
		fmt.Printf("ETH: Setting allowAllTokens=%v (mock)\n", allowAll)
		return nil
	}

	if s.tronClient != nil {
		_, err := s.tronClient.SendAdminTransaction(ctx, "setAllowAllTokens", allowAll)
		return err
	}

	return fmt.Errorf("no blockchain client configured")
}

func (s *AdminService) SetNativeTokenEnabled(ctx context.Context, enabled bool) error {
	if s.ethClient != nil {
		fmt.Printf("ETH: Setting native token enabled=%v (mock)\n", enabled)
		return nil
	}

	if s.tronClient != nil {
		_, err := s.tronClient.SendAdminTransaction(ctx, "setNativeTokenEnabled", enabled)
		return err
	}

	return fmt.Errorf("no blockchain client configured")
}

func (s *AdminService) ParseTxEvents(ctx context.Context, txHash, chain string) (map[string]interface{}, error) {
	if chain == "tron" && s.tronClient != nil {
		return map[string]interface{}{
			"chain":   "tron",
			"tx_hash": txHash,
			"status":  "parsed",
			"note":    "TRON event parsing not fully implemented in this version",
		}, nil
	}

	if s.ethClient != nil {
		txHashBytes := common.HexToHash(txHash)
		events, err := s.ethClient.ParseTransactionReceipt(ctx, txHashBytes)
		if err != nil {
			return nil, err
		}

		eventList := make([]map[string]interface{}, len(events))
		for i, e := range events {
			eventList[i] = map[string]interface{}{
				"name": e.Name,
				"data": e.Data,
			}
		}

		return map[string]interface{}{
			"chain":   "eth",
			"tx_hash": txHash,
			"events":  eventList,
		}, nil
	}

	return nil, fmt.Errorf("no client available for chain: %s", chain)
}
