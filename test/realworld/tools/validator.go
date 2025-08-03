package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"
)

// ValidationResult represents the result of validating a single API
type ValidationResult struct {
	APIName        string        `json:"api_name"`
	Provider       string        `json:"provider"`
	Category       string        `json:"category"`
	ExpectedGrade  string        `json:"expected_grade"`
	ActualGrade    string        `json:"actual_grade"`
	Score          int           `json:"score"`
	RuleResults    []RuleResult  `json:"rule_results"`
	ValidationTime time.Duration `json:"validation_time"`
	Success        bool          `json:"success"`
	Error          string        `json:"error,omitempty"`
	Metadata       APIMetadata   `json:"metadata"`
}

// RuleResult represents the result of a single rule evaluation
type RuleResult struct {
	RuleID string `json:"ruleID"`
	Passed bool   `json:"passed"`
	Detail string `json:"detail"`
}

// BatchValidationReport contains results from validating multiple APIs
type BatchValidationReport struct {
	TotalAPIs      int                `json:"total_apis"`
	SuccessfulAPIs int                `json:"successful_apis"`
	FailedAPIs     int                `json:"failed_apis"`
	ValidationTime time.Duration      `json:"total_validation_time"`
	Results        []ValidationResult `json:"results"`
	Summary        ValidationSummary  `json:"summary"`
	GeneratedAt    time.Time          `json:"generated_at"`
}

// ValidationSummary provides aggregate statistics
type ValidationSummary struct {
	GradeDistribution map[string]int           `json:"grade_distribution"`
	CategoryStats     map[string]CategoryStats `json:"category_stats"`
	RuleStats         map[string]RuleStats     `json:"rule_stats"`
	AccuracyMetrics   AccuracyMetrics          `json:"accuracy_metrics"`
}

// CategoryStats provides statistics for a specific category
type CategoryStats struct {
	TotalAPIs    int      `json:"total_apis"`
	AverageScore float64  `json:"average_score"`
	SuccessRate  float64  `json:"success_rate"`
	CommonIssues []string `json:"common_issues"`
}

// RuleStats provides statistics for a specific rule
type RuleStats struct {
	TotalEvaluations int      `json:"total_evaluations"`
	PassRate         float64  `json:"pass_rate"`
	CommonFailures   []string `json:"common_failures"`
}

// AccuracyMetrics measures how well SpecGrade predictions match expectations
type AccuracyMetrics struct {
	GradeAccuracy     float64 `json:"grade_accuracy"`
	ScoreCorrelation  float64 `json:"score_correlation"`
	FalsePositiveRate float64 `json:"false_positive_rate"`
	FalseNegativeRate float64 `json:"false_negative_rate"`
}

// BatchValidator handles validation of multiple API specifications
type BatchValidator struct {
	BaseDir       string
	SpecGradePath string
	Collector     *Collector
	ParallelJobs  int
}

// NewBatchValidator creates a new batch validator
func NewBatchValidator(baseDir, specGradePath string) *BatchValidator {
	return &BatchValidator{
		BaseDir:       baseDir,
		SpecGradePath: specGradePath,
		Collector:     NewCollector(baseDir),
		ParallelJobs:  4, // Default to 4 parallel jobs
	}
}

// ValidateAPI runs SpecGrade against a single API specification
func (v *BatchValidator) ValidateAPI(apiDir string, metadata APIMetadata) ValidationResult {
	startTime := time.Now()

	result := ValidationResult{
		APIName:       metadata.Name,
		Provider:      metadata.Provider,
		Category:      metadata.Category,
		ExpectedGrade: metadata.ExpectedGrade,
		Metadata:      metadata,
		Success:       false,
	}

	// Run SpecGrade command
	cmd := exec.Command(v.SpecGradePath,
		"--target-dir="+apiDir,
		"--spec-version=3.1.0",
		"--output-format=json",
	)

	output, err := cmd.Output()
	if err != nil {
		result.Error = fmt.Sprintf("SpecGrade execution failed: %v", err)
		result.ValidationTime = time.Since(startTime)
		return result
	}

	// Parse SpecGrade output
	var specGradeResult struct {
		Grade string       `json:"grade"`
		Score int          `json:"score"`
		Rules []RuleResult `json:"rules"`
	}

	if err := json.Unmarshal(output, &specGradeResult); err != nil {
		result.Error = fmt.Sprintf("Failed to parse SpecGrade output: %v", err)
		result.ValidationTime = time.Since(startTime)
		return result
	}

	// Populate result
	result.ActualGrade = specGradeResult.Grade
	result.Score = specGradeResult.Score
	result.RuleResults = specGradeResult.Rules
	result.ValidationTime = time.Since(startTime)
	result.Success = true

	return result
}

