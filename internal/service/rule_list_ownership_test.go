package service

import (
	"testing"

	"featureflag/internal/model"
)

func TestRulePreviewDoesNotAlterCompactAudience(t *testing.T) {
	s := newEvaluatorService()
	f, _ := s.CreateFlag(model.Flag{Key: "toolbar-rollout", Name: "Toolbar rollout", Enabled: true})
	v, _ := s.CreateVariant(model.Variant{FlagID: f.ID, Name: "compact"})
	_, _ = s.CreateRolloutRule(model.RolloutRule{FlagID: f.ID, VariantID: v.ID, Percentage: 100})
	listed, _, _ := s.ListRolloutRules(model.RolloutRuleFilter{FlagID: f.ID}, 1, 20)
	listed[0].Percentage = 0
	listed[0].Status = model.RolloutDisabled
	result, _ := s.Evaluate("toolbar-rollout", "visitor-a", nil)
	if !result.Matched || result.Variant != "compact" {
		t.Fatal("editing a rollout rule list result must not change later evaluations")
	}
}
