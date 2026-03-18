package app

import (
	"fmt"
	"tgBotPlan/internal/config"
	"tgBotPlan/internal/handlers"
	"tgBotPlan/internal/storage"

	"github.com/go-telegram/bot"
)

func NewBot(cfg *config.Config, store storage.TaskStore) (*bot.Bot, error) {
	if cfg.Telegram.Token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is empty")
	}

	opts := []bot.Option{
		bot.WithMiddlewares(LogIncoming()),
		bot.WithDefaultHandler(handlers.DefaultHandler(store)),
	}

	opts = append(opts, Routes(store)...)

	return bot.New(cfg.Telegram.Token, opts...)
}
