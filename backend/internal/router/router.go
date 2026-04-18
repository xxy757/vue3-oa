package router

import (
	"oa-saas/internal/config"
	"oa-saas/internal/handler"
	"oa-saas/internal/middleware"
	"oa-saas/internal/pkg/cache"

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

	tenantHandler := handler.NewTenantHandler(db)
	r.POST("/api/v1/tenant/register", tenantHandler.Register)
	r.GET("/api/v1/plans", tenantHandler.ListPlans)

	tenantGroup := r.Group("/api/v1")
	tenantGroup.Use(tenantMw)
	{
		authHandler := handler.NewAuthHandler(db, cfg)
		tenantGroup.POST("/auth/login", authHandler.Login)
		tenantGroup.GET("/auth/info", authMw, authHandler.GetInfo)
		tenantGroup.PUT("/auth/password", authMw, authHandler.ChangePassword)
	}

	api := r.Group("/api/v1")
	api.Use(tenantMw, authMw)
	{
		userHandler := handler.NewUserHandler(db)
		api.GET("/user/list", userHandler.List)
		api.POST("/user", userHandler.Create)
		api.PUT("/user/:id", userHandler.Update)
		api.DELETE("/user/:id", userHandler.Delete)
		api.PUT("/user/:id/status", userHandler.UpdateStatus)
	}
	{
		deptHandler := handler.NewDeptHandler(db)
		api.GET("/dept/list", deptHandler.List)
		api.POST("/dept", deptHandler.Create)
		api.PUT("/dept/:id", deptHandler.Update)
		api.DELETE("/dept/:id", deptHandler.Delete)
	}
	{
		roleHandler := handler.NewRoleHandler(db)
		api.GET("/role/list", roleHandler.List)
		api.POST("/role", roleHandler.Create)
		api.PUT("/role/:id", roleHandler.Update)
		api.DELETE("/role/:id", roleHandler.Delete)
	}
	{
		approvalHandler := handler.NewApprovalHandler(db)
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
		noticeHandler := handler.NewNoticeHandler(db)
		api.GET("/notices", noticeHandler.List)
		api.GET("/notices/unread-count", noticeHandler.UnreadCount)
		api.GET("/notices/:id", noticeHandler.Detail)
		api.POST("/notices", noticeHandler.Create)
		api.POST("/notices/:id/read", noticeHandler.MarkRead)
	}
	{
		scheduleHandler := handler.NewScheduleHandler(db)
		api.GET("/schedules", scheduleHandler.List)
		api.GET("/schedules/week", scheduleHandler.WeekList)
		api.GET("/schedules/:id", scheduleHandler.Detail)
		api.POST("/schedules", scheduleHandler.Create)
		api.PUT("/schedules/:id", scheduleHandler.Update)
		api.DELETE("/schedules/:id", scheduleHandler.Delete)
	}
	{
		flowHandler := handler.NewFlowHandler(db)
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

	adminHandler := handler.NewAdminHandler(db)
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
