package handler

import (
	"net/http"
	"oa-saas/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.DefaultQuery("keyword", "")
	tid := getTenantID(c)

	list, total, err := h.userService.List(tid, keyword, page, pageSize)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"list": list, "total": total}})
}

func (h *UserHandler) Create(c *gin.Context) {
	tid := getTenantID(c)
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Nickname string `json:"nickname" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Phone    string `json:"phone"`
		DeptID   *uint  `json:"deptId"`
		RoleID   *uint  `json:"roleId"`
		Status   *int8  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	status := int8(1)
	if req.Status != nil {
		status = *req.Status
	}
	user, err := h.userService.Create(tid, req.Username, req.Password, req.Nickname, req.Email, req.Phone, req.DeptID, req.RoleID, status)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": user})
}

func (h *UserHandler) Update(c *gin.Context) {
	tid := getTenantID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Avatar   string `json:"avatar"`
		DeptID   *uint  `json:"deptId"`
		RoleID   *uint  `json:"roleId"`
		Status   *int8  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := h.userService.Update(uint(id), tid, map[string]interface{}{
		"nickname": req.Nickname, "email": req.Email, "phone": req.Phone,
		"avatar": req.Avatar, "dept_id": req.DeptID, "role_id": req.RoleID, "status": req.Status,
	}); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *UserHandler) Delete(c *gin.Context) {
	tid := getTenantID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.userService.Delete(uint(id), tid); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *UserHandler) UpdateStatus(c *gin.Context) {
	tid := getTenantID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Status int8 `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := h.userService.UpdateStatus(uint(id), tid, req.Status); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "操作成功"})
}
