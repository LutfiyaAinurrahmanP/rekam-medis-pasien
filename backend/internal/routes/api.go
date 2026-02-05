package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupAPIRouter sets up all API routes by calling individual router setups
func SetupAPIRouter(rg *gin.RouterGroup, cfg *RouteConfig) {
	// Setup Auth routes
	SetupAuthRouter(rg, cfg)

	// Setup Users routes
	SetupUsersRouter(rg, cfg)

	// Setup Departments routes
	if cfg.DepartmentHandler != nil {
		SetupDepartmentsRouter(rg, cfg, cfg.DepartmentHandler)
	}

	// Setup Patients routes
	if cfg.PatientHandler != nil {
		SetupPatientsRouter(rg, cfg, cfg.PatientHandler)
	}
}
