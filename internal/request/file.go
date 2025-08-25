package request

import "mime/multipart"

// FileUploadRequest 单文件上传请求
type FileUploadRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required" json:"-"` // 上传的文件
}

// MultiFileUploadRequest 多文件上传请求
type MultiFileUploadRequest struct {
	Files []*multipart.FileHeader `form:"files" binding:"required" json:"-"` // 上传的文件列表
}

// GetFileRequest 获取文件请求
type GetFileRequest struct {
	FileID string `uri:"fileId" binding:"required" json:"fileId"` // 文件ID
}