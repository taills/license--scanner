package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/taills/license--scanner/internal/analyzer"
	"github.com/taills/license--scanner/internal/license"
	"github.com/taills/license--scanner/internal/parsers"
	goparser "github.com/taills/license--scanner/internal/parsers/golang"
	javaparser "github.com/taills/license--scanner/internal/parsers/java"
	nodeparser "github.com/taills/license--scanner/internal/parsers/nodejs"
	pythonparser "github.com/taills/license--scanner/internal/parsers/python"
	"github.com/taills/license--scanner/internal/reporter"
)

var version = "dev"

func main() {
	root := &cobra.Command{
		Use:   "license-scanner",
		Short: "用于商业软件风险识别的开源协议扫描工具",
	}

	root.AddCommand(newScanCmd(), &cobra.Command{
		Use:   "version",
		Short: "显示版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), version)
		},
	})

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newScanCmd() *cobra.Command {
	var format string
	var output string
	var policyPath string
	var failOn string
	var ecosystemCSV string
	cmd := &cobra.Command{
		Use:   "scan [path]",
		Short: "扫描指定项目目录的依赖协议",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			viper.Set("policy", policyPath)
			policy, err := license.LoadPolicy(viper.GetString("policy"))
			if err != nil {
				return err
			}

			detector := license.NewDetector("")
			eng := analyzer.Engine{
				Parsers:  []parsers.Parser{&goparser.Parser{}, &nodeparser.Parser{}, &javaparser.Parser{}, &pythonparser.Parser{}},
				Detector: detector,
				Policy:   policy,
			}

			ecFilter := map[string]bool{}
			for _, e := range strings.Split(strings.TrimSpace(ecosystemCSV), ",") {
				e = strings.TrimSpace(e)
				if e != "" {
					ecFilter[e] = true
				}
			}
			result, err := eng.Analyze(cmd.Context(), path, ecFilter)
			if err != nil {
				return err
			}

			r, err := selectReporter(format)
			if err != nil {
				return err
			}
			var out = cmd.OutOrStdout()
			if output != "" {
				f, err := os.Create(output)
				if err != nil {
					return err
				}
				defer f.Close()
				out = f
			}
			if err := r.Report(result, out); err != nil {
				return err
			}
			if failOn != "" {
				threshold, err := license.ParseRiskLevel(failOn)
				if err != nil {
					return err
				}
				for _, d := range result.Dependencies {
					if d.Risk >= threshold {
						return fmt.Errorf("检测到风险等级 >= %s 的依赖: %s@%s(%s)", threshold.String(), d.Name, d.Version, d.Risk.String())
					}
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "text", "报告格式：text|json|html|sarif")
	cmd.Flags().StringVar(&output, "output", "", "输出文件路径")
	cmd.Flags().StringVar(&policyPath, "policy", "", "风险策略配置路径")
	cmd.Flags().StringVar(&failOn, "fail-on", "", "达到指定风险级别时返回非零退出码")
	cmd.Flags().StringVar(&ecosystemCSV, "ecosystem", "", "仅扫描指定生态（逗号分隔）：go,npm,maven,pypi")
	return cmd
}

func selectReporter(format string) (reporter.Reporter, error) {
	switch strings.ToLower(format) {
	case "text":
		return &reporter.TextReporter{}, nil
	case "json":
		return &reporter.JSONReporter{}, nil
	case "html":
		return &reporter.HTMLReporter{}, nil
	case "sarif":
		return &reporter.SARIFReporter{}, nil
	default:
		return nil, fmt.Errorf("不支持的格式: %s", format)
	}
}
