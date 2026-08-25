package model

import (
	"strings"
	"time"
)

// Flag 状态常量。
const (
	FlagActive   = "active"
	FlagPaused   = "paused"
	FlagArchived = "archived"
)

// flagTransitions 定义开关合法状态流转。
var flagTransitions = map[string]map[string]bool{
	FlagActive:   {FlagPaused: true, FlagArchived: true},
	FlagPaused:   {FlagActive: true, FlagArchived: true},
	FlagArchived: {FlagActive: true},
}

// CanTransitionFlag 判断开关状态能否从 from 流转到 to。
func CanTransitionFlag(from, to string) bool {
	if m, ok := flagTransitions[from]; ok {
		return m[to]
	}
	return false
}

// Flag 表示一个特性开关。
type Flag struct {
	ID          string    `json:"id"`
	Key         string    `json:"key"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate 校验并规范化开关字段。
func (f *Flag) Validate() error {
	f.Key = strings.ToLower(strings.TrimSpace(f.Key))
	f.Name = strings.TrimSpace(f.Name)
	f.Description = strings.TrimSpace(f.Description)
	if f.Key == "" {
		return NewValidationError("key", "开关标识不能为空")
	}
	for _, r := range f.Key {
		if !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' || r == '_' || r == '.') {
			return NewValidationError("key", "开关标识只能包含小写字母、数字、-、_、.")
		}
	}
	if f.Name == "" {
		return NewValidationError("name", "开关名称不能为空")
	}
	if f.Status == "" {
		f.Status = FlagActive
	}
	if f.Status != FlagActive && f.Status != FlagPaused && f.Status != FlagArchived {
		return NewValidationError("status", "开关状态不合法")
	}
	return nil
}

// FlagFilter 开关列表筛选条件。
type FlagFilter struct {
	Status  string
	Keyword string
}

// Match 判断开关是否匹配筛选条件。
func (f FlagFilter) Match(fl *Flag) bool {
	if f.Status != "" && fl.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(fl.Name), k) &&
			!strings.Contains(strings.ToLower(fl.Key), k) {
			return false
		}
	}
	return true
}
