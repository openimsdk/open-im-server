package chain

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// TronClient handles TRON blockchain interactions using HTTP JSON-RPC
type TronClient struct {
	fullNodeURL   string
	contractBase58 string
	ownerBase58   string
	privateKeyHex string
	feeLimit      int64
	abiJSON       string
	parsedABI     abi.ABI
}

// NewTronClient creates a new TRON client
func NewTronClient(fullNodeURL, contractBase58, ownerBase58, privateKeyHex string, abiJSON []byte, feeLimit int64) (*TronClient, error) {
	if fullNodeURL == "" {
		return nil, fmt.Errorf("fullNodeURL is required for TRON")
	}

	parsedABI, err := abi.JSON(bytes.NewReader(abiJSON))
	if err != nil {
		return nil, fmt.Errorf("parse TRON ABI failed: %w", err)
	}

	return &TronClient{
		fullNodeURL:    fullNodeURL,
		contractBase58: contractBase58,
		ownerBase58:    ownerBase58,
		privateKeyHex:  privateKeyHex,
		feeLimit:       feeLimit,
		abiJSON:        string(abiJSON),
		parsedABI:      parsedABI,
	}, nil
}

// SendAdminTransaction sends an admin transaction on TRON (setSigner, setToken, etc.)
func (t *TronClient) SendAdminTransaction(ctx context.Context, methodName string, args ...interface{}) (string, error) {
	if t.privateKeyHex == "" || t.ownerBase58 == "" {
		return "", fmt.Errorf("TRON admin credentials not configured")
	}

	// Build function selector like "setSigner(address)"
	selector := methodName
	if len(args) > 0 {
		// Simple selector generation - in production use full ABI encoding
		selector = fmt.Sprintf("%s(%s)", methodName, getParamTypes(args))
	}

	if _, encodeErr := encodeTronParams(t.abiJSON, methodName, args...); encodeErr != nil {
		return "", fmt.Errorf("encode params failed: %w", encodeErr)
	}

	return SendTronAdminTx(
		ctx,
		t.fullNodeURL,
		t.ownerBase58,
		t.contractBase58,
		selector,
		methodName,
		t.feeLimit,
		t.privateKeyHex,
		t.abiJSON,
		args...,
	)
}

// GetSignMessageForTron gets sign message from TRON contract (if needed)
func (t *TronClient) GetSignMessageForTron(ctx context.Context, packetID *big.Int, claimer, authNonce, randomSeed, deadline string) (string, error) {
	// TRON version would call triggersmartcontract with getSignMessage
	// For simplicity, we can reuse similar logic as ETH or implement full TRON trigger
	return "", fmt.Errorf("TRON getSignMessage not fully implemented yet - use ETH path for signing")
}

// Helper functions

func getParamTypes(args []interface{}) string {
	types := make([]string, len(args))
	for i, arg := range args {
		switch arg.(type) {
		case string, common.Address:
			types[i] = "address"
		case bool:
			types[i] = "bool"
		case int, int64, *big.Int:
			types[i] = "uint256"
		default:
			types[i] = "unknown"
		}
	}
	return strings.Join(types, ",")
}

// SendTronAdminTx implements TRON transaction broadcasting (from design doc)
func SendTronAdminTx(
	ctx context.Context,
	fullNodeURL, ownerBase58, contractBase58, selector, methodName string,
	feeLimit int64,
	privateKeyHex string,
	abiJSON string,
	args ...interface{},
) (string, error) {

	paramHex, err := encodeTronParams(abiJSON, methodName, args...)
	if err != nil {
		return "", err
	}

	// Trigger smart contract
	var triggerResp map[string]interface{}
	err = postJSON(ctx, fullNodeURL+"/wallet/triggersmartcontract", map[string]interface{}{
		"owner_address":     ownerBase58,
		"contract_address":  contractBase58,
		"function_selector": selector,
		"parameter":         paramHex,
		"fee_limit":         feeLimit,
		"call_value":        0,
		"visible":           true,
	}, &triggerResp)
	if err != nil {
		return "", fmt.Errorf("trigger contract failed: %w", err)
	}

	txObj, ok := triggerResp["transaction"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("transaction not found in trigger response")
	}

	// Sign transaction
	var signedResp map[string]interface{}
	err = postJSON(ctx, fullNodeURL+"/wallet/gettransactionsign", map[string]interface{}{
		"transaction": txObj,
		"privateKey":  privateKeyHex,
	}, &signedResp)
	if err != nil {
		return "", fmt.Errorf("sign transaction failed: %w", err)
	}

	// Broadcast
	var broadcastResp map[string]interface{}
	err = postJSON(ctx, fullNodeURL+"/wallet/broadcasttransaction", signedResp, &broadcastResp)
	if err != nil {
		return "", fmt.Errorf("broadcast failed: %w", err)
	}

	if result, _ := broadcastResp["result"].(bool); !result {
		return "", fmt.Errorf("broadcast failed: %v", broadcastResp)
	}

	txid, _ := broadcastResp["txid"].(string)
	return txid, nil
}

func encodeTronParams(abiJSON, method string, args ...interface{}) (string, error) {
	parsed, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return "", err
	}
	m, ok := parsed.Methods[method]
	if !ok {
		return "", fmt.Errorf("method not found: %s", method)
	}
	packed, err := m.Inputs.Pack(args...)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(packed), nil
}

func postJSON(ctx context.Context, url string, body interface{}, out interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("http %d: %s", resp.StatusCode, string(raw))
	}

	if err := json.Unmarshal(raw, out); err != nil {
		return err
	}
	return nil
}
