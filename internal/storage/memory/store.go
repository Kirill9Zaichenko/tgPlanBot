package memory

import (
	"sync"
	"tgBotPlan/internal/domain"
	"time"
)

type TaskStore struct {
	mu    sync.Mutex
	seq   map[int64]int
	tasks map[int64][]*domain.Task
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		seq:   make(map[int64]int),
		tasks: make(map[int64][]*domain.Task),
	}
}

func (s *TaskStore) Add(chatID int64, text string) *domain.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.seq[chatID]++
	t := &domain.Task{
		ID:        s.seq[chatID],
		Text:      text,
		CreatedAt: time.Now(),
	}
	s.tasks[chatID] = append(s.tasks[chatID], t)
	return t
}

func (s *TaskStore) List(chatID int64) []*domain.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]*domain.Task, 0, len(s.tasks[chatID]))
	out = append(out, s.tasks[chatID]...)

	return out
}

func (s *TaskStore) Done(chatID int64, id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, t := range s.tasks[chatID] {
		if t.ID == id {
			t.Done = true
			return true
		}
	}
	return false
}

func (s *TaskStore) Delete(chatID int64, id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	ts := s.tasks[chatID]
	for i, t := range ts {
		if t.ID == id {
			s.tasks[chatID] = append(ts[:i], ts[i+1:]...)
			return true
		}
	}
	return false
}

func (s *TaskStore) Clear(chatID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks[chatID] = nil
	s.seq[chatID] = 0
}
