package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/codetestcode/specgrade/core"
	"github.com/codetestcode/specgrade/fetcher"
	"github.com/codetestcode/specgrade/registry"
	"github.com/codetestcode/specgrade/reporter"
	"github.com/codetestcode/specgrade/rules"
	"github.com/codetestcode/specgrade/runner"
	"github.com/codetestcode/specgrade/test/generator"
)

// TestSuiteRunner runs comprehensive validation tests
type TestSuiteRunner struct {
	registry *registry.DefaultRuleRegistry
	grader   *reporter.DefaultGrader
}

// NewTestSuiteRunner creates a new test suite runner
func NewTestSuiteRunner() *TestSuiteRunner {
	reg := registry.NewDefaultRuleRegistry()
	
	// Register all rules
	reg.RegisterRule("3.1.0", &rules.InfoTitleRule{})
	reg.RegisterRule("3.1.0", &rules.InfoVersionRule{})
	reg.RegisterRule("3.1.0", &rules.PathsExistRule{})
	reg.RegisterRule("3.1.0", &rules.OperationIDRule{})
	reg.RegisterRule("3.1.0", &rules.SchemaExampleConsistencyRule{})
	reg.RegisterRule("3.1.0", &rules.OperationDescriptionRule{})
	reg.RegisterRule("3.1.0", &rules.ErrorResponseRule{})
	reg.RegisterRule("3.1.0", &rules.SecuritySchemeRule{})
	
	return &TestSuiteRunner{
		registry: reg,
		grader:   &reporter.DefaultGrader{},
	}
}

// TestGradeDistribution validates that generated specs achieve expected grades
func TestGradeDistribution(t *testing.T) {
	runner := NewTestSuiteRunner()
	gen := generator.NewSpecGenerator()
	profiles := generator.PredefinedProfiles()
	
	for profileName, profile := range profiles {
		t.Run(fmt.Sprintf("Profile_%s", profileName), func(t *testing.T) {
			// Generate spec with known quality profile
			specContent := gen.GenerateSpec(profile)
			
			// Write to temporary file
			tmpDir := t.TempDir()
			specFile := filepath.Join(tmpDir, "openapi.yaml")
			err := os.WriteFile(specFile, []byte(specContent), 0644)
			require.NoError(t, err)
			
			// Load and validate
			loader := &fetcher.LocalSpecLoader{TargetDir: tmpDir}
			spec, err := loader.Load("3.1.0")
			require.NoError(t, err)
			
			// Run validation
			ctx := &core.SpecContext{
				Spec:    spec,
				Version: "3.1.0",
			}
			
			testRunner := &runner.DefaultRunner{Registry: runner.registry}
			results := testRunner.Run(ctx, []string{})
			
			// Grade the results
			report := &core.Report{
				Version: "3.1.0",
				Rules:   results,
			}
			runner.grader.Grade(report)
			
			// Validate grade is in expected range
			expectedGrades := map[string][]string{
				"perfect":   {"A+"},
				"excellent": {"A+", "A"},
				"good":      {"A", "B"},
				"average":   {"B", "C"},
				"poor":      {"C", "D"},
				"failing":   {"D", "F"},
			}
			
			validGrades := expectedGrades[profileName]
			assert.Contains(t, validGrades, report.Grade, 
				"Profile %s should achieve grade in %v, got %s (score: %d)", 
				profileName, validGrades, report.Grade, report.Score)
			
			t.Logf("Profile %s: Grade %s, Score %d", profileName, report.Grade, report.Score)
		})
	}
}

