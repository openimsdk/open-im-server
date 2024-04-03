// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"bytes"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/system/program"
	"time"
)

type LocalCache struct {
	User struct {
		Topic         string `mapstructure:"topic"`
		SlotNum       int    `mapstructure:"slotNum"`
		SlotSize      int    `mapstructure:"slotSize"`
		SuccessExpire int    `mapstructure:"successExpire"`
		FailedExpire  int    `mapstructure:"failedExpire"`
	} `mapstructure:"user"`
	Group struct {
		Topic         string `mapstructure:"topic"`
		SlotNum       int    `mapstructure:"slotNum"`
		SlotSize      int    `mapstructure:"slotSize"`
		SuccessExpire int    `mapstructure:"successExpire"`
		FailedExpire  int    `mapstructure:"failedExpire"`
	} `mapstructure:"group"`
	Friend struct {
		Topic         string `mapstructure:"topic"`
		SlotNum       int    `mapstructure:"slotNum"`
		SlotSize      int    `mapstructure:"slotSize"`
		SuccessExpire int    `mapstructure:"successExpire"`
		FailedExpire  int    `mapstructure:"failedExpire"`
	} `mapstructure:"friend"`
	Conversation struct {
		Topic         string `mapstructure:"topic"`
		SlotNum       int    `mapstructure:"slotNum"`
		SlotSize      int    `mapstructure:"slotSize"`
		SuccessExpire int    `mapstructure:"successExpire"`
		FailedExpire  int    `mapstructure:"failedExpire"`
	} `mapstructure:"conversation"`
}

type Log struct {
	StorageLocation     string `mapstructure:"storageLocation"`
	RotationTime        int    `mapstructure:"rotationTime"`
	RemainRotationCount int    `mapstructure:"remainRotationCount"`
	RemainLogLevel      int    `mapstructure:"remainLogLevel"`
	IsStdout            bool   `mapstructure:"isStdout"`
	IsJson              bool   `mapstructure:"isJson"`
	WithStack           bool   `mapstructure:"withStack"`
}

type Minio struct {
	Bucket          string `mapstructure:"bucket"`
	Port            int    `mapstructure:"port"`
	AccessKeyID     string `mapstructure:"accessKeyID"`
	SecretAccessKey string `mapstructure:"secretAccessKey"`
	SessionToken    string `mapstructure:"sessionToken"`
	InternalIP      string `mapstructure:"internalIP"`
	ExternalIP      string `mapstructure:"externalIP"`
	URL             string `mapstructure:"url"`
	PublicRead      bool   `mapstructure:"publicRead"`
}

type Mongo struct {
	URI         string   `mapstructure:"uri"`
	Address     []string `mapstructure:"address"`
	Database    string   `mapstructure:"database"`
	Username    string   `mapstructure:"username"`
	Password    string   `mapstructure:"password"`
	MaxPoolSize int      `mapstructure:"maxPoolSize"`
	MaxRetry    int      `mapstructure:"maxRetry"`
}

type API struct {
	Api struct {
		ListenIP string `mapstructure:"listenIP"`
		Ports    []int  `mapstructure:"ports"`
	} `mapstructure:"api"`
	Prometheus struct {
		Enable     bool   `mapstructure:"enable"`
		Ports      []int  `mapstructure:"ports"`
		GrafanaURL string `mapstructure:"grafanaURL"`
	} `mapstructure:"prometheus"`
}

