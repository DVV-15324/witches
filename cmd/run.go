package cmd

import (
	"fmt"
	"log"
	"os"

	cmd_run "github.com/DVV-15324/witches/cmd/cmd_run"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var db string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "En: Create a project Vi: Tạo chương trình",
	Long: `En: Create a new project with database configuration Vi: Tạo chương trình"
En: Required: --db=mysql | mssql | postgres
Vi: Yêu cầu: --db=mysql | mssql | postgres
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("missing init project")
			return
		}
		if len(db) < 1 {
			fmt.Println("Error: missing init project, required --db")
			return
		}
		cmd_run.WitchesCreate(args[0], db)
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "En: Run the program Vi: Chay chuong trinh",
	Long:  "En: Run the program Vi: Chay chuong trinh",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Println("Error: missing run project")
			return
		}
		cmd_run.WitchesRun()
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "En: Creates template for project Vi: Hàm tạo templates cho dự án",
	Long:  "En: Creates template for project Vi: Hàm tạo templates cho dự án",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Println("Error: missing run project")
			return
		}
		// En: Load environment variables from witches.env
		// Vi: Kiểm tra biến môi trường từ file witches.env
		err := godotenv.Load("witches.env")
		if err != nil {
			log.Fatal("Error: missing load")
		}
		cmd_run.WitchesInit(os.Getenv("DB_URL"))
	},
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "En: This function installs dependencies Vi: Hàm cài đặt các thư viện cần dùng",
	Long:  "En: This function installs dependencies Vi: Hàm cài đặt các thư viện cần dùng",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Println("Error: missing run project")
			return
		}
		cmd_run.WitchesInstall()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "En: Display the version of The Witch Vi: Hiển thị phiên bản của Witches",
	Long:  "En: Display the version of The Witch Vi: Hiển thị phiên bản của Witches",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Witches version v1.0.0")
	},
}
