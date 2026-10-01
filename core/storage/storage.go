// Package storage 文件存储抽象，支持本地存储，可扩展对象存储(OSS/S3)
package storage

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"gin-web/core/config"
	"gin-web/internal/errcode"
	"github.com/google/uuid"
)

// Driver 存储驱动接口
type Driver interface {
	// Save 保存文件，返回可直接访问的URL
	Save(filename string, data []byte) (string, error)
}

// Validate 校验上传文件的扩展名与大小
func Validate(header *multipart.FileHeader) error {
	maxSizeMB := config.GetInt("upload.maxSize")
	if maxSizeMB <= 0 {
		maxSizeMB = 5
	}
	if header.Size > int64(maxSizeMB)<<20 {
		return errcode.NewCustomError(errcode.ErrFileInvalid)
	}

	allowExt := config.GetStringSlice("upload.allowExt")
	if len(allowExt) == 0 {
		allowExt = []string{"jpg", "jpeg", "png", "gif", "webp"}
	}

	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(header.Filename)), ".")
	for _, allow := range allowExt {
		if ext == strings.ToLower(allow) {
			return nil
		}
	}
	return errcode.NewCustomError(errcode.ErrFileInvalid)
}

// ValidateContent rejects files whose signature does not match an allowed
// image media type. Extension checks alone are insufficient for uploads.
func ValidateContent(data []byte) error {
	if len(data) == 0 {
		return errcode.NewCustomError(errcode.ErrFileInvalid)
	}
	contentType := http.DetectContentType(data)
	allowed := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	if !allowed[contentType] {
		return errcode.NewCustomError(errcode.ErrFileInvalid)
	}
	return nil
}

// NewDriver 根据配置创建存储驱动
func NewDriver() Driver {
	driver := config.GetString("upload.driver")
	if driver == "" {
		driver = "local"
	}
	switch driver {
	case "local":
		return NewLocalDriver()
	default:
		panic(fmt.Sprintf("unsupported upload driver: %s", driver))
	}
}

// randomFilename 生成存储文件名：日期目录/uuid.ext
func randomFilename(filename string) string {
	ext := filepath.Ext(filename)
	return filepath.Join(dateDir(), uuid.NewString()+ext)
}
