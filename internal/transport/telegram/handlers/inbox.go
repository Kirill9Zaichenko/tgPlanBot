package telegram

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot/models"
	"log"
	"strings"

	tgbot "github.com/go-telegram/bot"
)

func (b *Bot) handleInbox(ctx context.Context, bot *tgbot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.From == nil {
		return
	}

	userID := update.Message.From.ID

	items, err := b.moderationService.ListInbox(ctx, userID)
	if err != nil {
		_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Не удалось загрузить входящие запросы.",
		})
		log.Printf("list inbox: %v", err)
		return
	}

	if len(items) == 0 {
		_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Входящих запросов нет.",
		})
		return
	}

	var sb strings.Builder
	sb.WriteString("Входящие запросы:\n\n")

	for _, item := range items {
		sb.WriteString(fmt.Sprintf(
			"Task #%d\nОтправитель: %d\nСтатус: %s\n\n",
			item.TaskID,
			item.SenderUserID,
			item.Status,
		))
	}

	_, err = bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   sb.String(),
	})
	if err != nil {
		log.Printf("send /inbox response: %v", err)
	}
}
