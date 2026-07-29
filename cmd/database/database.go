package cmd

import (
	"fmt"
	content "github.com/DVV-15324/witches/cmd/cmd_utils"
	utils "github.com/DVV-15324/witches/cmd/cmd_utils"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// En: Function to create database containers
// Vi: Chức năng tạo database containers
func WitchesDatabaseDockerUp(DB_DRIVER string, TYPE string) {
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
	DB_PASSWORD := os.Getenv("DB_PASSWORD") // En: Database password            // Vi: Mật khẩu của database
	DB_PORT := os.Getenv("DB_PORT")         // En: Database Port                // Vi: Cổng Port của database
	DB_NAME := os.Getenv("DB_NAME")         // En: Database name                // Vi: Tên của database
	DB_HOST := os.Getenv("DB_HOST")         // En: Database HOST/IP address     // Vi: Địa chỉ HOST/IP của database

	// En: Check required variables for BOTH access and refresh
	// Vi: Kiểm tra biến bắt buộc cho CẢ access và refresh
	if APP_PORT == "" {
		log.Fatal("Error: APP_PORT is not set in environment")
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

	if DB_DRIVER != "" {
		log.Fatal("Error: DB_DRIVER is not set in environment")
	}

	// En: Check specific variables for ACCESS project
	// Vi: Kiểm tra biến đặc thù cho project ACCESS
	if TYPE == "access" {
		// En: Access only needs basic configs, no Redis required
		// Vi: Access chỉ cần cấu hình cơ bản, không cần Redis
		ACCESS_TOKEN_TTL := os.Getenv("ACCESS_TOKEN_TTL")
		if ACCESS_TOKEN_TTL == "" {
			log.Println("Warning: ACCESS_TOKEN_TTL not set, using default: 86400 (1 day)")
			// Set default value
			os.Setenv("ACCESS_TOKEN_TTL", "86400")
		}
		log.Println("Info: Access project - only database required, Redis not needed")

		var DB_USER string
		var DB_URL string

		//En: Verify and configure the database on the user side
		//Vi: Xác thực và thiết lập cấu hình database cho bên người dùng
		switch DB_DRIVER {
		case "mysql":
			DB_USER = "root"
			DB_URL = fmt.Sprintf(
				"mysql://%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				DB_USER,
				DB_PASSWORD,
				DB_HOST,
				DB_PORT,
				DB_NAME,
			)

		case "postgres":
			DB_USER = "postgres"
			DB_URL = fmt.Sprintf(
				"postgres://%s:%s@%s:%s/%s?sslmode=disable",
				DB_USER,
				DB_PASSWORD,
				DB_HOST,
				DB_PORT,
				DB_NAME,
			)

		case "mssql":
			DB_USER = "sa"
			DB_URL = fmt.Sprintf(
				"sqlserver://%s:%s@%s:%s?database=%s&encrypt=disable",
				DB_USER,
				DB_PASSWORD,
				DB_HOST,
				DB_PORT,
				DB_NAME,
			)

		default:
			log.Fatalf("Error: unsupported database: %s. supported : mysql, postgres, mssql", DB_DRIVER)
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
		content := content.CreateContentAccessUsed(DB_DRIVER, DB_HOST, DB_PORT, DB_NAME, DB_PASSWORD, DB_URL, ACCESS_TOKEN_TTL, TYPE)

		if _, err := file.WriteString(content); err != nil {
			log.Fatalf("Error: write to witches.env: %v", err)

		}

		log.Printf("Info: successfully created %s", envPath)

		//En: Get the resource path of the framework
		//Vi: Lấy đưởng dẫn của framework
		frameworkPath := utils.GetFrameworkPath()
		if frameworkPath == "" {
			log.Fatal("Error: framework path not found")
		}
		//En: Get the resource path of the database framework
		//Vi: Lấy đưởng dẫn của database framework
		dockerPath := filepath.Join(
			frameworkPath,
			"pkg",
			"core",
			"database",
		)

		//En :Check if the Docker Path exists.
		//Vi: Kiểm tra đường dẫn dockerPath có tồn tại
		if _, err := os.Stat(dockerPath); os.IsNotExist(err) {
			log.Fatal("Error: Docker config not found")
		}

		log.Printf("Info: starting docker compose: %s", DB_DRIVER)

		//En: Start executing
		//Vi: Bắt đầu thực thi
		cmd := exec.Command(
			"docker",
			"compose",
			"--env-file", envPath,
			"--profile", DB_DRIVER,
			"up",
			"--build",
			"-d",
		)

		cmd.Dir = dockerPath
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatalf("Error: docker compose failed: %v", err)
		}

		log.Printf("Info: successfully started database containers")
	}

	// En: Check specific variables for REFRESH project
	// Vi: Kiểm tra biến đặc thù cho project REFRESH
	if TYPE == "refresh" {
		// En: Refresh needs Redis and all token configs
		// Vi: Refresh cần Redis và tất cả cấu hình token
		REDIS_HOST := os.Getenv("REDIS_HOST")
		REDIS_PORT := os.Getenv("REDIS_PORT")
		REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")
		ACCESS_TOKEN_TTL := os.Getenv("ACCESS_TOKEN_TTL")
		REFRESH_TOKEN_TTL := os.Getenv("REFRESH_TOKEN_TTL")
		SESSION_TTL := os.Getenv("SESSION_TTL")
		IDLE_TIMEOUT := os.Getenv("IDLE_TIMEOUT")
		REVOKED_TTL := os.Getenv("REVOKED_TTL")
		GOOGLE_CLIENT_ID := os.Getenv("GOOGLE_CLIENT_ID")
		RATE_LIMIT_DRIVER := os.Getenv("RATE_LIMIT_DRIVER")

		if GOOGLE_CLIENT_ID == "" {
			log.Fatal("Error: GOOGLE_CLIENT_ID is required for refresh project")
		}

		// En: Check Redis required
		// Vi: Kiểm tra Redis bắt buộc
		if REDIS_HOST == "" {
			log.Fatal("Error: REDIS_HOST is required for refresh project")
		}
		if REDIS_PORT == "" {
			log.Fatal("Error: REDIS_PORT is required for refresh project")
		}

		// En: Check token configs with defaults
		// Vi: Kiểm tra cấu hình token với giá trị mặc định
		if ACCESS_TOKEN_TTL == "" {
			log.Fatal("Error: ACCESS_TOKEN_TTL not set, using default: 900 (15 minutes)")
		}

		if REFRESH_TOKEN_TTL == "" {
			log.Fatal("Error: REFRESH_TOKEN_TTL not set, using default: 604800 (7 days)")
		}

		if SESSION_TTL == "" {
			log.Fatal("Error: SESSION_TTL not set, using default: 604800 (7 days)")
		}

		if IDLE_TIMEOUT == "" {
			log.Fatal("Error: IDLE_TIMEOUT not set, using default: 1800 (30 minutes)")
		}

		if REVOKED_TTL == "" {
			log.Fatal("Error: REVOKED_TTL not set, using default: 900 (15 minutes)")
		}

		if RATE_LIMIT_DRIVER == "" {
			log.Fatal("Error: RATE_LIMIT_DRIVER not set, using default: memory")
		}

		var DB_USER string
		var DB_URL string

		//En: Verify and configure the database on the user side
		//Vi: Xác thực và thiết lập cấu hình database cho bên người dùng
		switch DB_DRIVER {
		case "mysql":
			DB_USER = "root"
			DB_URL = fmt.Sprintf(
				"mysql://%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				DB_USER,
				DB_PASSWORD,
				DB_HOST,
				DB_PORT,
				DB_NAME,
			)

		case "postgres":
			DB_USER = "postgres"
			DB_URL = fmt.Sprintf(
				"postgres://%s:%s@%s:%s/%s?sslmode=disable",
				DB_USER,
				DB_PASSWORD,
				DB_HOST,
				DB_PORT,
				DB_NAME,
			)

		case "mssql":
			DB_USER = "sa"
			DB_URL = fmt.Sprintf(
				"sqlserver://%s:%s@%s:%s?database=%s&encrypt=disable",
				DB_USER,
				DB_PASSWORD,
				DB_HOST,
				DB_PORT,
				DB_NAME,
			)

		default:
			log.Fatalf("Error: unsupported database: %s. supported : mysql, postgres, mssql", DB_DRIVER)
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

		content := content.CreateContentRefreshUsed(DB_DRIVER, DB_HOST, DB_PORT, DB_NAME, DB_PASSWORD, DB_URL, REDIS_HOST, REDIS_PORT, REDIS_PASSWORD, ACCESS_TOKEN_TTL, REFRESH_TOKEN_TTL, SESSION_TTL, IDLE_TIMEOUT, REVOKED_TTL, GOOGLE_CLIENT_ID, RATE_LIMIT_DRIVER, TYPE)

		if _, err := file.WriteString(content); err != nil {
			log.Fatalf("Error: write to witches.env: %v", err)

		}

		log.Printf("Info: successfully created %s", envPath)

		//En: Get the resource path of the framework
		//Vi: Lấy đưởng dẫn của framework
		frameworkPath := utils.GetFrameworkPath()
		if frameworkPath == "" {
			log.Fatal("Error: framework path not found")
		}
		//En: Get the resource path of the database framework
		//Vi: Lấy đưởng dẫn của database framework
		dockerPath := filepath.Join(
			frameworkPath,
			"pkg",
			"core",
			"database",
		)

		//En :Check if the Docker Path exists.
		//Vi: Kiểm tra đường dẫn dockerPath có tồn tại
		if _, err := os.Stat(dockerPath); os.IsNotExist(err) {
			log.Fatal("Error: Docker config not found")
		}

		log.Printf("Info: starting docker compose: %s", DB_DRIVER)

		//En: Start executing
		//Vi: Bắt đầu thực thi

		cmd := exec.Command(
			"docker",
			"compose",
			"--env-file", envPath,
			"--profile", DB_DRIVER,
			"--profile=redis",
			"up",
			"--build",
			"-d",
		)

		cmd.Dir = dockerPath
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatalf("Error: docker compose failed: %v", err)
		}

		log.Printf("Info: successfully started database containers")
	}

	// En: If TYPE is neither access nor refresh
	// Vi: Nếu TYPE không phải access hoặc refresh
	if TYPE != "access" && TYPE != "refresh" {
		log.Fatalf("Error: invalid TYPE '%s'. Must be 'access' or 'refresh'", TYPE)
	}

}
