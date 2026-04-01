package user

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"tgPlanBot/internal/domain"
	repoif "tgPlanBot/internal/repository/interfaces"
)

type SyncTelegramUserInput struct {
	TelegramID int64
	Username   string
	FirstName  string
	LastName   string
}

type Service struct {
	userRepo repoif.UserRepository
}

func NewService(userRepo repoif.UserRepository) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

func (s *Service) SyncTelegramUser(ctx context.Context, input SyncTelegramUserInput) (*domain.User, error) {
	if input.TelegramID <= 0 {
		return nil, fmt.Errorf("telegram id is required")
	}

	input.Username = strings.TrimSpace(input.Username)
	input.FirstName = strings.TrimSpace(input.FirstName)
	input.LastName = strings.TrimSpace(input.LastName)

	existingUser, err := s.userRepo.GetByTelegramID(ctx, input.TelegramID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("get user by telegram id: %w", err)
	}

	if err == sql.ErrNoRows {
		user := &domain.User{
			TelegramID: input.TelegramID,
			Username:   input.Username,
			FirstName:  input.FirstName,
			LastName:   input.LastName,
		}

		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("create user: %w", err)
		}

		return user, nil
	}

	existingUser.Username = input.Username
	existingUser.FirstName = input.FirstName
	existingUser.LastName = input.LastName

	if err := s.userRepo.UpdateTelegramProfile(ctx, existingUser); err != nil {
		return nil, fmt.Errorf("update telegram profile: %w", err)
	}

	return existingUser, nil
}

func (s *Service) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	if telegramID <= 0 {
		return nil, fmt.Errorf("telegram id is required")
	}

	user, err := s.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, fmt.Errorf("get user by telegram id: %w", err)
	}

	return user, nil
}
