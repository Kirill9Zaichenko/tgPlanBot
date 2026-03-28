package state

import "sync"

type Store struct {
	mu     sync.RWMutex
	states map[int64]NewTaskState
}

func NewStore() *Store {
	return &Store{
		states: make(map[int64]NewTaskState),
	}
}

func (s *Store) Get(userID int64) (NewTaskState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.states[userID]
	return state, ok
}

func (s *Store) Set(userID int64, state NewTaskState) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.states[userID] = state
}

func (s *Store) Delete(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.states, userID)
}
