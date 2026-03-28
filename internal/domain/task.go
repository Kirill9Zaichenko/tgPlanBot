package domain

import "time"

type TaskStatus string

const (
	TaskStatusPendingAcceptance TaskStatus = "pending_acceptance"
	TaskStatusAccepted          TaskStatus = "accepted"
	TaskStatusInProgress        TaskStatus = "in_progress"
	TaskStatusDone              TaskStatus = "done"
	TaskStatusRejected          TaskStatus = "rejected"
	TaskStatusCancelled         TaskStatus = "cancelled"
)

type Task struct {
	ID             int64
	Title          string
	Description    string
	CreatorUserID  int64
	AssigneeUserID int64
	Status         TaskStatus
	DueAt          *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
