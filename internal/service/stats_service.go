package service

import (
	"featureflag/internal/model"
)

// FlagStat 单个开关的统计。
type FlagStat struct {
	Key             string `json:"key"`
	Name            string `json:"name"`
	Status          string `json:"status"`
	Enabled         bool   `json:"enabled"`
	VariantCount    int    `json:"variant_count"`
	RuleCount       int    `json:"rule_count"`
	EvaluationCount int    `json:"evaluation_count"`
	MatchedCount    int    `json:"matched_count"`
}

// FlagStats 开关全局统计。
type FlagStats struct {
	FlagCount       int        `json:"flag_count"`
	ActiveCount     int        `json:"active_count"`
	PausedCount     int        `json:"paused_count"`
	ArchivedCount   int        `json:"archived_count"`
	EnabledCount    int        `json:"enabled_count"`
	VariantCount    int        `json:"variant_count"`
	GroupCount      int        `json:"group_count"`
	RuleCount       int        `json:"rule_count"`
	EvaluationCount int        `json:"evaluation_count"`
	MatchedCount    int        `json:"matched_count"`
	ByFlag          []FlagStat `json:"by_flag"`
}

// Overview 汇总开关全局统计。
func (s *Service) Overview() (*FlagStats, error) {
	stats := &FlagStats{ByFlag: []FlagStat{}}
	for _, f := range s.store.ListFlags() {
		stats.FlagCount++
		switch f.Status {
		case model.FlagActive:
			stats.ActiveCount++
		case model.FlagPaused:
			stats.PausedCount++
		case model.FlagArchived:
			stats.ArchivedCount++
		}
		if f.Enabled {
			stats.EnabledCount++
		}
		stats.ByFlag = append(stats.ByFlag, FlagStat{Key: f.Key, Name: f.Name, Status: f.Status, Enabled: f.Enabled})
	}
	stats.VariantCount = len(s.store.ListVariants())
	stats.GroupCount = len(s.store.ListTargetGroups())
	stats.RuleCount = len(s.store.ListRolloutRules())
	for _, e := range s.store.ListEvaluationRecords() {
		stats.EvaluationCount++
		if e.Matched {
			stats.MatchedCount++
		}
	}
	// 汇总每个开关的变体/规则/评估计数。
	for i := range stats.ByFlag {
		fs := &stats.ByFlag[i]
		for _, v := range s.store.ListVariants() {
			if flagIDByKey(v.FlagID, s) == fs.Key {
				fs.VariantCount++
			}
		}
		for _, r := range s.store.ListRolloutRules() {
			if flagIDByKey(r.FlagID, s) == fs.Key {
				fs.RuleCount++
			}
		}
		for _, e := range s.store.ListEvaluationRecords() {
			if flagIDByKey(e.FlagID, s) == fs.Key {
				fs.EvaluationCount++
				if e.Matched {
					fs.MatchedCount++
				}
			}
		}
	}
	return stats, nil
}

// flagIDByKey 将 flag ID 映射回 key（用于聚合），未知时返回空串。
func flagIDByKey(flagID string, s *Service) string {
	f, err := s.store.GetFlag(flagID)
	if err != nil {
		return ""
	}
	return f.Key
}
