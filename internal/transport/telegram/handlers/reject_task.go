package telegram

import (
	"context"
	"log"
	"strconv"
	"strings"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *Bot) handleRejectTask(ctx context.Context, bot *tgbot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.From == nil {
		return
	}

	parts := splitCommand(update.Message.Text)
	if len(parts) < 2 {
		_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Использование: /reject {task_id} {comment}",
		})
		return
	}

	taskID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || taskID <= 0 {
		_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Некорректный task_id.",
		})
		return
	}

	comment := "rejected"
	if len(parts) > 2 {
		comment = strings.TrimSpace(strings.Join(parts[2:], " "))
		if comment == "" {
			comment = "rejected"
		}
	}

	receiverUserID := update.Message.From.ID

	if err := b.moderationService.RejectTask(ctx, taskID, receiverUserID, comment); err != nil {
		_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Не удалось отклонить задачу: " + err.Error(),
		})
		log.Printf("reject task %d by user %d: %v", taskID, receiverUserID, err)
		return
	}

	_, err = bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Задача отклонена.",
	})
	if err != nil {
		log.Printf("send /reject response: %v", err)
	}
}
