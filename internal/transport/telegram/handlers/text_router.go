package handlers

import (
	"context"
	"log"
	"strings"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	taskapp "tgPlanBot/internal/app/task"
	tgstate "tgPlanBot/internal/transport/telegram/state"
)

type TextRouterHandler struct {
	taskService *taskapp.Service
	stateStore  *tgstate.Store
}

func NewTextRouterHandler(
	taskService *taskapp.Service,
	stateStore *tgstate.Store,
) *TextRouterHandler {
	return &TextRouterHandler{
		taskService: taskService,
		stateStore:  stateStore,
	}
}

func (h *TextRouterHandler) Handle(ctx context.Context, bot *tgbot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.From == nil {
		return
	}

	text := strings.TrimSpace(update.Message.Text)
	if text == "" {
		return
	}

	if strings.HasPrefix(text, "/") {
		return
	}

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	currentState, ok := h.stateStore.Get(userID)
	if !ok || currentState.Step == tgstate.StepIdle {
		return
	}

	switch currentState.Step {
	case tgstate.StepWaitingTaskTitle:
		currentState.Title = text
		currentState.Step = tgstate.StepWaitingDescription
		h.stateStore.Set(userID, currentState)

		_, err := bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: chatID,
			Text:   "Теперь введи описание задачи.",
		})
		if err != nil {
			log.Printf("send waiting description message: %v", err)
		}
		return

	case tgstate.StepWaitingDescription:
		currentState.Description = text

		task, err := h.taskService.Create(ctx, taskapp.CreateTaskInput{
			Title:          currentState.Title,
			Description:    currentState.Description,
			CreatorUserID:  userID,
			AssigneeUserID: userID,
			DueAt:          nil,
		})
		if err != nil {
			_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
				ChatID: chatID,
				Text:   "Не удалось создать задачу: " + err.Error(),
			})
			log.Printf("create task from telegram: %v", err)
			return
		}

		h.stateStore.Delete(userID)

		_, err = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: chatID,
			Text:   "Задача создана.\n\nID: #" + int64ToString(task.ID) + "\nНазвание: " + task.Title,
		})
		if err != nil {
			log.Printf("send task created message: %v", err)
		}
		return
	}
}
