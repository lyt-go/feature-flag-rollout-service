package model

import (
	"strings"
	"time"
)

// ChangeLog 动作常量。
const (
	ActionCreate  = "create"
	ActionUpdate  = "update"
	ActionToggle  = "toggle"
	ActionPause   = "pause"
	ActionResume  = "resume"
	ActionArchive = "archive"
	ActionRestore = "restore"
)

// ChangeLog 表示一次开关变更的审计记录。
type ChangeLog struct {
	ID        string    `json:"id"`
	FlagID    string    `json:"flag_id"`
	Operator  string    `json:"operator"`
	Action    string    `json:"action"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate 校验变更记录字段。
func (c *ChangeLog) Validate() error {
	c.FlagID = strings.TrimSpace(c.FlagID)
	c.Operator = strings.TrimSpace(c.Operator)
	c.Action = strings.TrimSpace(c.Action)
	c.Detail = strings.TrimSpace(c.Detail)
	if c.FlagID == "" {
		return NewValidationError("flag_id", "开关 ID 不能为空")
	}
	if c.Operator == "" {
		return NewValidationError("operator", "操作人不能为空")
	}
	switch c.Action {
	case ActionCreate, ActionUpdate, ActionToggle, ActionPause, ActionResume, ActionArchive, ActionRestore:
	default:
		return NewValidationError("action", "变更动作不合法")
	}
	return nil
}

// ChangeLogFilter 变更记录筛选条件。
type ChangeLogFilter struct {
	FlagID string
	Action string
}

// Match 判断记录是否匹配筛选条件。
func (f ChangeLogFilter) Match(c *ChangeLog) bool {
	if f.FlagID != "" && c.FlagID != f.FlagID {
		return false
	}
	if f.Action != "" && c.Action != f.Action {
		return false
	}
	return true
}
