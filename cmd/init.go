package cmd

import (
	run "github.com/DVV-15324/witches/cmd/run"
	utils "github.com/DVV-15324/witches/cmd/utils"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:  "init",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := utils.PreloadNotDBURL()
		utils.LoadDbUrl(cfg)
		run.WitchesInit(cfg.DBDriver)
	},
}
