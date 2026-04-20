package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// MavenFetcher 通过 Maven Central 查询协议信息。
type MavenFetcher struct {
	client *http.Client
}

// NewMavenFetcher 创建 Maven 拉取器。
func NewMavenFetcher() *MavenFetcher {
	return &MavenFetcher{client: &http.Client{Timeout: 5 * time.Second}}
}

func (f *MavenFetcher) FetchLicense(ctx context.Context, name, version string) (string, error) {
	parts := strings.Split(name, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("maven 包名应为 group:artifact")
	}
	q := fmt.Sprintf("g:%s AND a:%s", parts[0], parts[1])
	u := "https://search.maven.org/solrsearch/select?q=" + url.QueryEscape(q) + "&rows=1&wt=json"
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	resp, err := f.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("maven central 返回状态码 %d", resp.StatusCode)
	}
	var body struct {
		Response struct {
			Docs []struct {
				L []string `json:"l"`
			} `json:"docs"`
		} `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	if len(body.Response.Docs) > 0 && len(body.Response.Docs[0].L) > 0 {
		return body.Response.Docs[0].L[0], nil
	}
	return "", fmt.Errorf("maven central 未返回 license")
}
