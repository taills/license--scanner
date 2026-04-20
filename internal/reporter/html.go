package reporter

import (
	"html/template"
	"io"

	"github.com/taills/license--scanner/internal/analyzer"
)

// HTMLReporter 输出 HTML 报告。
type HTMLReporter struct{}

const htmlTpl = `<!doctype html>
<html lang="zh-CN">
<head><meta charset="utf-8"><title>License Scanner Report</title>
<style>body{font-family:Arial,sans-serif;margin:24px}table{border-collapse:collapse;width:100%}th,td{border:1px solid #ddd;padding:8px}th{background:#f4f4f4}.high,.critical{color:#b00020;font-weight:bold}</style>
</head>
<body>
<h1>License Scanner Report</h1>
<p>项目：{{.ProjectPath}} | 策略：{{.PolicyName}} | 扫描时间：{{.ScannedAt}}</p>
<table>
<thead><tr><th>生态</th><th>依赖</th><th>版本</th><th>协议</th><th>风险</th><th>原因</th></tr></thead>
<tbody>
{{range .Dependencies}}
<tr><td>{{.Ecosystem}}</td><td>{{.Name}}</td><td>{{.Version}}</td><td>{{.LicenseID}}</td><td class="{{.Risk}}">{{.Risk}}</td><td>{{.RiskReason}}</td></tr>
{{end}}
</tbody>
</table>
</body></html>`

func (r *HTMLReporter) Report(result analyzer.ScanResult, w io.Writer) error {
	t, err := template.New("report").Parse(htmlTpl)
	if err != nil {
		return err
	}
	return t.Execute(w, result)
}
