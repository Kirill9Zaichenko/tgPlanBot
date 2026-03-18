package storage

import "tgBotPlan/internal/domain"

type TaskStore interface {
	Add(chatID int64, text string) *domain.Task
	List(chatID int64) []*domain.Task
	Done(chatID int64, id int) bool
	Delete(chatID int64, id int) bool
	Clear(chatID int64)
}
