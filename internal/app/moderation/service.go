package moderation

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"tgPlanBot/internal/db"
	"tgPlanBot/internal/domain"
	repoif "tgPlanBot/internal/repository/interfaces"
	sqliterepo "tgPlanBot/internal/repository/sqlite"
)

type Service struct {
	database        *sql.DB
	taskRepo        repoif.TaskRepository
	taskRequestRepo repoif.TaskRequestRepository
}

func NewService(
	database *sql.DB,
	taskRepo repoif.TaskRepository,
	taskRequestRepo repoif.TaskRequestRepository,
) *Service {
	return &Service{
		database:        database,
		taskRepo:        taskRepo,
		taskRequestRepo: taskRequestRepo,
	}
}

func (s *Service) ListInbox(ctx context.Context, receiverUserID int64) ([]domain.TaskRequest, error) {
	if receiverUserID <= 0 {
		return nil, fmt.Errorf("invalid receiver user id")
	}

	items, err := s.taskRequestRepo.ListPendingByReceiver(ctx, receiverUserID)
	if err != nil {
		return nil, fmt.Errorf("list inbox: %w", err)
	}

	return items, nil
}

func (s *Service) AcceptTask(ctx context.Context, taskID, receiverUserID int64) error {
	return db.WithTx(ctx, s.database, func(tx *sql.Tx) error {
		taskRepo := sqliterepo.NewTaskRepositoryTx(tx)
		taskRequestRepo := sqliterepo.NewTaskRequestRepositoryTx(tx)

		req, err := taskRequestRepo.GetByTaskID(ctx, taskID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("task request not found")
			}
			return fmt.Errorf("get task request: %w", err)
		}

		if req.ReceiverUserID != receiverUserID {
			return fmt.Errorf("you cannot accept another user's task")
		}

		if req.Status != domain.RequestStatusPending {
			return fmt.Errorf("task request is not pending")
		}

		if err := taskRequestRepo.UpdateDecision(ctx, taskID, domain.RequestStatusAccepted, "accepted"); err != nil {
			return fmt.Errorf("update task request decision: %w", err)
		}

		if err := taskRepo.UpdateStatus(ctx, taskID, domain.TaskStatusAccepted); err != nil {
			return fmt.Errorf("update task status: %w", err)
		}

		return nil
	})
}

func (s *Service) RejectTask(ctx context.Context, taskID, receiverUserID int64, comment string) error {
	return db.WithTx(ctx, s.database, func(tx *sql.Tx) error {
		taskRepo := sqliterepo.NewTaskRepositoryTx(tx)
		taskRequestRepo := sqliterepo.NewTaskRequestRepositoryTx(tx)

		req, err := taskRequestRepo.GetByTaskID(ctx, taskID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("task request not found")
			}
			return fmt.Errorf("get task request: %w", err)
		}

		if req.ReceiverUserID != receiverUserID {
			return fmt.Errorf("you cannot reject another user's task")
		}

		if req.Status != domain.RequestStatusPending {
			return fmt.Errorf("task request is not pending")
		}

		comment = strings.TrimSpace(comment)
		if comment == "" {
			comment = "rejected"
		}

		if err := taskRequestRepo.UpdateDecision(ctx, taskID, domain.RequestStatusRejected, comment); err != nil {
			return fmt.Errorf("update task request decision: %w", err)
		}

		if err := taskRepo.UpdateStatus(ctx, taskID, domain.TaskStatusRejected); err != nil {
			return fmt.Errorf("update task status: %w", err)
		}

		return nil
	})
}
