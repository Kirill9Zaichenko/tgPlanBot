package interfaces

import (
	"context"

	"tgPlanBot/internal/domain"
)

type UserRepository interface {
	GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	UpdateTelegramProfile(ctx context.Context, user *domain.User) error
}
