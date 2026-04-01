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

type RejectHandler struct {
	moderationService *moderationapp.Service
}

func NewRejectHandler(moderationService *moderationapp.Service) *RejectHandler {
	return &RejectHandler{moderationService: moderationService}
}

func (h *RejectHandler) Handle(
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
			Text:   messages.UsageReject(),
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

	comment := "rejected"
	if len(parts) > 2 {
		comment = strings.TrimSpace(strings.Join(parts[2:], " "))
		if comment == "" {
			comment = "rejected"
		}
	}

	if err := h.moderationService.RejectTask(ctx, taskID, user.ID, comment); err != nil {
		_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   messages.RejectFailed(err),
		})
		log.Printf("reject task %d by user %d: %v", taskID, user.ID, err)
		return
	}

	_, err = bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   messages.TaskRejected(),
	})
	if err != nil {
		log.Printf("send /reject response: %v", err)
	}
}
