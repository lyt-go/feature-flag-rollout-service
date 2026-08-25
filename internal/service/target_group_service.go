package service

import (
	"fmt"
	"sort"
	"time"

	"featureflag/internal/model"
	"featureflag/internal/store"
	"featureflag/pkg/idgen"
)

// CreateTargetGroup 创建目标分组。
func (s *Service) CreateTargetGroup(input model.TargetGroup) (*model.TargetGroup, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	input.ID = idgen.HexN(8)
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.store.CreateTargetGroup(&input); err != nil {
		return nil, err
	}
	s.log.Infof("创建目标分组 %s", input.ID)
	return &input, nil
}

// GetTargetGroup 按 ID 获取目标分组。
func (s *Service) GetTargetGroup(id string) (*model.TargetGroup, error) {
	return s.store.GetTargetGroup(id)
}

// ListTargetGroups 分页列出目标分组。
func (s *Service) ListTargetGroups(filter model.TargetGroupFilter, page, size int) ([]*model.TargetGroup, int, error) {
	all := s.store.ListTargetGroups()
	matched := make([]*model.TargetGroup, 0, len(all))
	for _, g := range all {
		if filter.Match(g) {
			matched = append(matched, g)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.TargetGroup{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateTargetGroup 更新目标分组。
func (s *Service) UpdateTargetGroup(id string, input model.TargetGroup) (*model.TargetGroup, error) {
	existing, err := s.store.GetTargetGroup(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.Description = input.Description
	existing.Rules = input.Rules
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateTargetGroup(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteTargetGroup 删除目标分组。若仍有灰度规则引用该分组，则拒绝删除以保持引用完整性。
func (s *Service) DeleteTargetGroup(id string) error {
	// 确保分组存在，不存在的分组直接返回 NotFound。
	if _, err := s.store.GetTargetGroup(id); err != nil {
		return err
	}
	// 检查是否仍有灰度规则引用该分组；存在引用时删除会留下悬挂引用，
	// 并使评估跳过对应规则，因此返回冲突。
	if refs := s.countRulesByTargetGroup(id); refs > 0 {
		return fmt.Errorf("%w: 目标分组仍被 %d 条灰度规则引用", store.ErrConflict, refs)
	}
	if err := s.store.DeleteTargetGroup(id); err != nil {
		return err
	}
	s.log.Infof("删除目标分组 %s", id)
	return nil
}

// countRulesByTargetGroup 统计引用指定目标分组的灰度规则数量。
func (s *Service) countRulesByTargetGroup(groupID string) int {
	count := 0
	for _, r := range s.store.ListRolloutRules() {
		if r.TargetGroupID == groupID {
			count++
		}
	}
	return count
}
