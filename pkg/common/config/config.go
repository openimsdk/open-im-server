package config

import (
	"OpenIM/pkg/discoveryregistry"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	_ "embed"
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
	CallbackTimeOut        int   `yaml:"callbackTimeOut"`
	CallbackFailedContinue *bool `yaml:"callbackFailedContinue"`
}

type PConversation struct {
	ReliabilityLevel int  `yaml:"reliabilityLevel"`
	UnreadCount      bool `yaml:"unreadCount"`
}

type POfflinePush struct {
	PushSwitch bool   `yaml:"switch"`
	Title      string `yaml:"title"`
	Desc       string `yaml:"desc"`
	Ext        string `yaml:"ext"`
}
type PDefaultTips struct {
	Tips string `yaml:"tips"`
}

type config struct {
	ServerIP string `yaml:"serverip"`

	RpcRegisterIP string `yaml:"rpcRegisterIP"`
	ListenIP      string `yaml:"listenIP"`

	ServerVersion string `yaml:"serverversion"`
	Api           struct {
		GinPort  []int  `yaml:"openImApiPort"`
		ListenIP string `yaml:"listenIP"`
	}
	Sdk struct {
		WsPort  []int    `yaml:"openImSdkWsPort"`
		DataDir []string `yaml:"dataDir"`
	}
	Credential struct {
		Tencent struct {
			AppID     string `yaml:"appID"`
			Region    string `yaml:"region"`
			Bucket    string `yaml:"bucket"`
			SecretID  string `yaml:"secretID"`
			SecretKey string `yaml:"secretKey"`
		}
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
		}
		Minio struct {
			Bucket              string `yaml:"bucket"`
			AppBucket           string `yaml:"appBucket"`
			Location            string `yaml:"location"`
			Endpoint            string `yaml:"endpoint"`
			AccessKeyID         string `yaml:"accessKeyID"`
			SecretAccessKey     string `yaml:"secretAccessKey"`
			EndpointInner       string `yaml:"endpointInner"`
			EndpointInnerEnable bool   `yaml:"endpointInnerEnable"`
			StorageTime         int    `yaml:"storageTime"`
			IsDistributedMod    bool   `yaml:"isDistributedMod"`
		} `yaml:"minio"`
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
	}

	Mysql struct {
		DBAddress      []string `yaml:"dbMysqlAddress"`
		DBUserName     string   `yaml:"dbMysqlUserName"`
		DBPassword     string   `yaml:"dbMysqlPassword"`
		DBDatabaseName string   `yaml:"dbMysqlDatabaseName"`
		DBTableName    string   `yaml:"DBTableName"`
		DBMsgTableNum  int      `yaml:"dbMsgTableNum"`
		DBMaxOpenConns int      `yaml:"dbMaxOpenConns"`
		DBMaxIdleConns int      `yaml:"dbMaxIdleConns"`
		DBMaxLifeTime  int      `yaml:"dbMaxLifeTime"`
		LogLevel       int      `yaml:"logLevel"`
		SlowThreshold  int      `yaml:"slowThreshold"`
	}
	Mongo struct {
		DBUri                string   `yaml:"dbUri"`
		DBAddress            []string `yaml:"dbAddress"`
		DBDirect             bool     `yaml:"dbDirect"`
		DBTimeout            int      `yaml:"dbTimeout"`
		DBDatabase           string   `yaml:"dbDatabase"`
		DBSource             string   `yaml:"dbSource"`
		DBUserName           string   `yaml:"dbUserName"`
		DBPassword           string   `yaml:"dbPassword"`
		DBMaxPoolSize        int      `yaml:"dbMaxPoolSize"`
		DBRetainChatRecords  int      `yaml:"dbRetainChatRecords"`
		ChatRecordsClearTime string   `yaml:"chatRecordsClearTime"`
	}
	Redis struct {
		DBAddress     []string `yaml:"dbAddress"`
		DBMaxIdle     int      `yaml:"dbMaxIdle"`
		DBMaxActive   int      `yaml:"dbMaxActive"`
		DBIdleTimeout int      `yaml:"dbIdleTimeout"`
		DBUserName    string   `yaml:"dbUserName"`
		DBPassWord    string   `yaml:"dbPassWord"`
		EnableCluster bool     `yaml:"enableCluster"`
	}
	RpcPort struct {
		OpenImUserPort           []int `yaml:"openImUserPort"`
		OpenImFriendPort         []int `yaml:"openImFriendPort"`
		OpenImMessagePort        []int `yaml:"openImMessagePort"`
		OpenImMessageGatewayPort []int `yaml:"openImMessageGatewayPort"`
		OpenImGroupPort          []int `yaml:"openImGroupPort"`
		OpenImAuthPort           []int `yaml:"openImAuthPort"`
		OpenImPushPort           []int `yaml:"openImPushPort"`
		OpenImConversationPort   []int `yaml:"openImConversationPort"`
		OpenImCachePort          []int `yaml:"openImCachePort"`
		OpenImRtcPort            []int `yaml:"openImRtcPort"`
		OpenImThirdPort          []int `yaml:"openImThirdPort"`
	}
	RpcRegisterName struct {
		OpenImUserName           string `yaml:"openImUserName"`
		OpenImFriendName         string `yaml:"openImFriendName"`
		OpenImMsgName            string `yaml:"openImMsgName"`
		OpenImPushName           string `yaml:"openImPushName"`
		OpenImMessageGatewayName string `yaml:"openImMessageGatewayName"`
		OpenImGroupName          string `yaml:"openImGroupName"`
		OpenImAuthName           string `yaml:"openImAuthName"`
		OpenImConversationName   string `yaml:"openImConversationName"`
		OpenImCacheName          string `yaml:"openImCacheName"`
		OpenImRtcName            string `yaml:"openImRtcName"`
		OpenImThirdName          string `yaml:"openImThirdName"`
	}
	Zookeeper struct {
		Schema   string   `yaml:"schema"`
		ZkAddr   []string `yaml:"zkAddr"`
		UserName string   `yaml:"userName"`
		Password string   `yaml:"password"`
	}
	Log struct {
		StorageLocation       string   `yaml:"storageLocation"`
		RotationTime          int      `yaml:"rotationTime"`
		RemainRotationCount   uint     `yaml:"remainRotationCount"`
		RemainLogLevel        uint     `yaml:"remainLogLevel"`
		ElasticSearchSwitch   bool     `yaml:"elasticSearchSwitch"`
		ElasticSearchAddr     []string `yaml:"elasticSearchAddr"`
		ElasticSearchUser     string   `yaml:"elasticSearchUser"`
		ElasticSearchPassword string   `yaml:"elasticSearchPassword"`
	}
	ModuleName struct {
		LongConnSvrName string `yaml:"longConnSvrName"`
		MsgTransferName string `yaml:"msgTransferName"`
		PushName        string `yaml:"pushName"`
	}
	LongConnSvr struct {
		WebsocketPort       []int `yaml:"openImWsPort"`
		WebsocketMaxConnNum int   `yaml:"websocketMaxConnNum"`
		WebsocketMaxMsgLen  int   `yaml:"websocketMaxMsgLen"`
		WebsocketTimeOut    int   `yaml:"websocketTimeOut"`
	}

	Push struct {
		Jpns struct {
			AppKey       string `yaml:"appKey"`
			MasterSecret string `yaml:"masterSecret"`
			PushUrl      string `yaml:"pushUrl"`
			PushIntent   string `yaml:"pushIntent"`
			Enable       bool   `yaml:"enable"`
		}
		Getui struct {
			PushUrl      string `yaml:"pushUrl"`
			AppKey       string `yaml:"appKey"`
			Enable       bool   `yaml:"enable"`
			Intent       string `yaml:"intent"`
			MasterSecret string `yaml:"masterSecret"`
			ChannelID    string `yaml:"channelID"`
			ChannelName  string `yaml:"channelName"`
		}
		Fcm struct {
			ServiceAccount string `yaml:"serviceAccount"`
			Enable         bool   `yaml:"enable"`
		}
	}
	Manager struct {
		AppManagerUid []string `yaml:"appManagerUid"`
		//	AppSysNotificationName string   `yaml:"appSysNotificationName"`
		Nickname []string `yaml:"nickname"`
	}

	Kafka struct {
		SASLUserName string `yaml:"SASLUserName"`
		SASLPassword string `yaml:"SASLPassword"`
		Ws2mschat    struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		}
		MsgToMongo struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		}
		Ms2pschat struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		}
		MsgToModify struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		}
		ConsumerGroupID struct {
			MsgToRedis  string `yaml:"msgToTransfer"`
			MsgToMongo  string `yaml:"msgToMongo"`
			MsgToMySql  string `yaml:"msgToMySql"`
			MsgToPush   string `yaml:"msgToPush"`
			MsgToModify string `yaml:"msgToModify"`
		}
	}
	Secret                            string `yaml:"secret"`
	MultiLoginPolicy                  int    `yaml:"multiloginpolicy"`
	ChatPersistenceMysql              bool   `yaml:"chatpersistencemysql"`
	MsgCacheTimeout                   int    `yaml:"msgCacheTimeout"`
	GroupMessageHasReadReceiptEnable  bool   `yaml:"groupMessageHasReadReceiptEnable"`
	SingleMessageHasReadReceiptEnable bool   `yaml:"singleMessageHasReadReceiptEnable"`

	TokenPolicy struct {
		AccessSecret string `yaml:"accessSecret"`
		AccessExpire int64  `yaml:"accessExpire"`
	}
	MessageVerify struct {
		FriendVerify *bool `yaml:"friendVerify"`
	}
	IOSPush struct {
		PushSound  string `yaml:"pushSound"`
		BadgeCount bool   `yaml:"badgeCount"`
		Production bool   `yaml:"production"`
	}
	Callback struct {
		CallbackUrl                        string         `yaml:"callbackUrl"`
		CallbackBeforeSendSingleMsg        CallBackConfig `yaml:"callbackBeforeSendSingleMsg"`
		CallbackAfterSendSingleMsg         CallBackConfig `yaml:"callbackAfterSendSingleMsg"`
		CallbackBeforeSendGroupMsg         CallBackConfig `yaml:"callbackBeforeSendGroupMsg"`
		CallbackAfterSendGroupMsg          CallBackConfig `yaml:"callbackAfterSendGroupMsg"`
		CallbackMsgModify                  CallBackConfig `yaml:"callbackMsgModify"`
		CallbackUserOnline                 CallBackConfig `yaml:"callbackUserOnline"`
		CallbackUserOffline                CallBackConfig `yaml:"callbackUserOffline"`
		CallbackUserKickOff                CallBackConfig `yaml:"callbackUserKickOff"`
		CallbackOfflinePush                CallBackConfig `yaml:"callbackOfflinePush"`
		CallbackOnlinePush                 CallBackConfig `yaml:"callbackOnlinePush"`
		CallbackBeforeSuperGroupOnlinePush CallBackConfig `yaml:"callbackSuperGroupOnlinePush"`
		CallbackBeforeAddFriend            CallBackConfig `yaml:"callbackBeforeAddFriend"`
		CallbackBeforeCreateGroup          CallBackConfig `yaml:"callbackBeforeCreateGroup"`
		CallbackBeforeMemberJoinGroup      CallBackConfig `yaml:"callbackBeforeMemberJoinGroup"`
		CallbackBeforeSetGroupMemberInfo   CallBackConfig `yaml:"callbackBeforeSetGroupMemberInfo"`
	} `yaml:"callback"`
	Notification Notification `yaml:"notification"`
	Rtc          struct {
		SignalTimeout string `yaml:"signalTimeout"`
	} `yaml:"rtc"`

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
	GroupCreated struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"groupCreated"`

	GroupInfoSet struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"groupInfoSet"`

	JoinGroupApplication struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"joinGroupApplication"`

	MemberQuit struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"memberQuit"`

	GroupApplicationAccepted struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"groupApplicationAccepted"`

	GroupApplicationRejected struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"groupApplicationRejected"`

	GroupOwnerTransferred struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"groupOwnerTransferred"`

	MemberKicked struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"memberKicked"`

	MemberInvited struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"memberInvited"`

	MemberEnter struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"memberEnter"`

	GroupDismissed struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"groupDismissed"`

	GroupMuted struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"groupMuted"`

	GroupCancelMuted struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"groupCancelMuted"`

	GroupMemberMuted struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"groupMemberMuted"`

	GroupMemberCancelMuted struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"groupMemberCancelMuted"`
	GroupMemberInfoSet struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"groupMemberInfoSet"`
	GroupMemberSetToAdmin struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"groupMemberSetToAdmin"`
	GroupMemberSetToOrdinary struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"groupMemberSetToOrdinaryUser"`

	////////////////////////user///////////////////////
	UserInfoUpdated struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"userInfoUpdated"`

	//////////////////////friend///////////////////////
	FriendApplication struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"friendApplicationAdded"`
	FriendApplicationApproved struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"friendApplicationApproved"`

	FriendApplicationRejected struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"friendApplicationRejected"`

	FriendAdded struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"friendAdded"`

	FriendDeleted struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"friendDeleted"`
	FriendRemarkSet struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"friendRemarkSet"`
	BlackAdded struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"blackAdded"`
	BlackDeleted struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"blackDeleted"`
	FriendInfoUpdated struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"friendInfoUpdated"`

	ConversationOptUpdate struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  PDefaultTips  `yaml:"defaultTips"`
	} `yaml:"conversationOptUpdate"`
	ConversationSetPrivate struct {
		Conversation PConversation `yaml:"conversation"`
		OfflinePush  POfflinePush  `yaml:"offlinePush"`
		DefaultTips  struct {
			OpenTips  string `yaml:"openTips"`
			CloseTips string `yaml:"closeTips"`
		} `yaml:"defaultTips"`
	} `yaml:"conversationSetPrivate"`
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
		fmt.Println(err.Error())
		if !os.IsNotExist(err) {
			return err
		}
		configPath = filepath.Join(Root, "config", configName)
		fmt.Println(configPath, "not exist, use", configPath)
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
