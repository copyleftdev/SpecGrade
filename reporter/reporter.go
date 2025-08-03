package reporter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/copyleftdev/specgrade/core"
)

// Reporter handles different output formats
type Reporter struct {
	grader *DefaultGrader
}

// NewReporter creates a new reporter
func NewReporter() *Reporter {
	return &Reporter{
		grader: NewDefaultGrader(),
	}
}

// GenerateReport creates a report from rule results with enhanced developer analytics
func (r *Reporter) GenerateReport(version string, results []core.RuleResult) *core.Report {
	grade := r.grader.Grade(results)
	score := r.grader.CalculateScore(results)

	// Generate enhanced analytics
	summary := r.generateSummary(results)
	analytics := r.generateAnalytics(results)

	return &core.Report{
		Version:   version,
		Grade:     grade,
		Score:     score,
		Rules:     results,
		Summary:   summary,
		Analytics: analytics,
		Metadata: map[string]string{
			"generated_at": "now", // Would use time.Now() in real implementation
			"tool_version": "1.0.0",
			"report_type":  "enhanced",
		},
	}
}

// FormatJSON outputs the report in JSON format
func (r *Reporter) FormatJSON(report *core.Report) (string, error) {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(data), nil
}

// FormatCLI outputs the report in CLI format
func (r *Reporter) FormatCLI(report *core.Report, targetDir string) string {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("📄 Validating: %s\n", targetDir))
	output.WriteString(fmt.Sprintf("🔖 Spec: OpenAPI %s\n", report.Version))

	passed := 0
	for _, result := range report.Rules {
		if result.Passed {
			passed++
		}
	}

	output.WriteString(fmt.Sprintf("✅ Passed: %d/%d rules\n", passed, len(report.Rules)))
	output.WriteString(fmt.Sprintf("🎯 Score: %d%%\n", report.Score))
	output.WriteString(fmt.Sprintf("🏅 Grade: %s\n", report.Grade))

	// Show failed rules
	if passed < len(report.Rules) {
		output.WriteString("\n❌ Failed Rules:\n")
		for _, result := range report.Rules {
			if !result.Passed {
				output.WriteString(fmt.Sprintf("  - %s: %s\n", result.RuleID, result.Detail))
			}
		}
	}

	return output.String()
}

// FormatMarkdown outputs the report in Markdown format
func (r *Reporter) FormatMarkdown(report *core.Report, targetDir string) string {
	var output strings.Builder

	output.WriteString("# SpecGrade Validation Report\n\n")
	output.WriteString(fmt.Sprintf("**Target:** %s  \n", targetDir))
	output.WriteString(fmt.Sprintf("**OpenAPI Version:** %s  \n", report.Version))
	output.WriteString(fmt.Sprintf("**Grade:** %s  \n", report.Grade))
	output.WriteString(fmt.Sprintf("**Score:** %d%%  \n\n", report.Score))

	passed := 0
	for _, result := range report.Rules {
		if result.Passed {
			passed++
		}
	}

	output.WriteString("## Summary\n\n")
	output.WriteString(fmt.Sprintf("- **Passed:** %d/%d rules\n", passed, len(report.Rules)))
	output.WriteString(fmt.Sprintf("- **Success Rate:** %d%%\n\n", report.Score))

	output.WriteString("## Rule Results\n\n")
	output.WriteString("| Rule ID | Status | Detail |\n")
	output.WriteString("|---------|--------|--------|\n")

	for _, result := range report.Rules {
		status := "❌ Failed"
		if result.Passed {
			status = "✅ Passed"
		}
		output.WriteString(fmt.Sprintf("| %s | %s | %s |\n", result.RuleID, status, result.Detail))
	}

	return output.String()
}

