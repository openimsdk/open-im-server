package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	kdisc "github.com/openimsdk/open-im-server/v3/pkg/common/discovery"
	disetcd "github.com/openimsdk/open-im-server/v3/pkg/common/discovery/etcd"
	"github.com/openimsdk/open-im-server/v3/version"
	"github.com/openimsdk/tools/discovery/etcd"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/runtimeenv"
	"github.com/spf13/cobra"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type RootCmd struct {
	Command        cobra.Command
	processName    string
	port           int
	prometheusPort int
	log            config.Log
	index          int
	configPath     string
	etcdClient     *clientv3.Client
}

func (r *RootCmd) ConfigPath() string {
	return r.configPath
}

func (r *RootCmd) Index() int {
	return r.index
}

func (r *RootCmd) Port() int {
	return r.port
}

type CmdOpts struct {
	loggerPrefixName string
	configMap        map[string]any
}

func WithCronTaskLogName() func(*CmdOpts) {
	return func(opts *CmdOpts) {
		opts.loggerPrefixName = "openim-crontask"
	}
}

func WithLogName(logName string) func(*CmdOpts) {
	return func(opts *CmdOpts) {
		opts.loggerPrefixName = logName
	}
}
func WithConfigMap(configMap map[string]any) func(*CmdOpts) {
	return func(opts *CmdOpts) {
		opts.configMap = configMap
	}
}

func NewRootCmd(processName string, opts ...func(*CmdOpts)) *RootCmd {
	rootCmd := &RootCmd{processName: processName}
	cmd := cobra.Command{
		Use:  "Start openIM application",
		Long: fmt.Sprintf(`Start %s `, processName),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.persistentPreRun(cmd, opts...)
		},
		SilenceUsage:  true,
		SilenceErrors: false,
	}
	cmd.Flags().StringP(config.FlagConf, "c", "", "path of config directory")
	cmd.Flags().IntP(config.FlagTransferIndex, "i", 0, "process startup sequence number")

	rootCmd.Command = cmd
	return rootCmd
}

func (r *RootCmd) initEtcd() error {
	configDirectory, _, err := r.getFlag(&r.Command)
	if err != nil {
		return err
	}
	disConfig := config.Discovery{}
	env := runtimeenv.PrintRuntimeEnvironment()
	err = config.Load(configDirectory, config.DiscoveryConfigFilename, config.EnvPrefixMap[config.DiscoveryConfigFilename],
		env, &disConfig)
	if err != nil {
		return err
	}
	if disConfig.Enable == config.ETCD {
		discov, _ := kdisc.NewDiscoveryRegister(&disConfig, env, nil)
		r.etcdClient = discov.(*etcd.SvcDiscoveryRegistryImpl).GetClient()
	}
	return nil
}

func (r *RootCmd) persistentPreRun(cmd *cobra.Command, opts ...func(*CmdOpts)) error {
	if err := r.initEtcd(); err != nil {
		return err
	}
	cmdOpts := r.applyOptions(opts...)
	if err := r.initializeConfiguration(cmd, cmdOpts); err != nil {
		return err
	}
	if err := r.updateConfigFromEtcd(cmdOpts); err != nil {
		return err
	}
	if err := r.initializeLogger(cmdOpts); err != nil {
		return errs.WrapMsg(err, "failed to initialize logger")
	}
	if err := r.etcdClient.Close(); err != nil {
		return errs.WrapMsg(err, "failed to close etcd client")
	}
	return nil
}

func (r *RootCmd) initializeConfiguration(cmd *cobra.Command, opts *CmdOpts) error {
	configDirectory, _, err := r.getFlag(cmd)
	if err != nil {
		return err
	}

	runtimeEnv := runtimeenv.PrintRuntimeEnvironment()

	// Load common configuration file
	//opts.configMap[ShareFileName] = StructEnvPrefix{EnvPrefix: shareEnvPrefix, ConfigStruct: &r.share}
	for configFileName, configStruct := range opts.configMap {
		err := config.Load(configDirectory, configFileName, config.EnvPrefixMap[configFileName], runtimeEnv, configStruct)
		if err != nil {
			return err
		}
	}
	// Load common log configuration file
	return config.Load(configDirectory, config.LogConfigFileName, config.EnvPrefixMap[config.LogConfigFileName], runtimeEnv, &r.log)
}

