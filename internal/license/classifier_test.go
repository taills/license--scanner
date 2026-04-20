package license

import "testing"

func TestNormalizeLicenseID(t *testing.T) {
	cases := map[string]string{
		"MIT":         "MIT",
		"Apache 2.0":  "Apache-2.0",
		"mit license": "MIT",
		"":            "UNKNOWN",
	}
	for in, want := range cases {
		if got := NormalizeLicenseID(in); got != want {
			t.Fatalf("NormalizeLicenseID(%q)=%q, want %q", in, got, want)
		}
	}
}

func TestIsKnownLicense(t *testing.T) {
	if !IsKnownLicense("MIT") {
		t.Fatal("MIT 应被识别为已知协议")
	}
	if IsKnownLicense("UNKNOWN-XYZ") {
		t.Fatal("UNKNOWN-XYZ 不应被识别为已知协议")
	}
}
