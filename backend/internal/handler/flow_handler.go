package handler

import (
	"net/http"
	"oa-saas/internal/model"
	"oa-saas/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FlowHandler struct {
	flowService *service.FlowService
}

func NewFlowHandler(flowService *service.FlowService) *FlowHandler {
	return &FlowHandler{flowService: flowService}
}

func (h *FlowHandler) List(c *gin.Context) {
	tid := getTenantID(c)
	flows, err := h.flowService.List(tid)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": flows})
}

func (h *FlowHandler) Create(c *gin.Context) {
	tid := getTenantID(c)
	var req struct {
		Name        string           `json:"name" binding:"required"`
		Code        string           `json:"code" binding:"required"`
		Description string           `json:"description"`
		Nodes       []model.FlowNode `json:"nodes" binding:"required"`
		Status      *int8            `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	status := int8(1)
	if req.Status != nil {
		status = *req.Status
	}
	flow, err := h.flowService.Create(tid, req.Name, req.Code, req.Description, req.Nodes, status)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": flow})
}

func (h *FlowHandler) Update(c *gin.Context) {
	tid := getTenantID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Name        string           `json:"name" binding:"required"`
		Description string           `json:"description"`
		Nodes       []model.FlowNode `json:"nodes" binding:"required"`
		Status      *int8            `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	updates := map[string]interface{}{
		"name": req.Name, "description": req.Description, "nodes": req.Nodes,
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if err := h.flowService.Update(uint(id), tid, updates); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *FlowHandler) Delete(c *gin.Context) {
	tid := getTenantID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.flowService.Delete(uint(id), tid); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}
