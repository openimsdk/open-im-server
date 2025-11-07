package config

import (
	"flag"
	"os"
)

// Flags is the flag set that AddOptions adds to. Test authors should
// also use it instead of directly adding to the global command line.
var Flags = flag.NewFlagSet("", flag.ContinueOnError)

// CopyFlags ensures that all flags that are defined in the source flag
// set appear in the target flag set as if they had been defined there
// directly. From the flag package it inherits the behavior that there
// is a panic if the target already contains a flag from the source.
func CopyFlags(source *flag.FlagSet, target *flag.FlagSet) {
	source.VisitAll(func(flag *flag.Flag) {
		// We don't need to copy flag.DefValue. The original
		// default (from, say, flag.String) was stored in
		// the value and gets extracted by Var for the help
		// message.
		target.Var(flag.Value, flag.Name, flag.Usage)
	})
}

// Config defines the configuration structure for the OpenIM components.
type Config struct {
	APIHost             string
	APIPort             string
	MsgGatewayHost      string
	MsgTransferHost     string
	PushHost            string
	RPCAuthHost         string
	RPCConversationHost string
	RPCFriendHost       string
	RPCGroupHost        string
	RPCMsgHost          string
	RPCThirdHost        string
	RPCUserHost         string
	// Add other configuration fields as needed
}

// LoadConfig loads the configurations from environment variables or default values.
func LoadConfig() *Config {
	return &Config{
		APIHost: getEnv("OPENIM_API_HOST", "127.0.0.1"),
		APIPort: getEnv("API_OPENIM_PORT", "10002"),

		// TODO: Set default variable
		MsgGatewayHost:      getEnv("OPENIM_MSGGATEWAY_HOST", "default-msggateway-host"),
		MsgTransferHost:     getEnv("OPENIM_MSGTRANSFER_HOST", "default-msgtransfer-host"),
		PushHost:            getEnv("OPENIM_PUSH_HOST", "default-push-host"),
		RPCAuthHost:         getEnv("OPENIM_RPC_AUTH_HOST", "default-rpc-auth-host"),
		RPCConversationHost: getEnv("OPENIM_RPC_CONVERSATION_HOST", "default-rpc-conversation-host"),
		RPCFriendHost:       getEnv("OPENIM_RPC_FRIEND_HOST", "default-rpc-friend-host"),
		RPCGroupHost:        getEnv("OPENIM_RPC_GROUP_HOST", "default-rpc-group-host"),
		RPCMsgHost:          getEnv("OPENIM_RPC_MSG_HOST", "default-rpc-msg-host"),
		RPCThirdHost:        getEnv("OPENIM_RPC_THIRD_HOST", "default-rpc-third-host"),
		RPCUserHost:         getEnv("OPENIM_RPC_USER_HOST", "default-rpc-user-host"),
	}
}

// getEnv is a helper function to read an environment variable or return a default value.
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
