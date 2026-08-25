package service

import (
	"featureflag/internal/config"
	"featureflag/internal/store"
	"featureflag/pkg/logger"
)

func newEvaluatorService() *Service {
	return New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 100})
}
