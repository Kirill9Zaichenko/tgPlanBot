package handlers

import (
	"context"
	"log"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"tgPlanBot/internal/domain"
	tgstate "tgPlanBot/internal/transport/telegram/state"
)

type NewTaskForHandler struct {
	stateStore *tgstate.Store
}

func NewNewTaskForHandler(stateStore *tgstate.Store) *NewTaskForHandler {
	return &NewTaskForHandler{stateStore: stateStore}
}

func (h *NewTaskForHandler) Handle(
	ctx context.Context,
	bot *tgbot.Bot,
	update *models.Update,
	user *domain.User,
) {
	if update.Message == nil {
		return
	}

	h.stateStore.Set(user.ID, tgstate.NewTaskState{
		UserID: user.ID,
		Flow:   tgstate.FlowNewTaskFor,
		Step:   tgstate.StepWaitingAssigneeTelegramID,
	})

	_, err := bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Введи telegram_id получателя задачи.",
	})
	if err != nil {
		log.Printf("send /newtaskfor response: %v", err)
	}
}
