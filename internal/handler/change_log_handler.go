package handler

import (
	"net/http"

	"featureflag/internal/model"
	"featureflag/pkg/httpx"
)

func (s *Server) registerChangeLogRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/change-logs", s.listChangeLogs)
	mux.HandleFunc("GET /api/change-logs/{id}", s.getChangeLog)
}

func (s *Server) listChangeLogs(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ChangeLogFilter{
		FlagID: r.URL.Query().Get("flag_id"),
		Action: r.URL.Query().Get("action"),
	}
	items, total, err := s.svc.ListChangeLogs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getChangeLog(w http.ResponseWriter, r *http.Request) {
	c, err := s.svc.GetChangeLog(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}
