package interfaces

import (
	"context"

	"tgPlanBot/internal/domain"
)

type OrganizationMembershipRepository interface {
	AddUser(ctx context.Context, membership *domain.OrganizationMembership) error
	IsMember(ctx context.Context, organizationID, userID int64) (bool, error)
	ListUsersByOrganizationID(ctx context.Context, organizationID int64) ([]domain.User, error)
}
