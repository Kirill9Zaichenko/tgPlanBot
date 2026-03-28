package handlers

import (
	"context"
	"fmt"
	"log"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	moderationapp "tgPlanBot/internal/app/moderation"
	"tgPlanBot/internal/transport/telegram/keyboards"
	"tgPlanBot/internal/transport/telegram/messages"
)

type InboxHandler struct {
	moderationService *moderationapp.Service
}

func NewInboxHandler(moderationService *moderationapp.Service) *InboxHandler {
	return &InboxHandler{moderationService: moderationService}
}

func (h *InboxHandler) Handle(ctx context.Context, bot *tgbot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.From == nil {
		return
	}

	userID := update.Message.From.ID

	items, err := h.moderationService.ListInbox(ctx, userID)
	if err != nil {
		_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   messages.InboxLoadFailed(),
		})
		log.Printf("list inbox: %v", err)
		return
	}

	if len(items) == 0 {
		_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   messages.NoInboxItems(),
		})
		return
	}

	for _, item := range items {
		text := fmt.Sprintf(
			"Task #%d\nОтправитель: %d\nСтатус: %s",
			item.TaskID,
			item.SenderUserID,
			item.Status,
		)

		_, err := bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: keyboards.InboxTaskActions(item.TaskID),
		})
		if err != nil {
			log.Printf("send inbox item %d: %v", item.TaskID, err)
		}
	}
}
