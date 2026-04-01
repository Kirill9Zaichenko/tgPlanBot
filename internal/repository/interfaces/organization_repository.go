package interfaces

import (
	"context"

	"tgPlanBot/internal/domain"
)

type OrganizationRepository interface {
	Create(ctx context.Context, org *domain.Organization) error
	GetByID(ctx context.Context, id int64) (*domain.Organization, error)
	ListByUserID(ctx context.Context, userID int64) ([]domain.Organization, error)
}
