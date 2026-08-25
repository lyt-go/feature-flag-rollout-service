package service

import (
	"sort"

	"featureflag/internal/model"
)

// GetChangeLog 按 ID 获取变更记录。
func (s *Service) GetChangeLog(id string) (*model.ChangeLog, error) {
	return s.store.GetChangeLog(id)
}

// ListChangeLogs 分页列出变更记录。
func (s *Service) ListChangeLogs(filter model.ChangeLogFilter, page, size int) ([]*model.ChangeLog, int, error) {
	all := s.store.ListChangeLogs()
	matched := make([]*model.ChangeLog, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.ChangeLog{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}
