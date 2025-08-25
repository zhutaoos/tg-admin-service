package router

import (
	"app/internal/controller/file"

	"github.com/gin-gonic/gin"
)

// FileRoute 文件路由结构
type FileRoute struct {
	fileController *file.FileController
}

// NewFileRoute 创建文件路由实例
func NewFileRoute(fileController *file.FileController) *FileRoute {
	return &FileRoute{
		fileController: fileController,
	}
}

// InitRoute 初始化文件路由
func (fr *FileRoute) InitRoute(engine *gin.Engine) {
	// 文件管理路由组
	fileGroup := engine.Group("/api/file")
	{
		// 单文件上传
		fileGroup.POST("/upload", fr.fileController.UploadFile)
		
		// 多文件上传
		fileGroup.POST("/uploads", fr.fileController.UploadFiles)
		
		// 下载文件（简化路径）
		fileGroup.GET("/:fileId", fr.fileController.DownloadFile)
	}
}