package reporter

import (
	"encoding/json"
	"io"

	"github.com/taills/license--scanner/internal/analyzer"
)

// SARIFReporter 输出 SARIF 格式报告。
type SARIFReporter struct{}

func (r *SARIFReporter) Report(result analyzer.ScanResult, w io.Writer) error {
	type sarifResult struct {
		RuleID  string `json:"ruleId"`
		Level   string `json:"level"`
		Message struct {
			Text string `json:"text"`
		} `json:"message"`
	}
	out := struct {
		Version string `json:"version"`
		Runs    []struct {
			Tool struct {
				Driver struct {
					Name string `json:"name"`
				} `json:"driver"`
			} `json:"tool"`
			Results []sarifResult `json:"results"`
		} `json:"runs"`
	}{Version: "2.1.0"}
	run := struct {
		Tool struct {
			Driver struct {
				Name string `json:"name"`
			} `json:"driver"`
		} `json:"tool"`
		Results []sarifResult `json:"results"`
	}{}
	run.Tool.Driver.Name = "license-scanner"
	for _, dep := range result.Dependencies {
		if dep.Risk.String() == "none" || dep.Risk.String() == "low" {
			continue
		}
		sr := sarifResult{RuleID: "license-risk-" + dep.Risk.String(), Level: toSARIFLevel(dep.Risk.String())}
		sr.Message.Text = dep.Name + " 使用 " + dep.LicenseID + "，风险: " + dep.RiskReason
		run.Results = append(run.Results, sr)
	}
	out.Runs = append(out.Runs, run)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func toSARIFLevel(risk string) string {
	switch risk {
	case "critical", "high":
		return "error"
	case "medium":
		return "warning"
	default:
		return "note"
	}
}
