package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"tgPlanBot/internal/domain"
)

type TaskRequestRepository struct {
	db queryer
}

func NewTaskRequestRepository(db *sql.DB) *TaskRequestRepository {
	return &TaskRequestRepository{db: db}
}

func NewTaskRequestRepositoryTx(tx *sql.Tx) *TaskRequestRepository {
	return &TaskRequestRepository{db: tx}
}

func (r *TaskRequestRepository) Create(ctx context.Context, req *domain.TaskRequest) error {
	const query = `
		INSERT INTO task_requests (task_id, sender_user_id, receiver_user_id, status, comment)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		req.TaskID,
		req.SenderUserID,
		req.ReceiverUserID,
		req.Status,
		req.Comment,
	)
	if err != nil {
		return fmt.Errorf("insert task request: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get inserted task request id: %w", err)
	}

	req.ID = id
	return nil
}

func (r *TaskRequestRepository) GetByTaskID(ctx context.Context, taskID int64) (*domain.TaskRequest, error) {
	const query = `
		SELECT id, task_id, sender_user_id, receiver_user_id, status, comment, decided_at, created_at
		FROM task_requests
		WHERE task_id = ?
	`

	row := r.db.QueryRowContext(ctx, query, taskID)

	req, err := scanTaskRequest(row)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func (r *TaskRequestRepository) ListPendingByReceiver(ctx context.Context, receiverUserID int64) ([]domain.TaskRequest, error) {
	const query = `
		SELECT id, task_id, sender_user_id, receiver_user_id, status, comment, decided_at, created_at
		FROM task_requests
		WHERE receiver_user_id = ? AND status = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, receiverUserID, domain.RequestStatusPending)
	if err != nil {
		return nil, fmt.Errorf("query pending task requests: %w", err)
	}
	defer rows.Close()

	requests := make([]domain.TaskRequest, 0)

	for rows.Next() {
		req, err := scanTaskRequest(rows)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate task request rows: %w", err)
	}

	return requests, nil
}

func (r *TaskRequestRepository) UpdateDecision(ctx context.Context, taskID int64, status domain.RequestStatus, comment string) error {
	const query = `
		UPDATE task_requests
		SET status = ?, comment = ?, decided_at = CURRENT_TIMESTAMP
		WHERE task_id = ?
	`

	result, err := r.db.ExecContext(ctx, query, status, comment, taskID)
	if err != nil {
		return fmt.Errorf("update task request decision: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("task request rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

type requestScanner interface {
	Scan(dest ...any) error
}

func scanTaskRequest(s requestScanner) (domain.TaskRequest, error) {
	var req domain.TaskRequest
	var decidedAt sql.NullTime
	var createdAt string

	err := s.Scan(
		&req.ID,
		&req.TaskID,
		&req.SenderUserID,
		&req.ReceiverUserID,
		&req.Status,
		&req.Comment,
		&decidedAt,
		&createdAt,
	)
	if err != nil {
		return domain.TaskRequest{}, err
	}

	if decidedAt.Valid {
		t := decidedAt.Time
		req.DecidedAt = &t
	}

	parsedCreatedAt, err := parseRequestTime(createdAt)
	if err != nil {
		return domain.TaskRequest{}, fmt.Errorf("parse created_at: %w", err)
	}

	req.CreatedAt = parsedCreatedAt
	return req, nil
}

func parseRequestTime(value string) (time.Time, error) {
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
