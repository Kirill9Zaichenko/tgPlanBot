package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"tgPlanBot/internal/domain"
)

type TaskRepository struct {
	db queryer
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func NewTaskRepositoryTx(tx *sql.Tx) *TaskRepository {
	return &TaskRepository{db: tx}
}

func (r *TaskRepository) ListByAssignee(ctx context.Context, assigneeUserID int64) ([]domain.Task, error) {
	const query = `
		SELECT id, title, description, creator_user_id, assignee_user_id, status, due_at, created_at, updated_at
		FROM tasks
		WHERE assignee_user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, assigneeUserID)
	if err != nil {
		return nil, fmt.Errorf("query tasks by assignee: %w", err)
	}
	defer rows.Close()

	tasks := make([]domain.Task, 0)

	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate task rows: %w", err)
	}

	return tasks, nil
}

func (r *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
	const insertQuery = `
		INSERT INTO tasks (title, description, creator_user_id, assignee_user_id, status, due_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	var dueAt any
	if task.DueAt != nil {
		dueAt = task.DueAt.UTC().Format(time.RFC3339)
	}

	result, err := r.db.ExecContext(
		ctx,
		insertQuery,
		task.Title,
		task.Description,
		task.CreatorUserID,
		task.AssigneeUserID,
		task.Status,
		dueAt,
	)
	if err != nil {
		return fmt.Errorf("insert task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get inserted task id: %w", err)
	}

	task.ID = id

	loadedTask, err := r.GetByID(ctx, task.ID)
	if err != nil {
		return fmt.Errorf("load inserted task: %w", err)
	}

	*task = *loadedTask
	return nil
}

func (r *TaskRepository) GetByID(ctx context.Context, taskID int64) (*domain.Task, error) {
	const query = `
		SELECT id, title, description, creator_user_id, assignee_user_id, status, due_at, created_at, updated_at
		FROM tasks
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, taskID)

	task, err := scanTask(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("get task by id: %w", err)
	}

	return &task, nil
}

func (r *TaskRepository) UpdateStatus(ctx context.Context, taskID int64, status domain.TaskStatus) error {
	const query = `
		UPDATE tasks
		SET status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, status, taskID)
	if err != nil {
		return fmt.Errorf("update task status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("task status rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanTask(s scanner) (domain.Task, error) {
	var task domain.Task
	var dueAt sql.NullTime
	var createdAt string
	var updatedAt string

	err := s.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.CreatorUserID,
		&task.AssigneeUserID,
		&task.Status,
		&dueAt,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return domain.Task{}, err
	}

	if dueAt.Valid {
		t := dueAt.Time
		task.DueAt = &t
	}

	task.CreatedAt, err = parseSQLiteTime(createdAt)
	if err != nil {
		return domain.Task{}, fmt.Errorf("parse created_at: %w", err)
	}

	task.UpdatedAt, err = parseSQLiteTime(updatedAt)
	if err != nil {
		return domain.Task{}, fmt.Errorf("parse updated_at: %w", err)
	}

	return task, nil
}

func parseSQLiteTime(value string) (time.Time, error) {
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
