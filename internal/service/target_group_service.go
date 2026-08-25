package service

import (
	"sort"
	"time"

	"featureflag/internal/model"
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

// DeleteTargetGroup 删除目标分组。
func (s *Service) DeleteTargetGroup(id string) error {
	if err := s.store.DeleteTargetGroup(id); err != nil {
		return err
	}
	s.log.Infof("删除目标分组 %s", id)
	return nil
}
