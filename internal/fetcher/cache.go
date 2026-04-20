package fetcher

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// FileCache 提供简单的文件缓存实现。
type FileCache struct {
	path string
	mu   sync.Mutex
	data map[string]string
}

// NewFileCache 创建文件缓存。
func NewFileCache(path string) *FileCache {
	if path == "" {
		if dir, err := os.UserCacheDir(); err == nil {
			path = filepath.Join(dir, "license-scanner", "cache.json")
		}
	}
	c := &FileCache{path: path, data: map[string]string{}}
	c.load()
	return c
}

// Get 读取缓存值。
func (c *FileCache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.data[key]
	return v, ok
}

// Set 写入缓存值。
func (c *FileCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	_ = c.save()
}

func (c *FileCache) load() {
	if c.path == "" {
		return
	}
	b, err := os.ReadFile(c.path)
	if err != nil {
		return
	}
	_ = json.Unmarshal(b, &c.data)
}

func (c *FileCache) save() error {
	if c.path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(c.path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(c.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, b, 0o644)
}
