package model

import (
	"encoding/json"
	"strings"
	"time"
)

// TargetGroup 表示一组目标用户，通过规则 JSON 描述匹配条件。
type TargetGroup struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Rules       string    `json:"rules"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate 校验目标分组字段，并确保规则是合法 JSON 对象。
func (g *TargetGroup) Validate() error {
	g.Name = strings.TrimSpace(g.Name)
	g.Description = strings.TrimSpace(g.Description)
	g.Rules = strings.TrimSpace(g.Rules)
	if g.Name == "" {
		return NewValidationError("name", "分组名称不能为空")
	}
	if g.Rules == "" {
		return NewValidationError("rules", "匹配规则不能为空")
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(g.Rules), &m); err != nil {
		return NewValidationError("rules", "匹配规则必须是合法的 JSON 对象")
	}
	return nil
}

// MatchContext 判断上下文是否命中该分组的全部规则（key=value 需逐一相等）。
func (g *TargetGroup) MatchContext(context map[string]string) bool {
	var rules map[string]string
	if err := json.Unmarshal([]byte(g.Rules), &rules); err != nil {
		return false
	}
	for k, v := range rules {
		if cv, ok := context[k]; !ok || cv != v {
			return false
		}
	}
	return true
}

// TargetGroupFilter 目标分组筛选条件。
type TargetGroupFilter struct {
	Keyword string
}

// Match 判断分组是否匹配筛选条件。
func (f TargetGroupFilter) Match(g *TargetGroup) bool {
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(g.Name), k) &&
			!strings.Contains(strings.ToLower(g.Description), k) {
			return false
		}
	}
	return true
}
