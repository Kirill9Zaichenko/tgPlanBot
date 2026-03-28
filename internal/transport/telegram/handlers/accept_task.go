package telegram

import (
	"context"
	"log"
	"strconv"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *Bot) handleAcceptTask(ctx context.Context, bot *tgbot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.From == nil {
		return
	}

	parts := splitCommand(update.Message.Text)
	if len(parts) < 2 {
		_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Использование: /accept {task_id}",
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

	receiverUserID := update.Message.From.ID

	if err := b.moderationService.AcceptTask(ctx, taskID, receiverUserID); err != nil {
		_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Не удалось принять задачу: " + err.Error(),
		})
		log.Printf("accept task %d by user %d: %v", taskID, receiverUserID, err)
		return
	}

	_, err = bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Задача успешно принята.",
	})
	if err != nil {
		log.Printf("send /accept response: %v", err)
	}
}
