package store

import (
	"featureflag/internal/model"
)

func (s *MemoryStore) CreateFlag(f *model.Flag) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.flags {
		if exist.Key == f.Key {
			return ErrConflict
		}
	}
	s.flags[f.ID] = f
	return nil
}

func (s *MemoryStore) GetFlag(id string) (*model.Flag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, ok := s.flags[id]
	if !ok {
		return nil, ErrNotFound
	}
	return f, nil
}

func (s *MemoryStore) GetFlagByKey(key string) (*model.Flag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, f := range s.flags {
		if f.Key == key {
			return f, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListFlags() []*model.Flag {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Flag, 0, len(s.flags))
	for _, f := range s.flags {
		list = append(list, f)
	}
	return list
}

func (s *MemoryStore) UpdateFlag(f *model.Flag) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.flags[f.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.flags {
		if exist.ID != f.ID && exist.Key == f.Key {
			return ErrConflict
		}
	}
	s.flags[f.ID] = f
	return nil
}

func (s *MemoryStore) DeleteFlag(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.flags[id]; !ok {
		return ErrNotFound
	}
	delete(s.flags, id)
	return nil
}
