package handlers

import (
	"context"
	"log"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"tgPlanBot/internal/transport/telegram/messages"
)

type MeHandler struct{}

func NewMeHandler() *MeHandler {
	return &MeHandler{}
}

func (h *MeHandler) Handle(ctx context.Context, bot *tgbot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.From == nil {
		return
	}

	user := update.Message.From

	_, err := bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: messages.Me(
			user.ID,
			user.Username,
			user.FirstName,
			user.LastName,
		),
	})
	if err != nil {
		log.Printf("send /me response: %v", err)
	}
}
