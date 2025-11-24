package v1

import (
	"net/http"

	"github.com/anton1ks96/college-app-core/internal/httpmw"
	"github.com/anton1ks96/college-app-core/internal/services"
	"github.com/anton1ks96/college-app-core/pkg/logger"
	"github.com/gin-gonic/gin"
)

type PerformanceHandler struct {
	performanceService *services.PerformanceService
}

func NewPerformanceHandler(svc *services.PerformanceService) *PerformanceHandler {
	return &PerformanceHandler{
		performanceService: svc,
	}
}

func (h *PerformanceHandler) GetSubjects(c *gin.Context) {
	login, _ := httpmw.GetUserID(c)

	subjects, err := h.performanceService.GetSubjects(login)
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Str("login", login).
			Msg("failed to get performance subjects")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subjects)
}

type scoreRequest struct {
	SuID      string `json:"SuID"`
	Datastart string `json:"datastart"`
	Dataend   string `json:"dataend"`
}

func (h *PerformanceHandler) GetScore(c *gin.Context) {
	var req scoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.SuID == "" || req.Datastart == "" || req.Dataend == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required fields: SuID, datastart, dataend"})
		return
	}

	login, _ := httpmw.GetUserID(c)

	scores, err := h.performanceService.GetScore(login, req.SuID, req.Datastart, req.Dataend)
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Str("login", login).
			Str("suID", req.SuID).
			Str("datastart", req.Datastart).
			Str("dataend", req.Dataend).
			Msg("failed to get performance score")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, scores)
}
