package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"tgPlanBot/internal/app"
	"tgPlanBot/internal/config"
	telegramtransport "tgPlanBot/internal/transport/telegram"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.NewConfig()

	appContainer, err := app.New(cfg)
	if err != nil {
		log.Fatalf("init app: %v", err)
	}
	defer appContainer.DB.Close()

	botApp, err := telegramtransport.NewBot(
		cfg,
		appContainer.TaskService,
		appContainer.ModerationService,
		appContainer.UserService,
		appContainer.OrganizationService,
	)
	if err != nil {
		log.Fatalf("init bot: %v", err)
	}

	log.Println("bot started")
	botApp.Start(ctx)
}
