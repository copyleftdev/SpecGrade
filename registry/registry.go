package registry

import (
	"github.com/codetestcode/specgrade/core"
)

// RuleRegistry manages the collection of validation rules
type RuleRegistry struct {
	rules []core.Rule
}

// NewRuleRegistry creates a new rule registry
func NewRuleRegistry() *RuleRegistry {
	return &RuleRegistry{
		rules: make([]core.Rule, 0),
	}
}

// Register adds a rule to the registry
func (r *RuleRegistry) Register(rule core.Rule) {
	r.rules = append(r.rules, rule)
}

// RulesForVersion returns all rules that apply to the given OpenAPI version
func (r *RuleRegistry) RulesForVersion(version string) []core.Rule {
	var applicableRules []core.Rule
	for _, rule := range r.rules {
		if rule.AppliesTo(version) {
			applicableRules = append(applicableRules, rule)
		}
	}
	return applicableRules
}

// AllRules returns all registered rules
func (r *RuleRegistry) AllRules() []core.Rule {
	return r.rules
}

// GetRule returns a rule by ID, or nil if not found
func (r *RuleRegistry) GetRule(id string) core.Rule {
	for _, rule := range r.rules {
		if rule.ID() == id {
			return rule
		}
	}
	return nil
}
