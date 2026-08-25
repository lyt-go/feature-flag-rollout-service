package handler

import (
	"net/http"

	"featureflag/internal/model"
	"featureflag/pkg/httpx"
)

func (s *Server) registerEvaluationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/evaluate", s.evaluate)
	mux.HandleFunc("GET /api/evaluation-records", s.listEvaluationRecords)
}

type evaluateRequest struct {
	FlagKey   string            `json:"flag_key"`
	TargetKey string            `json:"target_key"`
	Context   map[string]string `json:"context"`
}

func (s *Server) evaluate(w http.ResponseWriter, r *http.Request) {
	var req evaluateRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if req.Context == nil {
		req.Context = map[string]string{}
	}
	result, err := s.svc.Evaluate(req.FlagKey, req.TargetKey, req.Context)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) listEvaluationRecords(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.EvaluationFilter{
		FlagID:        r.URL.Query().Get("flag_id"),
		TargetKey:     r.URL.Query().Get("target_key"),
		ResultVariant: r.URL.Query().Get("result_variant"),
	}
	items, total, err := s.svc.ListEvaluationRecords(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}
