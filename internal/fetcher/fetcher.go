package fetcher

import "context"

// LicenseFetcher 定义远程协议拉取接口。
type LicenseFetcher interface {
	FetchLicense(ctx context.Context, name, version string) (string, error)
}
