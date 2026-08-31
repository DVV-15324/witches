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
		cfg := utils.PreloadNotDBURL()
		utils.LoadDbUrl(cfg)
		domainName := args[0]
		run.WitchesAdd(domainName, cfg.DBDriver)
	},
}
var removeCmd = &cobra.Command{
	Use:  "rm",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		run.WitchesRollback(args[0])
	},
}
var linkCmd = &cobra.Command{
	Use:  "link",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		domainName := args[0]
		url := args[1]
		run.WitchesLink(domainName, url)
	},
}
