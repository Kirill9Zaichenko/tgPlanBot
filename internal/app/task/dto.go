package task

import "time"

type CreateTaskInput struct {
	Title          string
	Description    string
	CreatorUserID  int64
	AssigneeUserID int64
	DueAt          *time.Time
}
