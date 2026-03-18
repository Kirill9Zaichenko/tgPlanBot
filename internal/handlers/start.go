package handlers

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"tgBotPlan/internal/storage"
	"tgBotPlan/internal/telegram"
	"tgBotPlan/internal/util"
)

func Start(store storage.TaskStore) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		telegram.Reply(ctx, b, update.Message.Chat.ID, util.HelpText())
	}
}
