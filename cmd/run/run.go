package cmd

import (
	"log"
	"os"
	"os/exec"
)

// En: This function runs app
// vi: Hàm chạy chương trình
func WitchesRun() {
	//En: Start executing
	//Vi: Bắt đầu thực thi
	cmd := exec.Command("go", "run", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
