package cmd

import (
	run "github.com/DVV-15324/witches/cmd/run"
	utils "github.com/DVV-15324/witches/cmd/utils"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:  "create",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		run.WitchesCreate(args[0])
	},
}

var runCmd = &cobra.Command{
	Use:  "run",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		run.WitchesRun()
	},
}

var initCmd = &cobra.Command{
	Use:  "init",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := utils.PreloadNotDBURL()
		utils.LoadDbUrl(cfg)
		run.WitchesInit(cfg.DBDriver)
	},
}

var installCmd = &cobra.Command{
	Use:  "install",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := utils.PreloadNotDBURL()
		utils.LoadDbUrl(cfg)
		run.WitchesInstall(cfg.DBDriver)
	},
}

var addCmd = &cobra.Command{
	Use:  "add",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domainName := args[0]
		run.WitchesAdd(domainName)
	},
}
var removeCmd = &cobra.Command{
	Use:  "rm",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		run.WitchesRollback(args[0])
	},
}
