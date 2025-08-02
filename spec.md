# 📘 SpecGrade: Full System Specification v1.0

## 🛍️ Mission Statement

**SpecGrade** is a modular, dynamic, and CICD-optimized conformance validator for OpenAPI specifications. It fetches versioned OpenAPI schema definitions, dynamically constructs validation rule sets based on that schema, and grades the conformance of target API specs against those rules. It is designed to be easy to maintain, extend, and integrate into automated pipelines.

---

## ✅ CLI Tool Design

### 🧑‍💻 Example Command

```bash
specgrade \
  --spec-version=3.1.0 \
  --target-dir=./specs/openai \
  --output-format=json \
  --fail-threshold=B
```

### Flags

| Flag               | Description                                                            |
| ------------------ | ---------------------------------------------------------------------- |
| `--spec-version`   | The official OpenAPI version to validate against (e.g., `3.1.0`)       |
| `--target-dir`     | Path to the local OpenAPI spec to validate                             |
| `--output-format`  | `json`, `cli`, `html`, or `markdown`                                   |
| `--fail-threshold` | Minimum acceptable grade (`A`, `B`, etc). Will exit non-zero if below. |
| `--config`         | Optional path to `specgrade.yaml` config file                          |
| `--skip`           | Comma-separated rule IDs to ignore                                     |
| `--docs`           | Generate rule documentation (markdown)                                 |

---

## 📂 Project Structure

```
specgrade/
├── cmd/             # CLI entry
├── core/            # Shared types/interfaces
├── registry/        # Rule discovery & filtering
├── rules/           # Pluggable rules
├── runner/          # Execution engine
├── fetcher/         # Downloads + parses OpenAPI spec
├── versions/        # Maps spec version to schema URL
├── reporter/        # Grading + output rendering
├── ci/              # Exit code strategy
├── utils/           # Helpers
├── test/            # Unit and property-based tests
└── go.mod
```

---

## 🔌 Rule Interface

```go
type Rule interface {
    ID() string
    Description() string
    AppliesTo(version string) bool
    Evaluate(ctx *SpecContext) RuleResult
}
```

## 📜 Rule Registry

```go
type RuleRegistry struct {
    rules []Rule
}
func (r *RuleRegistry) Register(rule Rule)
func (r *RuleRegistry) RulesForVersion(version string) []Rule
```

## 🧠 RuleResult

```go
type RuleResult struct {
    RuleID string
    Passed bool
    Detail string
}
```

## 📆 Spec Loader Strategy

```go
type SpecLoader interface {
    Load(version string) (*openapi3.T, error)
}
```

## �� Grader Strategy

```go
type Grader interface {
    Grade([]RuleResult) string // Returns A, B, C, etc
}
```

## 🚦 CI Exit Strategy

```go
type ExitHandler interface {
    Handle(grade string) int
}
```

## 🌍 Version → Schema URL Map

```go
var VersionToSchemaURL = map[string]string{
    "3.0.0": "https://spec.openapis.org/oas/3.0/schema/2019-04-02",
    "3.1.0": "https://spec.openapis.org/oas/3.1/schema/2022-10-07",
}
```

---

## 🗒️ Output Examples

### ✅ JSON

```json
{
  "version": "3.1.0",
  "grade": "B",
  "score": 83,
  "rules": [
    { "ruleID": "OPID-001", "passed": true, "detail": "OK" },
    { "ruleID": "RESP-002", "passed": false, "detail": "Missing 400 response" }
  ]
}
```

### ✅ CLI

```
📄 Validating: ./specs/openai
🔖 Spec: OpenAPI 3.1.0
✅ Passed: 25/30 rules
🎯 Score: 83%
🏅 Grade: B
```

---

## 🧪 Testing Strategy

* ✅ Unit tests per rule with coverage of pass/fail/edge
* ✅ Table-driven tests for runners, graders, fetchers
* ✅ Property-based testing of runner via `testing/quick`
* ✅ Schema error coverage (missing fields, malformed YAML)
* ✅ Threshold boundary tests (89.99% vs 90.00%)
* ✅ Config parsing + precedence (CLI vs YAML)
* ✅ End-to-end regression tests with known OpenAPI specs

---

## ⚙️ `specgrade.yaml` Config Support

```yaml
spec_version: 3.1.0
input_dir: ./specs/openai
fail_threshold: B
output_format: json
skip_rules:
  - RULE001
  - RULE012
```

* Declarative spec enforcement config
* Supports team-wide standardization
* Auto-loaded if in working directory or passed via `--config`

---

## 🏱 Optional Enhancements (Implemented)

* [x] `specgrade.yaml` config support
* [x] Rule skipping with `--skip`
* [x] Markdown/HTML report generation
* [x] Rule documentation output via `--docs`
* [x] GitHub Action support + semantic exit codes
* [x] `specgrade rules ls` for rule discovery
* [x] Dockerfile for containerized usage
* [x] Hypothesis-based property testing suite
* [x] Rule registry auto-generation (`go:generate`)

---

## 🔥 Summary

| Component              | Status                         |
| ---------------------- | ------------------------------ |
| Rule Plugin System     | ✅ Ready                        |
| Schema Loader Strategy | ✅ Modular                      |
| CLI & Flags            | ✅ Flexible                     |
| YAML Config            | ✅ Supported                    |
| CICD Output & Exit     | ✅ CI-native                    |
| Contributor Onramp     | ✅ Easy                         |
| Testing Coverage       | ✅ Thorough & Hypothesis-driven |

**SpecGrade is now fully specified, testable, and built for long-term maintainability, correctness, and extensibility.**

my github is codetestcode