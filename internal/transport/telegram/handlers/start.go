package telegram

import (
	"context"
	"log"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *Bot) handleStart(ctx context.Context, bot *tgbot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	text := "Привет! Я бот-планировщик задач.\n\n" +
		"Команды:\n" +
		"/help — помощь\n" +
		"/me — мой Telegram ID\n" +
		"/mytasks — мои задачи\n" +
		"/inbox — входящие запросы\n" +
		"/accept {task_id} — принять задачу\n" +
		"/reject {task_id} {comment} — отклонить задачу"

	_, err := bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})
	if err != nil {
		log.Printf("send /start response: %v", err)
	}
}
