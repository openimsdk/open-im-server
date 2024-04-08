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
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/openimsdk/tools/mq/kafka"
	"time"
)

type CacheConfig struct {
	Topic         string `mapstructure:"topic"`
	SlotNum       int    `mapstructure:"slotNum"`
	SlotSize      int    `mapstructure:"slotSize"`
	SuccessExpire int    `mapstructure:"successExpire"`
	FailedExpire  int    `mapstructure:"failedExpire"`
}

type LocalCache struct {
	User         CacheConfig `mapstructure:"user"`
	Group        CacheConfig `mapstructure:"group"`
	Friend       CacheConfig `mapstructure:"friend"`
	Conversation CacheConfig `mapstructure:"conversation"`
}

type Log struct {
	StorageLocation     string `mapstructure:"storageLocation"`
	RotationTime        uint   `mapstructure:"rotationTime"`
	RemainRotationCount uint   `mapstructure:"remainRotationCount"`
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
type Kafka struct {
	Username       string   `mapstructure:"username"`
	Password       string   `mapstructure:"password"`
	Address        []string `mapstructure:"address"`
	ToRedisTopic   string   `mapstructure:"toRedisTopic"`
	ToMongoTopic   string   `mapstructure:"toMongoTopic"`
	ToPushTopic    string   `mapstructure:"toPushTopic"`
	ToRedisGroupID string   `mapstructure:"toRedisGroupID"`
	ToMongoGroupID string   `mapstructure:"toMongoGroupID"`
	ToPushGroupID  string   `mapstructure:"toPushGroupID"`
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
	EnableCronLocker     bool   `yaml:"enableCronLocker"`
}

type OfflinePushConfig struct {
	Enable bool   `mapstructure:"enable"`
	Title  string `mapstructure:"title"`
	Desc   string `mapstructure:"desc"`
	Ext    string `mapstructure:"ext"`
}

type NotificationConfig struct {
	IsSendMsg        bool              `mapstructure:"isSendMsg"`
	ReliabilityLevel int               `mapstructure:"reliabilityLevel"`
	UnreadCount      bool              `mapstructure:"unreadCount"`
	OfflinePush      OfflinePushConfig `mapstructure:"offlinePush"`
}

type Notification struct {
	GroupCreated              NotificationConfig `mapstructure:"groupCreated"`
	GroupInfoSet              NotificationConfig `mapstructure:"groupInfoSet"`
	JoinGroupApplication      NotificationConfig `mapstructure:"joinGroupApplication"`
	MemberQuit                NotificationConfig `mapstructure:"memberQuit"`
	GroupApplicationAccepted  NotificationConfig `mapstructure:"groupApplicationAccepted"`
	GroupApplicationRejected  NotificationConfig `mapstructure:"groupApplicationRejected"`
	GroupOwnerTransferred     NotificationConfig `mapstructure:"groupOwnerTransferred"`
	MemberKicked              NotificationConfig `mapstructure:"memberKicked"`
	MemberInvited             NotificationConfig `mapstructure:"memberInvited"`
	MemberEnter               NotificationConfig `mapstructure:"memberEnter"`
	GroupDismissed            NotificationConfig `mapstructure:"groupDismissed"`
	GroupMuted                NotificationConfig `mapstructure:"groupMuted"`
	GroupCancelMuted          NotificationConfig `mapstructure:"groupCancelMuted"`
	GroupMemberMuted          NotificationConfig `mapstructure:"groupMemberMuted"`
	GroupMemberCancelMuted    NotificationConfig `mapstructure:"groupMemberCancelMuted"`
	GroupMemberInfoSet        NotificationConfig `mapstructure:"groupMemberInfoSet"`
	GroupInfoSetAnnouncement  NotificationConfig `mapstructure:"groupInfoSetAnnouncement"`
	GroupInfoSetName          NotificationConfig `mapstructure:"groupInfoSetName"`
	FriendApplicationAdded    NotificationConfig `mapstructure:"friendApplicationAdded"`
	FriendApplicationApproved NotificationConfig `mapstructure:"friendApplicationApproved"`
	FriendApplicationRejected NotificationConfig `mapstructure:"friendApplicationRejected"`
	FriendAdded               NotificationConfig `mapstructure:"friendAdded"`
	FriendDeleted             NotificationConfig `mapstructure:"friendDeleted"`
	FriendRemarkSet           NotificationConfig `mapstructure:"friendRemarkSet"`
	BlackAdded                NotificationConfig `mapstructure:"blackAdded"`
	BlackDeleted              NotificationConfig `mapstructure:"blackDeleted"`
	FriendInfoUpdated         NotificationConfig `mapstructure:"friendInfoUpdated"`
	UserInfoUpdated           NotificationConfig `mapstructure:"userInfoUpdated"`
	UserStatusChanged         NotificationConfig `mapstructure:"userStatusChanged"`
	ConversationChanged       NotificationConfig `mapstructure:"conversationChanged"`
	ConversationSetPrivate    NotificationConfig `mapstructure:"conversationSetPrivate"`
}

type Prometheus struct {
	Enable bool  `mapstructure:"enable"`
	Ports  []int `mapstructure:"ports"`
}

type MsgGateway struct {
	RPC struct {
		RegisterIP string `mapstructure:"registerIP"`
		Ports      []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus  Prometheus `mapstructure:"prometheus"`
	ListenIP    string     `mapstructure:"listenIP"`
	LongConnSvr struct {
		Ports               []int `mapstructure:"ports"`
		WebsocketMaxConnNum int   `mapstructure:"websocketMaxConnNum"`
		WebsocketMaxMsgLen  int   `mapstructure:"websocketMaxMsgLen"`
		WebsocketTimeout    int   `mapstructure:"websocketTimeout"`
	} `mapstructure:"longConnSvr"`
	MultiLoginPolicy int `mapstructure:"multiLoginPolicy"`
}

type MsgTransfer struct {
	Prometheus      Prometheus `mapstructure:"prometheus"`
	MsgCacheTimeout int        `mapstructure:"msgCacheTimeout"`
}

type Push struct {
	RPC struct {
		RegisterIP string `mapstructure:"registerIP"`
		ListenIP   string `mapstructure:"listenIP"`
		Ports      []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus Prometheus `mapstructure:"prometheus"`
	Enable     string     `mapstructure:"enable"`
	GeTui      struct {
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
	Prometheus  Prometheus `mapstructure:"prometheus"`
	TokenPolicy struct {
		Expire int64 `mapstructure:"expire"`
	} `mapstructure:"tokenPolicy"`
	Secret string `mapstructure:"secret"`
}

type Conversation struct {
	RPC struct {
		RegisterIP string `mapstructure:"registerIP"`
		ListenIP   string `mapstructure:"listenIP"`
		Ports      []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus Prometheus `mapstructure:"prometheus"`
}

type Friend struct {
	RPC struct {
		RegisterIP string `mapstructure:"registerIP"`
		ListenIP   string `mapstructure:"listenIP"`
		Ports      []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus Prometheus `mapstructure:"prometheus"`
}

type Group struct {
	RPC struct {
		RegisterIP string `mapstructure:"registerIP"`
		ListenIP   string `mapstructure:"listenIP"`
		Ports      []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus Prometheus `mapstructure:"prometheus"`
}

type Msg struct {
	RPC struct {
		RegisterIP string `mapstructure:"registerIP"`
		ListenIP   string `mapstructure:"listenIP"`
		Ports      []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus                        Prometheus `mapstructure:"prometheus"`
	FriendVerify                      bool       `mapstructure:"friendVerify"`
	GroupMessageHasReadReceiptEnable  bool       `mapstructure:"groupMessageHasReadReceiptEnable"`
	SingleMessageHasReadReceiptEnable bool       `mapstructure:"singleMessageHasReadReceiptEnable"`
}

type Third struct {
	RPC struct {
		RegisterIP string `mapstructure:"registerIP"`
		ListenIP   string `mapstructure:"listenIP"`
		Ports      []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus Prometheus `mapstructure:"prometheus"`
	Object     struct {
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
	Prometheus Prometheus `mapstructure:"prometheus"`
}

type Redis struct {
	Address     []string `mapstructure:"address"`
	Username    string   `mapstructure:"username"`
	Password    string   `mapstructure:"password"`
	ClusterMode bool     `mapstructure:"clusterMode"`
	DB          int      `mapstructure:"db"`
	MaxRetry    int      `mapstructure:"MaxRetry"`
}

type WebhookConfig struct {
	Enable         bool `mapstructure:"enable"`
	Timeout        int  `mapstructure:"timeout"`
	FailedContinue bool `mapstructure:"failedContinue"`
}

type Share struct {
	Env             string          `mapstructure:"env"`
	RpcRegisterName RpcRegisterName `mapstructure:"rpcRegisterName"`
	IMAdmin         IMAdmin         `mapstructure:"imAdmin"`
}
type RpcRegisterName struct {
	User           string `mapstructure:"user"`
	Friend         string `mapstructure:"friend"`
	Msg            string `mapstructure:"msg"`
	Push           string `mapstructure:"push"`
	MessageGateway string `mapstructure:"messageGateway"`
	Group          string `mapstructure:"group"`
	Auth           string `mapstructure:"auth"`
	Conversation   string `mapstructure:"conversation"`
	Third          string `mapstructure:"third"`
}

func (r *RpcRegisterName) GetServiceNames() []string {
	return []string{
		r.User,
		r.Friend,
		r.Msg,
		r.Push,
		r.MessageGateway,
		r.Group,
		r.Auth,
		r.Conversation,
		r.Third,
	}
}

type IMAdmin struct {
	UserID   []string `mapstructure:"userID"`
	Nickname []string `mapstructure:"nickname"`
}

type Webhooks struct {
	URL                      string        `mapstructure:"url"`
	BeforeSendSingleMsg      WebhookConfig `mapstructure:"beforeSendSingleMsg"`
	BeforeUpdateUserInfoEx   WebhookConfig `mapstructure:"beforeUpdateUserInfoEx"`
	AfterUpdateUserInfoEx    WebhookConfig `mapstructure:"afterUpdateUserInfoEx"`
	AfterSendSingleMsg       WebhookConfig `mapstructure:"afterSendSingleMsg"`
	BeforeSendGroupMsg       WebhookConfig `mapstructure:"beforeSendGroupMsg"`
	AfterSendGroupMsg        WebhookConfig `mapstructure:"afterSendGroupMsg"`
	AfterUserOnline          WebhookConfig `mapstructure:"afterUserOnline"`
	AfterUserOffline         WebhookConfig `mapstructure:"afterUserOffline"`
	AfterUserKickOff         WebhookConfig `mapstructure:"afterUserKickOff"`
	BeforeOfflinePush        WebhookConfig `mapstructure:"beforeOfflinePush"`
	BeforeOnlinePush         WebhookConfig `mapstructure:"beforeOnlinePush"`
	BeforeGroupOnlinePush    WebhookConfig `mapstructure:"beforeGroupOnlinePush"`
	BeforeAddFriend          WebhookConfig `mapstructure:"beforeAddFriend"`
	BeforeUpdateUserInfo     WebhookConfig `mapstructure:"beforeUpdateUserInfo"`
	BeforeCreateGroup        WebhookConfig `mapstructure:"beforeCreateGroup"`
	AfterCreateGroup         WebhookConfig `mapstructure:"afterCreateGroup"`
	BeforeMemberJoinGroup    WebhookConfig `mapstructure:"beforeMemberJoinGroup"`
	BeforeSetGroupMemberInfo WebhookConfig `mapstructure:"beforeSetGroupMemberInfo"`
	AfterSetGroupMemberInfo  WebhookConfig `mapstructure:"afterSetGroupMemberInfo"`
	AfterQuitGroup           WebhookConfig `mapstructure:"afterQuitGroup"`
	AfterKickGroupMember     WebhookConfig `mapstructure:"afterKickGroupMember"`
	AfterDismissGroup        WebhookConfig `mapstructure:"afterDismissGroup"`
	BeforeApplyJoinGroup     WebhookConfig `mapstructure:"beforeApplyJoinGroup"`
	AfterGroupMsgRead        WebhookConfig `mapstructure:"afterGroupMsgRead"`
	AfterSingleMsgRead       WebhookConfig `mapstructure:"afterSingleMsgRead"`
	BeforeUserRegister       WebhookConfig `mapstructure:"beforeUserRegister"`
	AfterUserRegister        WebhookConfig `mapstructure:"afterUserRegister"`
	AfterTransferGroupOwner  WebhookConfig `mapstructure:"afterTransferGroupOwner"`
	BeforeSetFriendRemark    WebhookConfig `mapstructure:"beforeSetFriendRemark"`
	AfterSetFriendRemark     WebhookConfig `mapstructure:"afterSetFriendRemark"`
	AfterGroupMsgRevoke      WebhookConfig `mapstructure:"afterGroupMsgRevoke"`
	AfterJoinGroup           WebhookConfig `mapstructure:"afterJoinGroup"`
	BeforeInviteUserToGroup  WebhookConfig `mapstructure:"beforeInviteUserToGroup"`
	AfterSetGroupInfo        WebhookConfig `mapstructure:"afterSetGroupInfo"`
	BeforeSetGroupInfo       WebhookConfig `mapstructure:"beforeSetGroupInfo"`
	AfterRevokeMsg           WebhookConfig `mapstructure:"afterRevokeMsg"`
	BeforeAddBlack           WebhookConfig `mapstructure:"beforeAddBlack"`
	AfterAddFriend           WebhookConfig `mapstructure:"afterAddFriend"`
	BeforeAddFriendAgree     WebhookConfig `mapstructure:"beforeAddFriendAgree"`
	AfterDeleteFriend        WebhookConfig `mapstructure:"afterDeleteFriend"`
	BeforeImportFriends      WebhookConfig `mapstructure:"beforeImportFriends"`
	AfterImportFriends       WebhookConfig `mapstructure:"afterImportFriends"`
	AfterRemoveBlack         WebhookConfig `mapstructure:"afterRemoveBlack"`
}

type ZooKeeper struct {
	Schema   string   `mapstructure:"schema"`
	Address  []string `mapstructure:"address"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
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

func (k *Kafka) Build() *kafka.Config {
	return &kafka.Config{}
}

func (l *CacheConfig) Failed() time.Duration {
	return time.Second * time.Duration(l.FailedExpire)
}

func (l *CacheConfig) Success() time.Duration {
	return time.Second * time.Duration(l.SuccessExpire)
}

func (l *CacheConfig) Enable() bool {
	return l.Topic != "" && l.SlotNum > 0 && l.SlotSize > 0
}
