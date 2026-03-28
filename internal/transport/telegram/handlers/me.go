package telegram

import (
	"context"
	"fmt"
	"log"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *Bot) handleMe(ctx context.Context, bot *tgbot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.From == nil {
		return
	}

	user := update.Message.From

	text := fmt.Sprintf(
		"👤 Твои данные:\n\n"+
			"ID: %d\n"+
			"Username: @%s\n"+
			"Имя: %s %s",
		user.ID,
		user.Username,
		user.FirstName,
		user.LastName,
	)

	_, err := bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})
	if err != nil {
		log.Printf("send /me response: %v", err)
	}
}
