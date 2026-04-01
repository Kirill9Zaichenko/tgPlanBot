package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"tgPlanBot/internal/domain"
)

type UserRepository struct {
	db queryer
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func NewUserRepositoryTx(tx *sql.Tx) *UserRepository {
	return &UserRepository{db: tx}
}

func (r *UserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	const query = `
		SELECT id, telegram_id, username, first_name, last_name, created_at, updated_at
		FROM users
		WHERE telegram_id = ?
	`

	row := r.db.QueryRowContext(ctx, query, telegramID)

	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by telegram id: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	const query = `
		INSERT INTO users (telegram_id, username, first_name, last_name)
		VALUES (?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.TelegramID,
		user.Username,
		user.FirstName,
		user.LastName,
	)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get inserted user id: %w", err)
	}

	user.ID = id

	loadedUser, err := r.GetByTelegramID(ctx, user.TelegramID)
	if err != nil {
		return fmt.Errorf("load inserted user: %w", err)
	}

	*user = *loadedUser
	return nil
}

func (r *UserRepository) UpdateTelegramProfile(ctx context.Context, user *domain.User) error {
	const query = `
		UPDATE users
		SET username = ?, first_name = ?, last_name = ?, updated_at = CURRENT_TIMESTAMP
		WHERE telegram_id = ?
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.Username,
		user.FirstName,
		user.LastName,
		user.TelegramID,
	)
	if err != nil {
		return fmt.Errorf("update telegram profile: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("user rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	loadedUser, err := r.GetByTelegramID(ctx, user.TelegramID)
	if err != nil {
		return fmt.Errorf("reload updated user: %w", err)
	}

	*user = *loadedUser
	return nil
}

type userScanner interface {
	Scan(dest ...any) error
}

func scanUser(s userScanner) (domain.User, error) {
	var user domain.User
	var createdAt string
	var updatedAt string

	err := s.Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return domain.User{}, err
	}

	user.CreatedAt, err = parseUserTime(createdAt)
	if err != nil {
		return domain.User{}, fmt.Errorf("parse created_at: %w", err)
	}

	user.UpdatedAt, err = parseUserTime(updatedAt)
	if err != nil {
		return domain.User{}, fmt.Errorf("parse updated_at: %w", err)
	}

	return user, nil
}

func parseUserTime(value string) (time.Time, error) {
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
