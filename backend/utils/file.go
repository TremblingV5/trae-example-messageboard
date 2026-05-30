package utils

import (
	"fmt"
	"messageboard/config"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SaveFile 保存上传的文件
func SaveFile(file *multipart.FileHeader, subDir string) (string, error) {
	// Validate file extension before saving
	if !IsImageFile(file.Filename) {
		return "", fmt.Errorf("file type not allowed: %s", file.Filename)
	}

	// Create upload directory if not exists
	uploadDir := filepath.Join(config.AppConfig.Upload.UploadDir, subDir)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join(uploadDir, filename)

	// Save file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(src); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	// Return relative URL path
	return fmt.Sprintf("/uploads/%s/%s", subDir, filename), nil
}

// GetFileExtension 获取文件扩展名
func GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// IsImageFile 检查是否为图片文件
func IsImageFile(filename string) bool {
	ext := GetFileExtension(filename)
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	return allowedExts[ext]
}

// EnsureDir 确保目录存在
func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}
