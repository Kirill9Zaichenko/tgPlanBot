package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"tgPlanBot/internal/domain"
)

type OrganizationRepository struct {
	db queryer
}

func NewOrganizationRepository(db *sql.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func NewOrganizationRepositoryTx(tx *sql.Tx) *OrganizationRepository {
	return &OrganizationRepository{db: tx}
}

func (r *OrganizationRepository) Create(ctx context.Context, org *domain.Organization) error {
	const query = `
		INSERT INTO organizations (name, slug)
		VALUES (?, ?)
	`

	result, err := r.db.ExecContext(ctx, query, org.Name, org.Slug)
	if err != nil {
		return fmt.Errorf("insert organizations: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get inserted organizations id: %w", err)
	}

	org.ID = id

	loaded, err := r.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("reload organizations: %w", err)
	}

	*org = *loaded
	return nil
}

func (r *OrganizationRepository) GetByID(ctx context.Context, id int64) (*domain.Organization, error) {
	const query = `
		SELECT id, name, slug, created_at, updated_at
		FROM organizations
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var org domain.Organization
	var createdAt string
	var updatedAt string

	err := row.Scan(&org.ID, &org.Name, &org.Slug, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("organizations not found")
		}
		return nil, fmt.Errorf("get organizations by id: %w", err)
	}

	org.CreatedAt, err = parseOrgTime(createdAt)
	if err != nil {
		return nil, err
	}

	org.UpdatedAt, err = parseOrgTime(updatedAt)
	if err != nil {
		return nil, err
	}

	return &org, nil
}

func (r *OrganizationRepository) ListByUserID(ctx context.Context, userID int64) ([]domain.Organization, error) {
	const query = `
		SELECT o.id, o.name, o.slug, o.created_at, o.updated_at
		FROM organizations o
		INNER JOIN organization_memberships om ON om.organization_id = o.id
		WHERE om.user_id = ?
		ORDER BY o.name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list organizations by user id: %w", err)
	}
	defer rows.Close()

	var items []domain.Organization

	for rows.Next() {
		var org domain.Organization
		var createdAt string
		var updatedAt string

		if err := rows.Scan(&org.ID, &org.Name, &org.Slug, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan organizations: %w", err)
		}

		org.CreatedAt, err = parseOrgTime(createdAt)
		if err != nil {
			return nil, err
		}

		org.UpdatedAt, err = parseOrgTime(updatedAt)
		if err != nil {
			return nil, err
		}

		items = append(items, org)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate organizations: %w", err)
	}

	return items, nil
}

func parseOrgTime(value string) (time.Time, error) {
	layouts := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
	}

	var lastErr error
	for _, layout := range layouts {
		t, err := time.Parse(layout, value)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}

	return time.Time{}, lastErr
}
