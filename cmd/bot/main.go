package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"tgPlanBot/internal/config"
	telegramtransport "tgPlanBot/internal/transport/telegram"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.NewConfig()

	botApp, err := telegramtransport.NewBot(cfg)
	if err != nil {
		log.Fatalf("failed to initialize bot: %v", err)
	}

	log.Printf("bot started in %s mode", cfg.App.Env)

	botApp.Start(ctx)
}
