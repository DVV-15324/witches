package database

import (
	"fmt"
	utils "github.com/DVV-15324/witches/cmd/utils"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
)

func WitchesDatabase(DB_DRIVER string) {
	//En: Get the current path
	//Vi: Lấy đường dẫn hiện tại
	currentPath := utils.GetCurrentPath()

	//En: Get the path to the witches.env file
	//Vi: Lấy đường dẫn witches.env file
	envPath := filepath.Join(currentPath, "witches.env")

	//Vi: Load để kiểm tra đường dẫn đến witches.env -> envPath
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Error: .env file not found: %v", err)
	}
	// En: Get the value of the environmental variable
	// Vi: Lấy giá trị của biến môi trường
	APP_PORT := os.Getenv("APP_PORT")       // En: Project Port                 // Vi: Cổng Port của dự án
	APP_HOST := os.Getenv("APP_HOST")       // En: Project Host                 // Vi: Cổng Port của dự án
	METRIC_PORT := os.Getenv("METRIC_PORT") // En: Mestrict Port                 // Vi: Cổng Mestrict Port của dự án
	DB_HOST := os.Getenv("DB_HOST")         // En: Database HOST/IP address     // Vi: Địa chỉ HOST/IP của database
	DB_PORT := os.Getenv("DB_PORT")         // En: Database Port                // Vi: Cổng Port của database
	DB_NAME := os.Getenv("DB_NAME")         // En: Database name                // Vi: Tên của database
	DB_USER := os.Getenv("DB_USER")         // En: Database user            	// Vi: user của database
	DB_PASSWORD := os.Getenv("DB_PASSWORD") // En: Database password            // Vi: Mật khẩu của database

	if APP_PORT == "" {
		log.Fatal("Error: APP_PORT is not set in environment")
	}
	if METRIC_PORT == "" {
		log.Fatal("Error: METRIC_PORT is not set in environment")
	}
	if DB_PASSWORD == "" {
		log.Fatal("Error: DB_PASSWORD is not set in environment")
	}
	if DB_NAME == "" {
		log.Fatal("Error: DB_NAME is not set in environment")
	}
	if DB_PORT == "" {
		log.Fatal("Error: DB_PORT is not set in environment")
	}
	if DB_HOST == "" {
		log.Fatal("Error: DB_HOST is not set in environment")
	}

	REDIS_HOST := os.Getenv("REDIS_HOST")
	REDIS_PORT := os.Getenv("REDIS_PORT")
	REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")
	ACCESS_TOKEN_TTL := os.Getenv("ACCESS_TOKEN_TTL")
	REFRESH_TOKEN_TTL := os.Getenv("REFRESH_TOKEN_TTL")
	SESSION_TTL := os.Getenv("SESSION_TTL")
	IDLE_TIMEOUT := os.Getenv("IDLE_TIMEOUT")
	REVOKED_TTL := os.Getenv("REVOKED_TTL")
	RATE_LIMIT_PERIOD := os.Getenv("RATE_LIMIT_PERIOD")
	RATE_LIMIT_MAX := os.Getenv("RATE_LIMIT_MAX")

	if REDIS_HOST == "" {
		log.Fatal("Error: REDIS_HOST is required for refresh project")
	}
	if REDIS_PORT == "" {
		log.Fatal("Error: REDIS_PORT is required for refresh project")
	}
	if REFRESH_TOKEN_TTL == "" {
		log.Fatal("Error: REFRESH_TOKEN_TTL not set")
	}
	if ACCESS_TOKEN_TTL == "" {
		log.Fatal("Error: ACCESS_TOKEN_TTL not set")
	}
	if SESSION_TTL == "" {
		log.Fatal("Error: SESSION_TTL not set")
	}
	if IDLE_TIMEOUT == "" {
		log.Fatal("Error: IDLE_TIMEOUT not set")
	}
	if REVOKED_TTL == "" {
		log.Fatal("Error: REVOKED_TTL not set")
	}

	if RATE_LIMIT_PERIOD == "" {
		log.Fatal("Error: RATE_LIMIT_PERIOD not set")
	}
	if RATE_LIMIT_MAX == "" {
		log.Fatal("Error: RATE_LIMIT_MAX not set")
	}
	var DB_URL string
	switch DB_DRIVER {
	case "mysql":
		DB_URL = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			DB_USER,
			DB_PASSWORD,
			DB_HOST,
			DB_PORT,
			DB_NAME,
		)

	case "postgresql", "postgres":
		DB_URL = fmt.Sprintf(
			"%s:%s@%s:%s/%s?sslmode=disable",
			DB_USER,
			DB_PASSWORD,
			DB_HOST,
			DB_PORT,
			DB_NAME,
		)

	case "mssql", "sqlserver":
		DB_URL = fmt.Sprintf(
			"%s:%s@%s:%s?database=%s&encrypt=disable",
			DB_USER,
			DB_PASSWORD,
			DB_HOST,
			DB_PORT,
			DB_NAME,
		)

	default:
		log.Fatalf("Error: unsupported database: %s. supported : mysql, postgresql, postgres, mssql, sqlserver", DB_DRIVER)
	}

	//En: O_CREATE:Permission to create a file when it doesn't exist
	//En: O_WRONLY:Write permissions only -> avoids making random edits during development.
	//En: O_TRUNC: Write over the old configuration -> avoid spam: witches database docker-up
	//Vi: O_CREATE: Quyền tạo khi file không tồn tại
	//Vi: O_WRONLY: Chỉ quyền ghi -> tránh trường hợp trong lúc phát triến sửa lung tung
	//Vi: O_TRUNC: Ghi đề lên cấu hình cũ -> tránh tình trạng spam: witches database docker-up
	file, err := os.OpenFile(
		envPath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0644,
	)

	//En: Check for errors when creating the file.
	//Vi: Kiểm tra lỗi tạo file
	if err != nil {
		log.Fatalf("Error: create witches.env: %v", err)
	}
	//En: Make sure the file is closed
	//Vi: Chắc chắn file đóng
	defer file.Close()

	content := utils.CreateContentRefreshUsed(APP_PORT, APP_HOST, METRIC_PORT, DB_DRIVER, DB_USER, DB_HOST, DB_PORT, DB_NAME, DB_PASSWORD, DB_URL, REDIS_HOST, REDIS_PORT, REDIS_PASSWORD, ACCESS_TOKEN_TTL, REFRESH_TOKEN_TTL, SESSION_TTL, IDLE_TIMEOUT, REVOKED_TTL, RATE_LIMIT_PERIOD, RATE_LIMIT_MAX)

	if _, err := file.WriteString(content); err != nil {
		log.Fatalf("Error: write to witches.env: %v", err)

	}

	log.Printf("Info: successfully created %s", envPath)

}
