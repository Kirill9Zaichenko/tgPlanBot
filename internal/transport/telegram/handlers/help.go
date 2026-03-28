package telegram

import (
	"context"
	"log"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *Bot) handleHelp(ctx context.Context, bot *tgbot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	text := "Доступные команды:\n\n" +
		"/start — запустить бота\n" +
		"/help — показать помощь\n" +
		"/me — показать мой Telegram ID\n" +
		"/mytasks — показать мои задачи\n" +
		"/inbox — показать входящие запросы\n" +
		"/accept {task_id} — принять задачу\n" +
		"/reject {task_id} {comment} — отклонить задачу"

	_, err := bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})
	if err != nil {
		log.Printf("send /help response: %v", err)
	}
}
