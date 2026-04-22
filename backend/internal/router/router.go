package router

import (
	"oa-saas/internal/config"
	"oa-saas/internal/handler"
	"oa-saas/internal/middleware"
	"oa-saas/internal/pkg/cache"
	"oa-saas/internal/repository"
	"oa-saas/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, c cache.Cache, cfg *config.Config) *gin.Engine {
	if cfg.Server.Mode != "" {
		gin.SetMode(cfg.Server.Mode)
	}

	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	tenantMw := middleware.Tenant(db)
	authMw := middleware.Auth(cfg.JWT.Secret)

	// Repositories
	planRepo := repository.NewPlanRepo(db)
	userRepo := repository.NewUserRepo(db)
	deptRepo := repository.NewDeptRepo(db)
	roleRepo := repository.NewRoleRepo(db)
	tenantRepo := repository.NewTenantRepo(db)
	approvalRepo := repository.NewApprovalRepo(db)
	flowRepo := repository.NewFlowRepo(db)
	noticeRepo := repository.NewNoticeRepo(db)
	scheduleRepo := repository.NewScheduleRepo(db)

	// Services
	authService := service.NewAuthService(userRepo, deptRepo, roleRepo, tenantRepo, planRepo, cfg.JWT.Secret, cfg.JWT.ExpireHours)
	userService := service.NewUserService(userRepo, tenantRepo, deptRepo, roleRepo, db)
	deptService := service.NewDeptService(deptRepo)
	roleService := service.NewRoleService(roleRepo)
	tenantService := service.NewTenantService(tenantRepo, roleRepo, userRepo, planRepo, db)
	approvalService := service.NewApprovalService(approvalRepo, flowRepo, db)
	flowService := service.NewFlowService(flowRepo)
	noticeService := service.NewNoticeService(noticeRepo, userRepo)
	scheduleService := service.NewScheduleService(scheduleRepo, db)
	adminService := service.NewAdminService(tenantRepo, userRepo, planRepo, db)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	deptHandler := handler.NewDeptHandler(deptService)
	roleHandler := handler.NewRoleHandler(roleService)
	tenantHandler := handler.NewTenantHandler(tenantService)
	approvalHandler := handler.NewApprovalHandler(approvalService)
	flowHandler := handler.NewFlowHandler(flowService)
	noticeHandler := handler.NewNoticeHandler(noticeService)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)
	adminHandler := handler.NewAdminHandler(adminService)

	// Public routes (no auth, no tenant)
	r.POST("/api/v1/tenant/register", tenantHandler.Register)
	r.GET("/api/v1/plans", tenantHandler.ListPlans)

	// Tenant-resolved routes (tenant middleware, some with auth)
	tenantGroup := r.Group("/api/v1")
	tenantGroup.Use(tenantMw)
	{
		tenantGroup.POST("/auth/login", authHandler.Login)
		tenantGroup.GET("/auth/info", authMw, authHandler.GetInfo)
		tenantGroup.PUT("/auth/password", authMw, authHandler.ChangePassword)
	}

	// Authenticated routes
	api := r.Group("/api/v1")
	api.Use(tenantMw, authMw)
	{
		api.GET("/user/list", userHandler.List)
		api.POST("/user", userHandler.Create)
		api.PUT("/user/:id", userHandler.Update)
		api.DELETE("/user/:id", userHandler.Delete)
		api.PUT("/user/:id/status", userHandler.UpdateStatus)
	}
	{
		api.GET("/dept/list", deptHandler.List)
		api.POST("/dept", deptHandler.Create)
		api.PUT("/dept/:id", deptHandler.Update)
		api.DELETE("/dept/:id", deptHandler.Delete)
	}
	{
		api.GET("/role/list", roleHandler.List)
		api.POST("/role", roleHandler.Create)
		api.PUT("/role/:id", roleHandler.Update)
		api.DELETE("/role/:id", roleHandler.Delete)
	}
	{
		api.POST("/approvals", approvalHandler.Create)
		api.GET("/approvals/my", approvalHandler.MyList)
		api.GET("/approvals/pending", approvalHandler.PendingList)
		api.GET("/approvals/done", approvalHandler.DoneList)
		api.GET("/approvals/stats", approvalHandler.Stats)
		api.GET("/approvals/:id", approvalHandler.Detail)
		api.POST("/approvals/:id/action", approvalHandler.Action)
		api.POST("/approvals/:id/withdraw", approvalHandler.Withdraw)
	}
	{
		api.GET("/notices", noticeHandler.List)
		api.GET("/notices/unread-count", noticeHandler.UnreadCount)
		api.GET("/notices/:id", noticeHandler.Detail)
		api.POST("/notices", noticeHandler.Create)
		api.POST("/notices/:id/read", noticeHandler.MarkRead)
	}
	{
		api.GET("/schedules", scheduleHandler.List)
		api.GET("/schedules/week", scheduleHandler.WeekList)
		api.GET("/schedules/:id", scheduleHandler.Detail)
		api.POST("/schedules", scheduleHandler.Create)
		api.PUT("/schedules/:id", scheduleHandler.Update)
		api.DELETE("/schedules/:id", scheduleHandler.Delete)
	}
	{
		api.GET("/flows", flowHandler.List)
		api.POST("/flows", flowHandler.Create)
		api.PUT("/flows/:id", flowHandler.Update)
		api.DELETE("/flows/:id", flowHandler.Delete)
	}
	{
		api.GET("/tenant/info", tenantHandler.GetInfo)
		api.PUT("/tenant/info", tenantHandler.UpdateInfo)
		api.POST("/tenant/plan/upgrade", tenantHandler.UpgradePlan)
		api.GET("/tenant/invoices", tenantHandler.ListInvoices)
	}

	// Admin routes
	admin := r.Group("/api/v1/admin")
	admin.Use(authMw)
	{
		admin.GET("/dashboard", adminHandler.Dashboard)
		admin.GET("/tenants", adminHandler.ListTenants)
		admin.POST("/tenants", adminHandler.CreateTenant)
		admin.PUT("/tenants/:id", adminHandler.UpdateTenant)
		admin.PUT("/tenants/:id/activate", adminHandler.ActivateTenant)
		admin.PUT("/tenants/:id/suspend", adminHandler.SuspendTenant)
		admin.GET("/plans", adminHandler.ListPlans)
		admin.POST("/plans", adminHandler.CreatePlan)
		admin.PUT("/plans/:id", adminHandler.UpdatePlan)
	}

	return r
}
