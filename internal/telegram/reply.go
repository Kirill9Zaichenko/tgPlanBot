package telegram

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
)

func Reply(ctx context.Context, b *bot.Bot, chatID int64, text string) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		log.Printf("send message error: %v", err)
	}
}
