package reporter

import (
	"fmt"
	"strings"

	"github.com/copyleftdev/specgrade/core"
)

// FormatDeveloperCLI outputs a practical developer-focused CLI format
// Focuses on maintainable features: file references and OpenAPI schema links
func (r *Reporter) FormatDeveloperCLI(report *core.Report, targetDir string) string {
	var output strings.Builder

	// Header
	output.WriteString("\n🚀 SpecGrade Developer Report\n")
	output.WriteString("=" + strings.Repeat("=", 35) + "\n")
	output.WriteString(fmt.Sprintf("📄 Target: %s\n", targetDir))
	output.WriteString(fmt.Sprintf("🔖 OpenAPI Version: %s\n", report.Version))
	output.WriteString(fmt.Sprintf("🏅 Grade: %s (%d%%)\n", report.Grade, report.Score))

	// Count issues by severity
	passed := 0
	failed := 0
	for _, result := range report.Rules {
		if result.Passed {
			passed++
		} else {
			failed++
		}
	}

	output.WriteString(fmt.Sprintf("✅ Passed: %d/%d rules\n", passed, len(report.Rules)))
	if failed > 0 {
		output.WriteString(fmt.Sprintf("❌ Failed: %d rules\n", failed))
	}

	// Detailed issues with file references and schema links
	if failed > 0 {
		output.WriteString("\n🔍 Issues Found\n")
		output.WriteString("=" + strings.Repeat("=", 15) + "\n")

		for _, result := range report.Rules {
			if !result.Passed {
				// Issue header with severity
				severityIcon := "ℹ️"
				if result.Severity == "error" {
					severityIcon = "❌"
				} else if result.Severity == "warning" {
					severityIcon = "⚠️"
				}

				output.WriteString(fmt.Sprintf("\n%s %s\n", severityIcon, result.RuleID))
				output.WriteString(fmt.Sprintf("   Problem: %s\n", result.Detail))

				// File location - most important for developers
				if result.Location != nil {
					if result.Location.FileRef != "" {
						output.WriteString(fmt.Sprintf("   📄 File: %s\n", result.Location.FileRef))
					} else if result.Location.File != "" {
						output.WriteString(fmt.Sprintf("   📄 File: %s\n", result.Location.File))
					}

					if result.Location.SpecSection != "" {
						output.WriteString(fmt.Sprintf("   📋 Section: %s\n", result.Location.SpecSection))
					}

					if result.Location.Path != "" {
						output.WriteString(fmt.Sprintf("   🔍 JSON Path: %s\n", result.Location.Path))
					}
				}

				// Simple fix suggestion with schema reference
				if result.Suggestion != nil {
					output.WriteString("   🔧 Fix:\n")
					output.WriteString(fmt.Sprintf("      %s\n", result.Suggestion.Description))

					// OpenAPI Schema Reference - key value add
					if result.Suggestion.SchemaRef != "" {
						output.WriteString(fmt.Sprintf("      📋 OpenAPI Schema: %s\n", result.Suggestion.SchemaRef))
					}

					// Simple example if available
					if result.Suggestion.Example != "" {
						output.WriteString("      Example:\n")
						exampleLines := strings.Split(result.Suggestion.Example, "\n")
						for _, line := range exampleLines {
							output.WriteString(fmt.Sprintf("        %s\n", line))
						}
					}

					// Documentation references
					if len(result.Suggestion.References) > 0 {
						output.WriteString("      📚 References:\n")
						for _, ref := range result.Suggestion.References {
							output.WriteString(fmt.Sprintf("        • %s\n", ref))
						}
					}
				}

				output.WriteString("\n" + strings.Repeat("-", 50) + "\n")
			}
		}
	}

	// Simple next steps
	output.WriteString("\n🎯 Next Steps\n")
	output.WriteString("=" + strings.Repeat("=", 12) + "\n")
	if failed > 0 {
		output.WriteString("1. Review the file locations and schema references above\n")
		output.WriteString("2. Fix issues using the provided examples and documentation\n")
		output.WriteString("3. Re-run SpecGrade to verify fixes\n")
	} else {
		output.WriteString("🎉 All validation rules passed! Your OpenAPI spec looks great.\n")
	}

	output.WriteString("\n" + strings.Repeat("=", 50) + "\n")
	output.WriteString("💡 Tip: Use 'specgrade --output-format json' for machine-readable results\n")

	return output.String()
}
