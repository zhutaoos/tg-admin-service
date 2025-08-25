package file

import (
	"app/internal/request"
	"app/internal/service"
	"app/internal/vo"
	"app/tools/resp"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FileController 文件控制器
type FileController struct {
	fileService service.FileService
}

// NewFileController 创建文件控制器实例
func NewFileController(fileService service.FileService) *FileController {
	return &FileController{
		fileService: fileService,
	}
}

// UploadFile 上传单个文件
func (fc *FileController) UploadFile(ctx *gin.Context) {
	var req request.FileUploadRequest
	if err := ctx.ShouldBind(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}

	// 上传文件
	fileID, originalName, err := fc.fileService.UploadFile(req.File)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "上传失败: " + err.Error()}).Response()
		return
	}

	// 构建响应对象
	fileVO := &vo.FileUploadVO{
		FileID:      fileID,
		FileName:    originalName,
		FileSize:    req.File.Size,
		DownloadURL: fmt.Sprintf("/api/file/%s", fileID),
	}

	resp.Ok(fileVO)
}

// UploadFiles 上传多个文件
func (fc *FileController) UploadFiles(ctx *gin.Context) {
	var req request.MultiFileUploadRequest
	if err := ctx.ShouldBind(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}

	var successFiles []*vo.FileUploadVO
	var failedFiles []*vo.FileError

	// 逐个上传文件
	for _, fileHeader := range req.Files {
		fileID, originalName, err := fc.fileService.UploadFile(fileHeader)
		if err != nil {
			failedFiles = append(failedFiles, &vo.FileError{
				FileName: fileHeader.Filename,
				Error:    err.Error(),
			})
			continue
		}

		successFiles = append(successFiles, &vo.FileUploadVO{
			FileID:      fileID,
			FileName:    originalName,
			FileSize:    fileHeader.Size,
			DownloadURL: fmt.Sprintf("/api/file/%s", fileID),
		})
	}

	// 构建响应对象
	result := &vo.MultiFileUploadVO{
		SuccessFiles: successFiles,
		FailedFiles:  failedFiles,
		SuccessCount: len(successFiles),
		FailedCount:  len(failedFiles),
	}

	resp.Ok(result)
}

// DownloadFile 下载文件（公开访问，无需权限验证）
func (fc *FileController) DownloadFile(ctx *gin.Context) {
	var req request.GetFileRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": resp.ReFail,
			"msg":  "参数错误: " + err.Error(),
			"data": nil,
		})
		return
	}

	// 直接根据文件ID获取文件
	content, _, err := fc.fileService.GetFileContent(req.FileID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code": resp.ReFail,
			"msg":  "获取文件失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	// 检测MIME类型
	mimeType := http.DetectContentType(content)

	// 设置响应头 - 移除attachment让浏览器直接显示文件
	ctx.Header("Content-Type", mimeType)
	ctx.Header("Content-Length", fmt.Sprintf("%d", len(content)))

	// 返回文件内容
	ctx.Data(http.StatusOK, mimeType, content)
}