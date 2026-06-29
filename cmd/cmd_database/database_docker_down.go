package cmd

import (
	utils "github.com/DVV-15324/witches/cmd/cmd_utils"

	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// En: Function to stop database containers
// Vi: Chức năng dừng database containers
func WitchesDatabaseDockerDown() {
	//En: Get the resource path of the framework
	//Vi: Lấy đường dẫn tài nguyên của framework
	frameworkPath := utils.GetFrameworkPath()

	//En: Get the path to the framework database
	//Vi: Lấy đường dẫn đến framework database
	dockerPath := filepath.Join(frameworkPath, "pkg", "core", "database")

	//En: Check if the folder exists.
	//Vi: Kiểm tra thư mục tồn tại
	if _, err := os.Stat(dockerPath); os.IsNotExist(err) {
		log.Fatalf("Error: docker config not found")
	}

	//En: Start executing
	//Vi: Bắt đầu thực thi
	cmd := exec.Command("docker", "compose", "down")
	cmd.Dir = dockerPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
