package handler

import (
	"net/http"
	"oa-saas/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DeptHandler struct {
	db *gorm.DB
}

func NewDeptHandler(db *gorm.DB) *DeptHandler {
	return &DeptHandler{db: db}
}

type DeptTreeNode struct {
	ID       uint            `json:"id"`
	ParentID *uint           `json:"parentId"`
	Name     string          `json:"name"`
	Sort     int             `json:"sort"`
	Leader   string          `json:"leader"`
	Phone    string          `json:"phone"`
	Email    string          `json:"email"`
	Status   int8            `json:"status"`
	Children []*DeptTreeNode `json:"children,omitempty"`
}

func (h *DeptHandler) List(c *gin.Context) {
	var depts []model.Department
	if err := h.db.Order("sort ASC").Find(&depts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取部门列表失败"})
		return
	}

	tree := buildDeptTree(depts, nil)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": tree})
}

func buildDeptTree(depts []model.Department, parentID *uint) []*DeptTreeNode {
	var nodes []*DeptTreeNode
	for i := range depts {
		if (depts[i].ParentID == nil && parentID == nil) || (depts[i].ParentID != nil && parentID != nil && *depts[i].ParentID == *parentID) {
			node := &DeptTreeNode{
				ID:       depts[i].ID,
				ParentID: depts[i].ParentID,
				Name:     depts[i].Name,
				Sort:     depts[i].Sort,
				Leader:   depts[i].Leader,
				Phone:    depts[i].Phone,
				Email:    depts[i].Email,
				Status:   depts[i].Status,
			}
			node.Children = buildDeptTree(depts, &depts[i].ID)
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (h *DeptHandler) Create(c *gin.Context) {
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

	dept := model.Department{
		ParentID: req.ParentID,
		Name:     req.Name,
		Sort:     sort,
		Leader:   req.Leader,
		Phone:    req.Phone,
		Email:    req.Email,
		Status:   status,
	}
	if err := h.db.Create(&dept).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建部门失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": dept})
}

func (h *DeptHandler) Update(c *gin.Context) {
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

	result := h.db.Model(&model.Department{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":      req.Name,
		"parent_id": req.ParentID,
		"sort":      req.Sort,
		"leader":    req.Leader,
		"phone":     req.Phone,
		"email":     req.Email,
		"status":    req.Status,
	})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "部门不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *DeptHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var count int64
	h.db.Model(&model.Department{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "该部门下有子部门，无法删除"})
		return
	}
	if err := h.db.Delete(&model.Department{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}
