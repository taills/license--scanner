package license

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/taills/license--scanner/internal/fetcher"
	"github.com/taills/license--scanner/internal/parsers"
)

// Detector 执行多源协议检测。
type Detector struct {
	Fetchers map[string]fetcher.LicenseFetcher
	GitHub   *fetcher.GitHubFetcher
	Cache    *fetcher.FileCache
}

// NewDetector 创建默认检测器。
func NewDetector(cachePath string) *Detector {
	cache := fetcher.NewFileCache(cachePath)
	return &Detector{
		Fetchers: map[string]fetcher.LicenseFetcher{
			"npm":   fetcher.NewNPMFetcher(),
			"pypi":  fetcher.NewPyPIFetcher(),
			"maven": fetcher.NewMavenFetcher(),
		},
		GitHub: fetcher.NewGitHubFetcher(),
		Cache:  cache,
	}
}

// Detect 检测依赖协议并返回协议与来源。
func (d *Detector) Detect(ctx context.Context, projectPath string, dep parsers.ParsedDependency) (string, string) {
	if id := detectFromProjectFiles(projectPath); id != "" {
		id = NormalizeLicenseID(id)
		return id, "project-license-file"
	}
	key := dep.Ecosystem + ":" + dep.Name + "@" + dep.Version
	if d.Cache != nil {
		if v, ok := d.Cache.Get(key); ok {
			return NormalizeLicenseID(v), "cache"
		}
	}
	if f, ok := d.Fetchers[dep.Ecosystem]; ok {
		if id, err := f.FetchLicense(ctx, dep.Name, dep.Version); err == nil && id != "" {
			id = NormalizeLicenseID(id)
			if d.Cache != nil {
				d.Cache.Set(key, id)
			}
			return id, dep.Ecosystem + "-registry"
		}
	}
	if d.GitHub != nil && (strings.HasPrefix(dep.Name, "github.com/") || strings.Count(dep.Name, "/") == 1) {
		if id, err := d.GitHub.FetchLicense(ctx, dep.Name, dep.Version); err == nil && id != "" {
			id = NormalizeLicenseID(id)
			if d.Cache != nil {
				d.Cache.Set(key, id)
			}
			return id, "github-api"
		}
	}
	return "UNKNOWN", "unresolved"
}

func detectFromProjectFiles(projectPath string) string {
	candidates := []string{"LICENSE", "LICENSE.txt", "COPYING", "COPYRIGHT"}
	for _, name := range candidates {
		p := filepath.Join(projectPath, name)
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		t := strings.ToLower(string(b))
		switch {
		case strings.Contains(t, "mit license"):
			return "MIT"
		case strings.Contains(t, "apache license"):
			return "Apache-2.0"
		case strings.Contains(t, "gnu general public license"):
			return "GPL-3.0"
		}
	}
	return ""
}
