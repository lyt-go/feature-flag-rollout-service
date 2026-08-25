package service

import (
	"testing"

	"featureflag/internal/model"
)

func TestRuleUpdateRejectsForeignVariant(t *testing.T) {
	s := newEvaluatorService()
	f1, _ := s.CreateFlag(model.Flag{Key: "editor-rollout", Name: "Editor rollout", Enabled: true})
	v1, _ := s.CreateVariant(model.Variant{FlagID: f1.ID, Name: "new-editor"})
	rule, _ := s.CreateRolloutRule(model.RolloutRule{FlagID: f1.ID, VariantID: v1.ID, Percentage: 100})
	f2, _ := s.CreateFlag(model.Flag{Key: "search-rollout", Name: "Search rollout", Enabled: true})
	v2, _ := s.CreateVariant(model.Variant{FlagID: f2.ID, Name: "new-search"})
	_, err := s.UpdateRolloutRule(rule.ID, model.RolloutRule{VariantID: v2.ID, Percentage: 100})
	if err == nil {
		t.Fatal("rule update must reject a variant owned by another flag")
	}
	result, _ := s.Evaluate("editor-rollout", "visitor-a", nil)
	if !result.Matched || result.Variant != "new-editor" {
		t.Fatal("rejected foreign variant update must preserve the editor rollout result")
	}
}
