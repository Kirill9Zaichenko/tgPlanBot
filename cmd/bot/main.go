package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	moderationapp "tgPlanBot/internal/app/moderation"
	taskapp "tgPlanBot/internal/app/task"
	"tgPlanBot/internal/config"
	"tgPlanBot/internal/db"
	sqliterepo "tgPlanBot/internal/repository/sqlite"
	telegramtransport "tgPlanBot/internal/transport/telegram"
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

	botApp, err := telegramtransport.NewBot(cfg, taskService, moderationService)
	if err != nil {
		log.Fatalf("failed to initialize bot: %v", err)
	}

	log.Printf("bot started in %s mode", cfg.App.Env)
	botApp.Start(ctx)
}

func closeDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("failed to close db: %v", err)
	}
}
