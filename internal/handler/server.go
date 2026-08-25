// Package handler 实现 HTTP 处理器层。
package handler

import (
	"errors"
	"net/http"
	"runtime/debug"
	"time"

	"featureflag/internal/config"
	"featureflag/internal/model"
	"featureflag/internal/service"
	"featureflag/internal/store"
	"featureflag/pkg/httpx"
	"featureflag/pkg/logger"
)

// Server 持有服务层并暴露路由。
type Server struct {
	svc *service.Service
	log *logger.Logger
	cfg *config.Config
}

// NewServer 构造处理器服务器。
func NewServer(svc *service.Service, log *logger.Logger, cfg *config.Config) *Server {
	return &Server{svc: svc, log: log, cfg: cfg}
}

// Routes 组装全部路由与中间件。
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	s.registerFlagRoutes(mux)
	s.registerVariantRoutes(mux)
	s.registerTargetGroupRoutes(mux)
	s.registerRolloutRuleRoutes(mux)
	s.registerChangeLogRoutes(mux)
	s.registerEvaluationRoutes(mux)
	s.registerStatsRoutes(mux)
	return s.loggingMiddleware(s.recoveryMiddleware(mux))
}

func (s *Server) maxPageSize() int {
	if s.cfg != nil && s.cfg.MaxPageSize > 0 {
		return s.cfg.MaxPageSize
	}
	return 100
}

func (s *Server) operator(r *http.Request) string {
	if v := r.URL.Query().Get("operator"); v != "" {
		return v
	}
	return "admin"
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.log.Infof("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				s.log.Errorf("panic: %v\n%s", rec, debug.Stack())
				httpx.InternalError(w, "服务器内部错误")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case model.IsValidationError(err):
		httpx.BadRequest(w, err.Error())
	case errors.Is(err, store.ErrNotFound):
		httpx.NotFound(w, err.Error())
	case errors.Is(err, store.ErrConflict):
		httpx.Conflict(w, err.Error())
	default:
		httpx.InternalError(w, err.Error())
	}
}
