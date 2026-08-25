package store

import (
	"featureflag/internal/model"
)

// cloneRolloutRule 复制一条灰度规则，确保读取方拿不到存储内部的指针。
// 上层（含灰度预览）若在返回值上改动 Percentage/Status，也不会回写到真实规则，
// 后续评估仍按原始规则命中变体。
func cloneRolloutRule(r *model.RolloutRule) *model.RolloutRule {
	if r == nil {
		return nil
	}
	cp := *r
	return &cp
}

func (s *MemoryStore) CreateRolloutRule(r *model.RolloutRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rules[r.ID]; ok {
		return ErrConflict
	}
	s.rules[r.ID] = r
	return nil
}

func (s *MemoryStore) GetRolloutRule(id string) (*model.RolloutRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.rules[id]
	if !ok {
		return nil, ErrNotFound
	}
	return cloneRolloutRule(r), nil
}

func (s *MemoryStore) ListRolloutRules() []*model.RolloutRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RolloutRule, 0, len(s.rules))
	for _, r := range s.rules {
		list = append(list, cloneRolloutRule(r))
	}
	return list
}

func (s *MemoryStore) UpdateRolloutRule(r *model.RolloutRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rules[r.ID]; !ok {
		return ErrNotFound
	}
	s.rules[r.ID] = r
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
