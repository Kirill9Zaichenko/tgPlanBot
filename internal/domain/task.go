package domain

import "time"

type Task struct {
	ID        int
	Text      string
	Done      bool
	CreatedAt time.Time
}
