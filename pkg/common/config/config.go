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

	"github.com/OpenIMSDK/tools/discoveryregistry"
	"gopkg.in/yaml.v3"
)

var Config configStruct

const ConfKey = "conf"

type CallBackConfig struct {
	Enable                 bool  `yaml:"enable"`
	CallbackTimeOut        int   `yaml:"timeout"`
	CallbackFailedContinue *bool `yaml:"failedContinue"`
}

type NotificationConf struct {
	IsSendMsg        bool         `yaml:"isSendMsg"`
	ReliabilityLevel int          `yaml:"reliabilityLevel"` // 1 online 2 persistent
	UnreadCount      bool         `yaml:"unreadCount"`
	OfflinePush      POfflinePush `yaml:"offlinePush"`
}

type POfflinePush struct {
	Enable bool   `yaml:"enable"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Ext    string `yaml:"ext"`
}

type configStruct struct {
	Envs struct {
		Discovery string `yaml:"discovery"`
	}
	Zookeeper struct {
		Schema   string   `yaml:"schema"`
		ZkAddr   []string `yaml:"address"`
		Username string   `yaml:"username"`
		Password string   `yaml:"password"`
	} `yaml:"zookeeper"`

	Mysql struct {
		Address       []string `yaml:"address"`
		Username      string   `yaml:"username"`
		Password      string   `yaml:"password"`
		Database      string   `yaml:"database"`
		MaxOpenConn   int      `yaml:"maxOpenConn"`
		MaxIdleConn   int      `yaml:"maxIdleConn"`
		MaxLifeTime   int      `yaml:"maxLifeTime"`
		LogLevel      int      `yaml:"logLevel"`
		SlowThreshold int      `yaml:"slowThreshold"`
	} `yaml:"mysql"`

	Mongo struct {
		Uri         string   `yaml:"uri"`
		Address     []string `yaml:"address"`
		Database    string   `yaml:"database"`
		Username    string   `yaml:"username"`
		Password    string   `yaml:"password"`
		MaxPoolSize int      `yaml:"maxPoolSize"`
	} `yaml:"mongo"`

	Redis struct {
		ClusterMode bool     `yaml:"clusterMode"`
		Address     []string `yaml:"address"`
		Username    string   `yaml:"username"`
		Password    string   `yaml:"password"`
	} `yaml:"redis"`

	Kafka struct {
		Username string   `yaml:"username"`
		Password string   `yaml:"password"`
		Addr     []string `yaml:"addr"`
		TLS      *struct {
			CACrt              string `yaml:"caCrt"`
			ClientCrt          string `yaml:"clientCrt"`
			ClientKey          string `yaml:"clientKey"`
			ClientKeyPwd       string `yaml:"clientKeyPwd"`
			InsecureSkipVerify bool   `yaml:"insecureSkipVerify"`
		} `yaml:"tls"`
		LatestMsgToRedis struct {
			Topic string `yaml:"topic"`
		} `yaml:"latestMsgToRedis"`
		MsgToMongo struct {
			Topic string `yaml:"topic"`
		} `yaml:"offlineMsgToMongo"`
		MsgToPush struct {
			Topic string `yaml:"topic"`
		} `yaml:"msgToPush"`
		ConsumerGroupID struct {
			MsgToRedis string `yaml:"msgToRedis"`
			MsgToMongo string `yaml:"msgToMongo"`
			MsgToMySql string `yaml:"msgToMySql"`
			MsgToPush  string `yaml:"msgToPush"`
		} `yaml:"consumerGroupID"`
	} `yaml:"kafka"`

	Rpc struct {
		RegisterIP string `yaml:"registerIP"`
		ListenIP   string `yaml:"listenIP"`
	} `yaml:"rpc"`

	Api struct {
		OpenImApiPort []int  `yaml:"openImApiPort"`
		ListenIP      string `yaml:"listenIP"`
	} `yaml:"api"`

	Object struct {
		Enable string `yaml:"enable"`
		ApiURL string `yaml:"apiURL"`
		Minio  struct {
			Bucket          string `yaml:"bucket"`
			Endpoint        string `yaml:"endpoint"`
			AccessKeyID     string `yaml:"accessKeyID"`
			SecretAccessKey string `yaml:"secretAccessKey"`
			SessionToken    string `yaml:"sessionToken"`
			SignEndpoint    string `yaml:"signEndpoint"`
			PublicRead      bool   `yaml:"publicRead"`
		} `yaml:"minio"`
		Cos struct {
			BucketURL    string `yaml:"bucketURL"`
			SecretID     string `yaml:"secretID"`
			SecretKey    string `yaml:"secretKey"`
			SessionToken string `yaml:"sessionToken"`
			PublicRead   bool   `yaml:"publicRead"`
		} `yaml:"cos"`
		Oss struct {
			Endpoint        string `yaml:"endpoint"`
			Bucket          string `yaml:"bucket"`
			BucketURL       string `yaml:"bucketURL"`
			AccessKeyID     string `yaml:"accessKeyID"`
			AccessKeySecret string `yaml:"accessKeySecret"`
			SessionToken    string `yaml:"sessionToken"`
			PublicRead      bool   `yaml:"publicRead"`
		} `yaml:"oss"`
	} `yaml:"object"`

	RpcPort struct {
		OpenImUserPort           []int `yaml:"openImUserPort"`
		OpenImFriendPort         []int `yaml:"openImFriendPort"`
		OpenImMessagePort        []int `yaml:"openImMessagePort"`
		OpenImMessageGatewayPort []int `yaml:"openImMessageGatewayPort"`
		OpenImGroupPort          []int `yaml:"openImGroupPort"`
		OpenImAuthPort           []int `yaml:"openImAuthPort"`
		OpenImPushPort           []int `yaml:"openImPushPort"`
		OpenImConversationPort   []int `yaml:"openImConversationPort"`
		OpenImRtcPort            []int `yaml:"openImRtcPort"`
		OpenImThirdPort          []int `yaml:"openImThirdPort"`
	} `yaml:"rpcPort"`

	RpcRegisterName struct {
		OpenImUserName           string `yaml:"openImUserName"`
		OpenImFriendName         string `yaml:"openImFriendName"`
		OpenImMsgName            string `yaml:"openImMsgName"`
		OpenImPushName           string `yaml:"openImPushName"`
		OpenImMessageGatewayName string `yaml:"openImMessageGatewayName"`
		OpenImGroupName          string `yaml:"openImGroupName"`
		OpenImAuthName           string `yaml:"openImAuthName"`
		OpenImConversationName   string `yaml:"openImConversationName"`
		OpenImThirdName          string `yaml:"openImThirdName"`
	} `yaml:"rpcRegisterName"`

	Log struct {
		StorageLocation     string `yaml:"storageLocation"`
		RotationTime        uint   `yaml:"rotationTime"`
		RemainRotationCount uint   `yaml:"remainRotationCount"`
		RemainLogLevel      int    `yaml:"remainLogLevel"`
		IsStdout            bool   `yaml:"isStdout"`
		IsJson              bool   `yaml:"isJson"`
		WithStack           bool   `yaml:"withStack"`
	} `yaml:"log"`

	LongConnSvr struct {
		OpenImMessageGatewayPort []int `yaml:"openImMessageGatewayPort"`
		OpenImWsPort             []int `yaml:"openImWsPort"`
		WebsocketMaxConnNum      int   `yaml:"websocketMaxConnNum"`
		WebsocketMaxMsgLen       int   `yaml:"websocketMaxMsgLen"`
		WebsocketTimeout         int   `yaml:"websocketTimeout"`
	} `yaml:"longConnSvr"`

	Push struct {
		Enable string `yaml:"enable"`
		GeTui  struct {
			PushUrl      string `yaml:"pushUrl"`
			AppKey       string `yaml:"appKey"`
			Intent       string `yaml:"intent"`
			MasterSecret string `yaml:"masterSecret"`
			ChannelID    string `yaml:"channelID"`
			ChannelName  string `yaml:"channelName"`
		} `yaml:"geTui"`
		Fcm struct {
			ServiceAccount string `yaml:"serviceAccount"`
		} `yaml:"fcm"`
		Jpns struct {
			AppKey       string `yaml:"appKey"`
			MasterSecret string `yaml:"masterSecret"`
			PushUrl      string `yaml:"pushUrl"`
			PushIntent   string `yaml:"pushIntent"`
		} `yaml:"jpns"`
	}
	Manager struct {
		UserID   []string `yaml:"userID"`
		Nickname []string `yaml:"nickname"`
	} `yaml:"manager"`

	MultiLoginPolicy                  int    `yaml:"multiLoginPolicy"`
	ChatPersistenceMysql              bool   `yaml:"chatPersistenceMysql"`
	MsgCacheTimeout                   int    `yaml:"msgCacheTimeout"`
	GroupMessageHasReadReceiptEnable  bool   `yaml:"groupMessageHasReadReceiptEnable"`
	SingleMessageHasReadReceiptEnable bool   `yaml:"singleMessageHasReadReceiptEnable"`
	RetainChatRecords                 int    `yaml:"retainChatRecords"`
	ChatRecordsClearTime              string `yaml:"chatRecordsClearTime"`
	MsgDestructTime                   string `yaml:"msgDestructTime"`
	Secret                            string `yaml:"secret"`
	TokenPolicy                       struct {
		Expire int64 `yaml:"expire"`
	} `yaml:"tokenPolicy"`
	MessageVerify struct {
		FriendVerify *bool `yaml:"friendVerify"`
	} `yaml:"messageVerify"`

	IOSPush struct {
		PushSound  string `yaml:"pushSound"`
		BadgeCount bool   `yaml:"badgeCount"`
		Production bool   `yaml:"production"`
	} `yaml:"iosPush"`
	Callback struct {
		CallbackUrl                        string         `yaml:"url"`
		CallbackBeforeSendSingleMsg        CallBackConfig `yaml:"beforeSendSingleMsg"`
		CallbackAfterSendSingleMsg         CallBackConfig `yaml:"afterSendSingleMsg"`
		CallbackBeforeSendGroupMsg         CallBackConfig `yaml:"beforeSendGroupMsg"`
		CallbackAfterSendGroupMsg          CallBackConfig `yaml:"afterSendGroupMsg"`
		CallbackMsgModify                  CallBackConfig `yaml:"msgModify"`
		CallbackUserOnline                 CallBackConfig `yaml:"userOnline"`
		CallbackUserOffline                CallBackConfig `yaml:"userOffline"`
		CallbackUserKickOff                CallBackConfig `yaml:"userKickOff"`
		CallbackOfflinePush                CallBackConfig `yaml:"offlinePush"`
		CallbackOnlinePush                 CallBackConfig `yaml:"onlinePush"`
		CallbackBeforeSuperGroupOnlinePush CallBackConfig `yaml:"superGroupOnlinePush"`
		CallbackBeforeAddFriend            CallBackConfig `yaml:"beforeAddFriend"`
		CallbackBeforeUpdateUserInfo       CallBackConfig `yaml:"beforeUpdateUserInfo"`
		CallbackBeforeCreateGroup          CallBackConfig `yaml:"beforeCreateGroup"`
		CallbackBeforeMemberJoinGroup      CallBackConfig `yaml:"beforeMemberJoinGroup"`
		CallbackBeforeSetGroupMemberInfo   CallBackConfig `yaml:"beforeSetGroupMemberInfo"`
	} `yaml:"callback"`

	Prometheus struct {
		Enable                        bool   `yaml:"enable"`
		PrometheusUrl                 string `yaml:"prometheusUrl"`
		ApiPrometheusPort             []int  `yaml:"apiPrometheusPort"`
		UserPrometheusPort            []int  `yaml:"userPrometheusPort"`
		FriendPrometheusPort          []int  `yaml:"friendPrometheusPort"`
		MessagePrometheusPort         []int  `yaml:"messagePrometheusPort"`
		MessageGatewayPrometheusPort  []int  `yaml:"messageGatewayPrometheusPort"`
		GroupPrometheusPort           []int  `yaml:"groupPrometheusPort"`
		AuthPrometheusPort            []int  `yaml:"authPrometheusPort"`
		PushPrometheusPort            []int  `yaml:"pushPrometheusPort"`
		ConversationPrometheusPort    []int  `yaml:"conversationPrometheusPort"`
		RtcPrometheusPort             []int  `yaml:"rtcPrometheusPort"`
		MessageTransferPrometheusPort []int  `yaml:"messageTransferPrometheusPort"`
		ThirdPrometheusPort           []int  `yaml:"thirdPrometheusPort"`
	} `yaml:"prometheus"`
	Notification notification `yaml:"notification"`
}

type notification struct {
	GroupCreated             NotificationConf `yaml:"groupCreated"`
	GroupInfoSet             NotificationConf `yaml:"groupInfoSet"`
	JoinGroupApplication     NotificationConf `yaml:"joinGroupApplication"`
	MemberQuit               NotificationConf `yaml:"memberQuit"`
	GroupApplicationAccepted NotificationConf `yaml:"groupApplicationAccepted"`
	GroupApplicationRejected NotificationConf `yaml:"groupApplicationRejected"`
	GroupOwnerTransferred    NotificationConf `yaml:"groupOwnerTransferred"`
	MemberKicked             NotificationConf `yaml:"memberKicked"`
	MemberInvited            NotificationConf `yaml:"memberInvited"`
	MemberEnter              NotificationConf `yaml:"memberEnter"`
	GroupDismissed           NotificationConf `yaml:"groupDismissed"`
	GroupMuted               NotificationConf `yaml:"groupMuted"`
	GroupCancelMuted         NotificationConf `yaml:"groupCancelMuted"`
	GroupMemberMuted         NotificationConf `yaml:"groupMemberMuted"`
	GroupMemberCancelMuted   NotificationConf `yaml:"groupMemberCancelMuted"`
	GroupMemberInfoSet       NotificationConf `yaml:"groupMemberInfoSet"`
	GroupMemberSetToAdmin    NotificationConf `yaml:"groupMemberSetToAdmin"`
	GroupMemberSetToOrdinary NotificationConf `yaml:"groupMemberSetToOrdinaryUser"`
	GroupInfoSetAnnouncement NotificationConf `yaml:"groupInfoSetAnnouncement"`
	GroupInfoSetName         NotificationConf `yaml:"groupInfoSetName"`
	////////////////////////user///////////////////////
	UserInfoUpdated   NotificationConf `yaml:"userInfoUpdated"`
	UserStatusChanged NotificationConf `yaml:"userStatusChanged"`
	//////////////////////friend///////////////////////
	FriendApplicationAdded    NotificationConf `yaml:"friendApplicationAdded"`
	FriendApplicationApproved NotificationConf `yaml:"friendApplicationApproved"`
	FriendApplicationRejected NotificationConf `yaml:"friendApplicationRejected"`
	FriendAdded               NotificationConf `yaml:"friendAdded"`
	FriendDeleted             NotificationConf `yaml:"friendDeleted"`
	FriendRemarkSet           NotificationConf `yaml:"friendRemarkSet"`
	BlackAdded                NotificationConf `yaml:"blackAdded"`
	BlackDeleted              NotificationConf `yaml:"blackDeleted"`
	FriendInfoUpdated         NotificationConf `yaml:"friendInfoUpdated"`
	//////////////////////conversation///////////////////////
	ConversationChanged    NotificationConf `yaml:"conversationChanged"`
	ConversationSetPrivate NotificationConf `yaml:"conversationSetPrivate"`
}

func (c *configStruct) GetServiceNames() []string {
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

func (c *configStruct) RegisterConf2Registry(registry discoveryregistry.SvcDiscoveryRegistry) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return registry.RegisterConf2Registry(ConfKey, data)
}

func (c *configStruct) GetConfFromRegistry(registry discoveryregistry.SvcDiscoveryRegistry) ([]byte, error) {
	return registry.GetConfFromRegistry(ConfKey)
}

func (c *configStruct) EncodeConfig() []byte {
	buf := bytes.NewBuffer(nil)
	if err := yaml.NewEncoder(buf).Encode(c); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
