package handlers

import (
	"context"
	"log"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	taskapp "tgPlanBot/internal/app/task"
	"tgPlanBot/internal/transport/telegram/messages"
)

type MyTasksHandler struct {
	taskService *taskapp.Service
}

func NewMyTasksHandler(taskService *taskapp.Service) *MyTasksHandler {
	return &MyTasksHandler{taskService: taskService}
}

func (h *MyTasksHandler) Handle(ctx context.Context, bot *tgbot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.From == nil {
		return
	}

	userID := update.Message.From.ID

	tasks, err := h.taskService.ListByAssignee(ctx, userID)
	if err != nil {
		_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   messages.TasksLoadFailed(),
		})
		log.Printf("list my tasks: %v", err)
		return
	}

	_, err = bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   messages.TasksList(tasks),
	})
	if err != nil {
		log.Printf("send /mytasks response: %v", err)
	}
}
