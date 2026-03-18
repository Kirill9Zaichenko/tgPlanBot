package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tgPlanBot/internal/config"
	httptransport "tgPlanBot/internal/transport/http"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.NewConfig()

	server := httptransport.NewServer(cfg)

	errCh := make(chan error, 1)

	go func() {
		log.Printf("http api started on %s", server.Address())
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")
	case err := <-errCh:
		log.Fatalf("http server failed: %v", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Stop(shutdownCtx); err != nil {
		log.Fatalf("failed to stop http server: %v", err)
	}

	log.Println("http api stopped")
}
