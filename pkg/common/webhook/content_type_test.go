package webhook

import (
	"testing"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/sdkws"
)

func TestFilterAfterMsg(t *testing.T) {
	tests := []struct {
		name            string
		msg             *sdkws.MsgData
		after           *config.AfterConfig
		attentionTarget string
		want            bool
	}{
		{
			name:  "empty filters allow message",
			msg:   &sdkws.MsgData{SendID: "sender", ContentType: 1},
			after: &config.AfterConfig{},
			want:  true,
		},
		{
			name:            "attention matches sender",
			msg:             &sdkws.MsgData{SendID: "sender", ContentType: 1},
			after:           &config.AfterConfig{AttentionIds: []string{"sender"}},
			attentionTarget: "other",
			want:            true,
		},
		{
			name:            "attention matches explicit target",
			msg:             &sdkws.MsgData{SendID: "sender", ContentType: 1},
			after:           &config.AfterConfig{AttentionIds: []string{"target"}},
			attentionTarget: "target",
			want:            true,
		},
		{
			name:            "attention rejects unrelated ids",
			msg:             &sdkws.MsgData{SendID: "sender", ContentType: 1},
			after:           &config.AfterConfig{AttentionIds: []string{"unrelated"}},
			attentionTarget: "target",
			want:            false,
		},
		{
			name:  "allowed type hit",
			msg:   &sdkws.MsgData{ContentType: 7},
			after: &config.AfterConfig{AllowedTypes: []string{"1-7"}},
			want:  true,
		},
		{
			name:  "allowed type miss",
			msg:   &sdkws.MsgData{ContentType: 8},
			after: &config.AfterConfig{AllowedTypes: []string{"1-7"}},
			want:  false,
		},
		{
			name:  "denied type hit",
			msg:   &sdkws.MsgData{ContentType: 8},
			after: &config.AfterConfig{DeniedTypes: []string{"8"}},
			want:  false,
		},
		{
			name: "denied wins over allowed",
			msg:  &sdkws.MsgData{ContentType: 8},
			after: &config.AfterConfig{
				AllowedTypes: []string{"1-10"},
				DeniedTypes:  []string{"8"},
			},
			want: false,
		},
		{
			name:  "malformed allowed interval rejects",
			msg:   &sdkws.MsgData{ContentType: 8},
			after: &config.AfterConfig{AllowedTypes: []string{"bad-interval"}},
			want:  false,
		},
		{
			name:  "malformed denied interval does not reject",
			msg:   &sdkws.MsgData{ContentType: 8},
			after: &config.AfterConfig{DeniedTypes: []string{"bad-interval"}},
			want:  true,
		},
		{
			name:  "typing is not hardcoded denied",
			msg:   &sdkws.MsgData{ContentType: constant.Typing},
			after: &config.AfterConfig{AllowedTypes: []string{"113"}},
			want:  true,
		},
		{
			name: "notification is not hardcoded denied",
			msg:  &sdkws.MsgData{ContentType: constant.GroupCreatedNotification},
			after: &config.AfterConfig{
				AllowedTypes: []string{"1501"},
			},
			want: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := FilterAfterMsg(test.msg, test.after, test.attentionTarget); got != test.want {
				t.Fatalf("FilterAfterMsg() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestFilterBeforeMsg(t *testing.T) {
	tests := []struct {
		name   string
		msg    *sdkws.MsgData
		before *config.BeforeConfig
		want   bool
	}{
		{
			name:   "empty filters allow message",
			msg:    &sdkws.MsgData{ContentType: 1},
			before: &config.BeforeConfig{},
			want:   true,
		},
		{
			name:   "allowed type hit",
			msg:    &sdkws.MsgData{ContentType: 7},
			before: &config.BeforeConfig{AllowedTypes: []string{"1-7"}},
			want:   true,
		},
		{
			name:   "allowed type miss",
			msg:    &sdkws.MsgData{ContentType: 8},
			before: &config.BeforeConfig{AllowedTypes: []string{"1-7"}},
			want:   false,
		},
		{
			name:   "denied type hit",
			msg:    &sdkws.MsgData{ContentType: 8},
			before: &config.BeforeConfig{DeniedTypes: []string{"8"}},
			want:   false,
		},
		{
			name: "denied wins over allowed",
			msg:  &sdkws.MsgData{ContentType: 8},
			before: &config.BeforeConfig{
				AllowedTypes: []string{"1-10"},
				DeniedTypes:  []string{"8"},
			},
			want: false,
		},
		{
			name:   "malformed allowed interval rejects",
			msg:    &sdkws.MsgData{ContentType: 8},
			before: &config.BeforeConfig{AllowedTypes: []string{"bad-interval"}},
			want:   false,
		},
		{
			name:   "malformed denied interval does not reject",
			msg:    &sdkws.MsgData{ContentType: 8},
			before: &config.BeforeConfig{DeniedTypes: []string{"bad-interval"}},
			want:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := FilterBeforeMsg(test.msg, test.before); got != test.want {
				t.Fatalf("FilterBeforeMsg() = %v, want %v", got, test.want)
			}
		})
	}
}
