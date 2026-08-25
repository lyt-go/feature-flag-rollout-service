package store

import (
	"featureflag/internal/model"
)

// cloneChangeLog 复制一条变更记录，确保读取方拿不到存储内部的指针。
func cloneChangeLog(c *model.ChangeLog) *model.ChangeLog {
	if c == nil {
		return nil
	}
	cp := *c
	return &cp
}

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
	return cloneChangeLog(c), nil
}

func (s *MemoryStore) ListChangeLogs() []*model.ChangeLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ChangeLog, 0, len(s.changeLogs))
	for _, c := range s.changeLogs {
		list = append(list, cloneChangeLog(c))
	}
	return list
}
