package rest

import (
	"net/http"
	"runtime"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/patrickmn/go-cache"
)

type HealthHandler struct {
	cfg   *config.Schema
	cache *cache.Cache
}

type HealthResponse struct {
	Status    string         `json:"status"`
	Timestamp time.Time      `json:"timestamp"`
	Version   string         `json:"version"`
	Uptime    time.Duration  `json:"uptime"`
	Memory    MemoryStats    `json:"memory"`
	Cache     CacheStats     `json:"cache"`
	Database  DatabaseStatus `json:"database"`
}

type MemoryStats struct {
	Allocated     uint64 `json:"allocated_mb"`
	TotalAlloc    uint64 `json:"total_alloc_mb"`
	Sys           uint64 `json:"sys_mb"`
	NumGC         uint32 `json:"num_gc"`
	NumGoroutines int    `json:"num_goroutines"`
}

type CacheStats struct {
	ItemCount int `json:"item_count"`
}

type DatabaseStatus struct {
	Connected bool   `json:"connected"`
	Error     string `json:"error,omitempty"`
}

var startTime = time.Now()

func NewHealthHandler(cfg *config.Schema, cache *cache.Cache) *HealthHandler {
	return &HealthHandler{
		cfg:   cfg,
		cache: cache,
	}
}

func (h *HealthHandler) HealthCheck(c echo.Context) error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   h.cfg.Version,
		Uptime:    time.Since(startTime),
		Memory: MemoryStats{
			Allocated:     bToMb(m.Alloc),
			TotalAlloc:    bToMb(m.TotalAlloc),
			Sys:           bToMb(m.Sys),
			NumGC:         m.NumGC,
			NumGoroutines: runtime.NumGoroutine(),
		},
		Cache: CacheStats{
			ItemCount: h.cache.ItemCount(),
		},
		Database: DatabaseStatus{
			Connected: true, // TODO: Add actual DB health check
		},
	}

	// Check for potential memory leaks
	if runtime.NumGoroutine() > 1000 {
		response.Status = "warning"
	}

	return c.JSON(http.StatusOK, response)
}

func (h *HealthHandler) RegisterRoutes(g *echo.Group) {
	g.GET("/health", h.HealthCheck)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
