package utils

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func ShutdownServer(ctx context.Context, engine http.Handler, address string, port string) {
	// Tạo addr = address:port
	addr := fmt.Sprintf("%s:%s", address, port)

	// Create HTTP server
	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	// Run server in goroutine
	go func() {
		log.Printf("Server running on http://%s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")
	// Graceful shutdown with timeout
	_ = ctx
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
