package license

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// RiskLevel 表示依赖协议风险等级。
type RiskLevel int

const (
	RiskNone RiskLevel = iota
	RiskLow
	RiskMedium
	RiskHigh
	RiskCritical
	RiskUnknown
)

func (r RiskLevel) String() string {
	switch r {
	case RiskNone:
		return "none"
	case RiskLow:
		return "low"
	case RiskMedium:
		return "medium"
	case RiskHigh:
		return "high"
	case RiskCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// ParseRiskLevel 将字符串解析为风险等级。
func ParseRiskLevel(s string) (RiskLevel, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "none":
		return RiskNone, nil
	case "low":
		return RiskLow, nil
	case "medium":
		return RiskMedium, nil
	case "high":
		return RiskHigh, nil
	case "critical":
		return RiskCritical, nil
	case "unknown":
		return RiskUnknown, nil
	default:
		return RiskUnknown, fmt.Errorf("未知风险等级: %s", s)
	}
}

// Policy 表示风险策略配置。
type Policy struct {
	Name          string
	Description   string
	Allowed       []string
	Conditional   []ConditionalRule
	Forbidden     []ForbiddenRule
	UnknownPolicy string `yaml:"unknown_policy"`
}

// ConditionalRule 表示有条件可用协议规则。
type ConditionalRule struct {
	ID        string `yaml:"id"`
	Condition string `yaml:"condition"`
	Note      string `yaml:"note"`
}

// ForbiddenRule 表示禁止协议规则。
type ForbiddenRule struct {
	ID     string `yaml:"id"`
	Reason string `yaml:"reason"`
}

type policyRoot struct {
	Policy Policy `yaml:"policy"`
}

// DefaultPolicy 返回内置默认策略。
func DefaultPolicy() *Policy {
	return &Policy{
		Name:          "commercial-software",
		Description:   "商业软件使用的开源协议风险策略",
		Allowed:       []string{"MIT", "Apache-2.0", "BSD-2-Clause", "BSD-3-Clause", "ISC", "0BSD", "Unlicense", "CC0-1.0", "Zlib"},
		Conditional:   []ConditionalRule{{ID: "LGPL-2.1", Condition: "dynamic-link-only", Note: "仅动态链接时可用于商业软件"}, {ID: "LGPL-3.0", Condition: "dynamic-link-only", Note: "仅动态链接时可用于商业软件"}, {ID: "MPL-2.0", Condition: "file-level-copyleft", Note: "修改的文件需开源，但可与私有代码结合"}},
		Forbidden:     []ForbiddenRule{{ID: "GPL-2.0", Reason: "强 Copyleft，会传染至整个项目"}, {ID: "GPL-3.0", Reason: "强 Copyleft，会传染至整个项目"}, {ID: "AGPL-3.0", Reason: "网络服务也需要开源，商业 SaaS 产品风险极高"}, {ID: "SSPL-1.0", Reason: "服务条款极为严格"}},
		UnknownPolicy: "warn",
	}
}

// LoadPolicy 从 YAML 文件加载策略。
func LoadPolicy(path string) (*Policy, error) {
	if path == "" {
		return DefaultPolicy(), nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var root policyRoot
	if err := yaml.Unmarshal(b, &root); err != nil {
		return nil, err
	}
	p := root.Policy
	if p.Name == "" {
		return nil, fmt.Errorf("策略名称不能为空")
	}
	return &p, nil
}

// EvaluateRisk 根据策略评估协议风险。
func EvaluateRisk(licenseID string, policy *Policy) (RiskLevel, string) {
	id := NormalizeLicenseID(licenseID)
	if policy == nil {
		policy = DefaultPolicy()
	}
	for _, a := range policy.Allowed {
		if NormalizeLicenseID(a) == id {
			return RiskNone, "策略允许的宽松协议"
		}
	}
	for _, c := range policy.Conditional {
		if NormalizeLicenseID(c.ID) == id {
			if strings.Contains(strings.ToLower(c.Condition), "dynamic") {
				return RiskLow, c.Note
			}
			return RiskMedium, c.Note
		}
	}
	for _, f := range policy.Forbidden {
		if NormalizeLicenseID(f.ID) == id {
			if id == "AGPL-3.0" || id == "SSPL-1.0" {
				return RiskCritical, f.Reason
			}
			return RiskHigh, f.Reason
		}
	}
	if id == "UNKNOWN" || !IsKnownLicense(id) {
		switch strings.ToLower(policy.UnknownPolicy) {
		case "allow":
			return RiskNone, "未知协议按策略放行"
		case "block":
			return RiskHigh, "未知协议按策略阻断"
		default:
			return RiskUnknown, "未知协议需人工审查"
		}
	}
	return RiskUnknown, "协议未命中策略规则"
}
