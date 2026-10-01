package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gin-web/core/config"
)

// LocalDriver 本地磁盘存储，文件保存在 upload.path 下，由 /uploads 静态路由或Nginx托管
type LocalDriver struct {
	root      string
	urlPrefix string
}

func NewLocalDriver() *LocalDriver {
	root := config.GetString("upload.path")
	if root == "" {
		root = "./runtime/uploads"
	}
	urlPrefix := config.GetString("upload.urlPrefix")
	if urlPrefix == "" {
		urlPrefix = "/uploads"
	}
	return &LocalDriver{root: root, urlPrefix: urlPrefix}
}

func (d *LocalDriver) Save(filename string, data []byte) (string, error) {
	name := randomFilename(filename)
	dst := filepath.Join(d.root, name)

	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return "", fmt.Errorf("create upload dir fail: %w", err)
	}

	if err := os.WriteFile(dst, data, 0o644); err != nil {
		return "", fmt.Errorf("write upload file fail: %w", err)
	}

	// filepath.Join在windows下会生成"\"，URL统一使用"/"
	return d.urlPrefix + "/" + filepath.ToSlash(name), nil
}

func dateDir() string {
	return time.Now().Format("20060102")
}
