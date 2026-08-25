package handler

import (
	"net/http"

	"featureflag/internal/model"
	"featureflag/pkg/httpx"
)

func (s *Server) registerVariantRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/variants", s.createVariant)
	mux.HandleFunc("GET /api/variants", s.listVariants)
	mux.HandleFunc("GET /api/variants/{id}", s.getVariant)
	mux.HandleFunc("PUT /api/variants/{id}", s.updateVariant)
	mux.HandleFunc("DELETE /api/variants/{id}", s.deleteVariant)
}

type createVariantRequest struct {
	FlagID  string `json:"flag_id"`
	Name    string `json:"name"`
	Payload string `json:"payload"`
	Weight  int    `json:"weight"`
}

func (s *Server) createVariant(w http.ResponseWriter, r *http.Request) {
	var req createVariantRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	v, err := s.svc.CreateVariant(model.Variant{
		FlagID: req.FlagID, Name: req.Name, Payload: req.Payload, Weight: req.Weight,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, v)
}

func (s *Server) listVariants(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.VariantFilter{
		FlagID:  r.URL.Query().Get("flag_id"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListVariants(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getVariant(w http.ResponseWriter, r *http.Request) {
	v, err := s.svc.GetVariant(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, v)
}

func (s *Server) updateVariant(w http.ResponseWriter, r *http.Request) {
	var req createVariantRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	v, err := s.svc.UpdateVariant(r.PathValue("id"), model.Variant{
		Name: req.Name, Payload: req.Payload, Weight: req.Weight,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, v)
}

func (s *Server) deleteVariant(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteVariant(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
