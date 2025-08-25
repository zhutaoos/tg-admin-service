package vo

// FileUploadVO 文件上传成功响应对象
type FileUploadVO struct {
	FileID      string `json:"fileId"`      // 文件ID
	FileName    string `json:"fileName"`    // 文件名
	FileSize    int64  `json:"fileSize"`    // 文件大小
	DownloadURL string `json:"downloadUrl"` // 下载链接
}

// MultiFileUploadVO 多文件上传响应对象
type MultiFileUploadVO struct {
	SuccessFiles []*FileUploadVO `json:"successFiles"` // 成功上传的文件
	FailedFiles  []*FileError    `json:"failedFiles"`  // 失败的文件
	SuccessCount int             `json:"successCount"` // 成功数量
	FailedCount  int             `json:"failedCount"`  // 失败数量
}

// FileError 文件错误信息
type FileError struct {
	FileName string `json:"fileName"` // 文件名
	Error    string `json:"error"`    // 错误信息
}