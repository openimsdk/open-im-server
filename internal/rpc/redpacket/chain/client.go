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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

//go:embed abi/RedPacket.json
var embeddedABI []byte

// ChainClient handles blockchain interactions for RedPacket.
type ChainClient struct {
	client         *ethclient.Client
	contractABI    abi.ABI
	contractAddr   common.Address
	signerKey      *ecdsa.PrivateKey
	configAdminKey *ecdsa.PrivateKey
	chainID        *big.Int
}

func NewClient(rpcURL, contractAddress string, chainID int64, signerPrivateKey, configAdminPrivateKey string) (*ChainClient, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ethereum: %w", err)
	}

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

func (c *ChainClient) SignClaim(digest [32]byte) ([]byte, error) {
	if c.signerKey == nil {
		return nil, fmt.Errorf("signer key not configured")
	}

	sig, err := crypto.Sign(digest[:], c.signerKey)
	if err != nil {
		return nil, fmt.Errorf("sign failed: %w", err)
	}

	if len(sig) == 65 && sig[64] < 27 {
		sig[64] += 27
	}

	return sig, nil
}

func (c *ChainClient) ParseTransactionReceipt(ctx context.Context, txHash common.Hash) ([]*ParsedEvent, error) {
	_, events, err := c.ParseTransactionReceiptWithStatus(ctx, txHash)
	return events, err
}

// ParseTransactionReceiptWithStatus fetches tx receipt once and returns both
// execution status and decoded contract events.
func (c *ChainClient) ParseTransactionReceiptWithStatus(ctx context.Context, txHash common.Hash) (bool, []*ParsedEvent, error) {
	receipt, err := c.client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return false, nil, fmt.Errorf("get receipt failed: %w", err)
	}
	events, err := ParseEventsFromLogs(receipt.Logs, c.contractABI)
	if err != nil {
		return false, nil, err
	}
	return receipt.Status == types.ReceiptStatusSuccessful, events, nil
}

// IsTransactionSuccessful reports whether the EVM transaction executed
// successfully according to receipt.status (1=success, 0=failure).
func (c *ChainClient) IsTransactionSuccessful(ctx context.Context, txHash common.Hash) (bool, error) {
	success, _, err := c.ParseTransactionReceiptWithStatus(ctx, txHash)
	return success, err
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

// EthClient exposes the underlying ethclient for indexers.
func (c *ChainClient) EthClient() *ethclient.Client {
	return c.client
}

// ContractABI exposes the parsed ABI for indexers.
func (c *ChainClient) ContractABI() abi.ABI {
	return c.contractABI
}

// RefundPacket submits an on-chain refund transaction for an expired red
// packet. It uses the configAdminKey to sign and broadcast the transaction.
// Returns the transaction hash on success.
func (c *ChainClient) RefundPacket(ctx context.Context, packetIDStr string) (string, error) {
	if c.configAdminKey == nil {
		return "", fmt.Errorf("config admin key not configured")
	}

	packetID, ok := new(big.Int).SetString(packetIDStr, 10)
	if !ok {
		return "", fmt.Errorf("invalid packetID: %s", packetIDStr)
	}

	data, err := c.contractABI.Pack("refundPacket", packetID)
	if err != nil {
		return "", fmt.Errorf("pack refundPacket failed: %w", err)
	}

	fromAddr := crypto.PubkeyToAddress(c.configAdminKey.PublicKey)
	nonce, err := c.client.PendingNonceAt(ctx, fromAddr)
	if err != nil {
		return "", fmt.Errorf("get nonce failed: %w", err)
	}

	gasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("suggest gas price failed: %w", err)
	}

	gasLimit, err := c.client.EstimateGas(ctx, ethereum.CallMsg{
		From: fromAddr,
		To:   &c.contractAddr,
		Data: data,
	})
	if err != nil {
		gasLimit = 200000 // fallback
	}

	tx := types.NewTransaction(nonce, c.contractAddr, big.NewInt(0), gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(c.chainID), c.configAdminKey)
	if err != nil {
		return "", fmt.Errorf("sign refund tx failed: %w", err)
	}

	if err := c.client.SendTransaction(ctx, signedTx); err != nil {
		return "", fmt.Errorf("send refund tx failed: %w", err)
	}

	return signedTx.Hash().Hex(), nil
}

func (c *ChainClient) Close() {
	if c.client != nil {
		c.client.Close()
	}
}

func ExtractABIFromEmbeddedArtifact() ([]byte, error) {
	if len(embeddedABI) == 0 {
		return nil, fmt.Errorf("embedded ABI is empty")
	}
	return embeddedABI, nil
}