// TestEdgeCases validates handling of edge cases and unknowns
func TestEdgeCases(t *testing.T) {
	runner := NewTestSuiteRunner()
	gen := generator.NewEdgeCaseGenerator()
	
	testCases := []struct {
		name     string
		specFunc func() string
		shouldPass bool
	}{
		{
			name:     "CircularReferences",
			specFunc: gen.GenerateCircularRef,
			shouldPass: true, // Should handle gracefully
		},
		{
			name:     "DeepNesting_10_Levels",
			specFunc: func() string { return gen.GenerateDeepNesting(10) },
			shouldPass: true,
		},
		{
			name:     "DeepNesting_50_Levels",
			specFunc: func() string { return gen.GenerateDeepNesting(50) },
			shouldPass: true, // Should not crash
		},
		{
			name:     "UnicodeContent",
			specFunc: gen.GenerateUnicodeContent,
			shouldPass: true,
		},
		{
			name:     "MassiveSpec_100_Endpoints",
			specFunc: func() string { return gen.GenerateMassiveSpec(100) },
			shouldPass: true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Generate spec
			specContent := tc.specFunc()
			
			// Write to temporary file
			tmpDir := t.TempDir()
			specFile := filepath.Join(tmpDir, "openapi.yaml")
			err := os.WriteFile(specFile, []byte(specContent), 0644)
			require.NoError(t, err)
			
			// Load and validate - should not panic
			loader := &fetcher.LocalSpecLoader{TargetDir: tmpDir}
			spec, err := loader.Load("3.1.0")
			
			if tc.shouldPass {
				require.NoError(t, err, "Edge case %s should load successfully", tc.name)
				
				// Run validation - should not panic
				ctx := &core.SpecContext{
					Spec:    spec,
					Version: "3.1.0",
				}
				
				testRunner := &runner.DefaultRunner{Registry: runner.registry}
				results := testRunner.Run(ctx, []string{})
				
				// Should produce valid results
				assert.NotEmpty(t, results, "Edge case %s should produce validation results", tc.name)
				
				// Grade should be valid
				report := &core.Report{
					Version: "3.1.0",
					Rules:   results,
				}
				runner.grader.Grade(report)
				
				validGrades := []string{"A+", "A", "B", "C", "D", "F"}
				assert.Contains(t, validGrades, report.Grade, 
					"Edge case %s should produce valid grade, got %s", tc.name, report.Grade)
				
				t.Logf("Edge case %s: Grade %s, Score %d", tc.name, report.Grade, report.Score)
			} else {
				// Some edge cases might fail to load, but should fail gracefully
				if err != nil {
					t.Logf("Edge case %s failed as expected: %v", tc.name, err)
				}
			}
		})
	}
}

// TestPropertyBased uses property-based testing to find unknown issues
func TestPropertyBased(t *testing.T) {
	runner := NewTestSuiteRunner()
	
	// Property: All valid specs should produce valid grades
	property := func(complexity uint8, missingDesc float32, typeMismatches uint8) bool {
		// Constrain inputs to reasonable ranges
		if complexity == 0 {
			complexity = 1
		}
		if complexity > 10 {
			complexity = 10
		}
		if missingDesc < 0 {
			missingDesc = 0
		}
		if missingDesc > 1 {
			missingDesc = 1
		}
		if typeMismatches > 20 {
			typeMismatches = 20
		}
		
		// Generate spec with random properties
		gen := generator.NewSpecGenerator()
		profile := generator.QualityProfile{
			TargetGrade:         "C", // Don't care about target
			MissingDescriptions: float64(missingDesc),
			TypeMismatches:      int(typeMismatches),
			MissingErrorCodes:   0.5,
			SecurityIssues:      false,
			ComplexityLevel:     int(complexity),
		}
		
		specContent := gen.GenerateSpec(profile)
		
		// Write to temporary file
		tmpDir, err := os.MkdirTemp("", "property_test")
		if err != nil {
			return false
		}
		defer os.RemoveAll(tmpDir)
		
		specFile := filepath.Join(tmpDir, "openapi.yaml")
		err = os.WriteFile(specFile, []byte(specContent), 0644)
		if err != nil {
			return false
		}
		
		// Load and validate
		loader := &fetcher.LocalSpecLoader{TargetDir: tmpDir}
		spec, err := loader.Load("3.1.0")
		if err != nil {
			return false // Invalid spec generation
		}
		
		// Run validation
		ctx := &core.SpecContext{
			Spec:    spec,
			Version: "3.1.0",
		}
		
		testRunner := &runner.DefaultRunner{Registry: runner.registry}
		results := testRunner.Run(ctx, []string{})
		
		// Grade the results
		report := &core.Report{
			Version: "3.1.0",
			Rules:   results,
		}
		runner.grader.Grade(report)
		
		// Properties that should always hold:
		// 1. Score should be 0-100
		if report.Score < 0 || report.Score > 100 {
			return false
		}
		
		// 2. Grade should be valid
		validGrades := []string{"A+", "A", "B", "C", "D", "F"}
		gradeValid := false
		for _, grade := range validGrades {
			if report.Grade == grade {
				gradeValid = true
				break
			}
		}
		if !gradeValid {
			return false
		}
		
		// 3. Should have results for all registered rules
		expectedRules := 8 // Number of rules we registered
		if len(results) != expectedRules {
			return false
		}
		
		return true
	}
	
	// Run property-based test
	err := quick.Check(property, &quick.Config{
		MaxCount: 100, // Run 100 random tests
	})
	assert.NoError(t, err, "Property-based testing should pass")
}

