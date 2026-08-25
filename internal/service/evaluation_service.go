package service

import (
	"encoding/json"
	"hash/fnv"
	"sort"
	"strings"
	"time"

	"featureflag/internal/model"
	"featureflag/pkg/idgen"
)

// EvaluationResult 一次开关评估的结果。
type EvaluationResult struct {
	FlagKey   string `json:"flag_key"`
	FlagID    string `json:"flag_id"`
	Enabled   bool   `json:"enabled"`
	Variant   string `json:"variant"`
	VariantID string `json:"variant_id"`
	Matched   bool   `json:"matched"`
}

// Evaluate 依据开关与灰度规则对目标做评估，并落一条评估记录。
func (s *Service) Evaluate(flagKey, targetKey string, context map[string]string) (*EvaluationResult, error) {
	flagKey = strings.TrimSpace(flagKey)
	targetKey = strings.TrimSpace(targetKey)
	if flagKey == "" {
		return nil, model.NewValidationError("flag_key", "开关标识不能为空")
	}
	if targetKey == "" {
		return nil, model.NewValidationError("target_key", "目标标识不能为空")
	}
	flag, err := s.store.GetFlagByKey(flagKey)
	if err != nil {
		return nil, err
	}

	result := &EvaluationResult{FlagKey: flag.Key, FlagID: flag.ID, Enabled: false}
	ctxJSON, _ := json.Marshal(context)

	// 状态为 active 且已启用时才继续走灰度规则。
	if flag.Status == model.FlagActive && flag.Enabled {
		result.Enabled = true
		variant, matched := s.matchRules(flag.ID, targetKey, context)
		if matched && variant != nil {
			result.Variant = variant.Name
			result.VariantID = variant.ID
			result.Matched = true
		}
	}

	s.recordEvaluation(result, targetKey, string(ctxJSON))
	return result, nil
}

// matchRules 按优先级依次匹配灰度规则，返回命中的变体。
func (s *Service) matchRules(flagID, targetKey string, context map[string]string) (*model.Variant, bool) {
	rules := s.enabledRules(flagID)
	for _, r := range rules {
		if r.TargetGroupID != "" {
			g, err := s.store.GetTargetGroup(r.TargetGroupID)
			if err != nil || !g.MatchContext(context) {
				continue
			}
		}
		if hashTarget(targetKey)%100 < r.Percentage {
			variant, err := s.store.GetVariant(r.VariantID)
			if err != nil {
				continue
			}
			return variant, true
		}
	}
	return nil, false
}

// enabledRules 获取某开关下启用的规则并按优先级降序排序。
func (s *Service) enabledRules(flagID string) []*model.RolloutRule {
	all := s.store.ListRolloutRules()
	rules := make([]*model.RolloutRule, 0, len(all))
	for _, r := range all {
		if r.FlagID == flagID && r.Status == model.RolloutEnabled {
			rules = append(rules, r)
		}
	}
	sort.Slice(rules, func(i, j int) bool {
		if rules[i].Priority != rules[j].Priority {
			return rules[i].Priority > rules[j].Priority
		}
		return rules[i].CreatedAt.After(rules[j].CreatedAt)
	})
	return rules
}

func (s *Service) recordEvaluation(result *EvaluationResult, targetKey, ctx string) {
	e := &model.EvaluationRecord{
		ID:            idgen.HexN(8),
		FlagID:        result.FlagID,
		TargetKey:     targetKey,
		Context:       ctx,
		ResultVariant: result.Variant,
		Matched:       result.Matched,
		CreatedAt:     time.Now(),
	}
	_ = s.store.CreateEvaluationRecord(e)
}

// ListEvaluationRecords 分页列出评估记录。
func (s *Service) ListEvaluationRecords(filter model.EvaluationFilter, page, size int) ([]*model.EvaluationRecord, int, error) {
	all := s.store.ListEvaluationRecords()
	matched := make([]*model.EvaluationRecord, 0, len(all))
	for _, e := range all {
		if filter.Match(e) {
			matched = append(matched, e)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.EvaluationRecord{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// hashTarget 计算目标标识的稳定哈希（0-99）。
func hashTarget(key string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return int(h.Sum32() % 100)
}
