package core

import "github.com/getkin/kin-openapi/openapi3"

// Rule represents a validation rule that can be applied to an OpenAPI spec
type Rule interface {
	ID() string
	Description() string
	AppliesTo(version string) bool
	Evaluate(ctx *SpecContext) RuleResult
}

// RuleResult represents the result of evaluating a rule with enhanced developer metadata
type RuleResult struct {
	RuleID     string            `json:"ruleID"`
	Passed     bool              `json:"passed"`
	Detail     string            `json:"detail"`
	Severity   string            `json:"severity"` // error, warning, info
	Category   string            `json:"category"` // documentation, structure, security, etc.
	Location   *RuleLocation     `json:"location,omitempty"`
	Suggestion *ActionableFix    `json:"suggestion,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	Impact     *ImpactAnalysis   `json:"impact,omitempty"`
}

// RuleLocation provides specific location information for issues
type RuleLocation struct {
	Path        string `json:"path,omitempty"`         // JSON path to the issue
	Line        int    `json:"line,omitempty"`         // Line number in YAML/JSON
	Column      int    `json:"column,omitempty"`       // Column number
	Component   string `json:"component,omitempty"`    // Component name (operation, schema, etc.)
	Method      string `json:"method,omitempty"`       // HTTP method if applicable
	Endpoint    string `json:"endpoint,omitempty"`     // API endpoint if applicable
	File        string `json:"file,omitempty"`         // Specific file name (for multi-file specs)
	FileRef     string `json:"file_ref,omitempty"`     // Full file reference for developers
	SpecSection string `json:"spec_section,omitempty"` // OpenAPI spec section (paths, components, etc.)
}

// ActionableFix provides practical, maintainable suggestions for developers
type ActionableFix struct {
	Title       string   `json:"title"`                // Short fix title
	Description string   `json:"description"`          // Brief explanation
	SchemaRef   string   `json:"schema_ref,omitempty"` // Direct link to OpenAPI schema definition
	References  []string `json:"references,omitempty"` // Links to documentation
	Example     string   `json:"example,omitempty"`    // Simple code example
}

// ImpactAnalysis provides simple, maintainable impact information
type ImpactAnalysis struct {
	Severity    string `json:"severity"`              // error, warning, info
	Category    string `json:"category"`              // documentation, structure, security, etc.
	Description string `json:"description,omitempty"` // Brief impact description
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
	SpecVersion   string   `yaml:"spec_version"`
	InputDir      string   `yaml:"input_dir"`
	FailThreshold string   `yaml:"fail_threshold"`
	OutputFormat  string   `yaml:"output_format"`
	SkipRules     []string `yaml:"skip_rules"`
	ConfigPath    string   `yaml:"-"`
}

// Report represents the final validation report with enhanced developer insights
type Report struct {
	Version   string            `json:"version"`
	Grade     string            `json:"grade"`
	Score     int               `json:"score"`
	Rules     []RuleResult      `json:"rules"`
	Summary   *ReportSummary    `json:"summary"`
	Analytics *ReportAnalytics  `json:"analytics"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// ReportSummary provides high-level insights for developers
type ReportSummary struct {
	TotalIssues      int                       `json:"total_issues"`
	CriticalIssues   int                       `json:"critical_issues"`
	QuickWins        int                       `json:"quick_wins"` // Issues that are easy to fix
	IssuesByCategory map[string]int            `json:"issues_by_category"`
	IssuesBySeverity map[string]int            `json:"issues_by_severity"`
	TopPriorities    []string                  `json:"top_priorities"`     // Rule IDs to fix first
	EstimatedFixTime string                    `json:"estimated_fix_time"` // Total estimated time to fix all issues
	ComplianceGaps   []string                  `json:"compliance_gaps"`    // Standards not being met
	Recommendations  []DeveloperRecommendation `json:"recommendations"`
}

// ReportAnalytics provides deeper insights and trends
type ReportAnalytics struct {
	SpecComplexity    *ComplexityAnalysis  `json:"spec_complexity"`
	QualityTrends     *QualityTrends       `json:"quality_trends,omitempty"`
	Benchmarks        *BenchmarkComparison `json:"benchmarks,omitempty"`
	RiskAssessment    *RiskAssessment      `json:"risk_assessment"`
	MaintenanceScore  int                  `json:"maintenance_score"`  // 0-100 how maintainable the spec is
	DeveloperFriendly int                  `json:"developer_friendly"` // 0-100 how dev-friendly the API is
}

// DeveloperRecommendation provides actionable recommendations
type DeveloperRecommendation struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    string   `json:"priority"` // high, medium, low
	Impact      string   `json:"impact"`   // Description of positive impact
	Effort      string   `json:"effort"`   // Time/effort required
	RuleIDs     []string `json:"rule_ids"` // Related rule IDs
}

// ComplexityAnalysis analyzes the complexity of the OpenAPI spec
type ComplexityAnalysis struct {
	EndpointCount   int `json:"endpoint_count"`
	SchemaCount     int `json:"schema_count"`
	ParameterCount  int `json:"parameter_count"`
	ResponseCount   int `json:"response_count"`
	ComplexityScore int `json:"complexity_score"` // 0-100
	NestingDepth    int `json:"nesting_depth"`    // Maximum schema nesting
	CircularRefs    int `json:"circular_refs"`    // Number of circular references
	ExternalRefs    int `json:"external_refs"`    // Number of external references
}

// QualityTrends tracks quality improvements over time (future enhancement)
type QualityTrends struct {
	PreviousScore  int    `json:"previous_score,omitempty"`
	ScoreChange    int    `json:"score_change,omitempty"`    // +/- change from previous
	TrendDirection string `json:"trend_direction,omitempty"` // improving, declining, stable
	FixedIssues    int    `json:"fixed_issues,omitempty"`
	NewIssues      int    `json:"new_issues,omitempty"`
}

// BenchmarkComparison compares against industry standards
type BenchmarkComparison struct {
	IndustryAverage int    `json:"industry_average"`
	TopPercentile   int    `json:"top_percentile"` // 90th percentile score
	YourRanking     string `json:"your_ranking"`   // "above average", "below average", etc.
	SimilarAPIs     int    `json:"similar_apis"`   // Number of similar APIs in comparison
}

// RiskAssessment identifies potential risks in the API specification
type RiskAssessment struct {
	SecurityRisks    []string `json:"security_risks,omitempty"`
	BreakingChanges  []string `json:"breaking_changes,omitempty"`
	MaintenanceRisks []string `json:"maintenance_risks,omitempty"`
	ComplianceRisks  []string `json:"compliance_risks,omitempty"`
	OverallRiskLevel string   `json:"overall_risk_level"` // low, medium, high, critical
}
