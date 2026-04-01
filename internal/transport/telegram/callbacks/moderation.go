package callbacks

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	moderationapp "tgPlanBot/internal/app/moderation"
	"tgPlanBot/internal/domain"
)

type ModerationHandler struct {
	moderationService *moderationapp.Service
}

func NewModerationHandler(moderationService *moderationapp.Service) *ModerationHandler {
	return &ModerationHandler{moderationService: moderationService}
}

func (h *ModerationHandler) Handle(
	ctx context.Context,
	bot *tgbot.Bot,
	update *models.Update,
	user *domain.User,
) {
	if update.CallbackQuery == nil {
		return
	}

	data := strings.TrimSpace(update.CallbackQuery.Data)
	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		h.answerCallback(ctx, bot, update, "Некорректное действие")
		return
	}

	action := parts[0]
	taskID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || taskID <= 0 {
		h.answerCallback(ctx, bot, update, "Некорректный task_id")
		return
	}

	switch action {
	case "accept":
		if err := h.moderationService.AcceptTask(ctx, taskID, user.ID); err != nil {
			log.Printf("accept from callback task %d by user %d: %v", taskID, user.ID, err)
			h.answerCallback(ctx, bot, update, "Не удалось принять: "+err.Error())
			return
		}

		h.answerCallback(ctx, bot, update, "Задача принята")
		h.sendStatusMessage(ctx, bot, update, fmt.Sprintf("Task #%d\nСтатус: accepted", taskID))

	case "reject":
		if err := h.moderationService.RejectTask(ctx, taskID, user.ID, "rejected from telegram"); err != nil {
			log.Printf("reject from callback task %d by user %d: %v", taskID, user.ID, err)
			h.answerCallback(ctx, bot, update, "Не удалось отклонить: "+err.Error())
			return
		}

		h.answerCallback(ctx, bot, update, "Задача отклонена")
		h.sendStatusMessage(ctx, bot, update, fmt.Sprintf("Task #%d\nСтатус: rejected", taskID))

	default:
		h.answerCallback(ctx, bot, update, "Неизвестное действие")
	}
}

func (h *ModerationHandler) answerCallback(ctx context.Context, bot *tgbot.Bot, update *models.Update, text string) {
	_, err := bot.AnswerCallbackQuery(ctx, &tgbot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		Text:            text,
		ShowAlert:       false,
	})
	if err != nil {
		log.Printf("answer callback query: %v", err)
	}
}

func (h *ModerationHandler) sendStatusMessage(ctx context.Context, bot *tgbot.Bot, update *models.Update, text string) {
	chatID := update.CallbackQuery.From.ID

	_, err := bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		log.Printf("send callback status message: %v", err)
	}
}
