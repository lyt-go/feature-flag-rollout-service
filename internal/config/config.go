// Package config 负责从环境变量加载服务配置。
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config 服务配置。
type Config struct {
	Addr           string
	MaxPageSize    int
	DefaultEnabled bool // 未找到开关时的默认返回
}

// Load 从环境变量加载配置。
func Load() *Config {
	cfg := &Config{
		Addr:           ":" + getenv("PORT", "8080"),
		MaxPageSize:    getenvInt("MAX_PAGE_SIZE", 100),
		DefaultEnabled: getenvBool("DEFAULT_ENABLED", false),
	}
	if v := os.Getenv("ADDR"); v != "" {
		cfg.Addr = v
	}
	return cfg
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return def
	}
	return n
}

func getenvBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

func (c *Config) String() string {
	return fmt.Sprintf("addr=%s max_page_size=%d default_enabled=%v", c.Addr, c.MaxPageSize, c.DefaultEnabled)
}
