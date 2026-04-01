package handlers

import (
	"context"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
	"tgPlanBot/internal/domain"
	"tgPlanBot/internal/transport/telegram/messages"
)

type MeHandler struct{}

func NewMeHandler() *MeHandler {
	return &MeHandler{}
}

func (h *MeHandler) Handle(
	ctx context.Context,
	bot *tgbot.Bot,
	update *models.Update,
	user *domain.User,
) {
	_ = ctx

	if update.Message == nil {
		return
	}

	_, err := bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: messages.Me(
			user.ID,
			user.TelegramID,
			user.Username,
			user.FirstName,
			user.LastName,
		),
	})
	if err != nil {
		log.Printf("send /me response: %v", err)
	}
}
