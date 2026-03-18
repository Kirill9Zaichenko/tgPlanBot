package telegram

import (
	"context"
	"log"

	tgbot "github.com/go-telegram/bot"
)

func (b *Bot) handleStart(ctx context.Context, bot *tgbot.Bot, update *tgbot.Update) {
	if update.Message == nil {
		return
	}

	_, err := bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Привет! Я бот-планировщик задач.\n\nСкоро здесь можно будет создавать задачи, отправлять их другим пользователям и модерировать входящие.",
	})
	if err != nil {
		log.Printf("send start message: %v", err)
	}
}
