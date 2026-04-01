package app

import (
	"database/sql"
	"fmt"

	moderationapp "tgPlanBot/internal/app/moderation"
	taskapp "tgPlanBot/internal/app/task"
	userapp "tgPlanBot/internal/app/user"
	"tgPlanBot/internal/config"
	"tgPlanBot/internal/db"
	sqliterepo "tgPlanBot/internal/repository/sqlite"
)

type App struct {
	DB                *sql.DB
	TaskService       *taskapp.Service
	ModerationService *moderationapp.Service
	UserService       *userapp.Service
}

func New(cfg *config.Config) (*App, error) {
	if err := db.RunMigrations(cfg.Database.SQLitePath, cfg.Database.MigrationsPath); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	sqliteDB, err := db.NewSQLite(cfg.Database.SQLitePath)
	if err != nil {
		return nil, fmt.Errorf("connect sqlite: %w", err)
	}

	taskRepo := sqliterepo.NewTaskRepository(sqliteDB)
	taskRequestRepo := sqliterepo.NewTaskRequestRepository(sqliteDB)
	userRepo := sqliterepo.NewUserRepository(sqliteDB)

	taskService := taskapp.NewService(sqliteDB, taskRepo, taskRequestRepo)
	moderationService := moderationapp.NewService(sqliteDB, taskRepo, taskRequestRepo)
	userService := userapp.NewService(userRepo)

	return &App{
		DB:                sqliteDB,
		TaskService:       taskService,
		ModerationService: moderationService,
		UserService:       userService,
	}, nil
}
