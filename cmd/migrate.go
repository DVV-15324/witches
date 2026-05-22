package cmd

import (
	"log"
	"os"
	"os/exec"
)

func WitchesMigrateUp(DB_URL string) {
	cmd := exec.Command("make", "migrate-up-docker", "DB_URL", DB_URL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func WitchesMigrateDown(DB_URL string) {
	cmd := exec.Command("make", "migrate-down-docker", "DB_URL", DB_URL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func WitchesMigrateVersion(DB_URL string) {
	cmd := exec.Command("make", "migrate-version-docker", "DB_URL", DB_URL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
func WitchesMigrateForce(DB_URL string) {
	cmd := exec.Command("make", "migrate-force-docker", "DB_URL", DB_URL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
