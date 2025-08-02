package reporter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/codetestcode/specgrade/core"
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

// GenerateReport creates a report from rule results
func (r *Reporter) GenerateReport(version string, results []core.RuleResult) *core.Report {
	grade := r.grader.Grade(results)
	score := r.grader.CalculateScore(results)

	return &core.Report{
		Version: version,
		Grade:   grade,
		Score:   score,
		Rules:   results,
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
	passed := 0
	for _, result := range results {
		if result.Passed {
			passed++
		}
	}
	return passed
}
