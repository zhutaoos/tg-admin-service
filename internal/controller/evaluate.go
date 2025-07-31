package controller

import (
	"app/internal/request"
	"app/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type EvaluateController struct {
	service.EvaluateService
}

// NewEvaluateController 创建评价控制器实例
func NewEvaluateController(evaluateService service.EvaluateService) *EvaluateController {
	return &EvaluateController{
		EvaluateService: evaluateService,
	}
}

func (e *EvaluateController) GetList(c *gin.Context) {
	var req request.EvaluateSearchRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	list, err := e.EvaluateService.GetList(req.GroupId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
