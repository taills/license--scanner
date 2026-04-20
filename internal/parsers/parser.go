package parsers

import "context"

// ParsedDependency 表示解析器输出的依赖信息。
type ParsedDependency struct {
	Name      string
	Version   string
	Ecosystem string
	Direct    bool
	DepPath   []string
}

// Parser 定义统一的项目依赖解析接口。
type Parser interface {
	Name() string
	Supports(projectPath string) bool
	Parse(ctx context.Context, projectPath string) ([]ParsedDependency, error)
}