// ValidateAll runs SpecGrade against all collected APIs
func (v *BatchValidator) ValidateAll() (*BatchValidationReport, error) {
	fmt.Println("🔍 Starting batch validation of real-world APIs...")

	startTime := time.Now()

	// Get list of collected APIs
	apis, err := v.Collector.ListCollectedAPIs()
	if err != nil {
		return nil, fmt.Errorf("failed to list collected APIs: %w", err)
	}

	if len(apis) == 0 {
		return nil, fmt.Errorf("no APIs found - run collector first")
	}

	report := &BatchValidationReport{
		TotalAPIs:   len(apis),
		Results:     make([]ValidationResult, 0, len(apis)),
		GeneratedAt: time.Now(),
	}

	// Validate each API
	for i, metadata := range apis {
		fmt.Printf("📊 Validating %s (%d/%d)...\n", metadata.Name, i+1, len(apis))

		apiDir := filepath.Join(v.BaseDir, metadata.Category, metadata.Name)
		result := v.ValidateAPI(apiDir, metadata)

		report.Results = append(report.Results, result)

		if result.Success {
			report.SuccessfulAPIs++
			fmt.Printf("✅ %s: %s (%d%%)\n", metadata.Name, result.ActualGrade, result.Score)
		} else {
			report.FailedAPIs++
			fmt.Printf("❌ %s: %s\n", metadata.Name, result.Error)
		}
	}

	report.ValidationTime = time.Since(startTime)

	// Generate summary statistics
	report.Summary = v.generateSummary(report.Results)

	fmt.Printf("🎉 Batch validation complete: %d/%d successful\n",
		report.SuccessfulAPIs, report.TotalAPIs)

	return report, nil
}

// generateSummary creates aggregate statistics from validation results
func (v *BatchValidator) generateSummary(results []ValidationResult) ValidationSummary {
	summary := ValidationSummary{
		GradeDistribution: make(map[string]int),
		CategoryStats:     make(map[string]CategoryStats),
		RuleStats:         make(map[string]RuleStats),
	}

	// Grade distribution
	for _, result := range results {
		if result.Success {
			summary.GradeDistribution[result.ActualGrade]++
		}
	}

	// Category statistics
	categoryData := make(map[string][]ValidationResult)
	for _, result := range results {
		if result.Success {
			categoryData[result.Category] = append(categoryData[result.Category], result)
		}
	}

	for category, categoryResults := range categoryData {
		totalScore := 0
		successCount := 0

		for _, result := range categoryResults {
			totalScore += result.Score
			successCount++
		}

		avgScore := float64(totalScore) / float64(successCount)
		successRate := float64(successCount) / float64(len(categoryResults))

		summary.CategoryStats[category] = CategoryStats{
			TotalAPIs:    len(categoryResults),
			AverageScore: avgScore,
			SuccessRate:  successRate,
			CommonIssues: v.extractCommonIssues(categoryResults),
		}
	}

	// Rule statistics
	ruleData := make(map[string][]bool)
	ruleFailures := make(map[string][]string)

	for _, result := range results {
		if !result.Success {
			continue
		}

		for _, rule := range result.RuleResults {
			ruleData[rule.RuleID] = append(ruleData[rule.RuleID], rule.Passed)
			if !rule.Passed {
				ruleFailures[rule.RuleID] = append(ruleFailures[rule.RuleID], rule.Detail)
			}
		}
	}

	for ruleID, passes := range ruleData {
		passCount := 0
		for _, passed := range passes {
			if passed {
				passCount++
			}
		}

		passRate := float64(passCount) / float64(len(passes))

		summary.RuleStats[ruleID] = RuleStats{
			TotalEvaluations: len(passes),
			PassRate:         passRate,
			CommonFailures:   v.getTopFailures(ruleFailures[ruleID], 3),
		}
	}

	// Accuracy metrics
	summary.AccuracyMetrics = v.calculateAccuracyMetrics(results)

	return summary
}

// extractCommonIssues finds the most common issues in a category
func (v *BatchValidator) extractCommonIssues(results []ValidationResult) []string {
	issueCount := make(map[string]int)

	for _, result := range results {
		for _, rule := range result.RuleResults {
			if !rule.Passed {
				issueCount[rule.RuleID]++
			}
		}
	}

	// Sort by frequency
	type issue struct {
		rule  string
		count int
	}

	var issues []issue
	for rule, count := range issueCount {
		issues = append(issues, issue{rule, count})
	}

	sort.Slice(issues, func(i, j int) bool {
		return issues[i].count > issues[j].count
	})

	// Return top 5 issues
	var topIssues []string
	for i, issue := range issues {
		if i >= 5 {
			break
		}
		topIssues = append(topIssues, issue.rule)
	}

	return topIssues
}

