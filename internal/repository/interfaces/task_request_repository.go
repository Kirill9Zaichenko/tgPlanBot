package interfaces

import (
	"context"

	"tgPlanBot/internal/domain"
)

type TaskRequestRepository interface {
	Create(ctx context.Context, req *domain.TaskRequest) error
	GetByTaskID(ctx context.Context, taskID int64) (*domain.TaskRequest, error)
	ListPendingByReceiver(ctx context.Context, receiverUserID int64) ([]domain.TaskRequest, error)
	ListInboxItemsByReceiver(ctx context.Context, receiverUserID int64) ([]domain.InboxItem, error)
	UpdateDecision(ctx context.Context, taskID int64, status domain.RequestStatus, comment string) error
}
