package store

import (
	"featureflag/internal/model"
)

func (s *MemoryStore) CreateTargetGroup(g *model.TargetGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.groups[g.ID]; ok {
		return ErrConflict
	}
	s.groups[g.ID] = clone(g)
	return nil
}

func (s *MemoryStore) GetTargetGroup(id string) (*model.TargetGroup, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	g, ok := s.groups[id]
	if !ok {
		return nil, ErrNotFound
	}
	return clone(g), nil
}

func (s *MemoryStore) ListTargetGroups() []*model.TargetGroup {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.TargetGroup, 0, len(s.groups))
	for _, g := range s.groups {
		list = append(list, clone(g))
	}
	return list
}

func (s *MemoryStore) UpdateTargetGroup(g *model.TargetGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.groups[g.ID]; !ok {
		return ErrNotFound
	}
	s.groups[g.ID] = clone(g)
	return nil
}

func (s *MemoryStore) DeleteTargetGroup(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.groups[id]; !ok {
		return ErrNotFound
	}
	delete(s.groups, id)
	return nil
}
