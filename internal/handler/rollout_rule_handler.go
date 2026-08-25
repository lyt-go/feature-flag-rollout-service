package handler

import (
	"net/http"

	"featureflag/internal/model"
	"featureflag/pkg/httpx"
)

func (s *Server) registerRolloutRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/rollout-rules", s.createRolloutRule)
	mux.HandleFunc("GET /api/rollout-rules", s.listRolloutRules)
	mux.HandleFunc("GET /api/rollout-rules/{id}", s.getRolloutRule)
	mux.HandleFunc("PUT /api/rollout-rules/{id}", s.updateRolloutRule)
	mux.HandleFunc("DELETE /api/rollout-rules/{id}", s.deleteRolloutRule)
}

type createRolloutRuleRequest struct {
	FlagID        string `json:"flag_id"`
	TargetGroupID string `json:"target_group_id"`
	VariantID     string `json:"variant_id"`
	Percentage    int    `json:"percentage"`
	Status        string `json:"status"`
	Priority      int    `json:"priority"`
}

func (s *Server) createRolloutRule(w http.ResponseWriter, r *http.Request) {
	var req createRolloutRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rule, err := s.svc.CreateRolloutRule(model.RolloutRule{
		FlagID: req.FlagID, TargetGroupID: req.TargetGroupID, VariantID: req.VariantID,
		Percentage: req.Percentage, Status: req.Status, Priority: req.Priority,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rule)
}

func (s *Server) listRolloutRules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RolloutRuleFilter{
		FlagID: r.URL.Query().Get("flag_id"),
		Status: r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListRolloutRules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRolloutRule(w http.ResponseWriter, r *http.Request) {
	rule, err := s.svc.GetRolloutRule(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rule)
}

func (s *Server) updateRolloutRule(w http.ResponseWriter, r *http.Request) {
	var req createRolloutRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rule, err := s.svc.UpdateRolloutRule(r.PathValue("id"), model.RolloutRule{
		TargetGroupID: req.TargetGroupID, VariantID: req.VariantID,
		Percentage: req.Percentage, Status: req.Status, Priority: req.Priority,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rule)
}

func (s *Server) deleteRolloutRule(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteRolloutRule(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
