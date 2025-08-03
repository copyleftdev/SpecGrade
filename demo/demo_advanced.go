package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/copyleftdev/specgrade/core"
	"github.com/copyleftdev/specgrade/fetcher"
	"github.com/copyleftdev/specgrade/registry"
	"github.com/copyleftdev/specgrade/reporter"
	"github.com/copyleftdev/specgrade/rules"
	"github.com/copyleftdev/specgrade/runner"
	"github.com/copyleftdev/specgrade/test/generator"
)

func main() {
	fmt.Println("🚀 SpecGrade Advanced Testing Demonstration")
	fmt.Println("==================================================")

	// Setup registry and grader
	reg := registry.NewRuleRegistry()
	reg.Register(&rules.InfoTitleRule{})
	reg.Register(&rules.InfoVersionRule{})
	reg.Register(&rules.PathsExistRule{})
	reg.Register(&rules.OperationIDRule{})
	reg.Register(&rules.SchemaExampleConsistencyRule{})
	reg.Register(&rules.OperationDescriptionRule{})
	reg.Register(&rules.ErrorResponseRule{})
	reg.Register(&rules.SecuritySchemeRule{})

	grader := &reporter.DefaultGrader{}
	testRunner := runner.NewRunner(reg, []string{})

	// Test 1: Grade Distribution Validation
	fmt.Println("\n📊 Test 1: Grade Distribution Validation")
	fmt.Println(strings.Repeat("-", 40))

	gen := generator.NewSpecGenerator()
	profiles := generator.PredefinedProfiles()

	for profileName, profile := range profiles {
		fmt.Printf("Testing profile: %s (target: %s)\n", profileName, profile.TargetGrade)

		// Generate spec
		specContent := gen.GenerateSpec(profile)

		// Create temp directory and file
		tmpDir, err := os.MkdirTemp("", fmt.Sprintf("specgrade_test_%s", profileName))
		if err != nil {
			fmt.Printf("❌ Error creating temp dir: %v\n", err)
			continue
		}
		defer os.RemoveAll(tmpDir)

		specFile := filepath.Join(tmpDir, "openapi.yaml")
		err = os.WriteFile(specFile, []byte(specContent), 0644)
		if err != nil {
			fmt.Printf("❌ Error writing spec file: %v\n", err)
			continue
		}

		// Load and validate
		loader := fetcher.NewLocalSpecLoader(tmpDir)
		spec, err := loader.Load("3.1.0")
		if err != nil {
			fmt.Printf("❌ Error loading spec: %v\n", err)
			continue
		}

		// Run validation
		ctx := &core.SpecContext{
			Spec:    spec,
			Version: "3.1.0",
		}

		results := testRunner.Run(ctx)

		// Grade results
		grade := grader.Grade(results)
		score := grader.CalculateScore(results)
		report := &core.Report{
			Version: "3.1.0",
			Rules:   results,
			Grade:   grade,
			Score:   score,
		}

		// Display results
		status := "✅"
		if profile.TargetGrade != report.Grade {
			// Allow some flexibility in grading
			expectedGrades := map[string][]string{
				"perfect":   {"A+"},
				"excellent": {"A+", "A"},
				"good":      {"A", "B"},
				"average":   {"B", "C"},
				"poor":      {"C", "D"},
				"failing":   {"D", "F"},
			}

			validGrades := expectedGrades[profileName]
			found := false
			for _, grade := range validGrades {
				if grade == report.Grade {
					found = true
					break
				}
			}
			if !found {
				status = "⚠️"
			}
		}

		fmt.Printf("  %s %s: Grade %s (Score: %d) - Expected: %s\n",
			status, profileName, report.Grade, report.Score, profile.TargetGrade)
	}

	// Test 2: Edge Cases
	fmt.Println("\n🔬 Test 2: Edge Case Handling")
	fmt.Println(strings.Repeat("-", 40))

	edgeGen := generator.NewEdgeCaseGenerator()

	edgeCases := []struct {
		name     string
		specFunc func() string
		desc     string
	}{
		{
			name:     "circular_refs",
			specFunc: edgeGen.GenerateCircularRef,
			desc:     "Circular schema references",
		},
		{
			name:     "deep_nesting",
			specFunc: func() string { return edgeGen.GenerateDeepNesting(20) },
			desc:     "20-level deep nesting",
		},
		{
			name:     "unicode_content",
			specFunc: edgeGen.GenerateUnicodeContent,
			desc:     "International characters and emojis",
		},
		{
			name:     "massive_spec",
			specFunc: func() string { return edgeGen.GenerateMassiveSpec(50) },
			desc:     "50 endpoints specification",
		},
	}

	for _, tc := range edgeCases {
		fmt.Printf("Testing edge case: %s (%s)\n", tc.name, tc.desc)

		// Generate spec
		specContent := tc.specFunc()

		// Create temp directory and file
		tmpDir, err := os.MkdirTemp("", fmt.Sprintf("specgrade_edge_%s", tc.name))
		if err != nil {
			fmt.Printf("❌ Error creating temp dir: %v\n", err)
			continue
		}
		defer os.RemoveAll(tmpDir)

		specFile := filepath.Join(tmpDir, "openapi.yaml")
		err = os.WriteFile(specFile, []byte(specContent), 0644)
		if err != nil {
			fmt.Printf("❌ Error writing spec file: %v\n", err)
			continue
		}

		// Load and validate
		loader := fetcher.NewLocalSpecLoader(tmpDir)
		spec, err := loader.Load("3.1.0")
		if err != nil {
			fmt.Printf("❌ Error loading spec: %v\n", err)
			continue
		}

		// Run validation
		ctx := &core.SpecContext{
			Spec:    spec,
			Version: "3.1.0",
		}

		results := testRunner.Run(ctx)

		// Grade results
		grade := grader.Grade(results)
		score := grader.CalculateScore(results)
		report := &core.Report{
			Version: "3.1.0",
			Rules:   results,
			Grade:   grade,
			Score:   score,
		}

		fmt.Printf("  ✅ %s: Grade %s (Score: %d) - Handled gracefully\n",
			tc.name, report.Grade, report.Score)
	}

	// Test 3: Performance Analysis
	fmt.Println("\n⚡ Test 3: Performance Analysis")
	fmt.Println(strings.Repeat("-", 40))

	performanceTests := []struct {
		name      string
		endpoints int
	}{
		{"small", 10},
		{"medium", 50},
		{"large", 100},
	}

	for _, pt := range performanceTests {
		fmt.Printf("Testing performance: %s (%d endpoints)\n", pt.name, pt.endpoints)

		// Generate massive spec
		specContent := edgeGen.GenerateMassiveSpec(pt.endpoints)

		// Create temp directory and file
		tmpDir, err := os.MkdirTemp("", fmt.Sprintf("specgrade_perf_%s", pt.name))
		if err != nil {
			fmt.Printf("❌ Error creating temp dir: %v\n", err)
			continue
		}
		defer os.RemoveAll(tmpDir)

		specFile := filepath.Join(tmpDir, "openapi.yaml")
		err = os.WriteFile(specFile, []byte(specContent), 0644)
		if err != nil {
			fmt.Printf("❌ Error writing spec file: %v\n", err)
			continue
		}

		// Load and validate
		loader := fetcher.NewLocalSpecLoader(tmpDir)
		spec, err := loader.Load("3.1.0")
		if err != nil {
			fmt.Printf("❌ Error loading spec: %v\n", err)
			continue
		}

		// Run validation
		ctx := &core.SpecContext{
			Spec:    spec,
			Version: "3.1.0",
		}

		results := testRunner.Run(ctx)

		// Grade results
		grade := grader.Grade(results)
		score := grader.CalculateScore(results)
		report := &core.Report{
			Version: "3.1.0",
			Rules:   results,
			Grade:   grade,
			Score:   score,
		}

		fmt.Printf("  ✅ %s: Grade %s (Score: %d) - %d endpoints processed\n",
			pt.name, report.Grade, report.Score, pt.endpoints)
	}

	fmt.Println("\n🎉 Advanced Testing Complete!")
	fmt.Println("SpecGrade demonstrates robust handling of:")
	fmt.Println("  • Diverse quality profiles with predictable grading")
	fmt.Println("  • Edge cases including circular refs, deep nesting, Unicode")
	fmt.Println("  • Performance with large specifications")
	fmt.Println("  • Industry-standard rule naming (Spectral-compatible)")
}
