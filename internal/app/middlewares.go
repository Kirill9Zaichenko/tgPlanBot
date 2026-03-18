package app

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
)

func LogIncoming() bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			if update.Message != nil && update.Message.From != nil {
				log.Printf("[chat:%d] %d (%s %s): %s",
					update.Message.Chat.ID,
					update.Message.From.ID,
					update.Message.From.FirstName,
					update.Message.From.LastName,
					update.Message.Text,
				)
			}
			next(ctx, b, update)
		}
	}
}
