package chain

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func TestTronLogsToEVMLogsAndParsePacketCreated(t *testing.T) {
	abiJSON, err := ExtractABIFromEmbeddedArtifact()
	if err != nil {
		t.Fatalf("ExtractABIFromEmbeddedArtifact() error = %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(abiJSON)))
	if err != nil {
		t.Fatalf("abi.JSON() error = %v", err)
	}

	eventDef := parsedABI.Events["PacketCreated"]
	packetID := big.NewInt(12)
	creator := common.HexToAddress("0x1111111111111111111111111111111111111111")
	packetType := big.NewInt(1)
	token := common.HexToAddress("0x2222222222222222222222222222222222222222")
	totalAmount := big.NewInt(1000)
	totalShares := big.NewInt(10)
	expiryAt := big.NewInt(1234567890)

	data, err := eventDef.Inputs.NonIndexed().Pack(token, totalAmount, totalShares, expiryAt)
	if err != nil {
		t.Fatalf("Pack() error = %v", err)
	}

	info := &tronTxInfoResp{
		ID:          "abc123",
		BlockNumber: 88,
		Log: []struct {
			Address string   `json:"address"`
			Topics  []string `json:"topics"`
			Data    string   `json:"data"`
		}{
			{
				Address: "41aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				Topics: []string{
					strings.TrimPrefix(eventDef.ID.Hex(), "0x"),
					strings.TrimPrefix(common.BigToHash(packetID).Hex(), "0x"),
					strings.TrimPrefix(common.BytesToHash(common.LeftPadBytes(creator.Bytes(), 32)).Hex(), "0x"),
					strings.TrimPrefix(common.BigToHash(packetType).Hex(), "0x"),
				},
				Data: common.Bytes2Hex(data),
			},
		},
	}

	logs, err := tronLogsToEVMLogs(info, info.ID)
	if err != nil {
		t.Fatalf("tronLogsToEVMLogs() error = %v", err)
	}

	events, err := ParseEventsFromLogs(logs, parsedABI)
	if err != nil {
		t.Fatalf("ParseEventsFromLogs() error = %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	event := events[0]
	if event.Name != "PacketCreated" {
		t.Fatalf("unexpected event name: %s", event.Name)
	}
	if got := GetPacketIDFromEvent(event).String(); got != "12" {
		t.Fatalf("packet id mismatch: got %s", got)
	}
	if got := GetAddressFromEvent(event, "creator").Hex(); got != creator.Hex() {
		t.Fatalf("creator mismatch: got %s want %s", got, creator.Hex())
	}
	if got := GetUintFromEvent(event, "packetType").String(); got != "1" {
		t.Fatalf("packetType mismatch: got %s", got)
	}
	if got := GetAddressFromEvent(event, "token").Hex(); got != token.Hex() {
		t.Fatalf("token mismatch: got %s want %s", got, token.Hex())
	}
	if event.BlockNumber != 88 {
		t.Fatalf("block number mismatch: got %d", event.BlockNumber)
	}
}
