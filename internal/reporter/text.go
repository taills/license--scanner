package reporter

import (
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/taills/license--scanner/internal/analyzer"
	"github.com/taills/license--scanner/internal/license"
)

// TextReporter 以终端彩色文本输出扫描结果。
type TextReporter struct{}

func (r *TextReporter) Report(result analyzer.ScanResult, w io.Writer) error {
	headline := color.New(color.FgCyan, color.Bold).SprintFunc()
	warn := color.New(color.FgYellow).SprintFunc()
	danger := color.New(color.FgRed, color.Bold).SprintFunc()
	ok := color.New(color.FgGreen).SprintFunc()
	fmt.Fprintf(w, "%s\n", headline("License Scanner Report"))
	fmt.Fprintf(w, "项目: %s\n策略: %s\n总依赖: %d（直接依赖: %d）\n\n", result.ProjectPath, result.PolicyName, result.Summary.Total, result.Summary.Direct)
	fmt.Fprintf(w, "风险汇总: %s=%d %s=%d %s=%d %s=%d %s=%d %s=%d\n\n",
		ok("none"), result.Summary.None,
		warn("low"), result.Summary.Low,
		warn("medium"), result.Summary.Medium,
		danger("high"), result.Summary.High,
		danger("critical"), result.Summary.Critical,
		warn("unknown"), result.Summary.Unknown,
	)
	for _, dep := range result.Dependencies {
		tag := dep.Risk.String()
		if dep.Risk == license.RiskHigh || dep.Risk == license.RiskCritical {
			tag = danger(tag)
		}
		fmt.Fprintf(w, "- [%s] %s@%s (%s) license=%s source=%s\n  原因: %s\n", tag, dep.Name, dep.Version, dep.Ecosystem, dep.LicenseID, dep.LicenseSource, dep.RiskReason)
	}
	return nil
}
