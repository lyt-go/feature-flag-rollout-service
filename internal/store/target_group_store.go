package store

import (
	"featureflag/internal/model"
)

// cloneTargetGroup 复制一条目标分组记录，确保读取方拿不到存储内部的指针。
func cloneTargetGroup(g *model.TargetGroup) *model.TargetGroup {
	if g == nil {
		return nil
	}
	cp := *g
	return &cp
}

func (s *MemoryStore) CreateTargetGroup(g *model.TargetGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.groups[g.ID]; ok {
		return ErrConflict
	}
	s.groups[g.ID] = g
	return nil
}

func (s *MemoryStore) GetTargetGroup(id string) (*model.TargetGroup, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	g, ok := s.groups[id]
	if !ok {
		return nil, ErrNotFound
	}
	return cloneTargetGroup(g), nil
}

func (s *MemoryStore) ListTargetGroups() []*model.TargetGroup {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.TargetGroup, 0, len(s.groups))
	for _, g := range s.groups {
		list = append(list, cloneTargetGroup(g))
	}
	return list
}

func (s *MemoryStore) UpdateTargetGroup(g *model.TargetGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.groups[g.ID]; !ok {
		return ErrNotFound
	}
	s.groups[g.ID] = g
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
