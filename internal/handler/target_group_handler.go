package handler

import (
	"net/http"

	"featureflag/internal/model"
	"featureflag/pkg/httpx"
)

func (s *Server) registerTargetGroupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/target-groups", s.createTargetGroup)
	mux.HandleFunc("GET /api/target-groups", s.listTargetGroups)
	mux.HandleFunc("GET /api/target-groups/{id}", s.getTargetGroup)
	mux.HandleFunc("PUT /api/target-groups/{id}", s.updateTargetGroup)
	mux.HandleFunc("DELETE /api/target-groups/{id}", s.deleteTargetGroup)
}

type createTargetGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Rules       string `json:"rules"`
}

func (s *Server) createTargetGroup(w http.ResponseWriter, r *http.Request) {
	var req createTargetGroupRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	g, err := s.svc.CreateTargetGroup(model.TargetGroup{
		Name: req.Name, Description: req.Description, Rules: req.Rules,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, g)
}

func (s *Server) listTargetGroups(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.TargetGroupFilter{
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListTargetGroups(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getTargetGroup(w http.ResponseWriter, r *http.Request) {
	g, err := s.svc.GetTargetGroup(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, g)
}

func (s *Server) updateTargetGroup(w http.ResponseWriter, r *http.Request) {
	var req createTargetGroupRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	g, err := s.svc.UpdateTargetGroup(r.PathValue("id"), model.TargetGroup{
		Name: req.Name, Description: req.Description, Rules: req.Rules,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, g)
}

func (s *Server) deleteTargetGroup(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteTargetGroup(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
