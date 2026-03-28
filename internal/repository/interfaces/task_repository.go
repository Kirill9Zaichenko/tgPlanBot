package interfaces

import (
	"context"

	"tgPlanBot/internal/domain"
)

type TaskRepository interface {
	ListByAssignee(ctx context.Context, assigneeUserID int64) ([]domain.Task, error)
	Create(ctx context.Context, task *domain.Task) error
	GetByID(ctx context.Context, taskID int64) (*domain.Task, error)
	UpdateStatus(ctx context.Context, taskID int64, status domain.TaskStatus) error
}
