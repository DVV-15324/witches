package cmd

import (
	"fmt"
	cmd_run "github.com/DVV-15324/witches/cmd/cmd_run"
	"github.com/spf13/cobra"
)

var db string
var runAll bool = false

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Thiết lập wiches cho dự án",
	Long: `Thiết lập wiches cho chương trình
Yêu cầu có ví dụ: --db=mysql | mssql | postgres
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("missing init project")
			return
		}
		if len(db) < 1 {
			fmt.Println("missing init project, required --db")
			return
		}
		cmd_run.WitchesInit(args[0], db)
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Chay chuong trinh",
	Long:  `Chạy chuong trinh`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Println("missing run project")
			return
		}
		cmd_run.WitchesRun()
	},
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Cài đặt các phụ thuộc cần thiết",
	Long:  `Cài đặt install easyjson && swag && migrate`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Println("missing run project")
			return
		}
		cmd_run.WitchesInstall()
	},
}
