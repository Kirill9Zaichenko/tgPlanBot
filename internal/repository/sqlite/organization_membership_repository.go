package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"tgPlanBot/internal/domain"
)

type OrganizationMembershipRepository struct {
	db queryer
}

func NewOrganizationMembershipRepository(db *sql.DB) *OrganizationMembershipRepository {
	return &OrganizationMembershipRepository{db: db}
}

func NewOrganizationMembershipRepositoryTx(tx *sql.Tx) *OrganizationMembershipRepository {
	return &OrganizationMembershipRepository{db: tx}
}

func (r *OrganizationMembershipRepository) AddUser(ctx context.Context, membership *domain.OrganizationMembership) error {
	const query = `
		INSERT INTO organization_memberships (organization_id, user_id, role)
		VALUES (?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query, membership.OrganizationID, membership.UserID, membership.Role)
	if err != nil {
		return fmt.Errorf("insert organizations membership: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get inserted membership id: %w", err)
	}

	membership.ID = id
	return nil
}

func (r *OrganizationMembershipRepository) IsMember(ctx context.Context, organizationID, userID int64) (bool, error) {
	const query = `
		SELECT 1
		FROM organization_memberships
		WHERE organization_id = ? AND user_id = ?
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, organizationID, userID)

	var exists int
	err := row.Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check organizations membership: %w", err)
	}

	return true, nil
}

func (r *OrganizationMembershipRepository) ListUsersByOrganizationID(ctx context.Context, organizationID int64) ([]domain.User, error) {
	const query = `
		SELECT u.id, u.telegram_id, u.username, u.first_name, u.last_name, u.created_at, u.updated_at
		FROM users u
		INNER JOIN organization_memberships om ON om.user_id = u.id
		WHERE om.organization_id = ?
		ORDER BY u.username ASC, u.first_name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, organizationID)
	if err != nil {
		return nil, fmt.Errorf("list users by organizations id: %w", err)
	}
	defer rows.Close()

	var users []domain.User

	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("scan org user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate org users: %w", err)
	}

	return users, nil
}
