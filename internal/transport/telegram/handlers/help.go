package handlers

import (
	"context"
	"log"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"tgPlanBot/internal/domain"
	"tgPlanBot/internal/transport/telegram/messages"
)

type HelpHandler struct{}

func NewHelpHandler() *HelpHandler {
	return &HelpHandler{}
}

func (h *HelpHandler) Handle(
	ctx context.Context,
	bot *tgbot.Bot,
	update *models.Update,
	user *domain.User,
) {
	_ = user // пока не используем

	if update.Message == nil {
		return
	}

	_, err := bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   messages.Help(),
	})
	if err != nil {
		log.Printf("send /help response: %v", err)
	}
}
