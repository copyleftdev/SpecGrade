package test

import (
	"testing"

	"github.com/copyleftdev/specgrade/core"
	"github.com/copyleftdev/specgrade/rules"
	"github.com/getkin/kin-openapi/openapi3"
)

func TestInfoTitleRule(t *testing.T) {
	rule := &rules.InfoTitleRule{}

	tests := []struct {
		name     string
		spec     *openapi3.T
		expected bool
		detail   string
	}{
		{
			name: "spec with title",
			spec: &openapi3.T{
				Info: &openapi3.Info{
					Title: "Test API",
				},
			},
			expected: true,
			detail:   "Title present: 'Test API'",
		},
		{
			name: "spec without title",
			spec: &openapi3.T{
				Info: &openapi3.Info{
					Title: "",
				},
			},
			expected: false,
			detail:   "Missing title in info section",
		},
		{
			name:     "spec without info section",
			spec:     &openapi3.T{},
			expected: false,
			detail:   "Missing info section",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &core.SpecContext{
				Spec:    tt.spec,
				Version: "3.1.0",
			}

			result := rule.Evaluate(ctx)

			if result.Passed != tt.expected {
				t.Errorf("Expected passed=%v, got %v", tt.expected, result.Passed)
			}

			if result.Detail != tt.detail {
				t.Errorf("Expected detail=%q, got %q", tt.detail, result.Detail)
			}

			if result.RuleID != rule.ID() {
				t.Errorf("Expected ruleID=%q, got %q", rule.ID(), result.RuleID)
			}
		})
	}
}

func TestInfoVersionRule(t *testing.T) {
	rule := &rules.InfoVersionRule{}

	tests := []struct {
		name     string
		spec     *openapi3.T
		expected bool
		detail   string
	}{
		{
			name: "spec with version",
			spec: &openapi3.T{
				Info: &openapi3.Info{
					Version: "1.0.0",
				},
			},
			expected: true,
			detail:   "Version present",
		},
		{
			name: "spec without version",
			spec: &openapi3.T{
				Info: &openapi3.Info{
					Version: "",
				},
			},
			expected: false,
			detail:   "Missing version in info section",
		},
		{
			name:     "spec without info section",
			spec:     &openapi3.T{},
			expected: false,
			detail:   "Missing info section",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &core.SpecContext{
				Spec:    tt.spec,
				Version: "3.1.0",
			}

			result := rule.Evaluate(ctx)

			if result.Passed != tt.expected {
				t.Errorf("Expected passed=%v, got %v", tt.expected, result.Passed)
			}

			if result.Detail != tt.detail {
				t.Errorf("Expected detail=%q, got %q", tt.detail, result.Detail)
			}
		})
	}
}

func TestPathsExistRule(t *testing.T) {
	rule := &rules.PathsExistRule{}

	tests := []struct {
		name     string
		spec     *openapi3.T
		expected bool
		detail   string
	}{
		{
			name: "spec with paths",
			spec: &openapi3.T{
				Paths: openapi3.Paths{
					"/users": &openapi3.PathItem{},
					"/posts": &openapi3.PathItem{},
				},
			},
			expected: true,
			detail:   "2 paths defined",
		},
		{
			name: "spec without paths",
			spec: &openapi3.T{
				Paths: openapi3.Paths{},
			},
			expected: false,
			detail:   "No paths defined",
		},
		{
			name:     "spec with nil paths",
			spec:     &openapi3.T{},
			expected: false,
			detail:   "No paths defined",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &core.SpecContext{
				Spec:    tt.spec,
				Version: "3.1.0",
			}

			result := rule.Evaluate(ctx)

			if result.Passed != tt.expected {
				t.Errorf("Expected passed=%v, got %v", tt.expected, result.Passed)
			}

			if result.Detail != tt.detail {
				t.Errorf("Expected detail=%q, got %q", tt.detail, result.Detail)
			}
		})
	}
}

func TestRuleAppliesTo(t *testing.T) {
	rules := []core.Rule{
		&rules.InfoTitleRule{},
		&rules.InfoVersionRule{},
		&rules.PathsExistRule{},
		&rules.OperationIDRule{},
	}

	testVersions := []struct {
		version  string
		expected bool
	}{
		{"3.0.0", true},
		{"3.1.0", true},
		{"2.0", false},
		{"4.0.0", false},
	}

	for _, rule := range rules {
		for _, tv := range testVersions {
			t.Run(rule.ID()+"_"+tv.version, func(t *testing.T) {
				result := rule.AppliesTo(tv.version)
				if result != tv.expected {
					t.Errorf("Rule %s with version %s: expected %v, got %v",
						rule.ID(), tv.version, tv.expected, result)
				}
			})
		}
	}
}
