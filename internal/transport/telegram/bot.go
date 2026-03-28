package telegram

import (
	"context"
	"fmt"

	tgbot "github.com/go-telegram/bot"

	moderationapp "tgPlanBot/internal/app/moderation"
	taskapp "tgPlanBot/internal/app/task"
	"tgPlanBot/internal/config"
	tgcallbacks "tgPlanBot/internal/transport/telegram/callbacks"
	tghandlers "tgPlanBot/internal/transport/telegram/handlers"
	tgstate "tgPlanBot/internal/transport/telegram/state"
)

type Bot struct {
	api *tgbot.Bot
}

func NewBot(
	cfg *config.Config,
	taskService *taskapp.Service,
	moderationService *moderationapp.Service,
) (*Bot, error) {
	api, err := tgbot.New(cfg.Telegram.Token)
	if err != nil {
		return nil, fmt.Errorf("create telegram bot: %w", err)
	}

	b := &Bot{
		api: api,
	}

	b.registerHandlers(taskService, moderationService)

	return b, nil
}

func (b *Bot) Start(ctx context.Context) {
	b.api.Start(ctx)
}

func (b *Bot) registerHandlers(
	taskService *taskapp.Service,
	moderationService *moderationapp.Service,
) {
	stateStore := tgstate.NewStore()

	startHandler := tghandlers.NewStartHandler()
	helpHandler := tghandlers.NewHelpHandler()
	meHandler := tghandlers.NewMeHandler()
	myTasksHandler := tghandlers.NewMyTasksHandler(taskService)
	inboxHandler := tghandlers.NewInboxHandler(moderationService)
	acceptHandler := tghandlers.NewAcceptHandler(moderationService)
	rejectHandler := tghandlers.NewRejectHandler(moderationService)
	newTaskHandler := tghandlers.NewNewTaskHandler(stateStore)
	textRouterHandler := tghandlers.NewTextRouterHandler(taskService, stateStore)

	callbackHandler := tgcallbacks.NewModerationHandler(moderationService)

	b.api.RegisterHandler(tgbot.HandlerTypeMessageText, "/start", tgbot.MatchTypeExact, startHandler.Handle)
	b.api.RegisterHandler(tgbot.HandlerTypeMessageText, "/help", tgbot.MatchTypeExact, helpHandler.Handle)
	b.api.RegisterHandler(tgbot.HandlerTypeMessageText, "/me", tgbot.MatchTypeExact, meHandler.Handle)
	b.api.RegisterHandler(tgbot.HandlerTypeMessageText, "/mytasks", tgbot.MatchTypeExact, myTasksHandler.Handle)
	b.api.RegisterHandler(tgbot.HandlerTypeMessageText, "/inbox", tgbot.MatchTypeExact, inboxHandler.Handle)
	b.api.RegisterHandler(tgbot.HandlerTypeMessageText, "/newtask", tgbot.MatchTypeExact, newTaskHandler.Handle)
	b.api.RegisterHandler(tgbot.HandlerTypeMessageText, "/accept", tgbot.MatchTypePrefix, acceptHandler.Handle)
	b.api.RegisterHandler(tgbot.HandlerTypeMessageText, "/reject", tgbot.MatchTypePrefix, rejectHandler.Handle)

	b.api.RegisterHandler(tgbot.HandlerTypeMessageText, "", tgbot.MatchTypePrefix, textRouterHandler.Handle)

	b.api.RegisterHandler(tgbot.HandlerTypeCallbackQueryData, "", tgbot.MatchTypePrefix, callbackHandler.Handle)
}
