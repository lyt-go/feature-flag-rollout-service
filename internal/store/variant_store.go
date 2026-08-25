package store

import (
	"featureflag/internal/model"
)

// cloneVariant 复制一条变体记录，确保读取方拿不到存储内部的指针。
func cloneVariant(v *model.Variant) *model.Variant {
	if v == nil {
		return nil
	}
	cp := *v
	return &cp
}

func (s *MemoryStore) CreateVariant(v *model.Variant) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.variants[v.ID]; ok {
		return ErrConflict
	}
	s.variants[v.ID] = v
	return nil
}

func (s *MemoryStore) GetVariant(id string) (*model.Variant, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.variants[id]
	if !ok {
		return nil, ErrNotFound
	}
	return cloneVariant(v), nil
}

func (s *MemoryStore) ListVariants() []*model.Variant {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Variant, 0, len(s.variants))
	for _, v := range s.variants {
		list = append(list, cloneVariant(v))
	}
	return list
}

func (s *MemoryStore) UpdateVariant(v *model.Variant) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.variants[v.ID]; !ok {
		return ErrNotFound
	}
	s.variants[v.ID] = v
	return nil
}

func (s *MemoryStore) DeleteVariant(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.variants[id]; !ok {
		return ErrNotFound
	}
	delete(s.variants, id)
	return nil
}
