package handlers

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strings"
	"tgBotPlan/internal/storage"
	"tgBotPlan/internal/telegram"
	"tgBotPlan/internal/util"
)

func FallbackHandler(store storage.TaskStore) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update == nil || update.Message == nil || strings.TrimSpace(update.Message.Text) == "" {
			return
		}
		chatID := update.Message.Chat.ID
		telegram.Reply(ctx, b, chatID, "Не понял команду. \n\n"+util.HelpText())
	}
}
