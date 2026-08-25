package store

import (
	"featureflag/internal/model"
)

func (s *MemoryStore) CreateChangeLog(c *model.ChangeLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.changeLogs[c.ID]; ok {
		return ErrConflict
	}
	s.changeLogs[c.ID] = c
	return nil
}

func (s *MemoryStore) GetChangeLog(id string) (*model.ChangeLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.changeLogs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *MemoryStore) ListChangeLogs() []*model.ChangeLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ChangeLog, 0, len(s.changeLogs))
	for _, c := range s.changeLogs {
		list = append(list, c)
	}
	return list
}
