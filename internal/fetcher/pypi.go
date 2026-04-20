package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// PyPIFetcher 通过 PyPI API 获取协议。
type PyPIFetcher struct {
	client *http.Client
}

// NewPyPIFetcher 创建 PyPI 拉取器。
func NewPyPIFetcher() *PyPIFetcher {
	return &PyPIFetcher{client: &http.Client{Timeout: 5 * time.Second}}
}

func (f *PyPIFetcher) FetchLicense(ctx context.Context, name, version string) (string, error) {
	u := fmt.Sprintf("https://pypi.org/pypi/%s/json", name)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	resp, err := f.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("pypi 返回状态码 %d", resp.StatusCode)
	}
	var body struct {
		Info struct {
			License string `json:"license"`
		} `json:"info"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	if body.Info.License == "" {
		return "", fmt.Errorf("pypi 未返回 license")
	}
	return body.Info.License, nil
}
