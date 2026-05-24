package cmd

// INFO: Phát Triển sau

// import (
// 	"log"
// 	"time"
// )

// func StartDB(done chan bool) {
// 	log.Println("Starting DB...")

// 	time.Sleep(5 * time.Second)

// 	log.Println("DB Ready")

// 	done <- true
// }

// func RunMigrate(done chan bool) {
// 	log.Println("Running migrate...")

// 	time.Sleep(2 * time.Second)

// 	log.Println("Migrate Done")

// 	done <- true
// }

// func StartApp() {
// 	log.Println("App Started")
// }

// func RunPipeline() {

// 	dbReady := make(chan bool)
// 	migrateDone := make(chan bool)

// 	go StartDB(dbReady)

// 	<-dbReady

// 	go RunMigrate(migrateDone)

// 	<-migrateDone

// 	StartApp()
// }
