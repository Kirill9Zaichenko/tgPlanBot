package handlers

import (
	"context"

	"github.com/go-telegram/bot/models"

	userapp "tgPlanBot/internal/app/user"
	"tgPlanBot/internal/domain"
)

func SyncTelegramUser(
	ctx context.Context,
	userService *userapp.Service,
	from *models.User,
) (*domain.User, error) {
	return userService.SyncTelegramUser(ctx, userapp.SyncTelegramUserInput{
		TelegramID: from.ID,
		Username:   from.Username,
		FirstName:  from.FirstName,
		LastName:   from.LastName,
	})
}
