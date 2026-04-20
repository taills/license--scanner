package spdx

// LicenseMeta 表示 SPDX 协议元数据。
type LicenseMeta struct {
	ID          string
	Name        string
	Copyleft    string
	Commercial  string
	Description string
}

// CommonLicenses 是常见协议元数据表。
var CommonLicenses = map[string]LicenseMeta{
	"MIT":          {ID: "MIT", Name: "MIT License", Copyleft: "none", Commercial: "allow", Description: "宽松协议，允许商业使用"},
	"Apache-2.0":   {ID: "Apache-2.0", Name: "Apache License 2.0", Copyleft: "none", Commercial: "allow", Description: "宽松协议，含专利授权"},
	"BSD-2-Clause": {ID: "BSD-2-Clause", Name: "BSD 2-Clause", Copyleft: "none", Commercial: "allow", Description: "宽松协议"},
	"BSD-3-Clause": {ID: "BSD-3-Clause", Name: "BSD 3-Clause", Copyleft: "none", Commercial: "allow", Description: "宽松协议"},
	"ISC":          {ID: "ISC", Name: "ISC License", Copyleft: "none", Commercial: "allow", Description: "宽松协议"},
	"0BSD":         {ID: "0BSD", Name: "Zero-Clause BSD", Copyleft: "none", Commercial: "allow", Description: "宽松协议"},
	"Unlicense":    {ID: "Unlicense", Name: "The Unlicense", Copyleft: "none", Commercial: "allow", Description: "公有领域协议"},
	"CC0-1.0":      {ID: "CC0-1.0", Name: "Creative Commons Zero", Copyleft: "none", Commercial: "allow", Description: "公有领域声明"},
	"Zlib":         {ID: "Zlib", Name: "zlib License", Copyleft: "none", Commercial: "allow", Description: "宽松协议"},
	"LGPL-2.1":     {ID: "LGPL-2.1", Name: "GNU LGPL v2.1", Copyleft: "weak", Commercial: "conditional", Description: "弱 Copyleft，通常要求动态链接"},
	"LGPL-3.0":     {ID: "LGPL-3.0", Name: "GNU LGPL v3.0", Copyleft: "weak", Commercial: "conditional", Description: "弱 Copyleft，通常要求动态链接"},
	"MPL-2.0":      {ID: "MPL-2.0", Name: "Mozilla Public License 2.0", Copyleft: "file", Commercial: "conditional", Description: "文件级 Copyleft"},
	"GPL-2.0":      {ID: "GPL-2.0", Name: "GNU GPL v2.0", Copyleft: "strong", Commercial: "high-risk", Description: "强 Copyleft"},
	"GPL-3.0":      {ID: "GPL-3.0", Name: "GNU GPL v3.0", Copyleft: "strong", Commercial: "high-risk", Description: "强 Copyleft"},
	"AGPL-3.0":     {ID: "AGPL-3.0", Name: "GNU AGPL v3.0", Copyleft: "network", Commercial: "critical", Description: "网络服务场景也需开源"},
	"SSPL-1.0":     {ID: "SSPL-1.0", Name: "Server Side Public License", Copyleft: "network", Commercial: "critical", Description: "服务条款严格"},
}