func (r *RootCmd) updateConfigFromEtcd(opts *CmdOpts) error {
	if r.etcdClient == nil {
		return nil
	}
	ctx := context.TODO()

	res, err := r.etcdClient.Get(ctx, disetcd.BuildKey(disetcd.EnableConfigCenterKey))
	if err != nil {
		log.ZWarn(ctx, "root cmd updateConfigFromEtcd, etcd Get EnableConfigCenterKey err: %v", errs.Wrap(err))
		return nil
	}
	if res.Count == 0 {
		return nil
	} else {
		if string(res.Kvs[0].Value) == disetcd.Disable {
			return nil
		} else if string(res.Kvs[0].Value) != disetcd.Enable {
			return errs.New("unknown EnableConfigCenter value").Wrap()
		}
	}

	update := func(configFileName string, configStruct any) error {
		key := disetcd.BuildKey(configFileName)
		etcdRes, err := r.etcdClient.Get(ctx, key)
		if err != nil {
			log.ZWarn(ctx, "root cmd updateConfigFromEtcd, etcd Get err: %v", errs.Wrap(err))
			return nil
		}
		if etcdRes.Count == 0 {
			data, err := json.Marshal(configStruct)
			if err != nil {
				return errs.ErrArgs.WithDetail(err.Error()).Wrap()
			}
			_, err = r.etcdClient.Put(ctx, disetcd.BuildKey(configFileName), string(data))
			if err != nil {
				log.ZWarn(ctx, "root cmd updateConfigFromEtcd, etcd Put err: %v", errs.Wrap(err))
			}
			return nil
		}
		err = json.Unmarshal(etcdRes.Kvs[0].Value, configStruct)
		if err != nil {
			return errs.WrapMsg(err, "failed to unmarshal config from etcd")
		}
		return nil
	}
	for configFileName, configStruct := range opts.configMap {
		if err := update(configFileName, configStruct); err != nil {
			return err
		}
	}
	if err := update(config.LogConfigFileName, &r.log); err != nil {
		return err
	}
	// Load common log configuration file
	return nil

}

func (r *RootCmd) applyOptions(opts ...func(*CmdOpts)) *CmdOpts {
	cmdOpts := defaultCmdOpts()
	for _, opt := range opts {
		opt(cmdOpts)
	}

	return cmdOpts
}

func (r *RootCmd) initializeLogger(cmdOpts *CmdOpts) error {
	err := log.InitLoggerFromConfig(

		cmdOpts.loggerPrefixName,
		r.processName,
		"", "",
		r.log.RemainLogLevel,
		r.log.IsStdout,
		r.log.IsJson,
		r.log.StorageLocation,
		r.log.RemainRotationCount,
		r.log.RotationTime,
		version.Version,
		r.log.IsSimplify,
	)
	if err != nil {
		return errs.Wrap(err)
	}
	return errs.Wrap(log.InitConsoleLogger(r.processName, r.log.RemainLogLevel, r.log.IsJson, version.Version))

}

func defaultCmdOpts() *CmdOpts {
	return &CmdOpts{
		loggerPrefixName: "openim-service-log",
	}
}

func (r *RootCmd) getFlag(cmd *cobra.Command) (string, int, error) {
	configDirectory, err := cmd.Flags().GetString(config.FlagConf)
	if err != nil {
		return "", 0, errs.Wrap(err)
	}
	r.configPath = configDirectory
	index, err := cmd.Flags().GetInt(config.FlagTransferIndex)
	if err != nil {
		return "", 0, errs.Wrap(err)
	}
	r.index = index
	return configDirectory, index, nil
}

func (r *RootCmd) Execute() error {
	return r.Command.Execute()
}
