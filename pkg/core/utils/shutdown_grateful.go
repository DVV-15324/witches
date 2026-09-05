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

var notifySignal = signal.Notify
var shutdownServerFunc = func(server *http.Server, ctx context.Context) error {
	return server.Shutdown(ctx)
}

func ShutdownServer(ctx context.Context, engine http.Handler, host string, port string) {
	addr := fmt.Sprintf("%s:%s", host, port)

	server := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	go func() {
		log.Printf("Server running on http://%s\n", addr)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	notifySignal(quit, os.Interrupt)
	<-quit

	shutdownServer(server)
}

func shutdownServer(server *http.Server) {
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		return
	}

	log.Println("Server exiting")
}
