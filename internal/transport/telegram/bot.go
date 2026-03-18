package telegram

import (
	"context"
	"fmt"

	tgbot "github.com/go-telegram/bot"

	"tgPlanBot/internal/config"
)

type Bot struct {
	api *tgbot.Bot
}

func NewBot(cfg *config.Config) (*Bot, error) {
	api, err := tgbot.New(cfg.Telegram.Token)
	if err != nil {
		return nil, fmt.Errorf("create telegram bot: %w", err)
	}

	b := &Bot{
		api: api,
	}

	b.registerHandlers()

	return b, nil
}

func (b *Bot) Start(ctx context.Context) {
	b.api.Start(ctx)
}

func (b *Bot) registerHandlers() {
	b.api.RegisterHandler(tgbot.HandlerTypeMessageText, "/start", tgbot.MatchTypePrefix, b.handleStart)
	b.api.RegisterHandler(tgbot.HandlerTypeMessageText, "/help", tgbot.MatchTypePrefix, b.handleHelp)
}
