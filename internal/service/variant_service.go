package service

import (
	"sort"
	"time"

	"featureflag/internal/model"
	"featureflag/pkg/idgen"
)

// CreateVariant 创建变体，校验所属开关存在。
func (s *Service) CreateVariant(input model.Variant) (*model.Variant, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetFlag(input.FlagID); err != nil {
		return nil, err
	}
	input.ID = idgen.HexN(8)
	input.CreatedAt = time.Now()
	if err := s.store.CreateVariant(&input); err != nil {
		return nil, err
	}
	s.log.Infof("创建变体 %s", input.ID)
	return &input, nil
}

// GetVariant 按 ID 获取变体。
func (s *Service) GetVariant(id string) (*model.Variant, error) {
	return s.store.GetVariant(id)
}

// ListVariants 分页列出变体。
func (s *Service) ListVariants(filter model.VariantFilter, page, size int) ([]*model.Variant, int, error) {
	all := s.store.ListVariants()
	matched := make([]*model.Variant, 0, len(all))
	for _, v := range all {
		if filter.Match(v) {
			matched = append(matched, v)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Variant{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateVariant 更新变体。
func (s *Service) UpdateVariant(id string, input model.Variant) (*model.Variant, error) {
	existing, err := s.store.GetVariant(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.Payload = input.Payload
	existing.Weight = input.Weight
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateVariant(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteVariant 删除变体。
func (s *Service) DeleteVariant(id string) error {
	if err := s.store.DeleteVariant(id); err != nil {
		return err
	}
	s.log.Infof("删除变体 %s", id)
	return nil
}
