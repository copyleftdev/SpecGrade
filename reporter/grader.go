package reporter

import (
	"github.com/copyleftdev/specgrade/core"
)

// DefaultGrader implements the standard grading logic
type DefaultGrader struct{}

// NewDefaultGrader creates a new default grader
func NewDefaultGrader() *DefaultGrader {
	return &DefaultGrader{}
}

// Grade calculates a letter grade based on rule results
func (g *DefaultGrader) Grade(results []core.RuleResult) string {
	if len(results) == 0 {
		return "F"
	}

	passed := 0
	for _, result := range results {
		if result.Passed {
			passed++
		}
	}

	percentage := float64(passed) / float64(len(results)) * 100

	switch {
	case percentage >= 95:
		return "A+"
	case percentage >= 90:
		return "A"
	case percentage >= 85:
		return "A-"
	case percentage >= 80:
		return "B+"
	case percentage >= 75:
		return "B"
	case percentage >= 70:
		return "B-"
	case percentage >= 65:
		return "C+"
	case percentage >= 60:
		return "C"
	case percentage >= 55:
		return "C-"
	case percentage >= 50:
		return "D"
	default:
		return "F"
	}
}

// CalculateScore returns the numeric score (0-100)
func (g *DefaultGrader) CalculateScore(results []core.RuleResult) int {
	if len(results) == 0 {
		return 0
	}

	passed := 0
	for _, result := range results {
		if result.Passed {
			passed++
		}
	}

	return int(float64(passed) / float64(len(results)) * 100)
}