type CronTask struct {
	ChatRecordsClearTime string `mapstructure:"chatRecordsClearTime"`
	MsgDestructTime      string `mapstructure:"msgDestructTime"`
	RetainChatRecords    int    `mapstructure:"retainChatRecords"`
}
type Notification struct {
	GroupCreated struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"groupCreated"`
	GroupInfoSet struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"groupInfoSet"`
	JoinGroupApplication struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"joinGroupApplication"`
	MemberQuit struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"memberQuit"`
	GroupApplicationAccepted struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"groupApplicationAccepted"`
	GroupApplicationRejected struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"groupApplicationRejected"`
	GroupOwnerTransferred struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"groupOwnerTransferred"`
	MemberKicked struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"memberKicked"`
	MemberInvited struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"memberInvited"`
	MemberEnter struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"memberEnter"`
	GroupDismissed struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"groupDismissed"`
	GroupMuted struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"groupMuted"`
	GroupCancelMuted struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
		DefaultTips struct {
			Tips string `mapstructure:"tips"`
		} `mapstructure:"defaultTips"`
	} `mapstructure:"

groupCancelMuted"`
	GroupMemberMuted struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"groupMemberMuted"`
	GroupMemberCancelMuted struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"groupMemberCancelMuted"`
	GroupMemberInfoSet struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"groupMemberInfoSet"`
	GroupInfoSetAnnouncement struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"groupInfoSetAnnouncement"`
	GroupInfoSetName struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"groupInfoSetName"`
	FriendApplicationAdded struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"friendApplicationAdded"`
	FriendApplicationApproved struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"friendApplicationApproved"`
	FriendApplicationRejected struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"friendApplicationRejected"`
	FriendAdded struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"friendAdded"`
	FriendDeleted struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"friendDeleted"`
	FriendRemarkSet struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"friendRemarkSet"`
	BlackAdded struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"blackAdded"`
	BlackDeleted struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"blackDeleted"`
	FriendInfoUpdated struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"re

liabilityLevel"`
		UnreadCount bool `mapstructure:"unreadCount"`
		OfflinePush struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"friendInfoUpdated"`
	UserInfoUpdated struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"userInfoUpdated"`
	UserStatusChanged struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"userStatusChanged"`
	ConversationChanged struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"conversationChanged"`
	ConversationSetPrivate struct {
		IsSendMsg        bool `mapstructure:"isSendMsg"`
		ReliabilityLevel int  `mapstructure:"reliabilityLevel"`
		UnreadCount      bool `mapstructure:"unreadCount"`
		OfflinePush      struct {
			Enable bool   `mapstructure:"enable"`
			Title  string `mapstructure:"title"`
			Desc   string `mapstructure:"desc"`
			Ext    string `mapstructure:"ext"`
		} `mapstructure:"offlinePush"`
	} `mapstructure:"conversationSetPrivate"`
}

