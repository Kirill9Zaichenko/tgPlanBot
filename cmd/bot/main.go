package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tgBotPlan/internal/app"
	"tgBotPlan/internal/config"
	"tgBotPlan/internal/repository/memory"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	cfg := config.NewConfig()
	store := memory.NewTaskStore()

	b, err := app.NewBot(cfg, store)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("bot started")

	b.Start(ctx)

	log.Println("bot stopped")
}
