package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"strings"
)

func LoadConfig(path string, prefix string, config any) error {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetEnvPrefix(prefix)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(config, func(config *mapstructure.DecoderConfig) {
		config.TagName = "mapstructure"
	}); err != nil {
		return err
	}

	return nil
}
