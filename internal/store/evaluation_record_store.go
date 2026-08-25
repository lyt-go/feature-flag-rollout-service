package store

import (
	"featureflag/internal/model"
)

// cloneEvaluationRecord 复制一条评估记录，确保读取方拿不到存储内部的指针。
func cloneEvaluationRecord(e *model.EvaluationRecord) *model.EvaluationRecord {
	if e == nil {
		return nil
	}
	cp := *e
	return &cp
}

func (s *MemoryStore) CreateEvaluationRecord(e *model.EvaluationRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.evaluations[e.ID]; ok {
		return ErrConflict
	}
	s.evaluations[e.ID] = e
	return nil
}

func (s *MemoryStore) GetEvaluationRecord(id string) (*model.EvaluationRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.evaluations[id]
	if !ok {
		return nil, ErrNotFound
	}
	return cloneEvaluationRecord(e), nil
}

func (s *MemoryStore) ListEvaluationRecords() []*model.EvaluationRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.EvaluationRecord, 0, len(s.evaluations))
	for _, e := range s.evaluations {
		list = append(list, cloneEvaluationRecord(e))
	}
	return list
}
