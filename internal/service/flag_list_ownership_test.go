package service

import (
	"testing"

	"featureflag/internal/model"
)

func TestReturnedFlagMutationLeavesNavigationEnabled(t *testing.T) {
	s := newEvaluatorService()
	_, _ = s.CreateFlag(model.Flag{Key: "navigation-rollout", Name: "Navigation rollout", Enabled: true})
	listed, _, _ := s.ListFlags(model.FlagFilter{}, 1, 20)
	listed[0].Enabled = false
	listed[0].Status = model.FlagArchived
	result, _ := s.Evaluate("navigation-rollout", "visitor-a", nil)
	if !result.Enabled {
		t.Fatal("editing a flag list result must not disable the stored rollout")
	}
}
