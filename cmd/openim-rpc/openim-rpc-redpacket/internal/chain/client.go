package chain

import (
	"context"
	"crypto/ecdsa"
	_ "embed"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

//go:embed abi/RedPacket.json
var embeddedABI []byte

// ChainClient handles blockchain interactions for RedPacket
type ChainClient struct {
	client         *ethclient.Client
	contractABI    abi.ABI
	contractAddr   common.Address
	signerKey      *ecdsa.PrivateKey
	configAdminKey *ecdsa.PrivateKey
	chainID        *big.Int
}

// NewClient creates a new ChainClient
func NewClient(rpcURL, contractAddress string, chainID int64, signerPrivateKey, configAdminPrivateKey string) (*ChainClient, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ethereum: %w", err)
	}

	// Load ABI
	abiJSON, err := ExtractABIFromEmbeddedArtifact()
	if err != nil {
		return nil, fmt.Errorf("failed to load ABI: %w", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(abiJSON)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	contractAddr := common.HexToAddress(contractAddress)

	var signerKey *ecdsa.PrivateKey
	if signerPrivateKey != "" {
		signerKey, err = crypto.HexToECDSA(strings.TrimPrefix(signerPrivateKey, "0x"))
		if err != nil {
			return nil, fmt.Errorf("invalid signer private key: %w", err)
		}
	}

	var adminKey *ecdsa.PrivateKey
	if configAdminPrivateKey != "" {
		adminKey, err = crypto.HexToECDSA(strings.TrimPrefix(configAdminPrivateKey, "0x"))
		if err != nil {
			return nil, fmt.Errorf("invalid config admin private key: %w", err)
		}
	}

	return &ChainClient{
		client:         client,
		contractABI:    parsedABI,
		contractAddr:   contractAddr,
		signerKey:      signerKey,
		configAdminKey: adminKey,
		chainID:        big.NewInt(chainID),
	}, nil
}

// GetSignMessage calls contract's getSignMessage view function
func (c *ChainClient) GetSignMessage(ctx context.Context, packetID *big.Int, claimer common.Address, authNonce, randomSeed, deadline *big.Int) ([32]byte, error) {
	var digest [32]byte

	data, err := c.contractABI.Pack("getSignMessage", packetID, claimer, authNonce, randomSeed, deadline)
	if err != nil {
		return digest, fmt.Errorf("failed to pack getSignMessage: %w", err)
	}

	msg := ethereum.CallMsg{
		To:   &c.contractAddr,
		Data: data,
	}

	result, err := c.client.CallContract(ctx, msg, nil)
	if err != nil {
		return digest, fmt.Errorf("call getSignMessage failed: %w", err)
	}

	copy(digest[:], result)
	return digest, nil
}

// SignClaim signs the digest using the signer key (naked signature as per contract)
func (c *ChainClient) SignClaim(digest [32]byte) ([]byte, error) {
	if c.signerKey == nil {
		return nil, fmt.Errorf("signer key not configured")
	}

	sig, err := crypto.Sign(digest[:], c.signerKey)
	if err != nil {
		return nil, fmt.Errorf("sign failed: %w", err)
	}

	// Adjust v from 0/1 to 27/28 as expected by EVM
	if len(sig) == 65 && sig[64] < 27 {
		sig[64] += 27
	}

	return sig, nil
}

// ParseTransactionReceipt parses events from a transaction receipt
func (c *ChainClient) ParseTransactionReceipt(ctx context.Context, txHash common.Hash) ([]*ParsedEvent, error) {
	receipt, err := c.client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("get receipt failed: %w", err)
	}

	return ParseEventsFromLogs(receipt.Logs, c.contractABI)
}

func (c *ChainClient) ContractAddress() common.Address {
	return c.contractAddr
}

func (c *ChainClient) ChainID() *big.Int {
	if c.chainID == nil {
		return nil
	}
	return new(big.Int).Set(c.chainID)
}

// Close closes the client connection
func (c *ChainClient) Close() {
	if c.client != nil {
		c.client.Close()
	}
}

// ExtractABIFromEmbeddedArtifact returns the embedded contract ABI
func ExtractABIFromEmbeddedArtifact() ([]byte, error) {
	if len(embeddedABI) == 0 {
		return nil, fmt.Errorf("embedded ABI is empty")
	}
	return embeddedABI, nil
}
