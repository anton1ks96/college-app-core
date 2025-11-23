package v1

import (
	"net/http"

	"github.com/anton1ks96/college-app-core/internal/httpmw"
	"github.com/anton1ks96/college-app-core/internal/services"
	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	attendanceService *services.AttendanceService
}

func NewAttendanceHandler(svc *services.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{
		attendanceService: svc,
	}
}

func (h *AttendanceHandler) GetAttendance(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")

	if start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required query params: start and end"})
		return
	}

	login, _ := httpmw.GetUserID(c)

	records, err := h.attendanceService.GetAttendance(login, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}
