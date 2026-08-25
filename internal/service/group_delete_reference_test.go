package service

import (
	"testing"

	"featureflag/internal/model"
)

func TestReferencedTargetGroupCannotBeDeleted(t *testing.T) {
	s := newEvaluatorService()
	f, _ := s.CreateFlag(model.Flag{Key: "region-rollout", Name: "Region rollout", Enabled: true})
	v, _ := s.CreateVariant(model.Variant{FlagID: f.ID, Name: "north"})
	g, _ := s.CreateTargetGroup(model.TargetGroup{Name: "North region", Rules: `{"region":"north"}`})
	_, _ = s.CreateRolloutRule(model.RolloutRule{FlagID: f.ID, VariantID: v.ID, TargetGroupID: g.ID, Percentage: 100})
	if err := s.DeleteTargetGroup(g.ID); err == nil {
		t.Fatal("referenced target group deletion must be rejected")
	}
	result, _ := s.Evaluate("region-rollout", "visitor-north", map[string]string{"region": "north"})
	if !result.Matched || result.Variant != "north" {
		t.Fatal("rejected target group deletion must preserve the regional rollout")
	}
}
