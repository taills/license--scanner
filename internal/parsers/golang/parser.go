package golang

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/mod/modfile"

	"github.com/taills/license--scanner/internal/parsers"
)

// Parser 实现 Go 项目依赖解析。
type Parser struct{}

func (p *Parser) Name() string { return "go" }

func (p *Parser) Supports(projectPath string) bool {
	_, err := os.Stat(filepath.Join(projectPath, "go.mod"))
	return err == nil
}

func (p *Parser) Parse(ctx context.Context, projectPath string) ([]parsers.ParsedDependency, error) {
	modPath := filepath.Join(projectPath, "go.mod")
	modData, err := os.ReadFile(modPath)
	if err != nil {
		return nil, err
	}
	mf, err := modfile.Parse("go.mod", modData, nil)
	if err != nil {
		return nil, err
	}
	direct := map[string]bool{}
	for _, r := range mf.Require {
		direct[r.Mod.Path] = !r.Indirect
	}
	mainModule := "main"
	if mf.Module != nil {
		mainModule = mf.Module.Mod.Path
	}

	cmd := exec.CommandContext(ctx, "go", "list", "-m", "-json", "all")
	cmd.Dir = projectPath
	out, err := cmd.Output()
	if err != nil {
		result := make([]parsers.ParsedDependency, 0, len(mf.Require))
		for _, r := range mf.Require {
			result = append(result, parsers.ParsedDependency{Name: r.Mod.Path, Version: r.Mod.Version, Ecosystem: "go", Direct: !r.Indirect, DepPath: []string{mainModule, r.Mod.Path}})
		}
		return result, nil
	}

	var deps []parsers.ParsedDependency
	dec := json.NewDecoder(bytes.NewReader(out))
	for {
		var m struct {
			Path    string
			Version string
			Main    bool
		}
		if err := dec.Decode(&m); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		if !m.Main {
			deps = append(deps, parsers.ParsedDependency{Name: m.Path, Version: m.Version, Ecosystem: "go", Direct: direct[m.Path], DepPath: []string{mainModule, m.Path}})
		}
	}
	return deps, nil
}
