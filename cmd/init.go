package cmd

import (
	run "github.com/DVV-15324/witches/cmd/run"
	utils "github.com/DVV-15324/witches/cmd/utils"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:  "init",
	Args: cobra.ExactArgs(1),
	Run:  func(cmd *cobra.Command, args []string) {},
}

var initCaptainCmd = &cobra.Command{
	Use:  "captain",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := utils.PreloadNotDBURL()
		utils.LoadDbUrl(cfg)
		run.WitchesInitCaptain(cfg.DBDriver)
	},
}

var initMemberCmd = &cobra.Command{
	Use:  "member",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := utils.PreloadNotDBURL()
		utils.LoadDbUrl(cfg)
		run.WitchesInitMember(cfg.DBDriver)
	},
}
