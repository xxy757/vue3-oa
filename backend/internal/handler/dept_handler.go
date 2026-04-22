package handler

import (
	"net/http"
	"oa-saas/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeptHandler struct {
	deptService *service.DeptService
}

func NewDeptHandler(deptService *service.DeptService) *DeptHandler {
	return &DeptHandler{deptService: deptService}
}

func (h *DeptHandler) List(c *gin.Context) {
	tid := getTenantID(c)
	tree, err := h.deptService.List(tid)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": tree})
}

func (h *DeptHandler) Create(c *gin.Context) {
	tid := getTenantID(c)
	var req struct {
		Name     string `json:"name" binding:"required"`
		ParentID *uint  `json:"parentId"`
		Sort     *int   `json:"sort"`
		Leader   string `json:"leader"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		Status   *int8  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	sort := 0
	if req.Sort != nil {
		sort = *req.Sort
	}
	status := int8(1)
	if req.Status != nil {
		status = *req.Status
	}
	dept, err := h.deptService.Create(tid, req.Name, req.ParentID, sort, req.Leader, req.Phone, req.Email, status)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": dept})
}

func (h *DeptHandler) Update(c *gin.Context) {
	tid := getTenantID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Name     string `json:"name" binding:"required"`
		ParentID *uint  `json:"parentId"`
		Sort     *int   `json:"sort"`
		Leader   string `json:"leader"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		Status   *int8  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := h.deptService.Update(uint(id), tid, map[string]interface{}{
		"name": req.Name, "parent_id": req.ParentID, "sort": req.Sort,
		"leader": req.Leader, "phone": req.Phone, "email": req.Email, "status": req.Status,
	}); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *DeptHandler) Delete(c *gin.Context) {
	tid := getTenantID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.deptService.Delete(uint(id), tid); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}
