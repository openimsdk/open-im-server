package rpcli

import (
	"github.com/openimsdk/protocol/rtc"
)

func NewRtcServiceClient(cli rtc.RtcServiceClient) *RtcServiceClient {
	return &RtcServiceClient{cli}
}

type RtcServiceClient struct {
	rtc.RtcServiceClient
}
