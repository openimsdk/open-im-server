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
	"github.com/ethereum/go-ethereum/core/types"
)

type TronClient struct {
	fullNodeURL    string
	contractBase58 string
	ownerBase58    string
	privateKeyHex  string
	feeLimit       int64
	abiJSON        string
	parsedABI      abi.ABI
}

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

func (t *TronClient) ContractAddress() string {
	return t.contractBase58
}

// ContractBase58 exposes the contract base58 address for indexers.
func (t *TronClient) ContractBase58() string {
	return t.contractBase58
}

// FullNodeURL exposes the full node URL for indexers.
func (t *TronClient) FullNodeURL() string {
	return t.fullNodeURL
}

func (t *TronClient) ParseTransactionReceipt(ctx context.Context, txID string) ([]*ParsedEvent, error) {
	info, err := t.getTransactionInfo(ctx, txID)
	if err != nil {
		return nil, err
	}

	logs, err := tronLogsToEVMLogs(info, txID)
	if err != nil {
		return nil, err
	}

	return ParseEventsFromLogs(logs, t.parsedABI)
}

func (t *TronClient) SendAdminTransaction(ctx context.Context, methodName string, args ...interface{}) (string, error) {
	if t.privateKeyHex == "" || t.ownerBase58 == "" {
		return "", fmt.Errorf("TRON admin credentials not configured")
	}

	selector := methodName
	if len(args) > 0 {
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

func (t *TronClient) GetSignMessageForTron(ctx context.Context, packetID *big.Int, claimer, authNonce, randomSeed, deadline string) (string, error) {
	return "", fmt.Errorf("TRON getSignMessage not fully implemented yet - use ETH path for signing")
}

type tronTxInfoResp struct {
	ID          string `json:"id"`
	BlockNumber uint64 `json:"blockNumber"`
	Log         []struct {
		Address string   `json:"address"`
		Topics  []string `json:"topics"`
		Data    string   `json:"data"`
	} `json:"log"`
}

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

	var signedResp map[string]interface{}
	err = postJSON(ctx, fullNodeURL+"/wallet/gettransactionsign", map[string]interface{}{
		"transaction": txObj,
		"privateKey":  privateKeyHex,
	}, &signedResp)
	if err != nil {
		return "", fmt.Errorf("sign transaction failed: %w", err)
	}

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

func (t *TronClient) getTransactionInfo(ctx context.Context, txID string) (*tronTxInfoResp, error) {
	var info tronTxInfoResp
	if err := postJSON(ctx, t.fullNodeURL+"/wallet/gettransactioninfobyid", map[string]interface{}{
		"value": txID,
	}, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func tronLogsToEVMLogs(info *tronTxInfoResp, txID string) ([]*types.Log, error) {
	if info == nil {
		return nil, fmt.Errorf("tron tx info is nil")
	}

	txHash := common.HexToHash(addHexPrefix(txID))
	logs := make([]*types.Log, 0, len(info.Log))
	for _, entry := range info.Log {
		topics := make([]common.Hash, 0, len(entry.Topics))
		for _, topic := range entry.Topics {
			topics = append(topics, common.HexToHash(addHexPrefix(topic)))
		}

		data, err := hex.DecodeString(strings.TrimPrefix(entry.Data, "0x"))
		if err != nil {
			return nil, fmt.Errorf("decode tron log data failed: %w", err)
		}

		logs = append(logs, &types.Log{
			Address:     tronLogAddressToCommonAddress(entry.Address),
			Topics:      topics,
			Data:        data,
			BlockNumber: info.BlockNumber,
			TxHash:      txHash,
		})
	}

	return logs, nil
}

func tronLogAddressToCommonAddress(raw string) common.Address {
	raw = strings.TrimPrefix(raw, "0x")
	raw = strings.TrimPrefix(raw, "41")
	if len(raw) > 40 {
		raw = raw[len(raw)-40:]
	}
	return common.HexToAddress(addHexPrefix(raw))
}

func addHexPrefix(value string) string {
	if strings.HasPrefix(value, "0x") || strings.HasPrefix(value, "0X") {
		return value
	}
	return "0x" + value
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
