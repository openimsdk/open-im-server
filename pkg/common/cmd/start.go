package cmd

//
//type StartFunc[C any] func(ctx context.Context, config *C, client discovery.Conn, server grpc.ServiceRegistrar) error
//
//func Start[C any](fn StartFunc[C]) {
//	var _ RootCmd
//	cmd := cobra.Command{
//		Use:  "Start openIM application",
//		Long: fmt.Sprintf(`Start %s `, program.GetProcessName()),
//		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
//			return rootCmd.persistentPreRun(cmd, opts...)
//		},
//		SilenceUsage:  true,
//		SilenceErrors: false,
//	}
//	cmd.Flags().StringP(config.FlagConf, "c", "", "path of config directory")
//	cmd.Flags().IntP(config.FlagTransferIndex, "i", 0, "process startup sequence number")
//
//
//
//}
//
//func start[C any](fn StartFunc[C]) error {
//
//
//	v := viper.New()
//	v.SetConfigType("yaml")
//	if err := v.ReadConfig(bytes.NewReader(confData)); err != nil {
//		return err
//	}
//	fn := func(conf *mapstructure.DecoderConfig) {
//		conf.TagName = config.StructTagName
//	}
//	if err := v.Unmarshal(val, fn); err != nil {
//		return err
//	}
//
//	return nil
//}
