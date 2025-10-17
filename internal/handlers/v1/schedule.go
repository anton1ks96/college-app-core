package v1

import (
	"net/http"

	"github.com/anton1ks96/college-app-core/internal/domain"
	"github.com/anton1ks96/college-app-core/internal/services"
	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	scheduleService *services.ScheduleService
}

func NewScheduleHandler(svc *services.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: svc,
	}
}

func (h *ScheduleHandler) GetSchedule(c *gin.Context) {
	group := c.Query("group")
	subgroup := c.Query("subgroup")
	englishGroup := c.Query("english_group")
	profileSubgroup := c.Query("profile_subgroup")
	start := c.Query("start")
	end := c.Query("end")

	if group == "" || start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required query params"})
		return
	}

	events, err := h.scheduleService.GetSchedule(group, subgroup, englishGroup, profileSubgroup, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := domain.ScheduleResponse{Events: events}
	c.JSON(http.StatusOK, resp)
}

func (h *ScheduleHandler) GetClassDetails(c *gin.Context) {
	clid := c.Query("id")
	if clid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}

	details, err := h.scheduleService.GetClassDetails(clid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, details)
}
