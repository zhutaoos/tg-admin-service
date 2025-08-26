package router

import (
	"app/internal/controller/task"

	"github.com/gin-gonic/gin"
)

// TaskRoute 任务路由结构
type TaskRoute struct {
	TaskController *task.TaskController
}

// NewTaskRoute 创建任务路由实例
func NewTaskRoute(taskController *task.TaskController) *TaskRoute {
	return &TaskRoute{
		TaskController: taskController,
	}
}

// InitRoute 初始化任务路由
func (tr *TaskRoute) InitRoute(r *gin.Engine) {
	// 任务管理路由组
	taskGroup := r.Group("/api/task")
	{
		// 创建任务
		taskGroup.POST("/create", tr.TaskController.CreateTask)
		
		// 更新任务
		taskGroup.PUT("/update", tr.TaskController.UpdateTask)
		
		// 删除任务
		taskGroup.DELETE("/delete/:id", tr.TaskController.DeleteTask)
		
		// 获取任务详情
		taskGroup.GET("/detail/:id", tr.TaskController.GetTaskDetail)
		
		// 任务列表
		taskGroup.GET("/list", tr.TaskController.TaskList)
		
		// 任务统计
		taskGroup.GET("/stats", tr.TaskController.TaskStats)
	}
}