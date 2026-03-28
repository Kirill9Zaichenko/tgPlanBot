package task

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

func (s *Service) ListByAssignee(ctx context.Context, assigneeUserID int64) ([]domain.Task, error) {
	if assigneeUserID <= 0 {
		return nil, fmt.Errorf("invalid assignee user id")
	}

	tasks, err := s.taskRepo.ListByAssignee(ctx, assigneeUserID)
	if err != nil {
		return nil, fmt.Errorf("list tasks by assignee: %w", err)
	}

	return tasks, nil
}

func (s *Service) Create(ctx context.Context, input CreateTaskInput) (*domain.Task, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if input.CreatorUserID <= 0 {
		return nil, fmt.Errorf("creator user id is required")
	}
	if input.AssigneeUserID <= 0 {
		return nil, fmt.Errorf("assignee user id is required")
	}

	var createdTask *domain.Task

	err := db.WithTx(ctx, s.database, func(tx *sql.Tx) error {
		taskRepo := sqliterepo.NewTaskRepositoryTx(tx)
		taskRequestRepo := sqliterepo.NewTaskRequestRepositoryTx(tx)

		task := &domain.Task{
			Title:          input.Title,
			Description:    input.Description,
			CreatorUserID:  input.CreatorUserID,
			AssigneeUserID: input.AssigneeUserID,
			Status:         domain.TaskStatusPendingAcceptance,
			DueAt:          input.DueAt,
		}

		if err := taskRepo.Create(ctx, task); err != nil {
			return fmt.Errorf("create task: %w", err)
		}

		req := &domain.TaskRequest{
			TaskID:         task.ID,
			SenderUserID:   input.CreatorUserID,
			ReceiverUserID: input.AssigneeUserID,
			Status:         domain.RequestStatusPending,
			Comment:        "",
		}

		if err := taskRequestRepo.Create(ctx, req); err != nil {
			return fmt.Errorf("create task request: %w", err)
		}

		createdTask = task
		return nil
	})
	if err != nil {
		return nil, err
	}

	return createdTask, nil
}
