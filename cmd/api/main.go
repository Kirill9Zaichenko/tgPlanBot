package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	moderationapp "tgPlanBot/internal/app/moderation"
	taskapp "tgPlanBot/internal/app/task"
	"tgPlanBot/internal/config"
	"tgPlanBot/internal/db"
	sqliterepo "tgPlanBot/internal/repository/sqlite"
	httptransport "tgPlanBot/internal/transport/http"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.NewConfig()

	if err := db.RunMigrations(cfg.Database.SQLitePath, cfg.Database.MigrationsPath); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	sqliteDB, err := db.NewSQLite(cfg.Database.SQLitePath)
	if err != nil {
		log.Fatalf("failed to connect sqlite: %v", err)
	}
	defer closeDB(sqliteDB)

	taskRepo := sqliterepo.NewTaskRepository(sqliteDB)
	taskRequestRepo := sqliterepo.NewTaskRequestRepository(sqliteDB)

	taskService := taskapp.NewService(sqliteDB, taskRepo, taskRequestRepo)
	moderationService := moderationapp.NewService(sqliteDB, taskRepo, taskRequestRepo)

	server := httptransport.NewServer(cfg, taskService, moderationService)

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

func closeDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("failed to close db: %v", err)
	}
}
