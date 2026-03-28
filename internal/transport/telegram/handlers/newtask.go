package handlers

import (
	"context"
	"log"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	tgstate "tgPlanBot/internal/transport/telegram/state"
)

type NewTaskHandler struct {
	stateStore *tgstate.Store
}

func NewNewTaskHandler(stateStore *tgstate.Store) *NewTaskHandler {
	return &NewTaskHandler{stateStore: stateStore}
}

func (h *NewTaskHandler) Handle(ctx context.Context, bot *tgbot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.From == nil {
		return
	}

	userID := update.Message.From.ID

	h.stateStore.Set(userID, tgstate.NewTaskState{
		UserID: userID,
		Step:   tgstate.StepWaitingTaskTitle,
	})

	_, err := bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Введи название задачи.",
	})
	if err != nil {
		log.Printf("send /newtask response: %v", err)
	}
}
