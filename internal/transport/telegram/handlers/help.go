package telegram

import (
	"context"
	"log"

	tgbot "github.com/go-telegram/bot"
)

func (b *Bot) handleHelp(ctx context.Context, bot *tgbot.Bot, update *tgbot.Update) {
	if update.Message == nil {
		return
	}

	_, err := bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: "/start — запустить бота\n" +
			"/help — показать помощь",
	})
	if err != nil {
		log.Printf("send help message: %v", err)
	}
}
