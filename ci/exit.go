package ci

import "strings"

// ExitHandler determines exit codes for CI/CD integration
type ExitHandler struct {
	failThreshold string
}

// NewExitHandler creates a new exit handler
func NewExitHandler(failThreshold string) *ExitHandler {
	return &ExitHandler{
		failThreshold: strings.ToUpper(failThreshold),
	}
}

// Handle returns the appropriate exit code based on the grade and threshold
func (e *ExitHandler) Handle(grade string) int {
	grade = strings.ToUpper(grade)

	// Grade hierarchy (higher is better)
	gradeValues := map[string]int{
		"A+": 12,
		"A":  11,
		"A-": 10,
		"B+": 9,
		"B":  8,
		"B-": 7,
		"C+": 6,
		"C":  5,
		"C-": 4,
		"D":  3,
		"F":  0,
	}

	actualValue, exists := gradeValues[grade]
	if !exists {
		return 1 // Unknown grade, fail
	}

	thresholdValue, exists := gradeValues[e.failThreshold]
	if !exists {
		return 1 // Unknown threshold, fail
	}

	if actualValue >= thresholdValue {
		return 0 // Success
	}

	return 1 // Failure
}
