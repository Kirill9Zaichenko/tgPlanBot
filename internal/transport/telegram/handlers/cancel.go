package handlers

import (
	"context"
	"log"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"tgPlanBot/internal/domain"
	"tgPlanBot/internal/transport/telegram/messages"
	tgstate "tgPlanBot/internal/transport/telegram/state"
)

type CancelHandler struct {
	stateStore *tgstate.Store
}

func NewCancelHandler(stateStore *tgstate.Store) *CancelHandler {
	return &CancelHandler{stateStore: stateStore}
}

func (h *CancelHandler) Handle(
	ctx context.Context,
	bot *tgbot.Bot,
	update *models.Update,
	user *domain.User,
) {
	if update.Message == nil {
		return
	}

	currentState, ok := h.stateStore.Get(user.ID)
	if !ok || currentState.Step == tgstate.StepIdle {
		_, err := bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   messages.NoActiveFlow(),
		})
		if err != nil {
			log.Printf("send /cancel no active flow response: %v", err)
		}
		return
	}

	h.stateStore.Delete(user.ID)

	_, err := bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   messages.FlowCancelled(),
	})
	if err != nil {
		log.Printf("send /cancel response: %v", err)
	}
}
