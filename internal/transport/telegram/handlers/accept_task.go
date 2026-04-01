package handlers

import (
	"context"
	"log"
	"strconv"
	"strings"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	moderationapp "tgPlanBot/internal/app/moderation"
	"tgPlanBot/internal/domain"
	"tgPlanBot/internal/transport/telegram/messages"
)

type AcceptHandler struct {
	moderationService *moderationapp.Service
}

func NewAcceptHandler(moderationService *moderationapp.Service) *AcceptHandler {
	return &AcceptHandler{moderationService: moderationService}
}

func (h *AcceptHandler) Handle(
	ctx context.Context,
	bot *tgbot.Bot,
	update *models.Update,
	user *domain.User,
) {
	if update.Message == nil {
		return
	}

	parts := strings.Fields(strings.TrimSpace(update.Message.Text))
	if len(parts) < 2 {
		_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   messages.UsageAccept(),
		})
		return
	}

	taskID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || taskID <= 0 {
		_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   messages.InvalidTaskID(),
		})
		return
	}

	if err := h.moderationService.AcceptTask(ctx, taskID, user.ID); err != nil {
		_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   messages.AcceptFailed(err),
		})
		log.Printf("accept task %d by user %d: %v", taskID, user.ID, err)
		return
	}

	_, err = bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   messages.TaskAccepted(),
	})
	if err != nil {
		log.Printf("send /accept response: %v", err)
	}
}