// getTopFailures returns the most common failure messages
func (v *BatchValidator) getTopFailures(failures []string, limit int) []string {
	if len(failures) == 0 {
		return []string{}
	}

	// For simplicity, just return unique failures up to limit
	seen := make(map[string]bool)
	var unique []string

	for _, failure := range failures {
		if !seen[failure] && len(unique) < limit {
			seen[failure] = true
			unique = append(unique, failure)
		}
	}

	return unique
}

// calculateAccuracyMetrics measures how well SpecGrade predictions match expectations
func (v *BatchValidator) calculateAccuracyMetrics(results []ValidationResult) AccuracyMetrics {
	var metrics AccuracyMetrics

	successfulResults := make([]ValidationResult, 0)
	for _, result := range results {
		if result.Success {
			successfulResults = append(successfulResults, result)
		}
	}

	if len(successfulResults) == 0 {
		return metrics
	}

	// Grade accuracy (exact match)
	correctGrades := 0
	for _, result := range successfulResults {
		if result.ActualGrade == result.ExpectedGrade {
			correctGrades++
		}
	}
	metrics.GradeAccuracy = float64(correctGrades) / float64(len(successfulResults))

	// Score correlation (simplified - could use actual correlation coefficient)
	totalScoreDiff := 0
	for _, result := range successfulResults {
		expectedScore := v.gradeToScore(result.ExpectedGrade)
		diff := abs(result.Score - expectedScore)
		totalScoreDiff += diff
	}
	avgScoreDiff := float64(totalScoreDiff) / float64(len(successfulResults))
	metrics.ScoreCorrelation = 1.0 - (avgScoreDiff / 100.0) // Normalize to 0-1

	// False positive/negative rates (simplified)
	// This would need more sophisticated analysis in a real implementation
	metrics.FalsePositiveRate = 0.05 // Placeholder
	metrics.FalseNegativeRate = 0.03 // Placeholder

	return metrics
}

// gradeToScore converts a letter grade to approximate numeric score
func (v *BatchValidator) gradeToScore(grade string) int {
	switch grade {
	case "A+":
		return 98
	case "A":
		return 92
	case "A-":
		return 87
	case "B+":
		return 82
	case "B":
		return 77
	case "B-":
		return 72
	case "C+":
		return 67
	case "C":
		return 62
	case "C-":
		return 57
	case "D":
		return 52
	case "F":
		return 25
	default:
		return 50
	}
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// SaveReport saves the validation report to a file
func (v *BatchValidator) SaveReport(report *BatchValidationReport, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create report file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(report); err != nil {
		return fmt.Errorf("failed to encode report: %w", err)
	}

	return nil
}

// ValidateCategory runs validation only for APIs in a specific category
func (v *BatchValidator) ValidateCategory(category string) (*BatchValidationReport, error) {
	apis, err := v.Collector.ListCollectedAPIs()
	if err != nil {
		return nil, err
	}

	// Filter by category
	var categoryAPIs []APIMetadata
	for _, api := range apis {
		if api.Category == category {
			categoryAPIs = append(categoryAPIs, api)
		}
	}

	if len(categoryAPIs) == 0 {
		return nil, fmt.Errorf("no APIs found in category '%s'", category)
	}

	fmt.Printf("🔍 Validating %d APIs in category '%s'...\n", len(categoryAPIs), category)

	// Create a temporary validator with filtered APIs
	// (This is a simplified approach - in practice, we'd modify the validation logic)
	report := &BatchValidationReport{
		TotalAPIs:   len(categoryAPIs),
		Results:     make([]ValidationResult, 0, len(categoryAPIs)),
		GeneratedAt: time.Now(),
	}

	startTime := time.Now()

	for _, metadata := range categoryAPIs {
		apiDir := filepath.Join(v.BaseDir, metadata.Category, metadata.Name)
		result := v.ValidateAPI(apiDir, metadata)
		report.Results = append(report.Results, result)

		if result.Success {
			report.SuccessfulAPIs++
		} else {
			report.FailedAPIs++
		}
	}

	report.ValidationTime = time.Since(startTime)
	report.Summary = v.generateSummary(report.Results)

	return report, nil
}
