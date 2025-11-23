package v1

import (
	"github.com/anton1ks96/college-app-core/internal/config"
	"github.com/anton1ks96/college-app-core/internal/httpmw"
	"github.com/anton1ks96/college-app-core/internal/repository"
	"github.com/anton1ks96/college-app-core/internal/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	cfg        *config.Config
	schedule   *ScheduleHandler
	attendance *AttendanceHandler
	auth       gin.HandlerFunc
}

func NewHandler(cfg *config.Config) *Handler {
	portalRepo := repository.NewPortalRepository(
		cfg.Portal.URL,
		cfg.Portal.AttendanceURL,
	)
	scheduleService := services.NewScheduleService(portalRepo)
	scheduleHandler := NewScheduleHandler(scheduleService)

	attendanceService := services.NewAttendanceService(portalRepo)
	attendanceHandler := NewAttendanceHandler(attendanceService)

	authMiddleware := httpmw.NewAuthMiddleware(cfg.Auth.ServiceURL, cfg.Auth.Timeout)

	return &Handler{
		cfg:        cfg,
		schedule:   scheduleHandler,
		attendance: attendanceHandler,
		auth:       authMiddleware.ValidateToken(),
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	api.GET("/schedule", h.schedule.GetSchedule)
	api.GET("/classdetails", h.schedule.GetClassDetails)
	api.GET("/attendance", h.auth, h.attendance.GetAttendance)
}
