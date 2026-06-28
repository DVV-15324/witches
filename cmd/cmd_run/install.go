package cmd

import (
	"log"
	"os"
	"os/exec"
)

func WitchesInstall() {
	tools := []string{
		"github.com/mailru/easyjson/...@latest",
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("Warning: failed to installcd  %v", err)
	}

	for _, tool := range tools {
		cmd := exec.Command("go", "install", tool)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Printf("Warning: failed to install %s: %v", tool, err)
		}
	}
}
