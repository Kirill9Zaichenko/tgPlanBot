package handlers

import (
	"context"
	"log"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"tgPlanBot/internal/domain"
	"tgPlanBot/internal/transport/telegram/messages"
)

type StartHandler struct{}

func NewStartHandler() *StartHandler {
	return &StartHandler{}
}

func (h *StartHandler) Handle(
	ctx context.Context,
	bot *tgbot.Bot,
	update *models.Update,
	user *domain.User,
) {
	_ = user

	if update.Message == nil {
		return
	}

	_, err := bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   messages.Start(),
	})
	if err != nil {
		log.Printf("send /start response: %v", err)
	}
}
