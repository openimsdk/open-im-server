package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../../..")
)

var Config config

type callBackConfig struct {
	Enable                 bool `yaml:"enable"`
	CallbackTimeOut        int  `yaml:"callbackTimeOut"`
	CallbackFailedContinue bool `yaml:"callbackFailedContinue"`
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
	CmsApi struct {
		GinPort  []int  `yaml:"openImCmsApiPort"`
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
		OpenImAdminCmsPort       []int `yaml:"openImAdminCmsPort"`
		OpenImOfficePort         []int `yaml:"openImOfficePort"`
		OpenImOrganizationPort   []int `yaml:"openImOrganizationPort"`
		OpenImConversationPort   []int `yaml:"openImConversationPort"`
		OpenImCachePort          []int `yaml:"openImCachePort"`
		OpenImRealTimeCommPort   []int `yaml:"openImRealTimeCommPort"`
	}
	RpcRegisterName struct {
		OpenImUserName   string `yaml:"openImUserName"`
		OpenImFriendName string `yaml:"openImFriendName"`
		//	OpenImOfflineMessageName     string `yaml:"openImOfflineMessageName"`
		OpenImMsgName          string `yaml:"openImMsgName"`
		OpenImPushName         string `yaml:"openImPushName"`
		OpenImRelayName        string `yaml:"openImRelayName"`
		OpenImGroupName        string `yaml:"openImGroupName"`
		OpenImAuthName         string `yaml:"openImAuthName"`
		OpenImAdminCMSName     string `yaml:"openImAdminCMSName"`
		OpenImOfficeName       string `yaml:"openImOfficeName"`
		OpenImOrganizationName string `yaml:"openImOrganizationName"`
		OpenImConversationName string `yaml:"openImConversationName"`
		OpenImCacheName        string `yaml:"openImCacheName"`
		OpenImRealTimeCommName string `yaml:"openImRealTimeCommName"`
	}
	Etcd struct {
		EtcdSchema string   `yaml:"etcdSchema"`
		EtcdAddr   []string `yaml:"etcdAddr"`
		UserName   string   `yaml:"userName"`
		Password   string   `yaml:"password"`
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
		Tpns struct {
			Ios struct {
				AccessID  string `yaml:"accessID"`
				SecretKey string `yaml:"secretKey"`
			}
			Android struct {
				AccessID  string `yaml:"accessID"`
				SecretKey string `yaml:"secretKey"`
			}
			Enable bool `yaml:"enable"`
		}
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
			Enable       *bool  `yaml:"enable"`
			Intent       string `yaml:"intent"`
			MasterSecret string `yaml:"masterSecret"`
			ChannelID    string `yaml:"channelID"`
			ChannelName  string `yaml:"channelName"`
		}
		Fcm struct {
			ServiceAccount string `yaml:"serviceAccount"`
			Enable         bool   `yaml:"enable"`
		}
		Mob struct {
			AppKey    string `yaml:"appKey"`
			PushUrl   string `yaml:"pushUrl"`
			Scheme    string `yaml:"scheme"`
			AppSecret string `yaml:"appSecret"`
			Enable    bool   `yaml:"enable"`
		}
	}
	Manager struct {
		AppManagerUid          []string `yaml:"appManagerUid"`
		Secrets                []string `yaml:"secrets"`
		AppSysNotificationName string   `yaml:"appSysNotificationName"`
	}

	Kafka struct {
		SASLUserName string `yaml:"SASLUserName"`
		SASLPassword string `yaml:"SASLPassword"`
		Ws2mschat    struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		}
		//Ws2mschatOffline struct {
		//	Addr  []string `yaml:"addr"`
		//	Topic string   `yaml:"topic"`
		//}
		MsgToMongo struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		}
		Ms2pschat struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		}
		ConsumerGroupID struct {
			MsgToRedis string `yaml:"msgToTransfer"`
			MsgToMongo string `yaml:"msgToMongo"`
			MsgToMySql string `yaml:"msgToMySql"`
			MsgToPush  string `yaml:"msgToPush"`
		}
	}
	Secret                            string `yaml:"secret"`
	MultiLoginPolicy                  int    `yaml:"multiloginpolicy"`
	ChatPersistenceMysql              bool   `yaml:"chatpersistencemysql"`
	ReliableStorage                   bool   `yaml:"reliablestorage"`
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
		CallbackBeforeSendSingleMsg        callBackConfig `yaml:"callbackBeforeSendSingleMsg"`
		CallbackAfterSendSingleMsg         callBackConfig `yaml:"callbackAfterSendSingleMsg"`
		CallbackBeforeSendGroupMsg         callBackConfig `yaml:"callbackBeforeSendGroupMsg"`
		CallbackAfterSendGroupMsg          callBackConfig `yaml:"callbackAfterSendGroupMsg"`
		CallbackMsgModify                  callBackConfig `yaml:"callbackMsgModify"`
		CallbackUserOnline                 callBackConfig `yaml:"callbackUserOnline"`
		CallbackUserOffline                callBackConfig `yaml:"callbackUserOffline"`
		CallbackUserKickOff                callBackConfig `yaml:"callbackUserKickOff"`
		CallbackOfflinePush                callBackConfig `yaml:"callbackOfflinePush"`
		CallbackOnlinePush                 callBackConfig `yaml:"callbackOnlinePush"`
		CallbackBeforeSuperGroupOnlinePush callBackConfig `yaml:"callbackSuperGroupOnlinePush"`
		CallbackBeforeAddFriend            callBackConfig `yaml:"callbackBeforeAddFriend"`
		CallbackBeforeCreateGroup          callBackConfig `yaml:"callbackBeforeCreateGroup"`
	} `yaml:"callback"`
	Notification struct {
		///////////////////////group/////////////////////////////
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
		OrganizationChanged struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"organizationChanged"`

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
		WorkMomentsNotification struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"workMomentsNotification"`
		JoinDepartmentNotification struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"joinDepartmentNotification"`
		Signal struct {
			OfflinePush struct {
				Title string `yaml:"title"`
			} `yaml:"offlinePush"`
		} `yaml:"signal"`
	}
	Demo struct {
		Port         []int  `yaml:"openImDemoPort"`
		ListenIP     string `yaml:"listenIP"`
		AliSMSVerify struct {
			AccessKeyID                  string `yaml:"accessKeyId"`
			AccessKeySecret              string `yaml:"accessKeySecret"`
			SignName                     string `yaml:"signName"`
			VerificationCodeTemplateCode string `yaml:"verificationCodeTemplateCode"`
			Enable                       bool   `yaml:"enable"`
		}
		TencentSMS struct {
			AppID                        string `yaml:"appID"`
			Region                       string `yaml:"region"`
			SecretID                     string `yaml:"secretID"`
			SecretKey                    string `yaml:"secretKey"`
			SignName                     string `yaml:"signName"`
			VerificationCodeTemplateCode string `yaml:"verificationCodeTemplateCode"`
			Enable                       bool   `yaml:"enable"`
		}
		SuperCode    string `yaml:"superCode"`
		CodeTTL      int    `yaml:"codeTTL"`
		UseSuperCode bool   `yaml:"useSuperCode"`
		Mail         struct {
			Title                   string `yaml:"title"`
			SenderMail              string `yaml:"senderMail"`
			SenderAuthorizationCode string `yaml:"senderAuthorizationCode"`
			SmtpAddr                string `yaml:"smtpAddr"`
			SmtpPort                int    `yaml:"smtpPort"`
		}
		TestDepartMentID                        string   `yaml:"testDepartMentID"`
		ImAPIURL                                string   `yaml:"imAPIURL"`
		NeedInvitationCode                      bool     `yaml:"needInvitationCode"`
		OnboardProcess                          bool     `yaml:"onboardProcess"`
		JoinDepartmentIDList                    []string `yaml:"joinDepartmentIDList"`
		JoinDepartmentGroups                    bool     `yaml:"joinDepartmentGroups"`
		OaNotification                          bool     `yaml:"oaNotification"`
		CreateOrganizationUserAndJoinDepartment bool     `yaml:"createOrganizationUserAndJoinDepartment"`
	}
	WorkMoment struct {
		OnlyFriendCanSee bool `yaml:"onlyFriendCanSee"`
	} `yaml:"workMoment"`
	Rtc struct {
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
		AdminCmsPrometheusPort        []int `yaml:"adminCmsPrometheusPort"`
		OfficePrometheusPort          []int `yaml:"officePrometheusPort"`
		OrganizationPrometheusPort    []int `yaml:"organizationPrometheusPort"`
		ConversationPrometheusPort    []int `yaml:"conversationPrometheusPort"`
		CachePrometheusPort           []int `yaml:"cachePrometheusPort"`
		RealTimeCommPrometheusPort    []int `yaml:"realTimeCommPrometheusPort"`
		MessageTransferPrometheusPort []int `yaml:"messageTransferPrometheusPort"`
	} `yaml:"prometheus"`
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

type usualConfig struct {
	Etcd struct {
		UserName string `yaml:"userName"`
		Password string `yaml:"password"`
	} `yaml:"etcd"`
	Mysql struct {
		DBUserName string `yaml:"dbMysqlUserName"`
		DBPassword string `yaml:"dbMysqlPassword"`
	} `yaml:"mysql"`
	Mongo struct {
		DBUserName string `yaml:"dbUserName"`
		DBPassword string `yaml:"dbPassword"`
	} `yaml:"mongo"`
	Redis struct {
		DBUserName string `yaml:"dbUserName"`
		DBPassword string `yaml:"dbPassWord"`
	} `yaml:"redis"`
	Kafka struct {
		SASLUserName string `yaml:"SASLUserName"`
		SASLPassword string `yaml:"SASLPassword"`
	} `yaml:"kafka"`

	Credential struct {
		Minio struct {
			AccessKeyID     string `yaml:"accessKeyID"`
			SecretAccessKey string `yaml:"secretAccessKey"`
			Endpoint        string `yaml:"endPoint"`
		} `yaml:"minio"`
	} `yaml:"credential"`

	Secret string `yaml:"secret"`

	Tokenpolicy struct {
		AccessSecret string `yaml:"accessSecret"`
		AccessExpire int64  `yaml:"accessExpire"`
	} `yaml:"tokenpolicy"`

	Messageverify struct {
		FriendVerify bool `yaml:"friendVerify"`
	} `yaml:"messageverify"`

	Push struct {
		Getui struct {
			PushUrl      string `yaml:"pushUrl"`
			MasterSecret string `yaml:"masterSecret"`
			AppKey       string `yaml:"appKey"`
			Enable       bool   `yaml:"enable"`
		} `yaml:"getui"`
	} `yaml:"push"`
}

var UsualConfig usualConfig

func unmarshalConfig(config interface{}, configName string) {
	var env string
	if configName == "config.yaml" {
		env = "CONFIG_NAME"
	} else if configName == "usualConfig.yaml" {
		env = "USUAL_CONFIG_NAME"
	}
	cfgName := os.Getenv(env)
	if len(cfgName) != 0 {
		bytes, err := ioutil.ReadFile(filepath.Join(cfgName, "config", configName))
		if err != nil {
			bytes, err = ioutil.ReadFile(filepath.Join(Root, "config", configName))
			if err != nil {
				panic(err.Error() + " config: " + filepath.Join(cfgName, "config", configName))
			}
		} else {
			Root = cfgName
		}
		if err = yaml.Unmarshal(bytes, config); err != nil {
			panic(err.Error())
		}
	} else {
		bytes, err := ioutil.ReadFile(fmt.Sprintf("../config/%s", configName))
		if err != nil {
			panic(err.Error())
		}
		if err = yaml.Unmarshal(bytes, config); err != nil {
			panic(err.Error())
		}
	}
}

func init() {
	unmarshalConfig(&Config, "config.yaml")
	unmarshalConfig(&UsualConfig, "usualConfig.yaml")
	fmt.Println(UsualConfig)
	if Config.Etcd.UserName == "" {
		Config.Etcd.UserName = UsualConfig.Etcd.UserName
	}
	if Config.Etcd.Password == "" {
		Config.Etcd.Password = UsualConfig.Etcd.Password
	}

	if Config.Mysql.DBUserName == "" {
		Config.Mysql.DBUserName = UsualConfig.Mysql.DBUserName
	}
	if Config.Mysql.DBPassword == "" {
		Config.Mysql.DBPassword = UsualConfig.Mysql.DBPassword
	}

	if Config.Redis.DBUserName == "" {
		Config.Redis.DBUserName = UsualConfig.Redis.DBUserName
	}
	if Config.Redis.DBPassWord == "" {
		Config.Redis.DBPassWord = UsualConfig.Redis.DBPassword
	}

	if Config.Mongo.DBUserName == "" {
		Config.Mongo.DBUserName = UsualConfig.Mongo.DBUserName
	}
	if Config.Mongo.DBPassword == "" {
		Config.Mongo.DBPassword = UsualConfig.Mongo.DBPassword
	}

	if Config.Kafka.SASLUserName == "" {
		Config.Kafka.SASLUserName = UsualConfig.Kafka.SASLUserName
	}
	if Config.Kafka.SASLPassword == "" {
		Config.Kafka.SASLPassword = UsualConfig.Kafka.SASLPassword
	}

	if Config.Credential.Minio.AccessKeyID == "" {
		Config.Credential.Minio.AccessKeyID = UsualConfig.Credential.Minio.AccessKeyID
	}
	if Config.Credential.Minio.SecretAccessKey == "" {
		Config.Credential.Minio.SecretAccessKey = UsualConfig.Credential.Minio.SecretAccessKey
	}
	if Config.Credential.Minio.Endpoint == "" {
		Config.Credential.Minio.Endpoint = UsualConfig.Credential.Minio.Endpoint
	}

	if Config.MessageVerify.FriendVerify == nil {
		Config.MessageVerify.FriendVerify = &UsualConfig.Messageverify.FriendVerify
	}

	if Config.Push.Getui.MasterSecret == "" {
		Config.Push.Getui.MasterSecret = UsualConfig.Push.Getui.MasterSecret
	}
	if Config.Push.Getui.AppKey == "" {
		Config.Push.Getui.AppKey = UsualConfig.Push.Getui.AppKey
	}
	if Config.Push.Getui.PushUrl == "" {
		Config.Push.Getui.PushUrl = UsualConfig.Push.Getui.PushUrl
	}
	if Config.Push.Getui.Enable == nil {
		Config.Push.Getui.Enable = &UsualConfig.Push.Getui.Enable
	}

	if Config.Secret == "" {
		Config.Secret = UsualConfig.Secret
	}

	if Config.TokenPolicy.AccessExpire == 0 {
		Config.TokenPolicy.AccessExpire = UsualConfig.Tokenpolicy.AccessExpire
	}
	if Config.TokenPolicy.AccessSecret == "" {
		Config.TokenPolicy.AccessSecret = UsualConfig.Tokenpolicy.AccessSecret
	}

}
