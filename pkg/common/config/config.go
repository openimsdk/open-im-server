package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../../..")
)

var Config config

type config struct {
	ServerIP      string `yaml:"serverip"`
	ServerVersion string `yaml:"serverversion"`
	Api           struct {
		GinPort []int `yaml:"openImApiPort"`
	}
	Sdk struct {
		WsPort []int `yaml:"openImSdkWsPort"`
	}
	Credential struct {
		Tencent struct {
			AppID     string `yaml:"appID"`
			Region    string `yaml:"region"`
			Bucket    string `yaml:"bucket"`
			SecretID  string `yaml:"secretID"`
			SecretKey string `yaml:"secretKey"`
		}
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
	}
	Mongo struct {
		DBAddress           []string `yaml:"dbAddress"`
		DBDirect            bool     `yaml:"dbDirect"`
		DBTimeout           int      `yaml:"dbTimeout"`
		DBDatabase          string   `yaml:"dbDatabase"`
		DBSource            string   `yaml:"dbSource"`
		DBUserName          string   `yaml:"dbUserName"`
		DBPassword          string   `yaml:"dbPassword"`
		DBMaxPoolSize       int      `yaml:"dbMaxPoolSize"`
		DBRetainChatRecords int      `yaml:"dbRetainChatRecords"`
	}
	Redis struct {
		DBAddress     string `yaml:"dbAddress"`
		DBMaxIdle     int    `yaml:"dbMaxIdle"`
		DBMaxActive   int    `yaml:"dbMaxActive"`
		DBIdleTimeout int    `yaml:"dbIdleTimeout"`
		DBPassWord    string `yaml:"dbPassWord"`
	}
	RpcPort struct {
		OpenImUserPort        []int `yaml:"openImUserPort"`
		openImFriendPort      []int `yaml:"openImFriendPort"`
		RpcMessagePort        []int `yaml:"rpcMessagePort"`
		RpcPushMessagePort    []int `yaml:"rpcPushMessagePort"`
		OpenImGroupPort       []int `yaml:"openImGroupPort"`
		RpcModifyUserInfoPort []int `yaml:"rpcModifyUserInfoPort"`
		RpcGetTokenPort       []int `yaml:"rpcGetTokenPort"`
	}
	RpcRegisterName struct {
		OpenImUserName               string `yaml:"openImUserName"`
		OpenImFriendName             string `yaml:"openImFriendName"`
		OpenImOfflineMessageName     string `yaml:"openImOfflineMessageName"`
		OpenImPushName               string `yaml:"openImPushName"`
		OpenImOnlineMessageRelayName string `yaml:"openImOnlineMessageRelayName"`
		OpenImGroupName              string `yaml:"openImGroupName"`
		OpenImAuthName               string `yaml:"openImAuthName"`
	}
	Etcd struct {
		EtcdSchema string   `yaml:"etcdSchema"`
		EtcdAddr   []string `yaml:"etcdAddr"`
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
		}
		Jpns struct {
			AppKey       string `yaml:"appKey"`
			MasterSecret string `yaml:"masterSecret"`
			PushUrl      string `yaml:"pushUrl"`
			PushIntent   string `yaml:"pushIntent"`
		}
	}
	Manager struct {
		AppManagerUid []string `yaml:"appManagerUid"`
		Secrets       []string `yaml:"secrets"`
	}
	Kafka struct {
		Ws2mschat struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		}
		Ms2pschat struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		}
		ConsumerGroupID struct {
			MsgToMongo string `yaml:"msgToMongo"`
			MsgToMySql string `yaml:"msgToMySql"`
			MsgToPush  string `yaml:"msgToPush"`
		}
	}
	Secret           string `yaml:"secret"`
	MultiLoginPolicy struct {
		OnlyOneTerminalAccess                                  bool `yaml:"onlyOneTerminalAccess"`
		MobileAndPCTerminalAccessButOtherTerminalKickEachOther bool `yaml:"mobileAndPCTerminalAccessButOtherTerminalKickEachOther"`
		AllTerminalAccess                                      bool `yaml:"allTerminalAccess"`
	}
	TokenPolicy struct {
		AccessSecret string `yaml:"accessSecret"`
		AccessExpire int64  `yaml:"accessExpire"`
	}
	MessageCallBack struct {
		CallbackSwitch  bool   `yaml:"callbackSwitch"`
		CallbackUrl     string `yaml:"callbackUrl"`
		CallBackTimeOut int    `yaml:"callbackTimeOut"`
	}
}

func init() {
	//path, _ := os.Getwd()
	//bytes, err := ioutil.ReadFile(path + "/config/config.yaml")
	// if we cd Open-IM-Server/src/utils and run go test
	// it will panic cannot find config/config.yaml
	bytes, err := ioutil.ReadFile(Root + "/config/config.yaml")
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(bytes, &Config); err != nil {
		panic(err)
	}

}
