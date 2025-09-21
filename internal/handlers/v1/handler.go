package v1

import (
	"github.com/anton1ks96/college-app-core/internal/config"
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

func (h *Handler) Init(api *gin.RouterGroup) {
}
