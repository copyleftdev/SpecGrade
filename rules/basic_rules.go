package rules

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/codetestcode/specgrade/core"
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
			RuleID: r.ID(),
			Passed: false,
			Detail: "Missing info section",
		}
	}

	if ctx.Spec.Info.Title == "" {
		return core.RuleResult{
			RuleID: r.ID(),
			Passed: false,
			Detail: "Missing title in info section",
		}
	}

	return core.RuleResult{
		RuleID: r.ID(),
		Passed: true,
		Detail: "Title present",
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
		operations := []struct{
			method string
			op *openapi3.Operation
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
