package reporter

import (
	"io"

	"github.com/taills/license--scanner/internal/analyzer"
)

// Reporter 定义结果输出接口。
type Reporter interface {
	Report(result analyzer.ScanResult, w io.Writer) error
}
