package cmd

import (
	"log"
	"os"
	"os/exec"
)

// En: This function installs dependencies
// Vi: Hàm cài đặt các thư viện cần dùng
func WitchesInstall() {
	tools := []string{
		"github.com/mailru/easyjson/...@latest",
	}
	//En: Start executing
	//Vi: Bắt đầu thực thi
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("Error: failed to install  %v", err)
	}

	//En: Start executing
	//Vi: Bắt đầu thực thi
	for _, tool := range tools {
		cmd := exec.Command("go", "install", tool)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Printf("Error: failed to install %s: %v", tool, err)
		}
	}
}
