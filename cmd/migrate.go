package cmd

import (
	"fmt"
	cmd_migrate "github.com/DVV-15324/witches/cmd/cmd_migrate"
	godotenv "github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// En: Migrate command group
// Vi: Nhóm lệnh quản lý Migrate

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "En: Migrate management Vi: Chức năng Migrate",
	Long:  `En: Manage migration execution tasks Vi: Quản lý các tác vụ thực thi migrate`,
	Run: func(cmd *cobra.Command, args []string) {
		// En: Load environment variables from witches.env
		// Vi: Kiểm tra biến môi trường từ file witches.env
		err := godotenv.Load("witches.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

	},
}

var migrateDropCmd = &cobra.Command{
	Use:   "drop",
	Short: "En: Drop 1 migrate Vi: Xóa toàn bộ migrate",
	Long:  `En: Manage migration execution tasks Vi: Quản lý các tác vụ thực thi migrate`,
	Run: func(cmd *cobra.Command, args []string) {
		// En: Load environment variables from witches.env
		// Vi: Kiểm tra biến môi trường từ file witches.env
		err := godotenv.Load("witches.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		db_url := os.Getenv("DB_URL")
		cmd_migrate.WitchesMigrateDrop(db_url)
	},
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "En: Up 1 migrate Vi: Triển khai migrate up 1",
	Long:  `En: Manage migration execution tasks Vi: Quản lý các tác vụ thực thi migrate`,
	Run: func(cmd *cobra.Command, args []string) {
		// En: Load environment variables from witches.env
		// Vi: Kiểm tra biến môi trường từ file witches.env
		err := godotenv.Load("witches.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		db_url := os.Getenv("DB_URL")
		cmd_migrate.WitchesMigrateUp(db_url)
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "En: Down 1 migrate Vi: Triển khai migrate down 1",
	Long:  `En: Manage migration execution tasks Vi: Quản lý các tác vụ thực thi migrate`,
	Run: func(cmd *cobra.Command, args []string) {
		// En: Load environment variables from witches.env
		// Vi: Kiểm tra biến môi trường từ file witches.env
		err := godotenv.Load("witches.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		db_url := os.Getenv("DB_URL")
		cmd_migrate.WitchesMigrateDown(db_url)
	},
}

var migrateVersionCmd = &cobra.Command{
	Use:   "force",
	Short: "En: Force migrate Vi: Thiết lập phiên bản Migrate",
	Long:  `En: Manage migration execution tasks Vi: Quản lý các tác vụ thực thi migrate`,
	Run: func(cmd *cobra.Command, args []string) {
		// En: Load environment variables from witches.env
		// Vi: Kiểm tra biến môi trường từ file witches.env
		err := godotenv.Load("witches.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		db_url := os.Getenv("DB_URL")
		if len(args) < 1 {
			fmt.Println("missing project")
			return
		}
		cmd_migrate.WitchesMigrateForce(db_url, args[0])
	},
}
var migrateForceCmd = &cobra.Command{
	Use:   "version",
	Short: "En: Version migrate Vi: Hàm hiển thị version migrate hiện tại",
	Long:  `En: Manage migration execution tasks Vi: Quản lý các tác vụ thực thi migrate`,
	Run: func(cmd *cobra.Command, args []string) {
		// En: Load environment variables from witches.env
		// Vi: Kiểm tra biến môi trường từ file witches.env
		err := godotenv.Load("witches.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		db_url := os.Getenv("DB_URL")
		cmd_migrate.WitchesMigrateVersion(db_url)
	},
}
