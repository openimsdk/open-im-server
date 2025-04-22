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
	"strings"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/kafka"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/openimsdk/tools/s3/aws"
	"github.com/openimsdk/tools/s3/cos"
	"github.com/openimsdk/tools/s3/kodo"
	"github.com/openimsdk/tools/s3/minio"
	"github.com/openimsdk/tools/s3/oss"
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
	IsSimplify          bool   `mapstructure:"isSimplify"`
	WithStack           bool   `mapstructure:"withStack"`
}

type Minio struct {
	Bucket          string `mapstructure:"bucket"`
	AccessKeyID     string `mapstructure:"accessKeyID"`
	SecretAccessKey string `mapstructure:"secretAccessKey"`
	SessionToken    string `mapstructure:"sessionToken"`
	InternalAddress string `mapstructure:"internalAddress"`
	ExternalAddress string `mapstructure:"externalAddress"`
	PublicRead      bool   `mapstructure:"publicRead"`
}

type Mongo struct {
	URI         string   `mapstructure:"uri"`
	Address     []string `mapstructure:"address"`
	Database    string   `mapstructure:"database"`
	Username    string   `mapstructure:"username"`
	Password    string   `mapstructure:"password"`
	AuthSource  string   `mapstructure:"authSource"`
	MaxPoolSize int      `mapstructure:"maxPoolSize"`
	MaxRetry    int      `mapstructure:"maxRetry"`
}
type Kafka struct {
	Username           string   `mapstructure:"username"`
	Password           string   `mapstructure:"password"`
	ProducerAck        string   `mapstructure:"producerAck"`
	CompressType       string   `mapstructure:"compressType"`
	Address            []string `mapstructure:"address"`
	ToRedisTopic       string   `mapstructure:"toRedisTopic"`
	ToMongoTopic       string   `mapstructure:"toMongoTopic"`
	ToPushTopic        string   `mapstructure:"toPushTopic"`
	ToOfflinePushTopic string   `mapstructure:"toOfflinePushTopic"`
	ToRedisGroupID     string   `mapstructure:"toRedisGroupID"`
	ToMongoGroupID     string   `mapstructure:"toMongoGroupID"`
	ToPushGroupID      string   `mapstructure:"toPushGroupID"`
	ToOfflineGroupID   string   `mapstructure:"toOfflinePushGroupID"`

	Tls TLSConfig `mapstructure:"tls"`
}
type TLSConfig struct {
	EnableTLS          bool   `mapstructure:"enableTLS"`
	CACrt              string `mapstructure:"caCrt"`
	ClientCrt          string `mapstructure:"clientCrt"`
	ClientKey          string `mapstructure:"clientKey"`
	ClientKeyPwd       string `mapstructure:"clientKeyPwd"`
	InsecureSkipVerify bool   `mapstructure:"insecureSkipVerify"`
}

type API struct {
	Api struct {
		ListenIP         string `mapstructure:"listenIP"`
		Ports            []int  `mapstructure:"ports"`
		CompressionLevel int    `mapstructure:"compressionLevel"`
	} `mapstructure:"api"`
	Prometheus struct {
		Enable       bool   `mapstructure:"enable"`
		AutoSetPorts bool   `mapstructure:"autoSetPorts"`
		Ports        []int  `mapstructure:"ports"`
		GrafanaURL   string `mapstructure:"grafanaURL"`
	} `mapstructure:"prometheus"`
}

type CronTask struct {
	CronExecuteTime   string   `mapstructure:"cronExecuteTime"`
	RetainChatRecords int      `mapstructure:"retainChatRecords"`
	FileExpireTime    int      `mapstructure:"fileExpireTime"`
	DeleteObjectType  []string `mapstructure:"deleteObjectType"`
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
	GroupMemberSetToAdmin     NotificationConfig `yaml:"groupMemberSetToAdmin"`
	GroupMemberSetToOrdinary  NotificationConfig `yaml:"groupMemberSetToOrdinaryUser"`
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
		RegisterIP   string `mapstructure:"registerIP"`
		AutoSetPorts bool   `mapstructure:"autoSetPorts"`
		Ports        []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus  Prometheus `mapstructure:"prometheus"`
	ListenIP    string     `mapstructure:"listenIP"`
	LongConnSvr struct {
		Ports               []int `mapstructure:"ports"`
		WebsocketMaxConnNum int   `mapstructure:"websocketMaxConnNum"`
		WebsocketMaxMsgLen  int   `mapstructure:"websocketMaxMsgLen"`
		WebsocketTimeout    int   `mapstructure:"websocketTimeout"`
	} `mapstructure:"longConnSvr"`
}

