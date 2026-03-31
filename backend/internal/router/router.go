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

	api := r.Group("/api/v1")
	authMw := middleware.Auth(cfg.JWT.Secret)
	{
		authHandler := handler.NewAuthHandler(db, cfg)
		api.POST("/auth/login", authHandler.Login)
		api.GET("/auth/info", authMw, authHandler.GetInfo)
	}
	{
		userHandler := handler.NewUserHandler(db)
		api.GET("/user/list", authMw, userHandler.List)
		api.POST("/user", authMw, userHandler.Create)
		api.PUT("/user/:id", authMw, userHandler.Update)
		api.DELETE("/user/:id", authMw, userHandler.Delete)
		api.PUT("/user/:id/status", authMw, userHandler.UpdateStatus)
	}
	{
		deptHandler := handler.NewDeptHandler(db)
		api.GET("/dept/list", authMw, deptHandler.List)
		api.POST("/dept", authMw, deptHandler.Create)
		api.PUT("/dept/:id", authMw, deptHandler.Update)
		api.DELETE("/dept/:id", authMw, deptHandler.Delete)
	}
	{
		roleHandler := handler.NewRoleHandler(db)
		api.GET("/role/list", authMw, roleHandler.List)
		api.POST("/role", authMw, roleHandler.Create)
		api.PUT("/role/:id", authMw, roleHandler.Update)
		api.DELETE("/role/:id", authMw, roleHandler.Delete)
	}
	{
		approvalHandler := handler.NewApprovalHandler(db)
		api.POST("/approvals", authMw, approvalHandler.Create)
		api.GET("/approvals/my", authMw, approvalHandler.MyList)
		api.GET("/approvals/pending", authMw, approvalHandler.PendingList)
		api.GET("/approvals/done", authMw, approvalHandler.DoneList)
		api.GET("/approvals/stats", authMw, approvalHandler.Stats)
		api.GET("/approvals/:id", authMw, approvalHandler.Detail)
		api.POST("/approvals/:id/action", authMw, approvalHandler.Action)
		api.POST("/approvals/:id/withdraw", authMw, approvalHandler.Withdraw)
	}
	{
		noticeHandler := handler.NewNoticeHandler(db)
		api.GET("/notices", authMw, noticeHandler.List)
		api.GET("/notices/unread-count", authMw, noticeHandler.UnreadCount)
		api.GET("/notices/:id", authMw, noticeHandler.Detail)
		api.POST("/notices", authMw, noticeHandler.Create)
		api.POST("/notices/:id/read", authMw, noticeHandler.MarkRead)
	}
	{
		scheduleHandler := handler.NewScheduleHandler(db)
		api.GET("/schedules", authMw, scheduleHandler.List)
		api.GET("/schedules/week", authMw, scheduleHandler.WeekList)
		api.GET("/schedules/:id", authMw, scheduleHandler.Detail)
		api.POST("/schedules", authMw, scheduleHandler.Create)
		api.PUT("/schedules/:id", authMw, scheduleHandler.Update)
		api.DELETE("/schedules/:id", authMw, scheduleHandler.Delete)
	}
	{
		flowHandler := handler.NewFlowHandler(db)
		api.GET("/flows", authMw, flowHandler.List)
		api.POST("/flows", authMw, flowHandler.Create)
		api.PUT("/flows/:id", authMw, flowHandler.Update)
		api.DELETE("/flows/:id", authMw, flowHandler.Delete)
	}

	return r
}
