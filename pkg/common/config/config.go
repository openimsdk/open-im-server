package config

import (
	_ "embed"
)

//go:embed version
var Version string

var Config config

type CallBackConfig struct {
	Enable                 bool  `yaml:"enable"`
	CallbackTimeOut        int   `yaml:"timeout"`
	CallbackFailedContinue *bool `yaml:"failedContinue"`
}

type NotificationConf struct {
	IsSendMsg        bool         `yaml:"isSendMsg"`
	ReliabilityLevel int          `yaml:"reliabilityLevel"` // 1 online 2 presistent
	UnreadCount      bool         `yaml:"unreadCount"`
	OfflinePush      POfflinePush `yaml:"offlinePush"`
}

type POfflinePush struct {
	Enable bool   `yaml:"enable"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Ext    string `yaml:"ext"`
}

type config struct {
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
		Address  []string `yaml:"address"`
		Username string   `yaml:"username"`
		Password string   `yaml:"password"`
	} `yaml:"redis"`

	Kafka struct {
		Username         string   `yaml:"username"`
		Password         string   `yaml:"password"`
		Addr             []string `yaml:"addr"`
		LatestMsgToRedis struct {
			Topic string `yaml:"topic"`
		} `yaml:"latestMsgToRedis"`
		MsgToMongo struct {
			Topic string `yaml:"topic"`
		} `yaml:"offlineMsgToMongo"`
		MsgToPush struct {
			Topic string `yaml:"topic"`
		} `yaml:"msgToPush"`
		MsgToModify struct {
			Topic string `yaml:"topic"`
		} `yaml:"msgToModify"`
		ConsumerGroupID struct {
			MsgToRedis  string `yaml:"msgToRedis"`
			MsgToMongo  string `yaml:"msgToMongo"`
			MsgToMySql  string `yaml:"msgToMySql"`
			MsgToPush   string `yaml:"msgToPush"`
			MsgToModify string `yaml:"msgToModify"`
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
			TempBucket       string `yaml:"tempBucket"`
			DataBucket       string `yaml:"dataBucket"`
			Location         string `yaml:"location"`
			Endpoint         string `yaml:"endpoint"`
			AccessKeyID      string `yaml:"accessKeyID"`
			SecretAccessKey  string `yaml:"secretAccessKey"`
			IsDistributedMod bool   `yaml:"isDistributedMod"`
		} `yaml:"minio"`
		Tencent struct {
			AppID     string `yaml:"appID"`
			Region    string `yaml:"region"`
			Bucket    string `yaml:"bucket"`
			SecretID  string `yaml:"secretID"`
			SecretKey string `yaml:"secretKey"`
		} `yaml:"tencent"`
		Ali struct {
			RegionID           string `yaml:"regionID"`
			AccessKeyID        string `yaml:"accessKeyID"`
			AccessKeySecret    string `yaml:"accessKeySecret"`
			StsEndpoint        string `yaml:"stsEndpoint"`
			OssEndpoint        string `yaml:"ossEndpoint"`
			Bucket             string `yaml:"bucket"`
			FinalHost          string `yaml:"finalHost"`
			StsDurationSeconds int64  `yaml:"stsDurationSeconds"`
			OssRoleArn         string `yaml:"OssRoleArn"`
		} `yaml:"ali"`
		Aws struct {
			AccessKeyID     string `yaml:"accessKeyID"`
			AccessKeySecret string `yaml:"accessKeySecret"`
			Region          string `yaml:"region"`
			Bucket          string `yaml:"bucket"`
			FinalHost       string `yaml:"finalHost"`
			RoleArn         string `yaml:"roleArn"`
			ExternalId      string `yaml:"externalId"`
			RoleSessionName string `yaml:"roleSessionName"`
		} `yaml:"aws"`
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
		RotationTime        int    `yaml:"rotationTime"`
		RemainRotationCount uint   `yaml:"remainRotationCount"`
		RemainLogLevel      int    `yaml:"remainLogLevel"`
		IsStdout            bool   `yaml:"isStdout"`
		IsJson              bool   `yaml:"isJson"`
		WithStack           bool   `yaml:"withStack"`
	} `yaml:"log"`

	LongConnSvr struct {
		OpenImWsPort        []int `yaml:"openImWsPort"`
		WebsocketMaxConnNum int   `yaml:"websocketMaxConnNum"`
		WebsocketMaxMsgLen  int   `yaml:"websocketMaxMsgLen"`
		WebsocketTimeout    int   `yaml:"websocketTimeout"`
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
	TokenPolicy                       struct {
		AccessSecret string `yaml:"accessSecret"`
		AccessExpire int64  `yaml:"accessExpire"`
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
		CallbackBeforeCreateGroup          CallBackConfig `yaml:"beforeCreateGroup"`
		CallbackBeforeMemberJoinGroup      CallBackConfig `yaml:"beforeMemberJoinGroup"`
		CallbackBeforeSetGroupMemberInfo   CallBackConfig `yaml:"beforeSetGroupMemberInfo"`
	} `yaml:"callback"`

	Prometheus struct {
		Enable                        bool  `yaml:"enable"`
		UserPrometheusPort            []int `yaml:"userPrometheusPort"`
		FriendPrometheusPort          []int `yaml:"friendPrometheusPort"`
		MessagePrometheusPort         []int `yaml:"messagePrometheusPort"`
		MessageGatewayPrometheusPort  []int `yaml:"messageGatewayPrometheusPort"`
		GroupPrometheusPort           []int `yaml:"groupPrometheusPort"`
		AuthPrometheusPort            []int `yaml:"authPrometheusPort"`
		PushPrometheusPort            []int `yaml:"pushPrometheusPort"`
		ConversationPrometheusPort    []int `yaml:"conversationPrometheusPort"`
		RtcPrometheusPort             []int `yaml:"rtcPrometheusPort"`
		MessageTransferPrometheusPort []int `yaml:"messageTransferPrometheusPort"`
		ThirdPrometheusPort           []int `yaml:"thirdPrometheusPort"`
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
	UserInfoUpdated NotificationConf `yaml:"userInfoUpdated"`
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

func GetServiceNames() []string {
	return []string{Config.RpcRegisterName.OpenImUserName, Config.RpcRegisterName.OpenImFriendName, Config.RpcRegisterName.OpenImMsgName, Config.RpcRegisterName.OpenImPushName, Config.RpcRegisterName.OpenImMessageGatewayName,
		Config.RpcRegisterName.OpenImGroupName, Config.RpcRegisterName.OpenImAuthName, Config.RpcRegisterName.OpenImConversationName, Config.RpcRegisterName.OpenImThirdName}
}
