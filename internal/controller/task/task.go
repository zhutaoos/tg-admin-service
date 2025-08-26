package task

import (
	"app/internal/controller"
	"app/internal/request"
	"app/internal/service"
	"app/tools/resp"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TaskController 任务控制器
type TaskController struct {
	controller.BaseController
	service.TaskService
}

// NewTaskController 创建任务控制器实例
func NewTaskController(taskService service.TaskService) *TaskController {
	return &TaskController{
		TaskService: taskService,
	}
}

// CreateTask 创建任务
func (tc *TaskController) CreateTask(ctx *gin.Context) {
	var req request.CreateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数缺失或格式错误: " + err.Error()}).Response()
		return
	}

	// 获取当前用户ID
	adminID := uint64(tc.CurrentUserId(ctx))

	// 调用服务层创建任务
	taskVO, err := tc.TaskService.CreateTask(&req, adminID)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReError, Msg: "创建任务失败: " + err.Error()}).Response()
		return
	}

	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "创建任务成功", Data: taskVO}).Response()
}

// UpdateTask 更新任务
func (tc *TaskController) UpdateTask(ctx *gin.Context) {
	var req request.UpdateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数缺失或格式错误: " + err.Error()}).Response()
		return
	}

	// 获取当前用户ID
	adminID := uint64(tc.CurrentUserId(ctx))

	// 调用服务层更新任务
	taskVO, err := tc.TaskService.UpdateTask(&req, adminID)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReError, Msg: "更新任务失败: " + err.Error()}).Response()
		return
	}

	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "更新任务成功", Data: taskVO}).Response()
}

// DeleteTask 删除任务
func (tc *TaskController) DeleteTask(ctx *gin.Context) {
	// 从URL参数获取任务ID
	idStr := ctx.Param("id")
	if idStr == "" {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "任务ID不能为空"}).Response()
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "任务ID格式错误"}).Response()
		return
	}

	req := request.DeleteTaskRequest{ID: id}

	// 获取当前用户ID
	adminID := uint64(tc.CurrentUserId(ctx))

	// 调用服务层删除任务
	if err := tc.TaskService.DeleteTask(&req, adminID); err != nil {
		(&resp.JsonResp{Code: resp.ReError, Msg: "删除任务失败: " + err.Error()}).Response()
		return
	}

	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "删除任务成功"}).Response()
}

// GetTaskDetail 获取任务详情
func (tc *TaskController) GetTaskDetail(ctx *gin.Context) {
	// 从URL参数获取任务ID
	idStr := ctx.Param("id")
	if idStr == "" {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "任务ID不能为空"}).Response()
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "任务ID格式错误"}).Response()
		return
	}

	// 获取当前用户ID
	adminID := uint64(tc.CurrentUserId(ctx))

	// 调用服务层获取任务详情
	taskVO, err := tc.TaskService.GetTaskByID(id, adminID)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReError, Msg: "获取任务详情失败: " + err.Error()}).Response()
		return
	}

	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "获取任务详情成功", Data: taskVO}).Response()
}

// TaskList 任务列表
func (tc *TaskController) TaskList(ctx *gin.Context) {
	var req request.TaskListRequest
	if err := ctx.ShouldBind(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数缺失或格式错误: " + err.Error()}).Response()
		return
	}

	// 获取当前用户ID
	adminID := uint64(tc.CurrentUserId(ctx))

	// 调用服务层获取任务列表
	taskListVO, err := tc.TaskService.ListTasks(&req, adminID)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReError, Msg: "查询任务列表失败: " + err.Error()}).Response()
		return
	}

	// 构建响应数据
	data := map[string]interface{}{
		"list":  taskListVO.List,
		"total": taskListVO.Total,
		"page":  req.Page,
		"limit": req.Limit,
		"pages": (taskListVO.Total + int64(req.Limit) - 1) / int64(req.Limit),
	}

	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "获取任务列表成功", Data: data}).Response()
}

// TaskStats 任务统计
func (tc *TaskController) TaskStats(ctx *gin.Context) {
	// 获取当前用户ID
	adminID := uint64(tc.CurrentUserId(ctx))

	// 调用服务层获取任务统计
	statsVO, err := tc.TaskService.GetTaskStats(adminID)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReError, Msg: "获取任务统计失败: " + err.Error()}).Response()
		return
	}

	(&resp.JsonResp{Code: resp.ReSuccess, Msg: "获取任务统计成功", Data: statsVO}).Response()
}