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

package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/openimsdk/open-im-server/v3/tools/component/checks"
	"github.com/openimsdk/open-im-server/v3/tools/component/util"
	"github.com/openimsdk/tools/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/s3/minio"
)

const defaultCfgPath = "./config.yaml"
const maxRetry = 100

func initConfig(cfgPath string) error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetConfigName("config")
	viper.AddConfigPath("../../../../../config")

	viper.SetEnvPrefix("openim")
	viper.AutomaticEnv()

	if cfgPath != "" {
		viper.SetConfigFile(cfgPath)
	} else if envPath, ok := os.LookupEnv("OPENIM_CONFIG"); ok && envPath != "" {
		viper.SetConfigFile(envPath)
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	fmt.Println("Using config file:", viper.ConfigFileUsed())
	return nil
}

func main() {
	var cfgFile string
	pflag.StringVarP(&cfgFile, "config", "c", "", "config file (default is ./config.yaml)")
	pflag.Parse()

	ctx := context.Background()

	if err := initConfig(cfgFile); err != nil {
		fmt.Fprintf(os.Stderr, "Initialization failed: %v\n", err)
		os.Exit(1)
	}

	// if err := util.ConfigGetEnv(conf); err != nil {
	// 	fmt.Fprintf(os.Stderr, "Environment variable override failed: %v\n", err)
	// 	os.Exit(1)
	// }

	// Define a slice of functions to perform each service check
	serviceChecks := []func(context.Context, *config.GlobalConfig) error{
		func(ctx context.Context, cfg *config.GlobalConfig) error {
			return checks.CheckMongo(ctx, &checks.MongoCheck{Mongo: &cfg.Mongo})
		},
		func(ctx context.Context, cfg *config.GlobalConfig) error {
			return checks.CheckRedis(ctx, &checks.RedisCheck{Redis: &cfg.Redis})
		},
		func(ctx context.Context, cfg *config.GlobalConfig) error {
			return checks.CheckZookeeper(ctx, &checks.ZookeeperCheck{Zookeeper: &cfg.Zookeeper})
		},
		func(ctx context.Context, cfg *config.GlobalConfig) error {
			return checks.CheckKafka(ctx, &checks.KafkaCheck{Kafka: &cfg.Kafka})
		},
	}

	if viper.GetString("object.enable") == "minio" {
		minioConfig := checks.MinioCheck{
			Config: minio.Config(viper.GetString("object.minio")),
			// UseSSL: conf.Minio.UseSSL,
			ApiURL: viper.GetString("object.apiURL"),
		}

		adjustUseSSL(&minioConfig)

		minioCheck := func(ctx context.Context, cfg *config.GlobalConfig) error {
			return checks.CheckMinio(ctx, minioConfig)
		}
		serviceChecks = append(serviceChecks, minioCheck)
	}

	// Execute checks with retry logic
	for i := 0; i < maxRetry; i++ {
		if i > 0 {
			time.Sleep(time.Second)
		}
		fmt.Printf("Checking components, attempt %d/%d\n", i+1, maxRetry)

		allSuccess := true
		for _, check := range serviceChecks {
			if err := check(ctx, &config.Config); err != nil {
				util.ColorErrPrint(fmt.Sprintf("Check failed: %v", err))
				allSuccess = false
				break
			}
		}

		if allSuccess {
			util.SuccessPrint("All components started successfully!")
			return
		}
	}

	util.ErrorPrint("Some components failed to start correctly.")
	os.Exit(-1)
}

// adjustUseSSL updates the UseSSL setting based on the MINIO_USE_SSL environment variable.
func adjustUseSSL(config *checks.MinioCheck) {
	useSSL := config.UseSSL
	if envSSL, exists := os.LookupEnv("MINIO_USE_SSL"); exists {
		parsedSSL, err := strconv.ParseBool(envSSL)
		if err == nil {
			useSSL = parsedSSL
		} else {
			log.CInfo(context.Background(), "Invalid MINIO_USE_SSL value; using config file setting.", "MINIO_USE_SSL", envSSL)
		}
	}
	config.UseSSL = useSSL
}
