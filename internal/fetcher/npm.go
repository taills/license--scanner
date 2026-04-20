package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// NPMFetcher 通过 npm registry 获取协议。
type NPMFetcher struct {
	client *http.Client
}

// NewNPMFetcher 创建 npm 拉取器。
func NewNPMFetcher() *NPMFetcher {
	return &NPMFetcher{client: &http.Client{Timeout: 5 * time.Second}}
}

func (f *NPMFetcher) FetchLicense(ctx context.Context, name, version string) (string, error) {
	u := fmt.Sprintf("https://registry.npmjs.org/%s", url.PathEscape(name))
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	resp, err := f.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("npm registry 返回状态码 %d", resp.StatusCode)
	}
	var body struct {
		License  any `json:"license"`
		Versions map[string]struct {
			License any `json:"license"`
		} `json:"versions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	if version != "" {
		if v, ok := body.Versions[version]; ok {
			switch x := v.License.(type) {
			case string:
				return x, nil
			}
		}
	}
	switch x := body.License.(type) {
	case string:
		return x, nil
	}
	return "", fmt.Errorf("npm 未返回 license 字段")
}
