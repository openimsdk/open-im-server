package main

import (
	"OpenIM/internal/tools"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)


var showSeqCmd = &cobra.Command{
	Use:   "show-seq",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		configFolderPath, _ := cmd.Flags().GetString(constant.FlagConf)
		config.InitConfig(configFolderPath)
	},
}

var


func init() {
	showSeqCmd.Flags().StringP("userID", "u", "", "openIM userID")
	showSeqCmd.Flags().StringP("groupID", "g", "", "openIM groupID")
	startCmd.Flags().StringP(constant.FlagConf, "c", "", "Path to config file folder")

}

func run(configFolderPath string, cmd *cobra.Command) error {
	if err := config.InitConfig(configFolderPath); err != nil {
		return err
	}


	return nil
}

func main() {
	rootCmd := cmd.NewRootCmd()
	rootCmd.AddCommand(showSeqCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}


//
func main() {
	if err := config.InitConfig(""); err != nil {
		panic(err.Error())
	}
	// clear msg by id
	var userIDClearMsg = flag.String("user_id_fix_seq", "", "userID to clear msg and reset seq")
	var superGroupIDClearMsg = flag.String("super_group_id_fix_seq", "", "superGroupID to clear msg and reset seq")
	// show seq by id
	var userIDShowSeq = flag.String("user_id_show_seq", "", "show userID")
	var superGroupIDShowSeq = flag.String("super_group_id_show_seq", "", "userID to clear msg and reset seq")
	// fix seq by id
	var userIDFixSeq = flag.String("user_id_fix_seq", "", "userID to Fix Seq")
	var superGroupIDFixSeq = flag.String("super_group_id_fix_seq", "", "super groupID to fix Seq")
	var fixAllSeq = flag.Bool("fix_all_seq", false, "fix seq")
	flag.Parse()
	msgTool, err := tools.InitMsgTool()
	if err != nil {
		panic(err.Error())
	}
	ctx := context.Background()
	if userIDFixSeq != nil {
		msgTool.GetAndFixUserSeqs(ctx, *userIDFixSeq)
	}
	if superGroupIDFixSeq != nil {
		msgTool.FixGroupSeq(ctx, *superGroupIDFixSeq)
	}
	if fixAllSeq != nil {
		msgTool.FixAllSeq(ctx)
	}
	if userIDClearMsg != nil {
		msgTool.ClearUsersMsg(ctx, []string{*userIDClearMsg})
	}

	if superGroupIDClearMsg != nil {
		msgTool.ClearSuperGroupMsg(ctx, []string{*superGroupIDClearMsg})
	}
	if userIDShowSeq != nil {
		msgTool.ShowUserSeqs(ctx, *userIDShowSeq)
	}

	if superGroupIDShowSeq != nil {
		msgTool.ShowSuperGroupSeqs(ctx, *superGroupIDShowSeq)
	}
}
