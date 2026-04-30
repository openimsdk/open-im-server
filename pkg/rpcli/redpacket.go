package rpcli

import (
	pbredpacket "github.com/openimsdk/protocol/redpacket"
	"google.golang.org/grpc"
)

func NewRedPacketClient(cc grpc.ClientConnInterface) *RedPacketClient {
	return &RedPacketClient{pbredpacket.NewRedPacketClient(cc)}
}

type RedPacketClient struct {
	pbredpacket.RedPacketClient
}
