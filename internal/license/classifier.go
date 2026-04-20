package license

import (
	"strings"

	"github.com/taills/license--scanner/pkg/spdx"
)

var aliases = map[string]string{
	"apache2":            "Apache-2.0",
	"apache-2":           "Apache-2.0",
	"apache 2.0":         "Apache-2.0",
	"mit license":        "MIT",
	"bsd-3":              "BSD-3-Clause",
	"bsd 3-clause":       "BSD-3-Clause",
	"bsd-2":              "BSD-2-Clause",
	"mozilla public 2.0": "MPL-2.0",
	"gplv2":              "GPL-2.0",
	"gplv3":              "GPL-3.0",
	"agplv3":             "AGPL-3.0",
	"lgplv2.1":           "LGPL-2.1",
	"lgplv3":             "LGPL-3.0",
}

// NormalizeLicenseID 归一化协议 ID 到 SPDX 标准形式。
func NormalizeLicenseID(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "UNKNOWN"
	}
	if _, ok := spdx.CommonLicenses[raw]; ok {
		return raw
	}
	k := strings.ToLower(raw)
	if v, ok := aliases[k]; ok {
		return v
	}
	for id := range spdx.CommonLicenses {
		if strings.EqualFold(id, raw) {
			return id
		}
	}
	return strings.ToUpper(strings.ReplaceAll(raw, " ", "-"))
}

// IsKnownLicense 判断协议是否在内置 SPDX 列表中。
func IsKnownLicense(id string) bool {
	_, ok := spdx.CommonLicenses[NormalizeLicenseID(id)]
	return ok
}
