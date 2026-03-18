package handlers

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"tgBotPlan/internal/storage"
)

func DefaultHandler(store storage.TaskStore) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		FallbackHandler(store)(ctx, b, update)
	}
}
