package license

import "testing"

func TestEvaluateRisk(t *testing.T) {
	p := DefaultPolicy()
	tests := []struct {
		id   string
		want RiskLevel
	}{
		{id: "MIT", want: RiskNone},
		{id: "LGPL-2.1", want: RiskLow},
		{id: "MPL-2.0", want: RiskMedium},
		{id: "GPL-3.0", want: RiskHigh},
		{id: "AGPL-3.0", want: RiskCritical},
		{id: "Unknown-License", want: RiskUnknown},
	}
	for _, tc := range tests {
		got, _ := EvaluateRisk(tc.id, p)
		if got != tc.want {
			t.Fatalf("EvaluateRisk(%s)=%s, want %s", tc.id, got.String(), tc.want.String())
		}
	}
}

func TestParseRiskLevel(t *testing.T) {
	v, err := ParseRiskLevel("high")
	if err != nil || v != RiskHigh {
		t.Fatalf("ParseRiskLevel(high) failed: %v, %v", v, err)
	}
	if _, err := ParseRiskLevel("nope"); err == nil {
		t.Fatal("ParseRiskLevel(nope) 应返回错误")
	}
}
