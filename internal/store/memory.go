package store

import (
	"sync"

	"featureflag/internal/model"
)

// MemoryStore 基于内存 map 的线程安全存储实现。
type MemoryStore struct {
	mu          sync.RWMutex
	flags       map[string]*model.Flag
	variants    map[string]*model.Variant
	groups      map[string]*model.TargetGroup
	rules       map[string]*model.RolloutRule
	changeLogs  map[string]*model.ChangeLog
	evaluations map[string]*model.EvaluationRecord
}

// NewMemoryStore 创建空的 MemoryStore。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		flags:       make(map[string]*model.Flag),
		variants:    make(map[string]*model.Variant),
		groups:      make(map[string]*model.TargetGroup),
		rules:       make(map[string]*model.RolloutRule),
		changeLogs:  make(map[string]*model.ChangeLog),
		evaluations: make(map[string]*model.EvaluationRecord),
	}
}

var _ Store = (*MemoryStore)(nil)
