package cmd

import (
	"fmt"
	run "github.com/DVV-15324/witches/cmd/run"
	utils "github.com/DVV-15324/witches/cmd/utils"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use: "create",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: missing create project")
			return
		}
		run.WitchesCreate(args[0])
	},
}

var runCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Println("Error: missing run project")
			return
		}
		run.WitchesRun()
	},
}

var initCmd = &cobra.Command{
	Use: "init",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Println("Error: missing run project")
			return
		}
		cfg := utils.PreloadNotDBURL()
		utils.Validate(cfg)
		utils.LoadDbUrl(cfg)
		run.WitchesInit(cfg.DBDriver)
	},
}

var installCmd = &cobra.Command{
	Use: "install",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Println("Error: missing install project")
			return
		}
		cfg := utils.PreloadNotDBURL()
		utils.Validate(cfg)
		utils.LoadDbUrl(cfg)
		run.WitchesInstall(cfg.DBDriver)
	},
}

var addCmd = &cobra.Command{
	Use: "add",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: missing service name")
			return
		}
		serviceName := args[0]
		run.WitchesAdd(serviceName)
	},
}
