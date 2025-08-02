# 📘 SpecGrade

A modular, dynamic, and CICD-optimized conformance validator for OpenAPI specifications.

## 🚀 Features

- **Modular Rule System**: Pluggable validation rules that can be easily extended
- **Multiple Output Formats**: JSON, CLI, HTML, and Markdown reporting
- **CI/CD Integration**: Semantic exit codes and configurable fail thresholds
- **YAML Configuration**: Team-wide standardization with `specgrade.yaml`
- **Rule Management**: Skip specific rules or generate documentation
- **Version Support**: OpenAPI 3.0.0 and 3.1.0 specifications

## 📦 Installation

```bash
go install github.com/codetestcode/specgrade@latest
```

Or build from source:

```bash
git clone https://github.com/codetestcode/specgrade.git
cd specgrade
go build -o specgrade
```

## 🧑‍💻 Usage

### Basic Usage

```bash
specgrade --target-dir=./specs/openai --spec-version=3.1.0
```

### Advanced Usage

```bash
specgrade \
  --spec-version=3.1.0 \
  --target-dir=./specs/openai \
  --output-format=json \
  --fail-threshold=B \
  --skip=RULE001,RULE012
```

### Using Configuration File

Create a `specgrade.yaml` file:

```yaml
spec_version: 3.1.0
input_dir: ./specs/openai
fail_threshold: B
output_format: json
skip_rules:
  - RULE001
  - RULE012
```

Then run:

```bash
specgrade --target-dir=./specs/openai
```

### Rule Management

List all available rules:

```bash
specgrade rules ls
```

Generate rule documentation:

```bash
specgrade --docs
```

## 📊 Output Formats

### CLI Output
```
📄 Validating: ./specs/openai
🔖 Spec: OpenAPI 3.1.0
✅ Passed: 25/30 rules
🎯 Score: 83%
🏅 Grade: B
```

### JSON Output
```json
{
  "version": "3.1.0",
  "grade": "B",
  "score": 83,
  "rules": [
    { "ruleID": "INFO-001", "passed": true, "detail": "Title present" },
    { "ruleID": "PATHS-001", "passed": false, "detail": "No paths defined" }
  ]
}
```

## 🔧 Configuration

### CLI Flags

| Flag               | Description                                                            |
| ------------------ | ---------------------------------------------------------------------- |
| `--spec-version`   | The official OpenAPI version to validate against (e.g., `3.1.0`)       |
| `--target-dir`     | Path to the local OpenAPI spec to validate                             |
| `--output-format`  | `json`, `cli`, `html`, or `markdown`                                   |
| `--fail-threshold` | Minimum acceptable grade (`A`, `B`, etc). Will exit non-zero if below. |
| `--config`         | Optional path to `specgrade.yaml` config file                          |
| `--skip`           | Comma-separated rule IDs to ignore                                     |
| `--docs`           | Generate rule documentation (markdown)                                 |

### Configuration Precedence

1. CLI flags (highest priority)
2. `specgrade.yaml` configuration file
3. Default values (lowest priority)

## 📋 Available Rules

| Rule ID   | Description                                    | Applies To    |
| --------- | ---------------------------------------------- | ------------- |
| INFO-001  | OpenAPI spec must have a title in info section | 3.0.0, 3.1.0 |
| INFO-002  | OpenAPI spec must have a version in info section | 3.0.0, 3.1.0 |
| PATHS-001 | OpenAPI spec must have at least one path defined | 3.0.0, 3.1.0 |
| OPID-001  | All operations should have unique operation IDs | 3.0.0, 3.1.0 |

## 🏅 Grading System

| Grade | Score Range | Description |
| ----- | ----------- | ----------- |
| A+    | 95-100%     | Excellent   |
| A     | 90-94%      | Very Good   |
| A-    | 85-89%      | Good        |
| B+    | 80-84%      | Above Average |
| B     | 75-79%      | Average     |
| B-    | 70-74%      | Below Average |
| C+    | 65-69%      | Poor        |
| C     | 60-64%      | Very Poor   |
| C-    | 55-59%      | Failing     |
| D     | 50-54%      | Very Failing |
| F     | 0-49%       | Complete Failure |

## 🔄 CI/CD Integration

SpecGrade returns appropriate exit codes for CI/CD integration:

- **Exit Code 0**: Grade meets or exceeds the fail threshold
- **Exit Code 1**: Grade is below the fail threshold or validation error

### GitHub Actions Example

```yaml
name: OpenAPI Validation
on: [push, pull_request]

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Install SpecGrade
        run: go install github.com/codetestcode/specgrade@latest
      - name: Validate OpenAPI Spec
        run: specgrade --target-dir=./api --fail-threshold=B
```

## 🧪 Development

### Running Tests

```bash
go test ./...
```

### Adding New Rules

1. Create a new rule struct implementing the `core.Rule` interface
2. Register the rule in `cmd/root.go`
3. Add tests for the new rule

Example rule:

```go
type MyCustomRule struct{}

func (r *MyCustomRule) ID() string {
    return "CUSTOM-001"
}

func (r *MyCustomRule) Description() string {
    return "My custom validation rule"
}

func (r *MyCustomRule) AppliesTo(version string) bool {
    return strings.HasPrefix(version, "3.")
}

func (r *MyCustomRule) Evaluate(ctx *core.SpecContext) core.RuleResult {
    // Your validation logic here
    return core.RuleResult{
        RuleID: r.ID(),
        Passed: true,
        Detail: "Validation passed",
    }
}
```

## 📄 License

MIT License - see LICENSE file for details.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for your changes
4. Ensure all tests pass
5. Submit a pull request

## 📞 Support

- GitHub Issues: [Report bugs or request features](https://github.com/codetestcode/specgrade/issues)
- Documentation: [Full documentation](https://github.com/codetestcode/specgrade/wiki)
