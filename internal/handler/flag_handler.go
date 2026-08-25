package handler

import (
	"net/http"

	"featureflag/internal/model"
	"featureflag/pkg/httpx"
)

func (s *Server) registerFlagRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/flags", s.createFlag)
	mux.HandleFunc("GET /api/flags", s.listFlags)
	mux.HandleFunc("GET /api/flags/{id}", s.getFlag)
	mux.HandleFunc("PUT /api/flags/{id}", s.updateFlag)
	mux.HandleFunc("DELETE /api/flags/{id}", s.deleteFlag)
	mux.HandleFunc("POST /api/flags/{id}/toggle", s.toggleFlag)
	mux.HandleFunc("POST /api/flags/{id}/status", s.changeFlagStatus)
}

type createFlagRequest struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

func (s *Server) createFlag(w http.ResponseWriter, r *http.Request) {
	var req createFlagRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	f, err := s.svc.CreateFlag(model.Flag{
		Key: req.Key, Name: req.Name, Description: req.Description, Enabled: req.Enabled,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, f)
}

func (s *Server) listFlags(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.FlagFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListFlags(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getFlag(w http.ResponseWriter, r *http.Request) {
	f, err := s.svc.GetFlag(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}

func (s *Server) updateFlag(w http.ResponseWriter, r *http.Request) {
	var req createFlagRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	f, err := s.svc.UpdateFlag(r.PathValue("id"), s.operator(r), model.Flag{
		Name: req.Name, Description: req.Description,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}

func (s *Server) deleteFlag(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteFlag(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) toggleFlag(w http.ResponseWriter, r *http.Request) {
	f, err := s.svc.ToggleFlag(r.PathValue("id"), s.operator(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}

type changeStatusRequest struct {
	Status string `json:"status"`
}

func (s *Server) changeFlagStatus(w http.ResponseWriter, r *http.Request) {
	var req changeStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	f, err := s.svc.ChangeFlagStatus(r.PathValue("id"), req.Status, s.operator(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}
