package rules

import (
	"fmt"
	"strings"

	"github.com/copyleftdev/specgrade/core"
	"github.com/getkin/kin-openapi/openapi3"
)

// InfoTitleRule checks if the OpenAPI spec has a title in the info section
type InfoTitleRule struct{}

func (r *InfoTitleRule) ID() string {
	return "info-title"
}

func (r *InfoTitleRule) Description() string {
	return "OpenAPI spec must have a title in the info section"
}

func (r *InfoTitleRule) AppliesTo(version string) bool {
	return strings.HasPrefix(version, "3.")
}

func (r *InfoTitleRule) Evaluate(ctx *core.SpecContext) core.RuleResult {
	if ctx.Spec.Info == nil {
		return core.RuleResult{
			RuleID:   r.ID(),
			Passed:   false,
			Detail:   "Missing info section",
			Severity: "error",
			Category: "documentation",
			Location: &core.RuleLocation{
				Path:        "$.info",
				Component:   "info",
				File:        "openapi.yaml", // or openapi.json
				FileRef:     "openapi.yaml:1-10 (info section)",
				SpecSection: "info",
			},
			Suggestion: &core.ActionableFix{
				Title:       "Add info section",
				Description: "OpenAPI specifications must include an info section with basic metadata",
				Example:     "info:\n  title: \"My API\"\n  version: \"1.0.0\"\n  description: \"A description of my API\"",
				SchemaRef:   "https://spec.openapis.org/oas/v3.1.0#info-object",
				References: []string{
					"https://swagger.io/specification/#info-object",
					"https://spec.openapis.org/oas/v3.1.0#info-object",
				},
			},
			Impact: &core.ImpactAnalysis{
				Severity:    "error",
				Category:    "documentation",
				Description: "Missing info section violates OpenAPI specification requirements",
			},
			Metadata: map[string]string{
				"required_fields": "title, version",
				"spec_section":    "info",
				"fix_priority":    "high",
			},
		}
	}

	if ctx.Spec.Info.Title == "" {
		return core.RuleResult{
			RuleID:   r.ID(),
			Passed:   false,
			Detail:   "Missing title in info section",
			Severity: "error",
			Category: "documentation",
			Location: &core.RuleLocation{
				Path:      "$.info.title",
				Component: "info.title",
			},
			Suggestion: &core.ActionableFix{
				Description: "Add a descriptive title to the info section of your OpenAPI specification",
				Example: `info:
  title: "My API"
  version: "1.0.0"
  description: "A comprehensive API for managing resources"`,
				References: []string{
					"https://spec.openapis.org/oas/v3.0.3#info-object",
					"https://swagger.io/specification/#info-object",
				},
				SchemaRef: "https://spec.openapis.org/oas/v3.0.3#info-object",
			},
			Impact: &core.ImpactAnalysis{
				Severity:    "error",
				Category:    "documentation",
				Description: "Missing info section violates OpenAPI specification requirements",
			},
			Metadata: map[string]string{
				"field_name":     "title",
				"parent_section": "info",
				"is_required":    "true",
				"fix_priority":   "high",
			},
		}
	}

	// Check for title quality
	title := ctx.Spec.Info.Title
	if len(title) < 5 {
		return core.RuleResult{
			RuleID:   r.ID(),
			Passed:   false,
			Detail:   fmt.Sprintf("Title too short: '%s' (minimum 5 characters recommended)", title),
			Severity: "warning",
			Category: "documentation",
			Location: &core.RuleLocation{
				Path:      "$.info.title",
				Component: "info.title",
			},
			Suggestion: &core.ActionableFix{
				Description: "Use a more descriptive title that clearly explains your API's purpose",
				Example:     fmt.Sprintf("# Instead of: '%s'\n# Consider: '%s Management API' or '%s Service'", title, title, title),
				References: []string{
					"https://spec.openapis.org/oas/v3.0.3#info-object",
				},
				SchemaRef: "https://spec.openapis.org/oas/v3.0.3#info-object",
			},
			Impact: &core.ImpactAnalysis{
				Severity:    "warning",
				Category:    "documentation",
				Description: "Vague titles make it harder for users to understand the API's purpose",
			},
			Metadata: map[string]string{
				"current_length":  fmt.Sprintf("%d", len(title)),
				"recommended_min": "5",
				"fix_priority":    "medium",
			},
		}
	}

	return core.RuleResult{
		RuleID:   r.ID(),
		Passed:   true,
		Detail:   fmt.Sprintf("Title present: '%s'", title),
		Severity: "info",
		Category: "documentation",
		Metadata: map[string]string{
			"title_length": fmt.Sprintf("%d", len(title)),
			"status":       "compliant",
		},
	}
}

