package service

import (
	"app/internal/config"
	"app/tools/logger"
	"app/tools/random"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// FileService 文件服务接口
type FileService interface {
	// UploadFile 上传单个文件，返回文件ID和原始文件名
	UploadFile(fileHeader *multipart.FileHeader) (string, string, error) // fileID, originalName, error
	// GetFileContent 根据文件ID获取文件内容
	GetFileContent(fileID string) ([]byte, string, error) // content, originalName, error
}

type fileService struct {
	config *config.Config
}

// NewFileService 创建文件服务实例
func NewFileService(config *config.Config) FileService {
	return &fileService{
		config: config,
	}
}

// 允许的文件类型
var allowedMimeTypes = map[string]bool{
	"image/jpeg":    true,
	"image/jpg":     true,
	"image/png":     true,
	"image/gif":     true,
	"image/webp":    true,
	"video/mp4":     true,
	"video/avi":     true,
	"video/mov":     true,
	"video/wmv":     true,
	"video/flv":     true,
	"audio/mp3":     true,
	"audio/wav":     true,
	"audio/aac":     true,
	"audio/ogg":     true,
}

// 最大文件大小限制 (100MB)
const maxFileSize = 100 * 1024 * 1024

// UploadFile 上传单个文件
func (fs *fileService) UploadFile(fileHeader *multipart.FileHeader) (string, string, error) {
	// 验证文件大小
	if fileHeader.Size > maxFileSize {
		return "", "", errors.New("文件大小超过限制（100MB）")
	}

	// 获取文件扩展名
	fileExt := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if fileExt == "" {
		return "", "", errors.New("无效的文件格式")
	}

	// 打开上传的文件
	file, err := fileHeader.Open()
	if err != nil {
		logger.Error("打开上传文件失败", "error", err)
		return "", "", errors.New("打开文件失败")
	}
	defer file.Close()

	// 读取文件内容以检测MIME类型
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", "", errors.New("读取文件失败")
	}

	// 重置文件指针
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", "", errors.New("重置文件指针失败")
	}

	// 检测MIME类型
	mimeType := http.DetectContentType(buffer)
	if !allowedMimeTypes[mimeType] {
		return "", "", errors.New("不支持的文件类型")
	}

	// 生成安全的文件ID
	fileID, err := random.GenerateURLSafeSalt(32)
	if err != nil {
		logger.Error("生成文件ID失败", "error", err)
		return "", "", errors.New("生成文件ID失败")
	}

	// 直接使用根目录，不创建子目录
	basePath := config.Get[string](fs.config, "server", "filePath")
	
	// 创建根目录（如果不存在）
	err = os.MkdirAll(basePath, 0755)
	if err != nil {
		logger.Error("创建目录失败", "path", basePath, "error", err)
		return "", "", errors.New("创建存储目录失败")
	}

	// 生成文件名：fileID + 扩展名
	fileName := fileID + fileExt
	fullPath := filepath.Join(basePath, fileName)

	// 保存文件
	dst, err := os.Create(fullPath)
	if err != nil {
		logger.Error("创建目标文件失败", "path", fullPath, "error", err)
		return "", "", errors.New("创建目标文件失败")
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		logger.Error("保存文件失败", "error", err)
		// 删除已创建的文件
		os.Remove(fullPath)
		return "", "", errors.New("保存文件失败")
	}

	logger.System("文件上传成功", "fileID", fileID, "filename", fileHeader.Filename)
	return fileID, fileHeader.Filename, nil
}

// GetFileContent 根据文件ID获取文件内容
func (fs *fileService) GetFileContent(fileID string) ([]byte, string, error) {
	basePath := config.Get[string](fs.config, "server", "filePath")
	
	// 遍历文件目录中的所有文件，查找以fileID开头的文件
	files, err := os.ReadDir(basePath)
	if err != nil {
		logger.Error("读取文件目录失败", "path", basePath, "error", err)
		return nil, "", errors.New("读取文件目录失败")
	}
	
	var foundFile string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		
		// 检查文件名是否以fileID开头（格式：fileID.ext）
		fileName := file.Name()
		if strings.HasPrefix(fileName, fileID+".") {
			foundFile = fileName
			break
		}
	}
	
	if foundFile == "" {
		return nil, "", errors.New("文件不存在")
	}
	
	// 构建完整文件路径
	fullPath := filepath.Join(basePath, foundFile)
	
	// 读取文件内容
	content, err := os.ReadFile(fullPath)
	if err != nil {
		logger.Error("读取文件内容失败", "path", fullPath, "error", err)
		return nil, "", errors.New("读取文件失败")
	}

	// 从文件名中提取扩展名作为原始文件名（简化处理）
	// 实际使用中可以考虑将原始文件名作为元数据单独存储
	originalName := foundFile

	return content, originalName, nil
}