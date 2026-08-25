package model

import (
	"strings"
	"time"
)

// EvaluationRecord 表示一次开关评估的命中记录。
type EvaluationRecord struct {
	ID            string    `json:"id"`
	FlagID        string    `json:"flag_id"`
	TargetKey     string    `json:"target_key"`
	Context       string    `json:"context"`
	ResultVariant string    `json:"result_variant"`
	Matched       bool      `json:"matched"`
	CreatedAt     time.Time `json:"created_at"`
}

// Validate 校验评估记录字段。
func (e *EvaluationRecord) Validate() error {
	e.FlagID = strings.TrimSpace(e.FlagID)
	e.TargetKey = strings.TrimSpace(e.TargetKey)
	e.Context = strings.TrimSpace(e.Context)
	e.ResultVariant = strings.TrimSpace(e.ResultVariant)
	if e.FlagID == "" {
		return NewValidationError("flag_id", "开关 ID 不能为空")
	}
	if e.TargetKey == "" {
		return NewValidationError("target_key", "目标标识不能为空")
	}
	return nil
}

// EvaluationFilter 评估记录筛选条件。
type EvaluationFilter struct {
	FlagID        string
	TargetKey     string
	ResultVariant string
}

// Match 判断评估记录是否匹配筛选条件。
func (f EvaluationFilter) Match(e *EvaluationRecord) bool {
	if f.FlagID != "" && e.FlagID != f.FlagID {
		return false
	}
	if f.TargetKey != "" && e.TargetKey != f.TargetKey {
		return false
	}
	if f.ResultVariant != "" && e.ResultVariant != f.ResultVariant {
		return false
	}
	return true
}
