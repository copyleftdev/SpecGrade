package rules

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/copyleftdev/specgrade/core"
	"github.com/getkin/kin-openapi/openapi3"
)

// SchemaExampleConsistencyRule checks if examples match their declared types
type SchemaExampleConsistencyRule struct{}

func (r *SchemaExampleConsistencyRule) ID() string {
	return "oas3-valid-schema-example"
}

func (r *SchemaExampleConsistencyRule) Description() string {
	return "Schema examples must match their declared types"
}

func (r *SchemaExampleConsistencyRule) AppliesTo(version string) bool {
	return strings.HasPrefix(version, "3.")
}

func (r *SchemaExampleConsistencyRule) Evaluate(ctx *core.SpecContext) core.RuleResult {
	issues := []string{}

	if ctx.Spec.Components != nil && ctx.Spec.Components.Schemas != nil {
		for schemaName, schemaRef := range ctx.Spec.Components.Schemas {
			if schemaRef.Value != nil {
				issues = append(issues, r.checkSchemaExamples(schemaName, schemaRef.Value, "")...)
			}
		}
	}

	if len(issues) == 0 {
		return core.RuleResult{
			RuleID: r.ID(),
			Passed: true,
			Detail: "All schema examples match their declared types",
		}
	}

	return core.RuleResult{
		RuleID: r.ID(),
		Passed: false,
		Detail: fmt.Sprintf("Found %d type/example mismatches: %s", len(issues), strings.Join(issues[:min(3, len(issues))], "; ")),
	}
}

func (r *SchemaExampleConsistencyRule) checkSchemaExamples(schemaName string, schema *openapi3.Schema, path string) []string {
	return r.checkSchemaExamplesWithDepth(schemaName, schema, path, 0, make(map[*openapi3.Schema]bool))
}

func (r *SchemaExampleConsistencyRule) checkSchemaExamplesWithDepth(schemaName string, schema *openapi3.Schema, path string, depth int, visited map[*openapi3.Schema]bool) []string {
	// Prevent infinite recursion with depth limit and cycle detection
	if depth > 10 || visited[schema] {
		return []string{}
	}
	visited[schema] = true
	defer func() { visited[schema] = false }()

	issues := []string{}
	currentPath := schemaName
	if path != "" {
		currentPath = fmt.Sprintf("%s.%s", schemaName, path)
	}

	// Check direct example
	if schema.Example != nil && schema.Type != "" {
		if !r.isExampleValidForType(schema.Example, schema.Type) {
			issues = append(issues, fmt.Sprintf("%s: %s example %v doesn't match type %s",
				currentPath, schema.Type, schema.Example, schema.Type))
		}
	}

	// Check properties recursively
	if schema.Properties != nil {
		for propName, propRef := range schema.Properties {
			if propRef.Value != nil {
				propPath := propName
				if path != "" {
					propPath = fmt.Sprintf("%s.%s", path, propName)
				}
				issues = append(issues, r.checkSchemaExamplesWithDepth(schemaName, propRef.Value, propPath, depth+1, visited)...)
			}
		}
	}

	return issues
}

func (r *SchemaExampleConsistencyRule) isExampleValidForType(example interface{}, schemaType string) bool {
	switch schemaType {
	case "string":
		_, ok := example.(string)
		return ok
	case "integer":
		switch example.(type) {
		case int, int32, int64, float64:
			// Check if float64 is actually an integer
			if f, ok := example.(float64); ok {
				return f == float64(int64(f))
			}
			return true
		default:
			return false
		}
	case "number":
		switch example.(type) {
		case int, int32, int64, float32, float64:
			return true
		default:
			return false
		}
	case "boolean":
		_, ok := example.(bool)
		return ok
	case "array":
		return reflect.TypeOf(example).Kind() == reflect.Slice
	case "object":
		return reflect.TypeOf(example).Kind() == reflect.Map
	default:
		return true // Unknown type, assume valid
	}
}

// OperationDescriptionRule checks if operations have meaningful descriptions
type OperationDescriptionRule struct{}

func (r *OperationDescriptionRule) ID() string {
	return "operation-description"
}

func (r *OperationDescriptionRule) Description() string {
	return "All operations should have meaningful descriptions"
}

func (r *OperationDescriptionRule) AppliesTo(version string) bool {
	return strings.HasPrefix(version, "3.")
}

func (r *OperationDescriptionRule) Evaluate(ctx *core.SpecContext) core.RuleResult {
	if len(ctx.Spec.Paths) == 0 {
		return core.RuleResult{
			RuleID: r.ID(),
			Passed: true,
			Detail: "No paths to check",
		}
	}

	totalOps := 0
	missingDesc := 0
	shortDesc := 0

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
			totalOps++

			if op.op.Description == "" {
				missingDesc++
			} else if len(strings.TrimSpace(op.op.Description)) < 10 {
				shortDesc++
			}
		}
	}

	if missingDesc == 0 && shortDesc == 0 {
		return core.RuleResult{
			RuleID: r.ID(),
			Passed: true,
			Detail: "All operations have meaningful descriptions",
		}
	}

	issues := []string{}
	if missingDesc > 0 {
		issues = append(issues, fmt.Sprintf("%d missing descriptions", missingDesc))
	}
	if shortDesc > 0 {
		issues = append(issues, fmt.Sprintf("%d too short (< 10 chars)", shortDesc))
	}

	return core.RuleResult{
		RuleID: r.ID(),
		Passed: false,
		Detail: fmt.Sprintf("Description issues: %s", strings.Join(issues, ", ")),
	}
}