type MsgTransfer struct {
	Prometheus struct {
		Enable       bool  `mapstructure:"enable"`
		AutoSetPorts bool  `mapstructure:"autoSetPorts"`
		Ports        []int `mapstructure:"ports"`
	} `mapstructure:"prometheus"`
}

type Push struct {
	RPC struct {
		RegisterIP   string `mapstructure:"registerIP"`
		ListenIP     string `mapstructure:"listenIP"`
		AutoSetPorts bool   `mapstructure:"autoSetPorts"`
		Ports        []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus           Prometheus `mapstructure:"prometheus"`
	MaxConcurrentWorkers int        `mapstructure:"maxConcurrentWorkers"`
	Enable               string     `mapstructure:"enable"`
	GeTui                struct {
		PushUrl      string `mapstructure:"pushUrl"`
		MasterSecret string `mapstructure:"masterSecret"`
		AppKey       string `mapstructure:"appKey"`
		Intent       string `mapstructure:"intent"`
		ChannelID    string `mapstructure:"channelID"`
		ChannelName  string `mapstructure:"channelName"`
	} `mapstructure:"geTui"`
	FCM struct {
		FilePath string `mapstructure:"filePath"`
		AuthURL  string `mapstructure:"authURL"`
	} `mapstructure:"fcm"`
	JPush struct {
		AppKey       string `mapstructure:"appKey"`
		MasterSecret string `mapstructure:"masterSecret"`
		PushURL      string `mapstructure:"pushURL"`
		PushIntent   string `mapstructure:"pushIntent"`
	} `mapstructure:"jpush"`
	IOSPush struct {
		PushSound  string `mapstructure:"pushSound"`
		BadgeCount bool   `mapstructure:"badgeCount"`
		Production bool   `mapstructure:"production"`
	} `mapstructure:"iosPush"`
	FullUserCache bool `mapstructure:"fullUserCache"`
}

type Auth struct {
	RPC struct {
		RegisterIP   string `mapstructure:"registerIP"`
		ListenIP     string `mapstructure:"listenIP"`
		AutoSetPorts bool   `mapstructure:"autoSetPorts"`
		Ports        []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus  Prometheus `mapstructure:"prometheus"`
	TokenPolicy struct {
		Expire int64 `mapstructure:"expire"`
	} `mapstructure:"tokenPolicy"`
}

type Conversation struct {
	RPC struct {
		RegisterIP   string `mapstructure:"registerIP"`
		ListenIP     string `mapstructure:"listenIP"`
		AutoSetPorts bool   `mapstructure:"autoSetPorts"`
		Ports        []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus Prometheus `mapstructure:"prometheus"`
}

type Friend struct {
	RPC struct {
		RegisterIP   string `mapstructure:"registerIP"`
		ListenIP     string `mapstructure:"listenIP"`
		AutoSetPorts bool   `mapstructure:"autoSetPorts"`
		Ports        []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus Prometheus `mapstructure:"prometheus"`
}

type Group struct {
	RPC struct {
		RegisterIP   string `mapstructure:"registerIP"`
		ListenIP     string `mapstructure:"listenIP"`
		AutoSetPorts bool   `mapstructure:"autoSetPorts"`
		Ports        []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus                 Prometheus `mapstructure:"prometheus"`
	EnableHistoryForNewMembers bool       `mapstructure:"enableHistoryForNewMembers"`
}

type Msg struct {
	RPC struct {
		RegisterIP   string `mapstructure:"registerIP"`
		ListenIP     string `mapstructure:"listenIP"`
		AutoSetPorts bool   `mapstructure:"autoSetPorts"`
		Ports        []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus   Prometheus `mapstructure:"prometheus"`
	FriendVerify bool       `mapstructure:"friendVerify"`
}

type Third struct {
	RPC struct {
		RegisterIP   string `mapstructure:"registerIP"`
		ListenIP     string `mapstructure:"listenIP"`
		AutoSetPorts bool   `mapstructure:"autoSetPorts"`
		Ports        []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus Prometheus `mapstructure:"prometheus"`
	Object     struct {
		Enable string `mapstructure:"enable"`
		Cos    Cos    `mapstructure:"cos"`
		Oss    Oss    `mapstructure:"oss"`
		Kodo   Kodo   `mapstructure:"kodo"`
		Aws    Aws    `mapstructure:"aws"`
	} `mapstructure:"object"`
}
type Cos struct {
	BucketURL    string `mapstructure:"bucketURL"`
	SecretID     string `mapstructure:"secretID"`
	SecretKey    string `mapstructure:"secretKey"`
	SessionToken string `mapstructure:"sessionToken"`
	PublicRead   bool   `mapstructure:"publicRead"`
}
type Oss struct {
	Endpoint        string `mapstructure:"endpoint"`
	Bucket          string `mapstructure:"bucket"`
	BucketURL       string `mapstructure:"bucketURL"`
	AccessKeyID     string `mapstructure:"accessKeyID"`
	AccessKeySecret string `mapstructure:"accessKeySecret"`
	SessionToken    string `mapstructure:"sessionToken"`
	PublicRead      bool   `mapstructure:"publicRead"`
}

type Kodo struct {
	Endpoint        string `mapstructure:"endpoint"`
	Bucket          string `mapstructure:"bucket"`
	BucketURL       string `mapstructure:"bucketURL"`
	AccessKeyID     string `mapstructure:"accessKeyID"`
	AccessKeySecret string `mapstructure:"accessKeySecret"`
	SessionToken    string `mapstructure:"sessionToken"`
	PublicRead      bool   `mapstructure:"publicRead"`
}

type Aws struct {
	Endpoint        string `mapstructure:"endpoint"`
	Region          string `mapstructure:"region"`
	Bucket          string `mapstructure:"bucket"`
	AccessKeyID     string `mapstructure:"accessKeyID"`
	SecretAccessKey string `mapstructure:"secretAccessKey"`
	SessionToken    string `mapstructure:"sessionToken"`
}

type User struct {
	RPC struct {
		RegisterIP   string `mapstructure:"registerIP"`
		ListenIP     string `mapstructure:"listenIP"`
		AutoSetPorts bool   `mapstructure:"autoSetPorts"`
		Ports        []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus Prometheus `mapstructure:"prometheus"`
}

type Redis struct {
	Address     []string `mapstructure:"address"`
	Username    string   `mapstructure:"username"`
	Password    string   `mapstructure:"password"`
	ClusterMode bool     `mapstructure:"clusterMode"`
	DB          int      `mapstructure:"storage"`
	MaxRetry    int      `mapstructure:"maxRetry"`
	PoolSize    int      `mapstructure:"poolSize"`
}

type BeforeConfig struct {
	Enable         bool    `yaml:"enable"`
	Timeout        int     `yaml:"timeout"`
	FailedContinue bool    `yaml:"failedContinue"`
	DeniedTypes    []int32 `yaml:"deniedTypes"`
}

type AfterConfig struct {
	Enable       bool     `yaml:"enable"`
	Timeout      int      `yaml:"timeout"`
	AttentionIds []string `yaml:"attentionIds"`
	DeniedTypes  []int32  `yaml:"deniedTypes"`
}

type Share struct {
	Secret         string         `yaml:"secret"`
	IMAdminUserID  []string       `yaml:"imAdminUserID"`
	MultiLogin     MultiLogin     `yaml:"multiLogin"`
	RPCMaxBodySize MaxRequestBody `yaml:"rpcMaxBodySize"`
}

type MaxRequestBody struct {
	RequestMaxBodySize  int `yaml:"requestMaxBodySize"`
	ResponseMaxBodySize int `yaml:"responseMaxBodySize"`
}

type MultiLogin struct {
	Policy       int `mapstructure:"policy"`
	MaxNumOneEnd int `mapstructure:"maxNumOneEnd"`
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

// FullConfig stores all configurations for before and after events
type Webhooks struct {
	URL                      string       `mapstructure:"url"`
	BeforeSendSingleMsg      BeforeConfig `mapstructure:"beforeSendSingleMsg"`
	BeforeUpdateUserInfoEx   BeforeConfig `mapstructure:"beforeUpdateUserInfoEx"`
	AfterUpdateUserInfoEx    AfterConfig  `mapstructure:"afterUpdateUserInfoEx"`
	AfterSendSingleMsg       AfterConfig  `mapstructure:"afterSendSingleMsg"`
	BeforeSendGroupMsg       BeforeConfig `mapstructure:"beforeSendGroupMsg"`
	BeforeMsgModify          BeforeConfig `mapstructure:"beforeMsgModify"`
	AfterSendGroupMsg        AfterConfig  `mapstructure:"afterSendGroupMsg"`
	AfterUserOnline          AfterConfig  `mapstructure:"afterUserOnline"`
	AfterUserOffline         AfterConfig  `mapstructure:"afterUserOffline"`
	AfterUserKickOff         AfterConfig  `mapstructure:"afterUserKickOff"`
	BeforeOfflinePush        BeforeConfig `mapstructure:"beforeOfflinePush"`
	BeforeOnlinePush         BeforeConfig `mapstructure:"beforeOnlinePush"`
	BeforeGroupOnlinePush    BeforeConfig `mapstructure:"beforeGroupOnlinePush"`
	BeforeAddFriend          BeforeConfig `mapstructure:"beforeAddFriend"`
	BeforeUpdateUserInfo     BeforeConfig `mapstructure:"beforeUpdateUserInfo"`
	AfterUpdateUserInfo      AfterConfig  `mapstructure:"afterUpdateUserInfo"`
	BeforeCreateGroup        BeforeConfig `mapstructure:"beforeCreateGroup"`
	AfterCreateGroup         AfterConfig  `mapstructure:"afterCreateGroup"`
	BeforeMemberJoinGroup    BeforeConfig `mapstructure:"beforeMemberJoinGroup"`
	BeforeSetGroupMemberInfo BeforeConfig `mapstructure:"beforeSetGroupMemberInfo"`
	AfterSetGroupMemberInfo  AfterConfig  `mapstructure:"afterSetGroupMemberInfo"`
	AfterQuitGroup           AfterConfig  `mapstructure:"afterQuitGroup"`
	AfterKickGroupMember     AfterConfig  `mapstructure:"afterKickGroupMember"`
	AfterDismissGroup        AfterConfig  `mapstructure:"afterDismissGroup"`
	BeforeApplyJoinGroup     BeforeConfig `mapstructure:"beforeApplyJoinGroup"`
	AfterGroupMsgRead        AfterConfig  `mapstructure:"afterGroupMsgRead"`
	AfterSingleMsgRead       AfterConfig  `mapstructure:"afterSingleMsgRead"`
	BeforeUserRegister       BeforeConfig `mapstructure:"beforeUserRegister"`
	AfterUserRegister        AfterConfig  `mapstructure:"afterUserRegister"`
	AfterTransferGroupOwner  AfterConfig  `mapstructure:"afterTransferGroupOwner"`
	BeforeSetFriendRemark    BeforeConfig `mapstructure:"beforeSetFriendRemark"`
	AfterSetFriendRemark     AfterConfig  `mapstructure:"afterSetFriendRemark"`
	AfterGroupMsgRevoke      AfterConfig  `mapstructure:"afterGroupMsgRevoke"`
	AfterJoinGroup           AfterConfig  `mapstructure:"afterJoinGroup"`
	BeforeInviteUserToGroup  BeforeConfig `mapstructure:"beforeInviteUserToGroup"`
	AfterSetGroupInfo        AfterConfig  `mapstructure:"afterSetGroupInfo"`
	BeforeSetGroupInfo       BeforeConfig `mapstructure:"beforeSetGroupInfo"`
	AfterSetGroupInfoEx      AfterConfig  `mapstructure:"afterSetGroupInfoEx"`
	BeforeSetGroupInfoEx     BeforeConfig `mapstructure:"beforeSetGroupInfoEx"`
	AfterRevokeMsg           AfterConfig  `mapstructure:"afterRevokeMsg"`
	BeforeAddBlack           BeforeConfig `mapstructure:"beforeAddBlack"`
	AfterAddFriend           AfterConfig  `mapstructure:"afterAddFriend"`
	BeforeAddFriendAgree     BeforeConfig `mapstructure:"beforeAddFriendAgree"`
	AfterAddFriendAgree      AfterConfig  `mapstructure:"afterAddFriendAgree"`
	AfterDeleteFriend        AfterConfig  `mapstructure:"afterDeleteFriend"`
	BeforeImportFriends      BeforeConfig `mapstructure:"beforeImportFriends"`
	AfterImportFriends       AfterConfig  `mapstructure:"afterImportFriends"`
	AfterRemoveBlack         AfterConfig  `mapstructure:"afterRemoveBlack"`
}

type ZooKeeper struct {
	Schema   string   `mapstructure:"schema"`
	Address  []string `mapstructure:"address"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
}

type Discovery struct {
	Enable     string     `mapstructure:"enable"`
	Etcd       Etcd       `mapstructure:"etcd"`
	Kubernetes Kubernetes `mapstructure:"kubernetes"`
	RpcService RpcService `mapstructure:"rpcService"`
}

type RpcService struct {
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

func (r *RpcService) GetServiceNames() []string {
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

type Kubernetes struct {
	Namespace string `yaml:"namespace"`
}
type Etcd struct {
	RootDirectory string   `mapstructure:"rootDirectory"`
	Address       []string `mapstructure:"address"`
	Username      string   `mapstructure:"username"`
	Password      string   `mapstructure:"password"`
}

func (m *Mongo) Build() *mongoutil.Config {
	return &mongoutil.Config{
		Uri:         m.URI,
		Address:     m.Address,
		Database:    m.Database,
		Username:    m.Username,
		Password:    m.Password,
		AuthSource:  m.AuthSource,
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
		PoolSize:    r.PoolSize,
	}
}

func (k *Kafka) Build() *kafka.Config {
	return &kafka.Config{
		Username:     k.Username,
		Password:     k.Password,
		ProducerAck:  k.ProducerAck,
		CompressType: k.CompressType,
		Addr:         k.Address,
		TLS: kafka.TLSConfig{
			EnableTLS:          k.Tls.EnableTLS,
			CACrt:              k.Tls.CACrt,
			ClientCrt:          k.Tls.ClientCrt,
			ClientKey:          k.Tls.ClientKey,
			ClientKeyPwd:       k.Tls.ClientKeyPwd,
			InsecureSkipVerify: k.Tls.InsecureSkipVerify,
		},
	}
}

func (m *Minio) Build() *minio.Config {
	formatEndpoint := func(address string) string {
		if strings.HasPrefix(address, "http://") || strings.HasPrefix(address, "https://") {
			return address
		}
		return "http://" + address
	}
	return &minio.Config{
		Bucket:          m.Bucket,
		AccessKeyID:     m.AccessKeyID,
		SecretAccessKey: m.SecretAccessKey,
		SessionToken:    m.SessionToken,
		PublicRead:      m.PublicRead,
		Endpoint:        formatEndpoint(m.InternalAddress),
		SignEndpoint:    formatEndpoint(m.ExternalAddress),
	}
}
func (c *Cos) Build() *cos.Config {
	return &cos.Config{
		BucketURL:    c.BucketURL,
		SecretID:     c.SecretID,
		SecretKey:    c.SecretKey,
		SessionToken: c.SessionToken,
		PublicRead:   c.PublicRead,
	}
}

func (o *Oss) Build() *oss.Config {
	return &oss.Config{
		Endpoint:        o.Endpoint,
		Bucket:          o.Bucket,
		BucketURL:       o.BucketURL,
		AccessKeyID:     o.AccessKeyID,
		AccessKeySecret: o.AccessKeySecret,
		SessionToken:    o.SessionToken,
		PublicRead:      o.PublicRead,
	}
}

func (o *Kodo) Build() *kodo.Config {
	return &kodo.Config{
		Endpoint:        o.Endpoint,
		Bucket:          o.Bucket,
		BucketURL:       o.BucketURL,
		AccessKeyID:     o.AccessKeyID,
		AccessKeySecret: o.AccessKeySecret,
		SessionToken:    o.SessionToken,
		PublicRead:      o.PublicRead,
	}
}

func (o *Aws) Build() *aws.Config {
	return &aws.Config{
		Region:          o.Region,
		Bucket:          o.Bucket,
		AccessKeyID:     o.AccessKeyID,
		SecretAccessKey: o.SecretAccessKey,
		SessionToken:    o.SessionToken,
	}
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

const (
	DiscoveryConfigFilename          = "discovery.yml"
	KafkaConfigFileName              = "kafka.yml"
	LocalCacheConfigFileName         = "local-cache.yml"
	LogConfigFileName                = "log.yml"
	MinioConfigFileName              = "minio.yml"
	MongodbConfigFileName            = "mongodb.yml"
	OpenIMAPICfgFileName             = "openim-api.yml"
	OpenIMCronTaskCfgFileName        = "openim-crontask.yml"
	OpenIMMsgGatewayCfgFileName      = "openim-msggateway.yml"
	OpenIMMsgTransferCfgFileName     = "openim-msgtransfer.yml"
	OpenIMPushCfgFileName            = "openim-push.yml"
	OpenIMRPCAuthCfgFileName         = "openim-rpc-auth.yml"
	OpenIMRPCConversationCfgFileName = "openim-rpc-conversation.yml"
	OpenIMRPCFriendCfgFileName       = "openim-rpc-friend.yml"
	OpenIMRPCGroupCfgFileName        = "openim-rpc-group.yml"
	OpenIMRPCMsgCfgFileName          = "openim-rpc-msg.yml"
	OpenIMRPCThirdCfgFileName        = "openim-rpc-third.yml"
	OpenIMRPCUserCfgFileName         = "openim-rpc-user.yml"
	RedisConfigFileName              = "redis.yml"
	ShareFileName                    = "share.yml"
	WebhooksConfigFileName           = "webhooks.yml"
	NotificationFileName             = "notification.yml"
)

func (d *Discovery) GetConfigFileName() string {
	return DiscoveryConfigFilename
}

func (k *Kafka) GetConfigFileName() string {
	return KafkaConfigFileName
}

func (lc *LocalCache) GetConfigFileName() string {
	return LocalCacheConfigFileName
}

func (l *Log) GetConfigFileName() string {
	return LogConfigFileName
}

func (m *Minio) GetConfigFileName() string {
	return MinioConfigFileName
}

func (m *Mongo) GetConfigFileName() string {
	return MongodbConfigFileName
}

func (n *Notification) GetConfigFileName() string {
	return NotificationFileName
}

func (a *API) GetConfigFileName() string {
	return OpenIMAPICfgFileName
}

func (ct *CronTask) GetConfigFileName() string {
	return OpenIMCronTaskCfgFileName
}

func (mg *MsgGateway) GetConfigFileName() string {
	return OpenIMMsgGatewayCfgFileName
}

func (mt *MsgTransfer) GetConfigFileName() string {
	return OpenIMMsgTransferCfgFileName
}

func (p *Push) GetConfigFileName() string {
	return OpenIMPushCfgFileName
}

func (a *Auth) GetConfigFileName() string {
	return OpenIMRPCAuthCfgFileName
}

func (c *Conversation) GetConfigFileName() string {
	return OpenIMRPCConversationCfgFileName
}

func (f *Friend) GetConfigFileName() string {
	return OpenIMRPCFriendCfgFileName
}

func (g *Group) GetConfigFileName() string {
	return OpenIMRPCGroupCfgFileName
}

func (m *Msg) GetConfigFileName() string {
	return OpenIMRPCMsgCfgFileName
}

func (t *Third) GetConfigFileName() string {
	return OpenIMRPCThirdCfgFileName
}

func (u *User) GetConfigFileName() string {
	return OpenIMRPCUserCfgFileName
}

func (r *Redis) GetConfigFileName() string {
	return RedisConfigFileName
}

func (s *Share) GetConfigFileName() string {
	return ShareFileName
}

func (w *Webhooks) GetConfigFileName() string {
	return WebhooksConfigFileName
}

func InitNotification(notification *Notification) {
	notification.GroupCreated.UnreadCount = false
	notification.GroupCreated.ReliabilityLevel = 1
	notification.GroupInfoSet.UnreadCount = false
	notification.GroupInfoSet.ReliabilityLevel = 1
	notification.JoinGroupApplication.UnreadCount = false
	notification.JoinGroupApplication.ReliabilityLevel = 1
	notification.MemberQuit.UnreadCount = false
	notification.MemberQuit.ReliabilityLevel = 1
	notification.GroupApplicationAccepted.UnreadCount = false
	notification.GroupApplicationAccepted.ReliabilityLevel = 1
	notification.GroupApplicationRejected.UnreadCount = false
	notification.GroupApplicationRejected.ReliabilityLevel = 1
	notification.GroupOwnerTransferred.UnreadCount = false
	notification.GroupOwnerTransferred.ReliabilityLevel = 1
	notification.MemberKicked.UnreadCount = false
	notification.MemberKicked.ReliabilityLevel = 1
	notification.MemberInvited.UnreadCount = false
	notification.MemberInvited.ReliabilityLevel = 1
	notification.MemberEnter.UnreadCount = false
	notification.MemberEnter.ReliabilityLevel = 1
	notification.GroupDismissed.UnreadCount = false
	notification.GroupDismissed.ReliabilityLevel = 1
	notification.GroupMuted.UnreadCount = false
	notification.GroupMuted.ReliabilityLevel = 1
	notification.GroupCancelMuted.UnreadCount = false
	notification.GroupCancelMuted.ReliabilityLevel = 1
	notification.GroupMemberMuted.UnreadCount = false
	notification.GroupMemberMuted.ReliabilityLevel = 1
	notification.GroupMemberCancelMuted.UnreadCount = false
	notification.GroupMemberCancelMuted.ReliabilityLevel = 1
	notification.GroupMemberInfoSet.UnreadCount = false
	notification.GroupMemberInfoSet.ReliabilityLevel = 1
	notification.GroupMemberSetToAdmin.UnreadCount = false
	notification.GroupMemberSetToAdmin.ReliabilityLevel = 1
	notification.GroupMemberSetToOrdinary.UnreadCount = false
	notification.GroupMemberSetToOrdinary.ReliabilityLevel = 1
	notification.GroupInfoSetAnnouncement.UnreadCount = false
	notification.GroupInfoSetAnnouncement.ReliabilityLevel = 1
	notification.GroupInfoSetName.UnreadCount = false
	notification.GroupInfoSetName.ReliabilityLevel = 1
	notification.FriendApplicationAdded.UnreadCount = false
	notification.FriendApplicationAdded.ReliabilityLevel = 1
	notification.FriendApplicationApproved.UnreadCount = false
	notification.FriendApplicationApproved.ReliabilityLevel = 1
	notification.FriendApplicationRejected.UnreadCount = false
	notification.FriendApplicationRejected.ReliabilityLevel = 1
	notification.FriendAdded.UnreadCount = false
	notification.FriendAdded.ReliabilityLevel = 1
	notification.FriendDeleted.UnreadCount = false
	notification.FriendDeleted.ReliabilityLevel = 1
	notification.FriendRemarkSet.UnreadCount = false
	notification.FriendRemarkSet.ReliabilityLevel = 1
	notification.BlackAdded.UnreadCount = false
	notification.BlackAdded.ReliabilityLevel = 1
	notification.BlackDeleted.UnreadCount = false
	notification.BlackDeleted.ReliabilityLevel = 1
	notification.FriendInfoUpdated.UnreadCount = false
	notification.FriendInfoUpdated.ReliabilityLevel = 1
	notification.UserInfoUpdated.UnreadCount = false
	notification.UserInfoUpdated.ReliabilityLevel = 1
	notification.UserStatusChanged.UnreadCount = false
	notification.UserStatusChanged.ReliabilityLevel = 1
	notification.ConversationChanged.UnreadCount = false
	notification.ConversationChanged.ReliabilityLevel = 1
	notification.ConversationSetPrivate.UnreadCount = false
	notification.ConversationSetPrivate.ReliabilityLevel = 1
}

type AllConfig struct {
	Discovery    Discovery
	Kafka        Kafka
	LocalCache   LocalCache
	Log          Log
	Minio        Minio
	Mongo        Mongo
	Notification Notification
	API          API
	CronTask     CronTask
	MsgGateway   MsgGateway
	MsgTransfer  MsgTransfer
	Push         Push
	Auth         Auth
	Conversation Conversation
	Friend       Friend
	Group        Group
	Msg          Msg
	Third        Third
	User         User
	Redis        Redis
	Share        Share
	Webhooks     Webhooks
}

func (a *AllConfig) Name2Config(name string) any {
	switch name {
	case a.Discovery.GetConfigFileName():
		return a.Discovery
	case a.Kafka.GetConfigFileName():
		return a.Kafka
	case a.LocalCache.GetConfigFileName():
		return a.LocalCache
	case a.Log.GetConfigFileName():
		return a.Log
	case a.Minio.GetConfigFileName():
		return a.Minio
	case a.Mongo.GetConfigFileName():
		return a.Mongo
	case a.Notification.GetConfigFileName():
		return a.Notification
	case a.API.GetConfigFileName():
		return a.API
	case a.CronTask.GetConfigFileName():
		return a.CronTask
	case a.MsgGateway.GetConfigFileName():
		return a.MsgGateway
	case a.MsgTransfer.GetConfigFileName():
		return a.MsgTransfer
	case a.Push.GetConfigFileName():
		return a.Push
	case a.Auth.GetConfigFileName():
		return a.Auth
	case a.Conversation.GetConfigFileName():
		return a.Conversation
	case a.Friend.GetConfigFileName():
		return a.Friend
	case a.Group.GetConfigFileName():
		return a.Group
	case a.Msg.GetConfigFileName():
		return a.Msg
	case a.Third.GetConfigFileName():
		return a.Third
	case a.User.GetConfigFileName():
		return a.User
	case a.Redis.GetConfigFileName():
		return a.Redis
	case a.Share.GetConfigFileName():
		return a.Share
	case a.Webhooks.GetConfigFileName():
		return a.Webhooks
	default:
		return nil
	}
}

func (a *AllConfig) GetConfigNames() []string {
	return []string{
		a.Discovery.GetConfigFileName(),
		a.Kafka.GetConfigFileName(),
		a.LocalCache.GetConfigFileName(),
		a.Log.GetConfigFileName(),
		a.Minio.GetConfigFileName(),
		a.Mongo.GetConfigFileName(),
		a.Notification.GetConfigFileName(),
		a.API.GetConfigFileName(),
		a.CronTask.GetConfigFileName(),
		a.MsgGateway.GetConfigFileName(),
		a.MsgTransfer.GetConfigFileName(),
		a.Push.GetConfigFileName(),
		a.Auth.GetConfigFileName(),
		a.Conversation.GetConfigFileName(),
		a.Friend.GetConfigFileName(),
		a.Group.GetConfigFileName(),
		a.Msg.GetConfigFileName(),
		a.Third.GetConfigFileName(),
		a.User.GetConfigFileName(),
		a.Redis.GetConfigFileName(),
		a.Share.GetConfigFileName(),
		a.Webhooks.GetConfigFileName(),
	}
}
