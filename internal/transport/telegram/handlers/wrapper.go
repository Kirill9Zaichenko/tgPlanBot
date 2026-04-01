package handlers

import (
	"context"
	"log"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	userapp "tgPlanBot/internal/app/user"
	"tgPlanBot/internal/domain"
)

type MessageHandlerFunc func(
	ctx context.Context,
	bot *tgbot.Bot,
	update *models.Update,
	user *domain.User,
)

type CallbackHandlerFunc func(
	ctx context.Context,
	bot *tgbot.Bot,
	update *models.Update,
	user *domain.User,
)

func WithSyncedMessageUser(
	userService *userapp.Service,
	next MessageHandlerFunc,
) func(ctx context.Context, bot *tgbot.Bot, update *models.Update) {
	return func(ctx context.Context, bot *tgbot.Bot, update *models.Update) {
		if update.Message == nil || update.Message.From == nil {
			return
		}

		user, err := SyncTelegramUser(ctx, userService, update.Message.From)
		if err != nil {
			log.Printf("sync telegram message user: %v", err)
			return
		}

		next(ctx, bot, update, user)
	}
}

func WithSyncedCallbackUser(
	userService *userapp.Service,
	next CallbackHandlerFunc,
) func(ctx context.Context, bot *tgbot.Bot, update *models.Update) {
	return func(ctx context.Context, bot *tgbot.Bot, update *models.Update) {
		if update.CallbackQuery == nil {
			return
		}

		user, err := SyncTelegramUser(ctx, userService, &update.CallbackQuery.From)
		if err != nil {
			log.Printf("sync telegram callback user: %v", err)
			return
		}

		next(ctx, bot, update, user)
	}
}
