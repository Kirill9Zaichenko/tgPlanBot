package organization

import (
	"context"
	"fmt"
	"strings"

	"tgPlanBot/internal/domain"
	repoif "tgPlanBot/internal/repository/interfaces"
)

type Service struct {
	orgRepo        repoif.OrganizationRepository
	membershipRepo repoif.OrganizationMembershipRepository
}

func NewService(
	orgRepo repoif.OrganizationRepository,
	membershipRepo repoif.OrganizationMembershipRepository,
) *Service {
	return &Service{
		orgRepo:        orgRepo,
		membershipRepo: membershipRepo,
	}
}

func (s *Service) ListByUserID(ctx context.Context, userID int64) ([]domain.Organization, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user id")
	}

	items, err := s.orgRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list organizations: %w", err)
	}

	return items, nil
}

func (s *Service) ListMembers(ctx context.Context, organizationID, userID int64) ([]domain.User, error) {
	if organizationID <= 0 || userID <= 0 {
		return nil, fmt.Errorf("invalid organization or user id")
	}

	isMember, err := s.membershipRepo.IsMember(ctx, organizationID, userID)
	if err != nil {
		return nil, fmt.Errorf("check membership: %w", err)
	}
	if !isMember {
		return nil, fmt.Errorf("user is not a member of organization")
	}

	users, err := s.membershipRepo.ListUsersByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("list organization members: %w", err)
	}

	return users, nil
}

func (s *Service) Create(ctx context.Context, ownerUserID int64, name, slug string) (*domain.Organization, error) {
	name = strings.TrimSpace(name)
	slug = strings.TrimSpace(slug)

	if ownerUserID <= 0 {
		return nil, fmt.Errorf("owner user id is required")
	}
	if name == "" {
		return nil, fmt.Errorf("organization name is required")
	}
	if slug == "" {
		return nil, fmt.Errorf("organization slug is required")
	}

	org := &domain.Organization{
		Name: name,
		Slug: slug,
	}

	if err := s.orgRepo.Create(ctx, org); err != nil {
		return nil, fmt.Errorf("create organization: %w", err)
	}

	membership := &domain.OrganizationMembership{
		OrganizationID: org.ID,
		UserID:         ownerUserID,
		Role:           domain.OrganizationRoleAdmin,
	}

	if err := s.membershipRepo.AddUser(ctx, membership); err != nil {
		return nil, fmt.Errorf("create owner membership: %w", err)
	}

	return org, nil
}

func (s *Service) GetByIDForUser(ctx context.Context, organizationID, userID int64) (*domain.Organization, error) {
	if organizationID <= 0 || userID <= 0 {
		return nil, fmt.Errorf("invalid organization or user id")
	}

	isMember, err := s.membershipRepo.IsMember(ctx, organizationID, userID)
	if err != nil {
		return nil, fmt.Errorf("check membership: %w", err)
	}
	if !isMember {
		return nil, fmt.Errorf("user is not a member of organization")
	}

	org, err := s.orgRepo.GetByID(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("get organization by id: %w", err)
	}

	return org, nil
}
