package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`

	DB struct {
		Driver string `yaml:"driver"`
		DSN    string `yaml:"dsn"`
	} `yaml:"db"`

	Chain struct {
		RPCURL             string `yaml:"rpc_url"`
		ContractAddress    string `yaml:"contract_address"`
		ChainID            int64  `yaml:"chain_id"`
		SignerPrivateKey   string `yaml:"signer_private_key"`
		ConfigAdminPrivateKey string `yaml:"config_admin_private_key"`
	} `yaml:"chain"`

	Tron struct {
		FullNodeURL   string `yaml:"full_node_url"`
		ContractBase58 string `yaml:"contract_base58"`
		OwnerBase58   string `yaml:"owner_base58"`
		PrivateKeyHex string `yaml:"private_key_hex"`
		FeeLimit      int64  `yaml:"fee_limit"`
	} `yaml:"tron"`

	Indexer struct {
		PollInterval int `yaml:"poll_interval"`
	} `yaml:"indexer"`
}

var Cfg Config

// Load loads configuration from YAML file
func Load(configPath string) {
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Warning: could not read config file %s: %v, using defaults\n", configPath, err)
		setDefaults()
		return
	}

	if err := yaml.Unmarshal(data, &Cfg); err != nil {
		fmt.Printf("Warning: could not parse config: %v, using defaults\n", err)
		setDefaults()
		return
	}

	fmt.Printf("Loaded config from %s\n", configPath)
}

func setDefaults() {
	Cfg.Server.Port = 8080
	Cfg.DB.Driver = "sqlite"
	Cfg.DB.DSN = "redpacket.db"
	Cfg.Chain.ChainID = 1
	Cfg.Indexer.PollInterval = 5
}
