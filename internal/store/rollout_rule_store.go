package store

import (
	"featureflag/internal/model"
)

func (s *MemoryStore) CreateRolloutRule(r *model.RolloutRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rules[r.ID]; ok {
		return ErrConflict
	}
	s.rules[r.ID] = clone(r)
	return nil
}

func (s *MemoryStore) GetRolloutRule(id string) (*model.RolloutRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.rules[id]
	if !ok {
		return nil, ErrNotFound
	}
	return clone(r), nil
}

func (s *MemoryStore) ListRolloutRules() []*model.RolloutRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RolloutRule, 0, len(s.rules))
	for _, r := range s.rules {
		list = append(list, clone(r))
	}
	return list
}

func (s *MemoryStore) UpdateRolloutRule(r *model.RolloutRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rules[r.ID]; !ok {
		return ErrNotFound
	}
	s.rules[r.ID] = clone(r)
	return nil
}

func (s *MemoryStore) DeleteRolloutRule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rules[id]; !ok {
		return ErrNotFound
	}
	delete(s.rules, id)
	return nil
}
