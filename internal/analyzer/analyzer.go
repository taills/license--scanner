package analyzer

import (
	"context"
	"sort"
	"time"

	"github.com/taills/license--scanner/internal/license"
	"github.com/taills/license--scanner/internal/parsers"
)

// Engine 协调 parser、detector、risk 评估流程。
type Engine struct {
	Parsers  []parsers.Parser
	Detector *license.Detector
	Policy   *license.Policy
}

// Analyze 执行项目扫描。
func (e *Engine) Analyze(ctx context.Context, projectPath string, ecosystems map[string]bool) (ScanResult, error) {
	result := ScanResult{ProjectPath: projectPath, ScannedAt: time.Now(), PolicyName: e.Policy.Name}
	merged := map[string]parsers.ParsedDependency{}

	for _, p := range e.Parsers {
		if !p.Supports(projectPath) {
			continue
		}
		ds, err := p.Parse(ctx, projectPath)
		if err != nil {
			return result, err
		}
		for _, d := range ds {
			if len(ecosystems) > 0 && !ecosystems[d.Ecosystem] {
				continue
			}
			k := d.Ecosystem + ":" + d.Name + "@" + d.Version
			if old, ok := merged[k]; ok {
				old.Direct = old.Direct || d.Direct
				if len(old.DepPath) == 0 {
					old.DepPath = d.DepPath
				}
				merged[k] = old
				continue
			}
			merged[k] = d
		}
	}

	for _, d := range merged {
		licenseID, source := e.Detector.Detect(ctx, projectPath, d)
		risk, reason := license.EvaluateRisk(licenseID, e.Policy)
		dep := Dependency{
			Name:          d.Name,
			Version:       d.Version,
			Ecosystem:     d.Ecosystem,
			Direct:        d.Direct,
			DepPath:       d.DepPath,
			LicenseID:     licenseID,
			LicenseSource: source,
			Risk:          risk,
			RiskReason:    reason,
		}
		result.Dependencies = append(result.Dependencies, dep)
		result.Summary.Total++
		if dep.Direct {
			result.Summary.Direct++
		}
		switch dep.Risk {
		case license.RiskNone:
			result.Summary.None++
		case license.RiskLow:
			result.Summary.Low++
		case license.RiskMedium:
			result.Summary.Medium++
		case license.RiskHigh:
			result.Summary.High++
		case license.RiskCritical:
			result.Summary.Critical++
		default:
			result.Summary.Unknown++
		}
	}

	sort.Slice(result.Dependencies, func(i, j int) bool {
		if result.Dependencies[i].Ecosystem == result.Dependencies[j].Ecosystem {
			return result.Dependencies[i].Name < result.Dependencies[j].Name
		}
		return result.Dependencies[i].Ecosystem < result.Dependencies[j].Ecosystem
	})

	return result, nil
}