// ErrorResponseRule checks if operations define proper error responses
type ErrorResponseRule struct{}

func (r *ErrorResponseRule) ID() string {
	return "operation-success-response"
}

func (r *ErrorResponseRule) Description() string {
	return "Operations should define proper error responses (400, 500)"
}

func (r *ErrorResponseRule) AppliesTo(version string) bool {
	return strings.HasPrefix(version, "3.")
}

func (r *ErrorResponseRule) Evaluate(ctx *core.SpecContext) core.RuleResult {
	if len(ctx.Spec.Paths) == 0 {
		return core.RuleResult{
			RuleID: r.ID(),
			Passed: true,
			Detail: "No paths to check",
		}
	}

	totalOps := 0
	missing400 := 0
	missing500 := 0

	for _, pathItem := range ctx.Spec.Paths {
		operations := []*openapi3.Operation{
			pathItem.Get, pathItem.Post, pathItem.Put,
			pathItem.Delete, pathItem.Patch, pathItem.Head, pathItem.Options,
		}

		for _, op := range operations {
			if op == nil {
				continue
			}
			totalOps++

			if op.Responses != nil {
				if _, has400 := op.Responses["400"]; !has400 {
					missing400++
				}
				if _, has500 := op.Responses["500"]; !has500 {
					missing500++
				}
			} else {
				missing400++
				missing500++
			}
		}
	}

	if missing400 == 0 && missing500 == 0 {
		return core.RuleResult{
			RuleID: r.ID(),
			Passed: true,
			Detail: "All operations define proper error responses",
		}
	}

	if missing400 > 0 || missing500 > 0 {
		detail := "Missing error responses: "
		if missing400 > 0 {
			detail += fmt.Sprintf("%d missing 400 responses", missing400)
		}
		if missing500 > 0 {
			if missing400 > 0 {
				detail += ", "
			}
			detail += fmt.Sprintf("%d missing 500 responses", missing500)
		}

		return core.RuleResult{
			RuleID:   r.ID(),
			Passed:   false,
			Detail:   detail,
			Severity: "warning",
			Category: "error_handling",
			Location: &core.RuleLocation{
				Path:        "$.paths",
				Component:   "operations",
				File:        "openapi.yaml",
				FileRef:     "openapi.yaml:paths section (operations missing error responses)",
				SpecSection: "paths",
			},
			Suggestion: &core.ActionableFix{
				Title:       "Add Error Response Definitions",
				Description: "Define proper error responses to help API consumers handle failures gracefully",
				Example: `responses:
  '200':
    description: Success
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/User'
  '400':
    description: Bad Request - Invalid input
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/Error'
  '500':
    description: Internal Server Error`,
				References: []string{
					"https://spec.openapis.org/oas/v3.0.3#responses-object",
					"https://tools.ietf.org/html/rfc7231#section-6",
				},
				SchemaRef: "https://spec.openapis.org/oas/v3.0.3#responses-object",
			},
			Impact: &core.ImpactAnalysis{
				Severity:    "warning",
				Category:    "error_handling",
				Description: "Missing error responses reduce API usability and debugging capability",
			},
			Metadata: map[string]string{
				"missing_400":      fmt.Sprintf("%d", missing400),
				"missing_500":      fmt.Sprintf("%d", missing500),
				"total_operations": fmt.Sprintf("%d", totalOps),
				"fix_priority":     "medium",
				"error_type":       "missing_responses",
			},
		}
	}

	return core.RuleResult{
		RuleID: r.ID(),
		Passed: false,
		Detail: fmt.Sprintf("Missing error responses: %s", strings.Join([]string{fmt.Sprintf("%d missing 400 responses", missing400), fmt.Sprintf("%d missing 500 responses", missing500)}, ", ")),
	}
}

// SecuritySchemeRule checks if the API defines security schemes
type SecuritySchemeRule struct{}

func (r *SecuritySchemeRule) ID() string {
	return "oas3-security-defined"
}

func (r *SecuritySchemeRule) Description() string {
	return "API should define security schemes for authentication"
}

func (r *SecuritySchemeRule) AppliesTo(version string) bool {
	return strings.HasPrefix(version, "3.")
}

func (r *SecuritySchemeRule) Evaluate(ctx *core.SpecContext) core.RuleResult {
	if ctx.Spec.Components == nil || ctx.Spec.Components.SecuritySchemes == nil || len(ctx.Spec.Components.SecuritySchemes) == 0 {
		return core.RuleResult{
			RuleID: r.ID(),
			Passed: false,
			Detail: "No security schemes defined",
		}
	}

	return core.RuleResult{
		RuleID: r.ID(),
		Passed: true,
		Detail: fmt.Sprintf("%d security schemes defined", len(ctx.Spec.Components.SecuritySchemes)),
	}
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