// TestPerformance validates performance with various spec sizes
func TestPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}
	
	runner := NewTestSuiteRunner()
	gen := generator.NewEdgeCaseGenerator()
	
	performanceTests := []struct {
		name        string
		specFunc    func() string
		maxDuration string // Not used in this simple version, but could be
	}{
		{
			name:     "Small_10_Endpoints",
			specFunc: func() string { return gen.GenerateMassiveSpec(10) },
		},
		{
			name:     "Medium_100_Endpoints",
			specFunc: func() string { return gen.GenerateMassiveSpec(100) },
		},
		{
			name:     "Large_500_Endpoints",
			specFunc: func() string { return gen.GenerateMassiveSpec(500) },
		},
	}
	
	for _, tc := range performanceTests {
		t.Run(tc.name, func(t *testing.T) {
			// Generate spec
			specContent := tc.specFunc()
			
			// Write to temporary file
			tmpDir := t.TempDir()
			specFile := filepath.Join(tmpDir, "openapi.yaml")
			err := os.WriteFile(specFile, []byte(specContent), 0644)
			require.NoError(t, err)
			
			// Load and validate
			loader := &fetcher.LocalSpecLoader{TargetDir: tmpDir}
			spec, err := loader.Load("3.1.0")
			require.NoError(t, err)
			
			// Run validation and measure performance
			ctx := &core.SpecContext{
				Spec:    spec,
				Version: "3.1.0",
			}
			
			testRunner := &runner.DefaultRunner{Registry: runner.registry}
			results := testRunner.Run(ctx, []string{})
			
			// Should complete successfully
			assert.NotEmpty(t, results)
			
			// Grade the results
			report := &core.Report{
				Version: "3.1.0",
				Rules:   results,
			}
			runner.grader.Grade(report)
			
			t.Logf("Performance test %s: Grade %s, Score %d", tc.name, report.Grade, report.Score)
		})
	}
}

// TestRegressionSuite ensures consistent behavior across versions
func TestRegressionSuite(t *testing.T) {
	// This would contain known specs with expected grades
	// For now, we'll test our existing samples
	
	runner := NewTestSuiteRunner()
	
	regressionCases := []struct {
		name          string
		specDir       string
		expectedGrade string
		tolerance     int // Score tolerance
	}{
		{
			name:          "Perfect_Sample",
			specDir:       "../sample-spec",
			expectedGrade: "A+",
			tolerance:     5,
		},
		{
			name:          "Bad_Example",
			specDir:       "../sample-spec/bad-example",
			expectedGrade: "C",
			tolerance:     10,
		},
	}
	
	for _, tc := range regressionCases {
		t.Run(tc.name, func(t *testing.T) {
			// Load spec
			loader := &fetcher.LocalSpecLoader{TargetDir: tc.specDir}
			spec, err := loader.Load("3.1.0")
			require.NoError(t, err)
			
			// Run validation
			ctx := &core.SpecContext{
				Spec:    spec,
				Version: "3.1.0",
			}
			
			testRunner := &runner.DefaultRunner{Registry: runner.registry}
			results := testRunner.Run(ctx, []string{})
			
			// Grade the results
			report := &core.Report{
				Version: "3.1.0",
				Rules:   results,
			}
			runner.grader.Grade(report)
			
			// Check grade matches expectation
			assert.Equal(t, tc.expectedGrade, report.Grade, 
				"Regression test %s should maintain grade %s, got %s", 
				tc.name, tc.expectedGrade, report.Grade)
			
			t.Logf("Regression test %s: Grade %s, Score %d", tc.name, report.Grade, report.Score)
		})
	}
}
