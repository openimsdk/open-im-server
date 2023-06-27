package config

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"

	"gopkg.in/yaml.v3"
)

//go:embed version
var Version string

var (
	_, b, _, _ = runtime.Caller(0)
	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../../..")
)

const (
	FileName             = "config.yaml"
	NotificationFileName = "notification.yaml"
	ENV                  = "CONFIG_NAME"
	DefaultFolderPath    = "../config/"
	ConfKey              = "conf"
)

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
		UserName string   `yaml:"userName"`
		Password string   `yaml:"password"`
	} `yaml:"zookeeper"`

	Mysql struct {
		DBAddress      []string `yaml:"address"`
		DBUserName     string   `yaml:"userName"`
		DBPassword     string   `yaml:"password"`
		DBDatabaseName string   `yaml:"databaseName"`
		DBMaxOpenConns int      `yaml:"maxOpenConns"`
		DBMaxIdleConns int      `yaml:"maxIdleConns"`
		DBMaxLifeTime  int      `yaml:"maxLifeTime"`
		LogLevel       int      `yaml:"logLevel"`
		SlowThreshold  int      `yaml:"slowThreshold"`
	} `yaml:"mysql"`

	Mongo struct {
		DBUri                string   `yaml:"uri"`
		DBAddress            []string `yaml:"address"`
		DBTimeout            int      `yaml:"timeout"`
		DBDatabase           string   `yaml:"database"`
		DBSource             string   `yaml:"source"`
		DBUserName           string   `yaml:"userName"`
		DBPassword           string   `yaml:"password"`
		DBMaxPoolSize        int      `yaml:"maxPoolSize"`
		DBRetainChatRecords  int      `yaml:"retainChatRecords"`
		ChatRecordsClearTime string   `yaml:"chatRecordsClearTime"`
	} `yaml:"mongo"`

	Redis struct {
		DBAddress  []string `yaml:"address"`
		DBUserName string   `yaml:"userName"`
		DBPassWord string   `yaml:"passWord"`
	} `yaml:"redis"`

	Kafka struct {
		SASLUserName     string   `yaml:"SASLUserName"`
		SASLPassword     string   `yaml:"SASLPassword"`
		Addr             []string `yaml:"addr"`
		LatestMsgToRedis struct {
			Topic string `yaml:"topic"`
		} `yaml:"latestMsgToRedis"`
		OfflineMsgToMongoMysql struct {
			Topic string `yaml:"topic"`
		} `yaml:"offlineMsgToMongoMysql"`
		MsqToPush struct {
			Topic string `yaml:"topic"`
		} `yaml:"msqToPush"`
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
		ListenIP string `yaml:"listenIP"`
	} `yaml:"api"`

	Sdk struct {
		DataDir []string `yaml:"dataDir"`
	} `yaml:"sdk"`

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
		OpenImRtcName            string `yaml:"openImRtcName"`
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
		WebsocketMaxConnNum int `yaml:"websocketMaxConnNum"`
		WebsocketMaxMsgLen  int `yaml:"websocketMaxMsgLen"`
		WebsocketTimeOut    int `yaml:"websocketTimeOut"`
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
		AppManagerUserID []string `yaml:"appManagerUserID"`
		Nickname         []string `yaml:"nickname"`
	} `yaml:"manager"`

	MultiLoginPolicy                  int  `yaml:"multiLoginPolicy"`
	ChatPersistenceMysql              bool `yaml:"chatPersistenceMysql"`
	MsgCacheTimeout                   int  `yaml:"msgCacheTimeout"`
	GroupMessageHasReadReceiptEnable  bool `yaml:"groupMessageHasReadReceiptEnable"`
	SingleMessageHasReadReceiptEnable bool `yaml:"singleMessageHasReadReceiptEnable"`

	TokenPolicy struct {
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
	Notification Notification `yaml:"notification"`

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
}

type Notification struct {
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

func GetOptionsByNotification(cfg NotificationConf) utils.Options {
	opts := utils.NewOptions()
	if cfg.UnreadCount {
		opts = utils.WithOptions(opts, utils.WithUnreadCount(true))
	}
	if cfg.OfflinePush.Enable {
		opts = utils.WithOptions(opts, utils.WithOfflinePush(true))
	}
	switch cfg.ReliabilityLevel {
	case constant.UnreliableNotification:
	case constant.ReliableNotificationNoMsg:
		opts = utils.WithOptions(opts, utils.WithHistory(true), utils.WithPersistent())
	}
	opts = utils.WithOptions(opts, utils.WithSendMsg(cfg.IsSendMsg))
	return opts
}

func (c *config) unmarshalConfig(config interface{}, configPath string) error {
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(bytes, config); err != nil {
		return err
	}
	return nil
}

func (c *config) initConfig(config interface{}, configName, configFolderPath string) error {
	if configFolderPath == "" {
		configFolderPath = DefaultFolderPath
	}
	configPath := filepath.Join(configFolderPath, configName)
	defer func() {
		fmt.Println("use config", configPath)
	}()
	_, err := os.Stat(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		configPath = filepath.Join(Root, "config", configName)
	} else {
		Root = filepath.Dir(configPath)
	}
	return c.unmarshalConfig(config, configPath)
}

func (c *config) RegisterConf2Registry(registry discoveryregistry.SvcDiscoveryRegistry) error {
	bytes, err := yaml.Marshal(Config)
	if err != nil {
		return err
	}
	return registry.RegisterConf2Registry(ConfKey, bytes)
}

func (c *config) GetConfFromRegistry(registry discoveryregistry.SvcDiscoveryRegistry) ([]byte, error) {
	return registry.GetConfFromRegistry(ConfKey)
}

func InitConfig(configFolderPath string) error {
	err := Config.initConfig(&Config, FileName, configFolderPath)
	if err != nil {
		return err
	}
	err = Config.initConfig(&Config.Notification, NotificationFileName, configFolderPath)
	if err != nil {
		return err
	}
	return nil
}
