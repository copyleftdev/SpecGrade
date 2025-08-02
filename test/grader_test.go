package test

import (
	"testing"

	"github.com/codetestcode/specgrade/core"
	"github.com/codetestcode/specgrade/reporter"
)

func TestDefaultGrader(t *testing.T) {
	grader := reporter.NewDefaultGrader()

	tests := []struct {
		name          string
		results       []core.RuleResult
		expectedGrade string
		expectedScore int
	}{
		{
			name:          "no results",
			results:       []core.RuleResult{},
			expectedGrade: "F",
			expectedScore: 0,
		},
		{
			name: "perfect score",
			results: []core.RuleResult{
				{RuleID: "RULE1", Passed: true},
				{RuleID: "RULE2", Passed: true},
				{RuleID: "RULE3", Passed: true},
			},
			expectedGrade: "A+",
			expectedScore: 100,
		},
		{
			name: "90% score - A grade",
			results: []core.RuleResult{
				{RuleID: "RULE1", Passed: true},
				{RuleID: "RULE2", Passed: true},
				{RuleID: "RULE3", Passed: true},
				{RuleID: "RULE4", Passed: true},
				{RuleID: "RULE5", Passed: true},
				{RuleID: "RULE6", Passed: true},
				{RuleID: "RULE7", Passed: true},
				{RuleID: "RULE8", Passed: true},
				{RuleID: "RULE9", Passed: true},
				{RuleID: "RULE10", Passed: false},
			},
			expectedGrade: "A",
			expectedScore: 90,
		},
		{
			name: "75% score - B grade",
			results: []core.RuleResult{
				{RuleID: "RULE1", Passed: true},
				{RuleID: "RULE2", Passed: true},
				{RuleID: "RULE3", Passed: true},
				{RuleID: "RULE4", Passed: false},
			},
			expectedGrade: "B",
			expectedScore: 75,
		},
		{
			name: "50% score - D grade",
			results: []core.RuleResult{
				{RuleID: "RULE1", Passed: true},
				{RuleID: "RULE2", Passed: false},
			},
			expectedGrade: "D",
			expectedScore: 50,
		},
		{
			name: "0% score - F grade",
			results: []core.RuleResult{
				{RuleID: "RULE1", Passed: false},
				{RuleID: "RULE2", Passed: false},
			},
			expectedGrade: "F",
			expectedScore: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grade := grader.Grade(tt.results)
			score := grader.CalculateScore(tt.results)

			if grade != tt.expectedGrade {
				t.Errorf("Expected grade %s, got %s", tt.expectedGrade, grade)
			}

			if score != tt.expectedScore {
				t.Errorf("Expected score %d, got %d", tt.expectedScore, score)
			}
		})
	}
}

func TestGradeBoundaries(t *testing.T) {
	grader := reporter.NewDefaultGrader()

	// Test boundary conditions
	boundaryTests := []struct {
		name          string
		passedCount   int
		totalCount    int
		expectedGrade string
	}{
		{"95% - A+", 95, 100, "A+"},
		{"94% - A", 94, 100, "A"},
		{"90% - A", 90, 100, "A"},
		{"89% - A-", 89, 100, "A-"},
		{"85% - A-", 85, 100, "A-"},
		{"84% - B+", 84, 100, "B+"},
		{"80% - B+", 80, 100, "B+"},
		{"79% - B", 79, 100, "B"},
		{"75% - B", 75, 100, "B"},
		{"74% - B-", 74, 100, "B-"},
		{"70% - B-", 70, 100, "B-"},
		{"69% - C+", 69, 100, "C+"},
		{"65% - C+", 65, 100, "C+"},
		{"64% - C", 64, 100, "C"},
		{"60% - C", 60, 100, "C"},
		{"59% - C-", 59, 100, "C-"},
		{"55% - C-", 55, 100, "C-"},
		{"54% - D", 54, 100, "D"},
		{"50% - D", 50, 100, "D"},
		{"49% - F", 49, 100, "F"},
	}

	for _, tt := range boundaryTests {
		t.Run(tt.name, func(t *testing.T) {
			results := make([]core.RuleResult, tt.totalCount)
			for i := 0; i < tt.totalCount; i++ {
				results[i] = core.RuleResult{
					RuleID: "RULE" + string(rune(i)),
					Passed: i < tt.passedCount,
				}
			}

			grade := grader.Grade(results)
			if grade != tt.expectedGrade {
				t.Errorf("Expected grade %s, got %s", tt.expectedGrade, grade)
			}
		})
	}
}
