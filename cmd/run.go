package cmd

import (
	"fmt"
	"log"
	"os"

	run "github.com/DVV-15324/witches/cmd/run"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "En: Create a project Vi: Tạo chương trình",
	Long: `En: Create a new project with database configuration Vi: Tạo chương trình"
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("missing init project")
			return
		}
		run.WitchesCreate(args[0])
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
		run.WitchesRun()
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
		run.WitchesInit()
	},
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "En: This function installs dependencies Vi: Hàm cài đặt các thư viện cần dùng",
	Long:  "En: This function installs dependencies Vi: Hàm cài đặt các thư viện cần dùng",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Println("Error: missing install project")
			return
		}
		db_driver := os.Getenv("DB_DRIVER")
		run.WitchesInstall(db_driver)
	},
}

var addCmd = &cobra.Command{
	Use:   "add <service_name>",
	Short: "En: Add template for project Vi: Hàm thêm templates cho dự án",
	Long:  "En: Add template for project Vi: Hàm thêm templates cho dự án",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: missing service name")
			fmt.Println("Usage: witches add <service_name>")
			fmt.Println("Example: witches add book")
			return
		}
		serviceName := args[0]

		// En: Load environment variables from witches.env
		// Vi: Kiểm tra biến môi trường từ file witches.env
		err := godotenv.Load("witches.env")
		if err != nil {
			log.Fatal("Error: missing load")
		}

		run.WitchesAdd(serviceName)
	},
}
