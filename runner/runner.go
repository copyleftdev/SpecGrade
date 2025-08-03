package runner

import (
	"github.com/copyleftdev/specgrade/core"
	"github.com/copyleftdev/specgrade/registry"
)

// Runner executes validation rules against OpenAPI specs
type Runner struct {
	registry  *registry.RuleRegistry
	skipRules map[string]bool
}

// NewRunner creates a new rule runner
func NewRunner(registry *registry.RuleRegistry, skipRules []string) *Runner {
	skipMap := make(map[string]bool)
	for _, ruleID := range skipRules {
		skipMap[ruleID] = true
	}

	return &Runner{
		registry:  registry,
		skipRules: skipMap,
	}
}

// Run executes all applicable rules for the given spec and version
func (r *Runner) Run(spec *core.SpecContext) []core.RuleResult {
	rules := r.registry.RulesForVersion(spec.Version)
	results := make([]core.RuleResult, 0, len(rules))

	for _, rule := range rules {
		// Skip rules that are in the skip list
		if r.skipRules[rule.ID()] {
			continue
		}

		result := rule.Evaluate(spec)
		results = append(results, result)
	}

	return results
}

// RunRule executes a specific rule by ID
func (r *Runner) RunRule(spec *core.SpecContext, ruleID string) (*core.RuleResult, error) {
	rule := r.registry.GetRule(ruleID)
	if rule == nil {
		return nil, nil // Rule not found
	}

	if !rule.AppliesTo(spec.Version) {
		return nil, nil // Rule doesn't apply to this version
	}

	result := rule.Evaluate(spec)
	return &result, nil
}
