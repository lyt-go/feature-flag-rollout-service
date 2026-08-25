package service

import (
	"sort"
	"time"

	"featureflag/internal/model"
	"featureflag/pkg/idgen"
)

// CreateFlag 创建特性开关并记录审计日志。
func (s *Service) CreateFlag(input model.Flag) (*model.Flag, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetFlagByKey(input.Key); err == nil {
		return nil, model.NewValidationError("key", "开关标识已存在")
	}
	input.ID = idgen.HexN(8)
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if input.Status == "" {
		input.Status = model.FlagActive
	}
	if err := s.store.CreateFlag(&input); err != nil {
		return nil, err
	}
	s.recordChange(input.ID, "system", model.ActionCreate, "创建开关 "+input.Key)
	s.log.Infof("创建开关 %s (%s)", input.Key, input.ID)
	return &input, nil
}

// GetFlag 按 ID 获取开关。
func (s *Service) GetFlag(id string) (*model.Flag, error) {
	return s.store.GetFlag(id)
}

// GetFlagByKey 按标识获取开关。
func (s *Service) GetFlagByKey(key string) (*model.Flag, error) {
	return s.store.GetFlagByKey(key)
}

// ListFlags 分页列出开关。
func (s *Service) ListFlags(filter model.FlagFilter, page, size int) ([]*model.Flag, int, error) {
	all := s.store.ListFlags()
	matched := make([]*model.Flag, 0, len(all))
	for _, f := range all {
		if filter.Match(f) {
			matched = append(matched, f)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Flag{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateFlag 更新开关可编辑字段并记录审计日志。
func (s *Service) UpdateFlag(id, operator string, input model.Flag) (*model.Flag, error) {
	existing, err := s.store.GetFlag(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.Description = input.Description
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateFlag(existing); err != nil {
		return nil, err
	}
	s.recordChange(id, operator, model.ActionUpdate, "更新开关 "+existing.Key)
	return existing, nil
}

// ToggleFlag 翻转开关的 Enabled 状态。
func (s *Service) ToggleFlag(id, operator string) (*model.Flag, error) {
	existing, err := s.store.GetFlag(id)
	if err != nil {
		return nil, err
	}
	existing.Enabled = !existing.Enabled
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateFlag(existing); err != nil {
		return nil, err
	}
	s.recordChange(id, operator, model.ActionToggle, "切换启用状态")
	return existing, nil
}

// ChangeFlagStatus 变更开关状态（含状态机校验）。
func (s *Service) ChangeFlagStatus(id, status, operator string) (*model.Flag, error) {
	existing, err := s.store.GetFlag(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionFlag(existing.Status, status) {
		return nil, model.NewValidationError("status", "非法状态流转: "+existing.Status+" -> "+status)
	}
	existing.Status = status
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateFlag(existing); err != nil {
		return nil, err
	}
	action := statusToAction(status)
	s.recordChange(id, operator, action, "变更状态为 "+status)
	s.log.Infof("开关 %s 状态变更为 %s", id, status)
	return existing, nil
}

func statusToAction(status string) string {
	switch status {
	case model.FlagPaused:
		return model.ActionPause
	case model.FlagArchived:
		return model.ActionArchive
	case model.FlagActive:
		return model.ActionResume
	default:
		return model.ActionUpdate
	}
}

// DeleteFlag 删除开关。
func (s *Service) DeleteFlag(id string) error {
	if err := s.store.DeleteFlag(id); err != nil {
		return err
	}
	s.log.Infof("删除开关 %s", id)
	return nil
}

// recordChange 写入一条变更审计日志。
func (s *Service) recordChange(flagID, operator, action, detail string) {
	c := &model.ChangeLog{
		ID:        idgen.HexN(8),
		FlagID:    flagID,
		Operator:  operator,
		Action:    action,
		Detail:    detail,
		CreatedAt: time.Now(),
	}
	_ = s.store.CreateChangeLog(c)
}
