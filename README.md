# 📘 SpecGrade

A modular, dynamic, and CICD-optimized conformance validator for OpenAPI specifications.

## 🚀 Features

- **Modular Rule System**: Pluggable validation rules that can be easily extended
- **Enhanced Developer Reporting**: Actionable insights with file references and OpenAPI schema links
- **Multiple Output Formats**: JSON, CLI, HTML, Markdown, and Developer reporting
- **CI/CD Integration**: Semantic exit codes and configurable fail thresholds
- **YAML Configuration**: Team-wide standardization with `specgrade.yaml`
- **Rule Management**: Skip specific rules or generate documentation
- **Version Support**: OpenAPI 3.0.0 and 3.1.0 specifications
- **Advanced Rigor**: Real-world API testing, fuzzing, ML prediction, and community contributions

## 🎯 Grading System

SpecGrade uses a comprehensive grading system that evaluates your OpenAPI specification across multiple dimensions:

### Grade Scale
| Grade | Score Range | Description |
|-------|-------------|-------------|
| **A+** | 95-100% | Exceptional - Production-ready with best practices |
| **A**  | 90-94%  | Excellent - High quality with minor improvements needed |
| **B+** | 85-89%  | Very Good - Solid specification with some gaps |
| **B**  | 80-84%  | Good - Functional but needs attention |
| **C+** | 75-79%  | Fair - Multiple issues to address |
| **C**  | 70-74%  | Poor - Significant problems present |
| **D**  | 60-69%  | Very Poor - Major rework required |
| **F**  | 0-59%   | Failing - Specification has critical issues |

### Scoring Methodology

SpecGrade calculates your score using a **weighted rule evaluation system**:

```
Final Score = (Passed Rules / Total Rules) × 100%
```

#### Rule Categories & Weights

1. **Basic Validation Rules** (25% weight)
   - `info-title`: API must have a descriptive title
   - `info-description`: API must have a clear description
   - `paths-kebab-case`: Paths should use kebab-case naming
   - `operation-operationid`: All operations must have unique operation IDs

2. **Advanced Validation Rules** (75% weight)
   - `oas3-valid-schema-example`: Examples must match their schema types
   - `operation-description`: All operations must be documented
   - `no-trailing-slash`: Paths should not have trailing slashes
   - `security-defined`: Security schemes must be properly defined

#### Example Calculation

```yaml
# Your API has 8 total rules
# 7 rules passed, 1 rule failed

Score = (7 ÷ 8) × 100% = 87.5%
Grade = B+ (85-89% range)
```

### Quality Insights

Beyond the grade, SpecGrade provides:

- **Detailed Rule Results**: See exactly which rules passed/failed
- **Actionable Recommendations**: Specific steps to improve your score
- **Trend Analysis**: Track improvements over time
- **Best Practice Guidance**: Learn industry standards for API design

## 🚦 Implementation Status

### ✅ **Production Ready** (Fully Implemented)

These features are **100% functional** and ready for production use:

- **Core Validation Engine**: 8 comprehensive validation rules
- **Grading System**: Transparent A+ to F scoring with detailed breakdowns
- **Enhanced Developer Reporting**: Actionable insights with file references and OpenAPI schema links
- **Multi-Format Output**: JSON, CLI, HTML, Markdown, and Developer reporters
- **Configuration Management**: YAML config files with CLI flag precedence
- **Multi-File Support**: External `$ref` resolution for complex specs
- **Rule Management**: Skip rules, generate documentation, list available rules
- **CI/CD Integration**: Semantic exit codes and fail thresholds
- **Docker Support**: Containerized deployment ready
- **Comprehensive Testing**: Unit tests, integration tests, property-based tests

### 🚧 **Prototype/Demo** (Architectural Previews)

These features demonstrate **future capabilities** with simulated data:

- **🌍 Real-World API Collection**: Simulates downloading APIs from major providers
- **🔥 Fuzzing Tests**: Mock corruption strategies and robustness testing
- **🤖 ML Quality Prediction**: Prototype feature extraction and quality insights
- **🤝 Community Framework**: Simulated contribution statistics and edge case patterns

> **📝 Note**: The prototype features show the intended user experience and architecture. They use hardcoded/simulated data for demonstration purposes. The core validation system provides real, actionable results.

### 🎯 **Example: Real vs. Simulated**

```bash
# ✅ REAL - Actual validation with genuine results
./specgrade --target-dir test/sample-spec
# Output: 6/8 rules passed, 75% score, Grade B

# 🎭 SIMULATED - Architectural preview with mock data
./specgrade advanced community --action stats
# Output: "47 contributions, Alice Johnson (12 contributions)"
```

## 📦 Installation

```bash
go install github.com/copyleftdev/specgrade@latest
```

Or build from source:

