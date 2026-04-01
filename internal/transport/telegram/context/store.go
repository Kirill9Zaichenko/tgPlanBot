package context

import "sync"

type Store struct {
	mu     sync.RWMutex
	active map[int64]int64
}

func NewStore() *Store {
	return &Store{
		active: make(map[int64]int64),
	}
}

func (s *Store) SetActiveOrganization(userID, organizationID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.active[userID] = organizationID
}

func (s *Store) GetActiveOrganization(userID int64) (int64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	orgID, ok := s.active[userID]
	return orgID, ok
}

func (s *Store) ClearActiveOrganization(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.active, userID)
}
