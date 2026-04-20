package reporter

import (
	"encoding/json"
	"io"

	"github.com/taills/license--scanner/internal/analyzer"
)

// JSONReporter 输出 JSON 报告。
type JSONReporter struct{}

func (r *JSONReporter) Report(result analyzer.ScanResult, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}
