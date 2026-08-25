package model

import (
	"strings"
	"time"
)

// RolloutRule 状态常量。
const (
	RolloutEnabled  = "enabled"
	RolloutDisabled = "disabled"
)

// RolloutRule 表示一条灰度规则：命中目标分组时按比例返回指定变体。
type RolloutRule struct {
	ID            string    `json:"id"`
	FlagID        string    `json:"flag_id"`
	TargetGroupID string    `json:"target_group_id"`
	VariantID     string    `json:"variant_id"`
	Percentage    int       `json:"percentage"`
	Status        string    `json:"status"`
	Priority      int       `json:"priority"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Validate 校验灰度规则字段。
func (r *RolloutRule) Validate() error {
	r.FlagID = strings.TrimSpace(r.FlagID)
	r.TargetGroupID = strings.TrimSpace(r.TargetGroupID)
	r.VariantID = strings.TrimSpace(r.VariantID)
	if r.FlagID == "" {
		return NewValidationError("flag_id", "开关 ID 不能为空")
	}
	if r.VariantID == "" {
		return NewValidationError("variant_id", "变体 ID 不能为空")
	}
	if r.Percentage < 0 || r.Percentage > 100 {
		return NewValidationError("percentage", "灰度比例必须在 0-100 之间")
	}
	if r.Priority < 0 {
		return NewValidationError("priority", "优先级不能为负数")
	}
	if r.Status == "" {
		r.Status = RolloutEnabled
	}
	if r.Status != RolloutEnabled && r.Status != RolloutDisabled {
		return NewValidationError("status", "规则状态不合法")
	}
	return nil
}

// RolloutRuleFilter 灰度规则筛选条件。
type RolloutRuleFilter struct {
	FlagID string
	Status string
}

// Match 判断规则是否匹配筛选条件。
func (f RolloutRuleFilter) Match(r *RolloutRule) bool {
	if f.FlagID != "" && r.FlagID != f.FlagID {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	return true
}
