package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// GitHubFetcher 通过 GitHub API 获取仓库协议。
type GitHubFetcher struct {
	client *http.Client
	token  string
}

// NewGitHubFetcher 创建 GitHub 拉取器。
func NewGitHubFetcher() *GitHubFetcher {
	return &GitHubFetcher{client: &http.Client{Timeout: 5 * time.Second}, token: os.Getenv("GITHUB_TOKEN")}
}

func (f *GitHubFetcher) FetchLicense(ctx context.Context, name, version string) (string, error) {
	owner, repo, ok := parseRepo(name)
	if !ok {
		return "", fmt.Errorf("无法从依赖名推断 GitHub 仓库")
	}
	u := fmt.Sprintf("https://api.github.com/repos/%s/%s/license", owner, repo)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	if f.token != "" {
		req.Header.Set("Authorization", "Bearer "+f.token)
	}
	resp, err := f.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("github api 返回状态码 %d", resp.StatusCode)
	}
	var body struct {
		License struct {
			SPDXID string `json:"spdx_id"`
		} `json:"license"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	if body.License.SPDXID == "" {
		return "", fmt.Errorf("github api 未返回 spdx_id")
	}
	return body.License.SPDXID, nil
}

func parseRepo(name string) (string, string, bool) {
	n := strings.TrimPrefix(name, "github.com/")
	parts := strings.Split(n, "/")
	if len(parts) < 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}
