package service

import (
	"sort"
	"time"

	"featureflag/internal/model"
	"featureflag/pkg/idgen"
)

// CreateRolloutRule 创建灰度规则，校验开关与变体存在且归属一致。
func (s *Service) CreateRolloutRule(input model.RolloutRule) (*model.RolloutRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetFlag(input.FlagID); err != nil {
		return nil, err
	}
	variant, err := s.store.GetVariant(input.VariantID)
	if err != nil {
		return nil, err
	}
	if variant.FlagID != input.FlagID {
		return nil, model.NewValidationError("variant_id", "变体不属于该开关")
	}
	if input.TargetGroupID != "" {
		if _, err := s.store.GetTargetGroup(input.TargetGroupID); err != nil {
			return nil, err
		}
	}
	input.ID = idgen.HexN(8)
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.store.CreateRolloutRule(&input); err != nil {
		return nil, err
	}
	s.log.Infof("创建灰度规则 %s", input.ID)
	return &input, nil
}

// GetRolloutRule 按 ID 获取灰度规则。
func (s *Service) GetRolloutRule(id string) (*model.RolloutRule, error) {
	return s.store.GetRolloutRule(id)
}

// ListRolloutRules 分页列出灰度规则。
func (s *Service) ListRolloutRules(filter model.RolloutRuleFilter, page, size int) ([]*model.RolloutRule, int, error) {
	all := s.store.ListRolloutRules()
	matched := make([]*model.RolloutRule, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		if matched[i].Priority != matched[j].Priority {
			return matched[i].Priority > matched[j].Priority
		}
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.RolloutRule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateRolloutRule 更新灰度规则。
// 先在副本上完成全部校验，确认无误后再写回存储：校验或关联校验失败时，
// 原规则保持不变继续生效，绝不让跨开关的变体/分组污染现有规则。
func (s *Service) UpdateRolloutRule(id string, input model.RolloutRule) (*model.RolloutRule, error) {
	existing, err := s.store.GetRolloutRule(id)
	if err != nil {
		return nil, err
	}

	// 在副本上应用变更，避免校验失败时把脏数据写回存储里的原对象。
	updated := *existing
	updated.TargetGroupID = input.TargetGroupID
	updated.VariantID = input.VariantID
	updated.Percentage = input.Percentage
	updated.Priority = input.Priority
	if input.Status != "" {
		updated.Status = input.Status
	}
	if err := updated.Validate(); err != nil {
		return nil, err
	}
	// 关联校验：变体必须属于该规则所属的开关（与创建逻辑一致），跨开关的变体一律拒绝。
	variant, err := s.store.GetVariant(updated.VariantID)
	if err != nil {
		return nil, err
	}
	if variant.FlagID != existing.FlagID {
		return nil, model.NewValidationError("variant_id", "变体不属于该开关")
	}
	if updated.TargetGroupID != "" {
		if _, err := s.store.GetTargetGroup(updated.TargetGroupID); err != nil {
			return nil, err
		}
	}

	updated.UpdatedAt = time.Now()
	if err := s.store.UpdateRolloutRule(&updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// DeleteRolloutRule 删除灰度规则。
func (s *Service) DeleteRolloutRule(id string) error {
	if err := s.store.DeleteRolloutRule(id); err != nil {
		return err
	}
	s.log.Infof("删除灰度规则 %s", id)
	return nil
}