// InfoVersionRule checks if the OpenAPI spec has a version in the info section
type InfoVersionRule struct{}

func (r *InfoVersionRule) ID() string {
	return "info-version"
}

func (r *InfoVersionRule) Description() string {
	return "OpenAPI spec must have a version in the info section"
}

func (r *InfoVersionRule) AppliesTo(version string) bool {
	return strings.HasPrefix(version, "3.")
}

func (r *InfoVersionRule) Evaluate(ctx *core.SpecContext) core.RuleResult {
	if ctx.Spec.Info == nil {
		return core.RuleResult{
			RuleID: r.ID(),
			Passed: false,
			Detail: "Missing info section",
		}
	}

	if ctx.Spec.Info.Version == "" {
		return core.RuleResult{
			RuleID: r.ID(),
			Passed: false,
			Detail: "Missing version in info section",
		}
	}

	return core.RuleResult{
		RuleID: r.ID(),
		Passed: true,
		Detail: "Version present",
	}
}

// PathsExistRule checks if the OpenAPI spec has at least one path defined
type PathsExistRule struct{}

func (r *PathsExistRule) ID() string {
	return "paths-exist"
}

func (r *PathsExistRule) Description() string {
	return "OpenAPI spec must have at least one path defined"
}

func (r *PathsExistRule) AppliesTo(version string) bool {
	return strings.HasPrefix(version, "3.")
}

func (r *PathsExistRule) Evaluate(ctx *core.SpecContext) core.RuleResult {
	if len(ctx.Spec.Paths) == 0 {
		return core.RuleResult{
			RuleID: r.ID(),
			Passed: false,
			Detail: "No paths defined",
		}
	}

	return core.RuleResult{
		RuleID: r.ID(),
		Passed: true,
		Detail: fmt.Sprintf("%d paths defined", len(ctx.Spec.Paths)),
	}
}

// OperationIDRule checks if all operations have operation IDs
type OperationIDRule struct{}

func (r *OperationIDRule) ID() string {
	return "operation-operationId-unique"
}

func (r *OperationIDRule) Description() string {
	return "All operations should have unique operation IDs"
}

func (r *OperationIDRule) AppliesTo(version string) bool {
	return strings.HasPrefix(version, "3.")
}

func (r *OperationIDRule) Evaluate(ctx *core.SpecContext) core.RuleResult {
	if ctx.Spec.Paths == nil {
		return core.RuleResult{
			RuleID: r.ID(),
			Passed: true,
			Detail: "No paths to check",
		}
	}

	operationIDs := make(map[string]bool)
	missingCount := 0
	duplicateCount := 0

	for _, pathItem := range ctx.Spec.Paths {
		operations := []struct {
			method string
			op     *openapi3.Operation
		}{
			{"GET", pathItem.Get},
			{"POST", pathItem.Post},
			{"PUT", pathItem.Put},
			{"DELETE", pathItem.Delete},
			{"PATCH", pathItem.Patch},
			{"HEAD", pathItem.Head},
			{"OPTIONS", pathItem.Options},
		}

		for _, op := range operations {
			if op.op == nil {
				continue
			}

			if op.op.OperationID == "" {
				missingCount++
			} else {
				if operationIDs[op.op.OperationID] {
					duplicateCount++
				}
				operationIDs[op.op.OperationID] = true
			}
		}
	}

	if missingCount > 0 || duplicateCount > 0 {
		detail := ""
		if missingCount > 0 {
			detail += fmt.Sprintf("%d operations missing operation ID", missingCount)
		}
		if duplicateCount > 0 {
			if detail != "" {
				detail += ", "
			}
			detail += fmt.Sprintf("%d duplicate operation IDs", duplicateCount)
		}

		return core.RuleResult{
			RuleID: r.ID(),
			Passed: false,
			Detail: detail,
		}
	}

	return core.RuleResult{
		RuleID: r.ID(),
		Passed: true,
		Detail: "All operations have unique operation IDs",
	}
}