```bash
git clone https://github.com/copyleftdev/specgrade.git
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

### Enhanced Developer Reporting

🎯 **New!** SpecGrade now includes enhanced developer-focused reporting with actionable insights:

```bash
# Get enhanced developer report with file references and schema links
specgrade --target-dir ./specs --output-format developer
```

**Key Features:**
- 📄 **Precise File References** - Shows exact file locations for each issue
- 📋 **OpenAPI Schema Links** - Direct links to official OpenAPI specification sections
- 🔧 **Actionable Fix Guidance** - Practical examples and step-by-step instructions
- 📚 **Documentation References** - Links to relevant OpenAPI documentation
- 🎯 **Developer-Focused Output** - Clean, readable format optimized for developers

**Example Output:**
```
🚀 SpecGrade Developer Report
====================================
📄 Target: ./specs
🔖 OpenAPI Version: 3.1.0
🏅 Grade: B (75%)

🔍 Issues Found
================

⚠️ operation-success-response
   Problem: Missing error responses: 1 missing 400 responses
   📄 File: openapi.yaml:paths section
   📋 Section: paths
   🔍 JSON Path: $.paths
   🔧 Fix:
      Define proper error responses to help API consumers handle failures gracefully
      📋 OpenAPI Schema: https://spec.openapis.org/oas/v3.0.3#responses-object
      📚 References:
        • https://spec.openapis.org/oas/v3.0.3#responses-object
```

### Using Configuration File

Create a `specgrade.yaml` file:

```yaml
spec_version: 3.1.0
input_dir: ./specs/openai
fail_threshold: B
output_format: developer  # Use enhanced developer reporting
skip_rules:
  - RULE001
  - RULE012
```

Then run:

```bash
specgrade --config=specgrade.yaml
```

### Advanced Rigor Features

> **⚠️ Implementation Status**: The advanced rigor features below are currently **architectural prototypes** with simulated data for demonstration purposes. The core validation system is fully production-ready.

SpecGrade includes cutting-edge features for comprehensive API validation:

#### 🌍 Real-World API Collection
```bash
# Collect APIs from major providers
specgrade advanced collect --categories fintech,developer

# Validate all collected APIs in batch
specgrade advanced validate-batch --report batch_report.json
```

#### 🔥 Fuzzing Tests
```bash
# Test robustness with corrupted specs
specgrade advanced fuzz --input api.yaml --strategies structural,semantic --iterations 100
```

#### 🤖 ML-Based Quality Prediction
```bash
# Predict API quality using machine learning
specgrade advanced predict --input api.yaml --detailed
```

#### 🤝 Community Contributions
```bash
# View community statistics and edge case patterns
specgrade advanced community --action stats
specgrade advanced community --action patterns
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
        run: go install github.com/copyleftdev/specgrade@latest
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

## 🗺️ Development Roadmap

### Making Prototype Features Production-Ready

To convert the advanced rigor prototypes into fully functional features:

#### 🌍 **Real-World API Collection**
- [ ] Implement HTTP clients for major API providers (Stripe, GitHub, AWS, etc.)
- [ ] Add API discovery and metadata extraction
- [ ] Build automated update scheduling and version tracking
- [ ] Create API categorization and tagging system

#### 🔥 **Fuzzing Framework**
- [ ] Implement actual OpenAPI corruption algorithms
- [ ] Add crash detection and error analysis
- [ ] Build fuzzing campaign management
- [ ] Integrate with CI/CD for automated robustness testing

#### 🤖 **ML Quality Prediction**
- [ ] Train models on real OpenAPI specification datasets
- [ ] Implement feature extraction from actual specs
- [ ] Add model versioning and A/B testing
- [ ] Build feedback loop for continuous improvement

#### 🤝 **Community Framework**
- [ ] Create web interface for contribution submission
- [ ] Build review workflow and approval system
- [ ] Implement real contribution storage and analytics
- [ ] Add user authentication and reputation system

### 🎯 **Priority Order**
1. **Fuzzing Framework** - Highest impact for robustness testing
2. **Real-World API Collection** - Valuable for benchmarking
3. **ML Quality Prediction** - Advanced feature for insights
4. **Community Framework** - Long-term ecosystem building

## 🤝 Contributing

### Current Contribution Areas

**✅ Ready for Contributions:**
- New validation rules for the core engine
- Additional output formats (SARIF, XML, etc.)
- Performance optimizations
- Documentation improvements
- Bug fixes and edge case handling

**🚧 Advanced Features (Prototype → Production):**
- Help implement the roadmap items above
- Contribute real-world API specifications for testing
- Share edge cases and validation scenarios

### How to Contribute

1. Fork the repository
2. Create a feature branch
3. Add tests for your changes
4. Ensure all tests pass
5. Submit a pull request

## 📞 Support

- GitHub Issues: [Report bugs or request features](https://github.com/copyleftdev/specgrade/issues)
- Documentation: [Full documentation](https://github.com/copyleftdev/specgrade/wiki)
