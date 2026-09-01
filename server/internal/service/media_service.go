package service

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"leslie-blog-server/config"
	"leslie-blog-server/internal/model"
	"leslie-blog-server/internal/repository"
)

type MediaService struct {
	repo *repository.MediaRepo
	cfg  *config.Config
}

func NewMediaService(repo *repository.MediaRepo, cfg *config.Config) *MediaService {
	return &MediaService{repo: repo, cfg: cfg}
}

func (s *MediaService) List(keyword, mtype string, page, pageSize int) ([]model.Media, int64, error) {
	return s.repo.List(keyword, mtype, page, pageSize)
}

// Save 上传文件落盘 + 入库
func (s *MediaService) Save(filename string, content []byte) (*model.Media, error) {
	max := int64(s.cfg.Upload.MaxSizeMB) * 1024 * 1024
	if max > 0 && int64(len(content)) > max {
		return nil, fmt.Errorf("文件超过 %dMB 限制", s.cfg.Upload.MaxSizeMB)
	}
	// 生成存储路径：uploads/YYYY/MM/<unique>
	rel, url, err := s.buildPath(filename)
	if err != nil {
		return nil, err
	}
	// 写入 uploads 目录（由 main.go 静态服务）
	if err := writeFile(rel, content); err != nil {
		return nil, err
	}
	m := model.Media{
		Name:     filename,
		URL:      url,
		Type:     detectType(filename),
		Size:     int64(len(content)),
		MimeType: mime.TypeByExtension(filepath.Ext(filename)),
	}
	if err := s.repo.Create(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *MediaService) Delete(id uint) error { return s.repo.Delete(id) }

// buildPath 生成磁盘相对路径与可访问 URL
func (s *MediaService) buildPath(filename string) (string, string, error) {
	dir := filepath.Join(s.cfg.App.UploadDir, time.Now().Format("2006/01"))
	if err := mkdirAll(dir); err != nil {
		return "", "", err
	}
	name := strings.TrimSpace(filename)
	if name == "" {
		name = "file"
	}
	// 避免重名：加时间戳
	base := strings.TrimSuffix(name, filepath.Ext(name))
	ext := filepath.Ext(name)
	unique := fmt.Sprintf("%s_%d%s", base, time.Now().UnixMilli(), ext)
	rel := filepath.Join(dir, unique)
	url := fmt.Sprintf("%s/%s/%s",
		strings.TrimRight(s.cfg.App.BaseURL, "/"),
		strings.TrimPrefix(s.cfg.App.UploadDir, "./"),
		filepath.Join(time.Now().Format("2006/01"), unique),
	)
	url = strings.ReplaceAll(url, "\\", "/")
	return rel, url, nil
}

func detectType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg", ".bmp":
		return "image"
	case ".mp4", ".webm", ".mov", ".avi", ".mkv":
		return "video"
	default:
		return "file"
	}
}

// writeFile / mkdirAll 薄封装，便于测试替换
func writeFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0o644)
}
func mkdirAll(path string) error {
	return os.MkdirAll(path, 0o755)
}
