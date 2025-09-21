package handlers

import (
	"github.com/anton1ks96/college-app-core/internal/config"
	v1 "github.com/anton1ks96/college-app-core/internal/handlers/v1"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	cfg *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		cfg: cfg,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.New()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	router.GET("/health", h.healthCheck)
	router.GET("/ready", h.readinessCheck)

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	api := router.Group("/api")

	v1Handler := v1.NewHandler(h.cfg)
	v1Group := api.Group("/v1")

	v1Handler.Init(v1Group)
}

func (h *Handler) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "OK",
		"service": "college-app-core",
	})
}

func (h *Handler) readinessCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"ready":   true,
		"service": "college-app-core",
	})
}
