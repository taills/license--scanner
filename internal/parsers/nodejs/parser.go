package nodejs

import (
"context"
"encoding/json"
"os"
"path/filepath"
"strings"

"github.com/taills/license--scanner/internal/parsers"
)

// Parser 实现 Node.js 项目依赖解析。
type Parser struct{}

func (p *Parser) Name() string { return "npm" }

func (p *Parser) Supports(projectPath string) bool {
files := []string{"package.json", "package-lock.json", "yarn.lock", "pnpm-lock.yaml"}
for _, f := range files {
if _, err := os.Stat(filepath.Join(projectPath, f)); err == nil {
return true
}
}
return false
}

func (p *Parser) Parse(ctx context.Context, projectPath string) ([]parsers.ParsedDependency, error) {
_ = ctx
directMap, err := readPackageJSONDirect(projectPath)
if err != nil {
directMap = map[string]bool{}
}
lock := filepath.Join(projectPath, "package-lock.json")
if _, err := os.Stat(lock); err == nil {
return parsePackageLock(lock, directMap)
}
result := make([]parsers.ParsedDependency, 0, len(directMap))
for name := range directMap {
result = append(result, parsers.ParsedDependency{Name: name, Version: "", Ecosystem: "npm", Direct: true, DepPath: []string{name}})
}
return result, nil
}

func readPackageJSONDirect(projectPath string) (map[string]bool, error) {
b, err := os.ReadFile(filepath.Join(projectPath, "package.json"))
if err != nil {
return nil, err
}
var pkg struct {
Dependencies    map[string]string `json:"dependencies"`
DevDependencies map[string]string `json:"devDependencies"`
}
if err := json.Unmarshal(b, &pkg); err != nil {
return nil, err
}
result := map[string]bool{}
for k := range pkg.Dependencies {
result[k] = true
}
for k := range pkg.DevDependencies {
if _, ok := result[k]; !ok {
result[k] = false
}
}
return result, nil
}

func parsePackageLock(path string, direct map[string]bool) ([]parsers.ParsedDependency, error) {
b, err := os.ReadFile(path)
if err != nil {
return nil, err
}
var lock struct {
Packages map[string]struct {
Version string `json:"version"`
} `json:"packages"`
}
if err := json.Unmarshal(b, &lock); err != nil {
return nil, err
}
result := make([]parsers.ParsedDependency, 0, len(lock.Packages))
for k, v := range lock.Packages {
if k == "" {
continue
}
name := k
if idx := strings.LastIndex(k, "node_modules/"); idx >= 0 {
name = k[idx+len("node_modules/"):]
}
result = append(result, parsers.ParsedDependency{Name: name, Version: v.Version, Ecosystem: "npm", Direct: direct[name], DepPath: []string{name}})
}
return result, nil
}
