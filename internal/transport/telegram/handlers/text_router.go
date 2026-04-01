package handlers

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	taskapp "tgPlanBot/internal/app/task"
	userapp "tgPlanBot/internal/app/user"
	"tgPlanBot/internal/domain"
	tgstate "tgPlanBot/internal/transport/telegram/state"
)

type TextRouterHandler struct {
	taskService *taskapp.Service
	userService *userapp.Service
	stateStore  *tgstate.Store
}

func NewTextRouterHandler(
	taskService *taskapp.Service,
	userService *userapp.Service,
	stateStore *tgstate.Store,
) *TextRouterHandler {
	return &TextRouterHandler{
		taskService: taskService,
		userService: userService,
		stateStore:  stateStore,
	}
}

func (h *TextRouterHandler) Handle(
	ctx context.Context,
	bot *tgbot.Bot,
	update *models.Update,
	user *domain.User,
) {
	if update.Message == nil {
		return
	}

	text := strings.TrimSpace(update.Message.Text)
	if text == "" {
		return
	}

	if strings.HasPrefix(text, "/") {
		return
	}

	chatID := update.Message.Chat.ID

	currentState, ok := h.stateStore.Get(user.ID)
	if !ok || currentState.Step == tgstate.StepIdle {
		return
	}

	switch currentState.Step {
	case tgstate.StepWaitingAssigneeTelegramID:
		telegramID, err := strconv.ParseInt(text, 10, 64)
		if err != nil || telegramID <= 0 {
			_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
				ChatID: chatID,
				Text:   "Некорректный telegram_id. Введи число.",
			})
			return
		}

		assignee, err := h.userService.GetByTelegramID(ctx, telegramID)
		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
					ChatID: chatID,
					Text:   "Пользователь не найден. Пусть сначала напишет боту /start или /me.",
				})
				return
			}

			_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
				ChatID: chatID,
				Text:   "Не удалось найти пользователя.",
			})
			log.Printf("find assignee by telegram id: %v", err)
			return
		}

		currentState.AssigneeTelegramID = telegramID
		currentState.AssigneeUserID = assignee.ID
		currentState.Step = tgstate.StepWaitingTaskTitle
		h.stateStore.Set(user.ID, currentState)

		_, err = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: chatID,
			Text:   "Получатель найден. Теперь введи название задачи.",
		})
		if err != nil {
			log.Printf("send title prompt after assignee: %v", err)
		}
		return

	case tgstate.StepWaitingTaskTitle:
		currentState.Title = text
		currentState.Step = tgstate.StepWaitingDescription
		h.stateStore.Set(user.ID, currentState)

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

		assigneeUserID := currentState.AssigneeUserID
		if assigneeUserID == 0 {
			assigneeUserID = user.ID
		}

		task, err := h.taskService.Create(ctx, taskapp.CreateTaskInput{
			Title:          currentState.Title,
			Description:    currentState.Description,
			CreatorUserID:  user.ID,
			AssigneeUserID: assigneeUserID,
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

		h.stateStore.Delete(user.ID)

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
