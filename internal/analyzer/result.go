package analyzer

import (
	"time"

	"github.com/taills/license--scanner/internal/license"
)

// Dependency 表示最终输出的依赖与协议风险信息。
type Dependency struct {
	Name          string            `json:"name"`
	Version       string            `json:"version"`
	Ecosystem     string            `json:"ecosystem"`
	Direct        bool              `json:"direct"`
	DepPath       []string          `json:"dep_path,omitempty"`
	LicenseID     string            `json:"license_id"`
	LicenseSource string            `json:"license_source"`
	Risk          license.RiskLevel `json:"risk"`
	RiskReason    string            `json:"risk_reason"`
}

// Summary 表示扫描统计结果。
type Summary struct {
	Total    int `json:"total"`
	Direct   int `json:"direct"`
	None     int `json:"none"`
	Low      int `json:"low"`
	Medium   int `json:"medium"`
	High     int `json:"high"`
	Critical int `json:"critical"`
	Unknown  int `json:"unknown"`
}

// ScanResult 表示扫描完整结果。
type ScanResult struct {
	ProjectPath  string       `json:"project_path"`
	PolicyName   string       `json:"policy_name"`
	ScannedAt    time.Time    `json:"scanned_at"`
	Dependencies []Dependency `json:"dependencies"`
	Summary      Summary      `json:"summary"`
}
