// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"featureflag/internal/model"
)

var (
	// ErrNotFound 表示记录不存在。
	ErrNotFound = errors.New("记录不存在")
	// ErrConflict 表示记录已存在或状态冲突。
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// Flag
	CreateFlag(f *model.Flag) error
	GetFlag(id string) (*model.Flag, error)
	GetFlagByKey(key string) (*model.Flag, error)
	ListFlags() []*model.Flag
	UpdateFlag(f *model.Flag) error
	DeleteFlag(id string) error

	// Variant
	CreateVariant(v *model.Variant) error
	GetVariant(id string) (*model.Variant, error)
	ListVariants() []*model.Variant
	UpdateVariant(v *model.Variant) error
	DeleteVariant(id string) error

	// TargetGroup
	CreateTargetGroup(g *model.TargetGroup) error
	GetTargetGroup(id string) (*model.TargetGroup, error)
	ListTargetGroups() []*model.TargetGroup
	UpdateTargetGroup(g *model.TargetGroup) error
	DeleteTargetGroup(id string) error

	// RolloutRule
	CreateRolloutRule(r *model.RolloutRule) error
	GetRolloutRule(id string) (*model.RolloutRule, error)
	ListRolloutRules() []*model.RolloutRule
	UpdateRolloutRule(r *model.RolloutRule) error
	DeleteRolloutRule(id string) error

	// ChangeLog
	CreateChangeLog(c *model.ChangeLog) error
	GetChangeLog(id string) (*model.ChangeLog, error)
	ListChangeLogs() []*model.ChangeLog

	// EvaluationRecord
	CreateEvaluationRecord(e *model.EvaluationRecord) error
	GetEvaluationRecord(id string) (*model.EvaluationRecord, error)
	ListEvaluationRecords() []*model.EvaluationRecord
}