type MsgGateway struct {
	RPC struct {
		RegisterIP string `mapstructure:"registerIP"`
		Ports      []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus struct {
		Enable bool  `mapstructure:"enable"`
		Ports  []int `mapstructure:"ports"`
	} `mapstructure:"prometheus"`
	ListenIP    string `mapstructure:"listenIP"`
	LongConnSvr struct {
		Ports               []int `mapstructure:"ports"`
		WebsocketMaxConnNum int   `mapstructure:"websocketMaxConnNum"`
		WebsocketMaxMsgLen  int   `mapstructure:"websocketMaxMsgLen"`
		WebsocketTimeout    int   `mapstructure:"websocketTimeout"`
	} `mapstructure:"longConnSvr"`
	MultiLoginPolicy int `mapstructure:"multiLoginPolicy"`
}

type MsgTransfer struct {
	Prometheus struct {
		Enable bool  `mapstructure:"enable"`
		Ports  []int `mapstructure:"ports"`
	} `mapstructure:"prometheus"`
	MsgCacheTimeout int `mapstructure:"msgCacheTimeout"`
}

type Push struct {
	RPC struct {
		RegisterIP string `mapstructure:"registerIP"`
		ListenIP   string `mapstructure:"listenIP"`
		Ports      []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus struct {
		Enable bool  `mapstructure:"enable"`
		Ports  []int `mapstructure:"ports"`
	} `mapstructure:"prometheus"`
	Enable string `mapstructure:"enable"`
	GeTui  struct {
		PushUrl      string `mapstructure:"pushUrl"`
		MasterSecret string `mapstructure:"masterSecret"`
		AppKey       string `mapstructure:"appKey"`
		Intent       string `mapstructure:"intent"`
		ChannelID    string `mapstructure:"channelID"`
		ChannelName  string `mapstructure:"channelName"`
	} `mapstructure:"geTui"`
	FCM struct {
		ServiceAccount string `mapstructure:"serviceAccount"`
	} `mapstructure:"fcm"`
	JPNS struct {
		AppKey       string `mapstructure:"appKey"`
		MasterSecret string `mapstructure:"masterSecret"`
		PushURL      string `mapstructure:"pushURL"`
		PushIntent   string `mapstructure:"pushIntent"`
	} `mapstructure:"jpns"`
	IOSPush struct {
		PushSound  string `mapstructure:"pushSound"`
		BadgeCount bool   `mapstructure:"badgeCount"`
		Production bool   `mapstructure:"production"`
	} `mapstructure:"iosPush"`
}

type Auth struct {
	RPC struct {
		RegisterIP string `mapstructure:"registerIP"`
		ListenIP   string `mapstructure:"listenIP"`
		Ports      []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus struct {
		Enable bool  `mapstructure:"enable"`
		Ports  []int `mapstructure:"ports"`
	} `mapstructure:"prometheus"`
	TokenPolicy struct {
		Expire int `mapstructure:"expire"`
	} `mapstructure:"tokenPolicy"`
	Secret string `mapstructure:"secret"`
}

type Conversation struct {
	RPC struct {
		RegisterIP string `mapstructure:"registerIP"`
		ListenIP   string `mapstructure:"listenIP"`
		Ports      []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus struct {
		Enable bool  `mapstructure:"enable"`
		Ports  []int `mapstructure:"ports"`
	} `mapstructure:"prometheus"`
}

type Friend struct {
	RPC struct {
		RegisterIP string `mapstructure:"registerIP"`
		ListenIP   string `mapstructure:"listenIP"`
		Ports      []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus struct {
		Enable bool  `mapstructure:"enable"`
		Ports  []int `mapstructure:"ports"`
	} `mapstructure:"prometheus"`
}

type Group struct {
	RPC struct {
		RegisterIP string `mapstructure:"registerIP"`
		ListenIP   string `mapstructure:"listenIP"`
		Ports      []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus struct {
		Enable bool  `mapstructure:"enable"`
		Ports  []int `mapstructure:"ports"`
	} `mapstructure:"prometheus"`
}

type Msg struct {
	RPC struct {
		RegisterIP string `mapstructure:"registerIP"`
		ListenIP   string `mapstructure:"listenIP"`
		Ports      []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus struct {
		Enable bool  `mapstructure:"enable"`
		Ports  []int `mapstructure:"ports"`
	} `mapstructure:"prometheus"`
	FriendVerify                      bool `mapstructure:"friendVerify"`
	GroupMessageHasReadReceiptEnable  bool `mapstructure:"groupMessageHasReadReceiptEnable"`
	SingleMessageHasReadReceiptEnable bool `mapstructure:"singleMessageHasReadReceiptEnable"`
}

type Third struct {
	RPC struct {
		RegisterIP string `mapstructure:"registerIP"`
		ListenIP   string `mapstructure:"listenIP"`
		Ports      []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus struct {
		Enable bool  `mapstructure:"enable"`
		Ports  []int `mapstructure:"ports"`
	} `mapstructure:"prometheus"`
	Object struct {
		Enable string `mapstructure:"enable"`
		Cos    struct {
			BucketURL    string `mapstructure:"bucketURL"`
			SecretID     string `mapstructure:"secretID"`
			SecretKey    string `mapstructure:"secretKey"`
			SessionToken string `mapstructure:"sessionToken"`
			PublicRead   bool   `mapstructure:"publicRead"`
		} `mapstructure:"cos"`
		Oss struct {
			Endpoint        string `mapstructure:"endpoint"`
			Bucket          string `mapstructure:"bucket"`
			BucketURL       string `mapstructure:"bucketURL"`
			AccessKeyID     string `mapstructure:"accessKeyID"`
			AccessKeySecret string `mapstructure:"accessKeySecret"`
			SessionToken    string `mapstructure:"sessionToken"`
			PublicRead      bool   `mapstructure:"publicRead"`
		} `mapstructure:"oss"`
		Kodo struct {
			Endpoint        string `mapstructure:"endpoint"`
			Bucket          string `mapstructure:"bucket"`
			BucketURL       string `mapstructure:"bucketURL"`
			AccessKeyID     string `mapstructure:"accessKeyID"`
			AccessKeySecret string `mapstructure:"accessKeySecret"`
			SessionToken    string `mapstructure:"sessionToken"`
			PublicRead      bool   `mapstructure:"publicRead"`
		} `mapstructure:"kodo"`
		Aws struct {
			Endpoint        string `mapstructure:"endpoint"`
			Region          string `mapstructure:"region"`
			Bucket          string `mapstructure:"bucket"`
			AccessKeyID     string `mapstructure:"accessKeyID"`
			AccessKeySecret string `mapstructure:"accessKeySecret"`
			PublicRead      bool   `mapstructure:"publicRead"`
		} `mapstructure:"aws"`
	} `mapstructure:"object"`
}

type User struct {
	RPC struct {
		RegisterIP string `mapstructure:"registerIP"`
		ListenIP   string `mapstructure:"listenIP"`
		Ports      []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus struct {
		Enable bool  `mapstructure:"enable"`
		Ports  []int `mapstructure:"ports"`
	} `mapstructure:"prometheus"`
}

type Redis struct {
	Address  []string `mapstructure:"address"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
}

type Webhooks struct {
	URL                 string `mapstructure:"url"`
	BeforeSendSingleMsg struct {
		Enable         bool `mapstructure:"enable"`
		Timeout        int  `mapstructure:"timeout"`
		FailedContinue bool `mapstructure:"failedContinue"`
	} `mapstructure:"beforeSendSingleMsg"`
	BeforeUpdateUserInfoEx struct {
		Enable         bool `mapstructure:"enable"`
		Timeout        int  `mapstructure:"timeout"`
		FailedContinue bool `mapstructure:"failedContinue"`
	} `mapstructure:"beforeUpdateUserInfoEx"`
	AfterUpdateUserInfoEx struct {
		Enable         bool `mapstructure:"enable"`
		Timeout        int  `mapstructure:"timeout"`
		FailedContinue bool `mapstructure:"failedContinue"`
	} `mapstructure:"afterUpdateUserInfoEx"`
	AfterSendSingleMsg struct {
		Enable         bool `mapstructure:"enable"`
		Timeout        int  `mapstructure:"timeout"`
		FailedContinue bool `mapstructure:"failedContinue"`
	} `mapstructure:"afterSendSingleMsg"`
	BeforeSendGroupMsg struct {
		Enable         bool `mapstructure:"enable"`
		Timeout        int  `mapstructure:"timeout"`
		FailedContinue bool `mapstructure:"failedContinue"`
	} `mapstructure:"beforeSendGroupMsg"`
	AfterSendGroupMsg struct {
		Enable         bool `mapstructure:"enable"`
		Timeout        int  `mapstructure:"timeout"`
		FailedContinue bool `mapstructure:"failedContinue"`
	} `mapstructure:"afterSendGroupMsg"`
	MsgModify struct {
		Enable         bool `mapstructure:"enable"`
		Timeout        int  `mapstructure:"timeout"`
		FailedContinue bool `mapstructure:"failedContinue"`
	} `mapstructure:"msgModify"`
	UserOnline struct {
		Enable         bool `mapstructure:"enable"`
		Timeout        int  `mapstructure:"timeout"`
		FailedContinue bool `mapstructure:"failedContinue"`
	} `mapstructure:"userOnline"`
	UserOffline struct {
		Enable         bool `mapstructure:"enable"`
		Timeout        int  `mapstructure:"timeout"`
		FailedContinue bool `mapstructure:"failedContinue"`
	} `mapstructure:"userOffline"`
	UserKickOff struct {
		Enable         bool `mapstructure:"enable"`
		Timeout        int  `mapstructure:"timeout"`
		FailedContinue bool `mapstructure:"failedContinue"`
	} `mapstructure:"userKickOff"`
	OfflinePush struct {
		Enable         bool `mapstructure:"enable"`
		Timeout        int  `mapstructure:"timeout"`
		FailedContinue bool `mapstructure:"failedContinue"`
	} `mapstructure:"offlinePush"`
	OnlinePush struct {
		Enable         bool `mapstructure:"enable"`
		Timeout        int  `mapstructure:"timeout"`
		FailedContinue bool `mapstructure:"failedContinue"`
	} `mapstructure:"onlinePush"`
	SuperGroupOnlinePush struct {
		Enable         bool `mapstructure:"enable"`
		Timeout        int  `mapstructure:"timeout"`
		FailedContinue bool `mapstructure:"failedContinue"`
	} `mapstructure:"superGroupOnlinePush"`
	// Add additional fields here following the same pattern for other hooks
}

type ZooKeeper struct {
	Schema          string   `mapstructure:"schema"`
	Address         []string `mapstructure:"address"`
	Username        string   `mapstructure:"username"`
	Password        string   `mapstructure:"password"`
	RpcRegisterName struct {
		User           string `mapstructure:"User"`
		Friend         string `mapstructure:"Friend"`
		Msg            string `mapstructure:"Msg"`
		Push           string `mapstructure:"Push"`
		MessageGateway string `mapstructure:"MessageGateway"`
		Group          string `mapstructure:"Group"`
		Auth           string `mapstructure:"Auth"`
		Conversation   string `mapstructure:"Conversation"`
		Third          string `mapstructure:"Third"`
	} `mapstructure:"rpcRegisterName"`
}

func (m *Mongo) Build() *mongoutil.Config {
	return &mongoutil.Config{
		Uri:         m.URI,
		Address:     m.Address,
		Database:    m.Database,
		Username:    m.Username,
		Password:    m.Password,
		MaxPoolSize: m.MaxPoolSize,
		MaxRetry:    m.MaxRetry,
	}
}

func (r *Redis) Build() *redisutil.Config {
	return &redisutil.Config{
		ClusterMode: r.ClusterMode,
		Address:     r.Address,
		Username:    r.Username,
		Password:    r.Password,
		DB:          r.DB,
		MaxRetry:    r.MaxRetry,
	}
}

type LocalCache struct {
	Topic         string `yaml:"topic"`
	SlotNum       int    `yaml:"slotNum"`
	SlotSize      int    `yaml:"slotSize"`
	SuccessExpire int    `yaml:"successExpire"` // second
	FailedExpire  int    `yaml:"failedExpire"`  // second
}

func (l LocalCache) Failed() time.Duration {
	return time.Second * time.Duration(l.FailedExpire)
}

func (l LocalCache) Success() time.Duration {
	return time.Second * time.Duration(l.SuccessExpire)
}

func (l LocalCache) Enable() bool {
	return l.Topic != "" && l.SlotNum > 0 && l.SlotSize > 0
}

type localCache struct {
	User         LocalCache `yaml:"user"`
	Group        LocalCache `yaml:"group"`
	Friend       LocalCache `yaml:"friend"`
	Conversation LocalCache `yaml:"conversation"`
}

func (c *GlobalConfig) GetServiceNames() []string {
	return []string{
		c.RpcRegisterName.OpenImUserName,
		c.RpcRegisterName.OpenImFriendName,
		c.RpcRegisterName.OpenImMsgName,
		c.RpcRegisterName.OpenImPushName,
		c.RpcRegisterName.OpenImMessageGatewayName,
		c.RpcRegisterName.OpenImGroupName,
		c.RpcRegisterName.OpenImAuthName,
		c.RpcRegisterName.OpenImConversationName,
		c.RpcRegisterName.OpenImThirdName,
	}
}

func (c *GlobalConfig) RegisterConf2Registry(registry discovery.SvcDiscoveryRegistry) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return registry.RegisterConf2Registry(ConfKey, data)
}

func (c *GlobalConfig) GetConfFromRegistry(registry discovery.SvcDiscoveryRegistry) ([]byte, error) {
	return registry.GetConfFromRegistry(ConfKey)
}

func (c *GlobalConfig) EncodeConfig() []byte {
	buf := bytes.NewBuffer(nil)
	if err := yaml.NewEncoder(buf).Encode(c); err != nil {
		program.ExitWithError(err)
	}
	return buf.Bytes()
}
