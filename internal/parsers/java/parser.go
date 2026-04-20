package java

import (
"bufio"
"context"
"encoding/xml"
"os"
"path/filepath"
"regexp"
"strings"

"github.com/taills/license--scanner/internal/parsers"
)

// Parser 实现 Java 项目依赖解析。
type Parser struct{}

func (p *Parser) Name() string { return "maven" }

func (p *Parser) Supports(projectPath string) bool {
_, pomErr := os.Stat(filepath.Join(projectPath, "pom.xml"))
_, gradleErr := os.Stat(filepath.Join(projectPath, "build.gradle"))
return pomErr == nil || gradleErr == nil
}

func (p *Parser) Parse(ctx context.Context, projectPath string) ([]parsers.ParsedDependency, error) {
_ = ctx
pom := filepath.Join(projectPath, "pom.xml")
if _, err := os.Stat(pom); err == nil {
return parsePOM(pom)
}
gradle := filepath.Join(projectPath, "build.gradle")
return parseGradle(gradle)
}

func parsePOM(path string) ([]parsers.ParsedDependency, error) {
b, err := os.ReadFile(path)
if err != nil {
return nil, err
}
var pom struct {
Dependencies []struct {
GroupID    string `xml:"groupId"`
ArtifactID string `xml:"artifactId"`
Version    string `xml:"version"`
} `xml:"dependencies>dependency"`
}
if err := xml.Unmarshal(b, &pom); err != nil {
return nil, err
}
result := make([]parsers.ParsedDependency, 0, len(pom.Dependencies))
for _, d := range pom.Dependencies {
result = append(result, parsers.ParsedDependency{Name: d.GroupID + ":" + d.ArtifactID, Version: d.Version, Ecosystem: "maven", Direct: true, DepPath: []string{d.GroupID + ":" + d.ArtifactID}})
}
return result, nil
}

func parseGradle(path string) ([]parsers.ParsedDependency, error) {
f, err := os.Open(path)
if err != nil {
return nil, err
}
defer f.Close()
re := regexp.MustCompile(`['"]([^:'"]+:[^:'"]+:[^:'"]+)['"]`)
var result []parsers.ParsedDependency
s := bufio.NewScanner(f)
for s.Scan() {
line := strings.TrimSpace(s.Text())
if !strings.Contains(line, "implementation") && !strings.Contains(line, "api") {
continue
}
m := re.FindStringSubmatch(line)
if len(m) != 2 {
continue
}
parts := strings.Split(m[1], ":")
if len(parts) != 3 {
continue
}
name := parts[0] + ":" + parts[1]
result = append(result, parsers.ParsedDependency{Name: name, Version: parts[2], Ecosystem: "maven", Direct: true, DepPath: []string{name}})
}
return result, s.Err()
}
