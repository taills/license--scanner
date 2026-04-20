package python

import (
"bufio"
"context"
"encoding/json"
"os"
"path/filepath"
"strings"

"github.com/taills/license--scanner/internal/parsers"
)

// Parser 实现 Python 项目依赖解析。
type Parser struct{}

func (p *Parser) Name() string { return "pypi" }

func (p *Parser) Supports(projectPath string) bool {
files := []string{"requirements.txt", "Pipfile.lock", "pyproject.toml", "poetry.lock"}
for _, f := range files {
if _, err := os.Stat(filepath.Join(projectPath, f)); err == nil {
return true
}
}
return false
}

func (p *Parser) Parse(ctx context.Context, projectPath string) ([]parsers.ParsedDependency, error) {
_ = ctx
depMap := map[string]parsers.ParsedDependency{}
if deps, err := parseRequirements(filepath.Join(projectPath, "requirements.txt")); err == nil {
for _, d := range deps {
depMap[d.Name] = d
}
}
if deps, err := parsePipfileLock(filepath.Join(projectPath, "Pipfile.lock")); err == nil {
for _, d := range deps {
depMap[d.Name] = d
}
}
if deps, err := parsePyProject(filepath.Join(projectPath, "pyproject.toml")); err == nil {
for _, d := range deps {
if _, ok := depMap[d.Name]; !ok {
depMap[d.Name] = d
}
}
}
if deps, err := parsePoetryLock(filepath.Join(projectPath, "poetry.lock")); err == nil {
for _, d := range deps {
depMap[d.Name] = d
}
}
result := make([]parsers.ParsedDependency, 0, len(depMap))
for _, d := range depMap {
result = append(result, d)
}
return result, nil
}

func parseRequirements(path string) ([]parsers.ParsedDependency, error) {
f, err := os.Open(path)
if err != nil {
return nil, err
}
defer f.Close()
var result []parsers.ParsedDependency
s := bufio.NewScanner(f)
for s.Scan() {
line := strings.TrimSpace(s.Text())
if line == "" || strings.HasPrefix(line, "#") {
continue
}
parts := strings.SplitN(line, "==", 2)
name, version := parts[0], ""
if len(parts) == 2 {
version = parts[1]
}
result = append(result, parsers.ParsedDependency{Name: name, Version: version, Ecosystem: "pypi", Direct: true, DepPath: []string{name}})
}
return result, s.Err()
}

func parsePipfileLock(path string) ([]parsers.ParsedDependency, error) {
b, err := os.ReadFile(path)
if err != nil {
return nil, err
}
var root map[string]map[string]map[string]any
if err := json.Unmarshal(b, &root); err != nil {
return nil, err
}
sections := []string{"default", "develop"}
var result []parsers.ParsedDependency
for _, sec := range sections {
for name, meta := range root[sec] {
version, _ := meta["version"].(string)
version = strings.TrimPrefix(version, "==")
result = append(result, parsers.ParsedDependency{Name: name, Version: version, Ecosystem: "pypi", Direct: sec == "default", DepPath: []string{name}})
}
}
return result, nil
}

func parsePyProject(path string) ([]parsers.ParsedDependency, error) {
f, err := os.Open(path)
if err != nil {
return nil, err
}
defer f.Close()
var result []parsers.ParsedDependency
inPoetry := false
s := bufio.NewScanner(f)
for s.Scan() {
line := strings.TrimSpace(s.Text())
if strings.HasPrefix(line, "[") {
inPoetry = line == "[tool.poetry.dependencies]"
continue
}
if !inPoetry || line == "" || strings.HasPrefix(line, "#") {
continue
}
parts := strings.SplitN(line, "=", 2)
if len(parts) != 2 {
continue
}
name := strings.TrimSpace(parts[0])
if name == "python" {
continue
}
version := strings.Trim(strings.TrimSpace(parts[1]), "\"'")
result = append(result, parsers.ParsedDependency{Name: name, Version: version, Ecosystem: "pypi", Direct: true, DepPath: []string{name}})
}
return result, s.Err()
}

func parsePoetryLock(path string) ([]parsers.ParsedDependency, error) {
f, err := os.Open(path)
if err != nil {
return nil, err
}
defer f.Close()
var result []parsers.ParsedDependency
var name, version string
s := bufio.NewScanner(f)
for s.Scan() {
line := strings.TrimSpace(s.Text())
switch {
case line == "[[package]]":
if name != "" {
result = append(result, parsers.ParsedDependency{Name: name, Version: version, Ecosystem: "pypi", Direct: true, DepPath: []string{name}})
}
name, version = "", ""
case strings.HasPrefix(line, "name ="):
name = strings.Trim(strings.TrimSpace(strings.TrimPrefix(line, "name =")), "\"")
case strings.HasPrefix(line, "version ="):
version = strings.Trim(strings.TrimSpace(strings.TrimPrefix(line, "version =")), "\"")
}
}
if name != "" {
result = append(result, parsers.ParsedDependency{Name: name, Version: version, Ecosystem: "pypi", Direct: true, DepPath: []string{name}})
}
return result, s.Err()
}
