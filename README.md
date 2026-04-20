# license--scanner

`license--scanner` 是一个使用 Go 开发的项目依赖开源协议扫描工具，用于商业软件研发阶段的开源合规风险识别与预警。

## 1. 项目简介与使用场景

在商业软件中引入第三方依赖时，不同协议（MIT、Apache、GPL、AGPL 等）对应的义务不同。该工具支持对多生态依赖进行协议扫描与风险分级，适用于：

- 上线前开源合规检查
- CI/CD 自动阻断高风险依赖
- 法务与研发协同审计

支持生态：

- Go（`go.mod` / `go.sum`）
- Node.js（`package.json` / `package-lock.json` / `yarn.lock` / `pnpm-lock.yaml`）
- Java（`pom.xml` / `build.gradle`）
- Python（`requirements.txt` / `Pipfile.lock` / `pyproject.toml` / `poetry.lock`）

## 2. 安装方式

### 2.1 go install

```bash
go install github.com/taills/license--scanner/cmd/scanner@latest
```

### 2.2 预编译二进制

从 Releases 下载对应平台二进制后重命名为 `license-scanner` 并加入 `PATH`。

## 3. 快速开始示例

```bash
# 扫描当前目录
license-scanner scan .

# 输出 JSON 报告
license-scanner scan . --format json --output report.json

# 输出 HTML 报告
license-scanner scan . --format html --output report.html

# 输出 SARIF 报告
license-scanner scan . --format sarif --output results.sarif

# 使用自定义策略
license-scanner scan . --policy ./configs/risk_policy.yaml

# CI 模式：高风险退出非零
license-scanner scan . --fail-on high

# 只扫描 Go + npm
license-scanner scan . --ecosystem go,npm

# 显示版本
license-scanner version
```

## 4. CLI 参数说明

| 参数 | 说明 |
| --- | --- |
| `--format` | 报告格式，支持 `text/json/html/sarif` |
| `--output` | 输出文件路径，不指定则输出到终端 |
| `--policy` | 自定义风险策略 YAML |
| `--fail-on` | 风险阈值（`none/low/medium/high/critical/unknown`） |
| `--ecosystem` | 仅扫描指定生态（逗号分隔） |

## 5. 风险等级说明

| 等级 | 含义 | 典型协议 |
| --- | --- | --- |
| `RiskNone` | 安全，可直接使用 | MIT / Apache-2.0 / BSD |
| `RiskLow` | 低风险，有条件可用 | LGPL |
| `RiskMedium` | 中风险，文件级约束 | MPL-2.0 |
| `RiskHigh` | 高风险，强 Copyleft | GPL-2.0 / GPL-3.0 |
| `RiskCritical` | 严重风险，网络服务约束 | AGPL-3.0 / SSPL-1.0 |
| `RiskUnknown` | 未知，需人工审查 | 未识别协议 |

## 6. 商业使用风险说明

- **宽松协议（MIT/Apache/BSD）**：通常可闭源商用。
- **弱 Copyleft（LGPL/MPL）**：通常需满足动态链接或文件级开源义务。
- **强 Copyleft（GPL）**：可能要求衍生作品整体开源。
- **网络 Copyleft（AGPL/SSPL）**：SaaS 场景风险极高，建议避免。

## 7. 贡献指南

1. Fork 仓库并创建特性分支
2. 提交前运行：`make test && make lint`
3. 提交 PR 并补充变更说明
4. 通过 CI 后进行 Code Review 与合并