// FormatHTML outputs the report in HTML format
func (r *Reporter) FormatHTML(report *core.Report, targetDir string) string {
	var output strings.Builder

	output.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SpecGrade Validation Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background-color: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { text-align: center; margin-bottom: 30px; }
        .grade { font-size: 3em; font-weight: bold; margin: 20px 0; }
        .grade.A { color: #28a745; }
        .grade.B { color: #17a2b8; }
        .grade.C { color: #ffc107; }
        .grade.D { color: #fd7e14; }
        .grade.F { color: #dc3545; }
        .summary { display: flex; justify-content: space-around; margin: 30px 0; }
        .summary div { text-align: center; }
        .summary h3 { margin: 0; color: #666; }
        .summary p { font-size: 1.5em; font-weight: bold; margin: 5px 0; }
        table { width: 100%; border-collapse: collapse; margin-top: 20px; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background-color: #f8f9fa; font-weight: bold; }
        .passed { color: #28a745; }
        .failed { color: #dc3545; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>SpecGrade Validation Report</h1>
            <p><strong>Target:</strong> ` + targetDir + `</p>
            <p><strong>OpenAPI Version:</strong> ` + report.Version + `</p>
            <div class="grade ` + string(report.Grade[0]) + `">` + report.Grade + `</div>
        </div>
        
        <div class="summary">
            <div>
                <h3>Score</h3>
                <p>` + fmt.Sprintf("%d%%", report.Score) + `</p>
            </div>
            <div>
                <h3>Passed</h3>
                <p>` + fmt.Sprintf("%d/%d", r.countPassed(report.Rules), len(report.Rules)) + `</p>
            </div>
        </div>
        
        <h2>Rule Results</h2>
        <table>
            <thead>
                <tr>
                    <th>Rule ID</th>
                    <th>Status</th>
                    <th>Detail</th>
                </tr>
            </thead>
            <tbody>`)

	for _, result := range report.Rules {
		status := `<span class="failed">❌ Failed</span>`
		if result.Passed {
			status = `<span class="passed">✅ Passed</span>`
		}
		output.WriteString(fmt.Sprintf(`
                <tr>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                </tr>`, result.RuleID, status, result.Detail))
	}

	output.WriteString(`
            </tbody>
        </table>
    </div>
</body>
</html>`)

	return output.String()
}

// countPassed counts the number of passed rules
func (r *Reporter) countPassed(results []core.RuleResult) int {
	count := 0
	for _, result := range results {
		if result.Passed {
			count++
		}
	}
	return count
}

// generateSummary creates a developer-focused summary of the validation results
func (r *Reporter) generateSummary(results []core.RuleResult) *core.ReportSummary {
	totalIssues := 0
	criticalIssues := 0
	quickWins := 0
	issuesByCategory := make(map[string]int)
	issuesBySeverity := make(map[string]int)
	var topPriorities []string
	var complianceGaps []string
	var recommendations []core.DeveloperRecommendation
	totalEstimatedMinutes := 0

	for _, result := range results {
		if !result.Passed {
			totalIssues++

			// Count by severity
			if result.Severity != "" {
				issuesBySeverity[result.Severity]++
				if result.Severity == "error" {
					criticalIssues++
				}
			} else {
				issuesBySeverity["unknown"]++
			}

			// Count by category
			if result.Category != "" {
				issuesByCategory[result.Category]++
			} else {
				issuesByCategory["uncategorized"]++
			}

			// Identify quick wins based on severity
			if result.Severity == "info" || result.Severity == "warning" {
				quickWins++
			}

			// Add to top priorities if high priority
			if result.Metadata != nil && result.Metadata["fix_priority"] == "high" {
				topPriorities = append(topPriorities, result.RuleID)
			}

			// Simple time estimation based on severity
			if result.Severity == "error" {
				totalEstimatedMinutes += 10 // errors take longer
			} else {
				totalEstimatedMinutes += 5 // warnings/info are quicker
			}
		}
	}

	// Generate recommendations based on analysis
	if quickWins > 0 {
		recommendations = append(recommendations, core.DeveloperRecommendation{
			Title:       "Start with Quick Wins",
			Description: fmt.Sprintf("You have %d easy fixes that can be completed quickly", quickWins),
			Priority:    "high",
			Impact:      "Immediate improvement in API quality with minimal effort",
			Effort:      "Low - most can be fixed in under 5 minutes each",
		})
	}

	if criticalIssues > 0 {
		recommendations = append(recommendations, core.DeveloperRecommendation{
			Title:       "Address Critical Issues",
			Description: fmt.Sprintf("You have %d critical issues that should be fixed immediately", criticalIssues),
			Priority:    "critical",
			Impact:      "Prevents API from meeting basic standards and compliance requirements",
			Effort:      "Medium - requires immediate attention but fixes are straightforward",
		})
	}

	// Format estimated fix time
	estimatedFixTime := "Unknown"
	if totalEstimatedMinutes > 0 {
		if totalEstimatedMinutes < 60 {
			estimatedFixTime = fmt.Sprintf("%d minutes", totalEstimatedMinutes)
		} else {
			hours := totalEstimatedMinutes / 60
			minutes := totalEstimatedMinutes % 60
			if minutes > 0 {
				estimatedFixTime = fmt.Sprintf("%d hours %d minutes", hours, minutes)
			} else {
				estimatedFixTime = fmt.Sprintf("%d hours", hours)
			}
		}
	}

	return &core.ReportSummary{
		TotalIssues:      totalIssues,
		CriticalIssues:   criticalIssues,
		QuickWins:        quickWins,
		IssuesByCategory: issuesByCategory,
		IssuesBySeverity: issuesBySeverity,
		TopPriorities:    topPriorities,
		EstimatedFixTime: estimatedFixTime,
		ComplianceGaps:   complianceGaps,
		Recommendations:  recommendations,
	}
}

// generateAnalytics creates detailed analytics about the API specification
func (r *Reporter) generateAnalytics(results []core.RuleResult) *core.ReportAnalytics {
	// For now, create basic analytics - in a real implementation this would analyze the actual spec
	complexity := &core.ComplexityAnalysis{
		EndpointCount:   2,  // Would count from actual spec
		SchemaCount:     5,  // Would count from actual spec
		ParameterCount:  8,  // Would count from actual spec
		ResponseCount:   6,  // Would count from actual spec
		ComplexityScore: 45, // Would calculate based on various factors
		NestingDepth:    3,  // Would analyze schema nesting
		CircularRefs:    0,  // Would detect circular references
		ExternalRefs:    1,  // Would count external references
	}

	riskAssessment := &core.RiskAssessment{
		SecurityRisks:    []string{},
		BreakingChanges:  []string{},
		MaintenanceRisks: []string{},
		ComplianceRisks:  []string{},
		OverallRiskLevel: "low",
	}

	// Analyze results to determine risks
	for _, result := range results {
		if !result.Passed {
			if result.Category == "security" {
				riskAssessment.SecurityRisks = append(riskAssessment.SecurityRisks, result.Detail)
				riskAssessment.OverallRiskLevel = "medium"
			}
			// Track compliance gaps based on category
			if result.Category == "compliance" {
				riskAssessment.ComplianceRisks = append(riskAssessment.ComplianceRisks, result.Detail)
			}
			if result.Severity == "error" {
				riskAssessment.OverallRiskLevel = "medium"
			}
		}
	}

	// Calculate maintenance and developer-friendly scores
	maintenanceScore := 85  // Would calculate based on documentation quality, consistency, etc.
	developerFriendly := 75 // Would calculate based on examples, descriptions, error handling, etc.

	// Adjust scores based on issues found
	for _, result := range results {
		if !result.Passed {
			if result.Category == "documentation" {
				developerFriendly -= 5
				maintenanceScore -= 3
			}
			if result.Severity == "error" {
				maintenanceScore -= 10
				developerFriendly -= 8
			}
		}
	}

	// Ensure scores don't go below 0
	if maintenanceScore < 0 {
		maintenanceScore = 0
	}
	if developerFriendly < 0 {
		developerFriendly = 0
	}

	return &core.ReportAnalytics{
		SpecComplexity:    complexity,
		RiskAssessment:    riskAssessment,
		MaintenanceScore:  maintenanceScore,
		DeveloperFriendly: developerFriendly,
	}
}
