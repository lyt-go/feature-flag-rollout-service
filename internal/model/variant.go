package model

import (
	"strings"
	"time"
)

// Variant 表示开关的一个取值变体。
type Variant struct {
	ID        string    `json:"id"`
	FlagID    string    `json:"flag_id"`
	Name      string    `json:"name"`
	Payload   string    `json:"payload"`
	Weight    int       `json:"weight"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate 校验变体字段。
func (v *Variant) Validate() error {
	v.FlagID = strings.TrimSpace(v.FlagID)
	v.Name = strings.TrimSpace(v.Name)
	v.Payload = strings.TrimSpace(v.Payload)
	if v.FlagID == "" {
		return NewValidationError("flag_id", "开关 ID 不能为空")
	}
	if v.Name == "" {
		return NewValidationError("name", "变体名称不能为空")
	}
	if v.Weight < 0 {
		return NewValidationError("weight", "权重不能为负数")
	}
	return nil
}

// VariantFilter 变体列表筛选条件。
type VariantFilter struct {
	FlagID  string
	Keyword string
}

// Match 判断变体是否匹配筛选条件。
func (f VariantFilter) Match(v *Variant) bool {
	if f.FlagID != "" && v.FlagID != f.FlagID {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(v.Name), k) &&
			!strings.Contains(strings.ToLower(v.Payload), k) {
			return false
		}
	}
	return true
}
