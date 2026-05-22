package cmd

import (
	"log"
	"os"
	"os/exec"
)

func WitchesDatabaseUp(DB_PROFILE string) {
	cmd := exec.Command("make", "database-up-docker", "DB_PROFILE", DB_PROFILE)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func WitchesDatabaseDown() {
	cmd := exec.Command("make", "database-down-docker")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
