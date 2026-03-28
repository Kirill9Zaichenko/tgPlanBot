package telegram

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot/models"
	"log"
	"strings"

	tgbot "github.com/go-telegram/bot"
)

func (b *Bot) handleMyTasks(ctx context.Context, bot *tgbot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.From == nil {
		return
	}

	userID := update.Message.From.ID

	tasks, err := b.taskService.ListByAssignee(ctx, userID)
	if err != nil {
		_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Не удалось загрузить задачи.",
		})
		log.Printf("list my tasks: %v", err)
		return
	}

	if len(tasks) == 0 {
		_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "У тебя пока нет задач.",
		})
		return
	}

	var sb strings.Builder
	sb.WriteString("Твои задачи:\n\n")

	for _, task := range tasks {
		sb.WriteString(fmt.Sprintf(
			"#%d | %s\nСтатус: %s\n\n",
			task.ID,
			task.Title,
			task.Status,
		))
	}

	_, err = bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   sb.String(),
	})
	if err != nil {
		log.Printf("send /mytasks response: %v", err)
	}
}
