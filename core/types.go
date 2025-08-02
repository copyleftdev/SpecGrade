package core

import "github.com/getkin/kin-openapi/openapi3"

// Rule represents a validation rule that can be applied to an OpenAPI spec
type Rule interface {
	ID() string
	Description() string
	AppliesTo(version string) bool
	Evaluate(ctx *SpecContext) RuleResult
}

// RuleResult represents the result of evaluating a rule
type RuleResult struct {
	RuleID string `json:"ruleID"`
	Passed bool   `json:"passed"`
	Detail string `json:"detail"`
}

// SpecContext provides context for rule evaluation
type SpecContext struct {
	Spec    *openapi3.T
	Version string
}

// SpecLoader loads OpenAPI specifications
type SpecLoader interface {
	Load(version string) (*openapi3.T, error)
}

// Grader assigns grades based on rule results
type Grader interface {
	Grade(results []RuleResult) string // Returns A, B, C, etc
}

// ExitHandler determines exit codes for CI/CD integration
type ExitHandler interface {
	Handle(grade string) int
}

// Config represents the configuration for SpecGrade
type Config struct {
	SpecVersion    string   `yaml:"spec_version"`
	InputDir       string   `yaml:"input_dir"`
	FailThreshold  string   `yaml:"fail_threshold"`
	OutputFormat   string   `yaml:"output_format"`
	SkipRules      []string `yaml:"skip_rules"`
	ConfigPath     string   `yaml:"-"`
}

// Report represents the final validation report
type Report struct {
	Version string       `json:"version"`
	Grade   string       `json:"grade"`
	Score   int          `json:"score"`
	Rules   []RuleResult `json:"rules"`
}
