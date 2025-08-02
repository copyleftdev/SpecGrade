package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/copyleftdev/specgrade/ci"
	"github.com/copyleftdev/specgrade/core"
	"github.com/copyleftdev/specgrade/fetcher"
	"github.com/copyleftdev/specgrade/registry"
	"github.com/copyleftdev/specgrade/reporter"
	"github.com/copyleftdev/specgrade/rules"
	"github.com/copyleftdev/specgrade/runner"
	"github.com/copyleftdev/specgrade/utils"
	"github.com/copyleftdev/specgrade/versions"
)

var (
	specVersion   string
	targetDir     string
	outputFormat  string
	failThreshold string
	configPath    string
	skipRules     string
	generateDocs  bool
)

var rootCmd = &cobra.Command{
	Use:   "specgrade",
	Short: "A modular, dynamic, and CICD-optimized conformance validator for OpenAPI specifications",
	Long: `SpecGrade is a modular, dynamic, and CICD-optimized conformance validator for OpenAPI specifications. 
It fetches versioned OpenAPI schema definitions, dynamically constructs validation rule sets based on that schema, 
and grades the conformance of target API specs against those rules.`,
	RunE: runSpecGrade,
}

func init() {
	rootCmd.Flags().StringVar(&specVersion, "spec-version", "", "The official OpenAPI version to validate against (e.g., 3.1.0)")
	rootCmd.Flags().StringVar(&targetDir, "target-dir", "", "Path to the local OpenAPI spec to validate")
	rootCmd.Flags().StringVar(&outputFormat, "output-format", "", "Output format: json, cli, html, or markdown")
	rootCmd.Flags().StringVar(&failThreshold, "fail-threshold", "", "Minimum acceptable grade (A, B, etc). Will exit non-zero if below")
	rootCmd.Flags().StringVar(&configPath, "config", "", "Optional path to specgrade.yaml config file")
	rootCmd.Flags().StringVar(&skipRules, "skip", "", "Comma-separated rule IDs to ignore")
	rootCmd.Flags().BoolVar(&generateDocs, "docs", false, "Generate rule documentation (markdown)")
}

func Execute() error {
	return rootCmd.Execute()
}

func runSpecGrade(cmd *cobra.Command, args []string) error {
	// Load configuration
	config, err := utils.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create flags config for merging
	flagsConfig := &core.Config{
		SpecVersion:   specVersion,
		InputDir:      targetDir,
		OutputFormat:  outputFormat,
		FailThreshold: failThreshold,
	}

	// Parse skip rules
	if skipRules != "" {
		flagsConfig.SkipRules = strings.Split(skipRules, ",")
		for i, rule := range flagsConfig.SkipRules {
			flagsConfig.SkipRules[i] = strings.TrimSpace(rule)
		}
	}

	// Merge config with flags (flags take precedence)
	finalConfig := utils.MergeConfigWithFlags(config, flagsConfig)

	// Validate required fields
	if finalConfig.InputDir == "" {
		return fmt.Errorf("target directory is required (use --target-dir flag or input_dir in config file)")
	}

	// Set defaults if not specified
	if finalConfig.SpecVersion == "" {
		finalConfig.SpecVersion = "3.1.0"
	}
	if finalConfig.OutputFormat == "" {
		finalConfig.OutputFormat = "cli"
	}
	if finalConfig.FailThreshold == "" {
		finalConfig.FailThreshold = "B"
	}

	// Validate spec version
	if !versions.IsValidVersion(finalConfig.SpecVersion) {
		return fmt.Errorf("unsupported OpenAPI version: %s", finalConfig.SpecVersion)
	}

	// Generate documentation if requested
	if generateDocs {
		return generateRuleDocumentation()
	}

	// Initialize components
	ruleRegistry := registry.NewRuleRegistry()
	registerRules(ruleRegistry)

	specLoader := fetcher.NewLocalSpecLoader(finalConfig.InputDir)
	ruleRunner := runner.NewRunner(ruleRegistry, finalConfig.SkipRules)
	rep := reporter.NewReporter()
	exitHandler := ci.NewExitHandler(finalConfig.FailThreshold)

	// Load the OpenAPI spec
	spec, err := specLoader.Load(finalConfig.SpecVersion)
	if err != nil {
		return fmt.Errorf("failed to load OpenAPI spec: %w", err)
	}

	// Create spec context
	specContext := &core.SpecContext{
		Spec:    spec,
		Version: finalConfig.SpecVersion,
	}

	// Run validation rules
	results := ruleRunner.Run(specContext)

	// Generate report
	report := rep.GenerateReport(finalConfig.SpecVersion, results)

	// Output report in requested format
	var output string
	switch strings.ToLower(finalConfig.OutputFormat) {
	case "json":
		output, err = rep.FormatJSON(report)
		if err != nil {
			return fmt.Errorf("failed to format JSON output: %w", err)
		}
	case "markdown":
		output = rep.FormatMarkdown(report, finalConfig.InputDir)
	case "html":
		output = rep.FormatHTML(report, finalConfig.InputDir)
	case "cli":
		output = rep.FormatCLI(report, finalConfig.InputDir)
	default:
		return fmt.Errorf("unsupported output format: %s", finalConfig.OutputFormat)
	}

	fmt.Print(output)

	// Exit with appropriate code for CI/CD
	exitCode := exitHandler.Handle(report.Grade)
	if exitCode != 0 {
		os.Exit(exitCode)
	}

	return nil
}

// registerRules registers all available validation rules
func registerRules(registry *registry.RuleRegistry) {
	// Basic structural rules
	registry.Register(&rules.InfoTitleRule{})
	registry.Register(&rules.InfoVersionRule{})
	registry.Register(&rules.PathsExistRule{})
	registry.Register(&rules.OperationIDRule{})
	
	// Advanced quality rules
	registry.Register(&rules.SchemaExampleConsistencyRule{})
	registry.Register(&rules.OperationDescriptionRule{})
	registry.Register(&rules.ErrorResponseRule{})
	registry.Register(&rules.SecuritySchemeRule{})
}

// generateRuleDocumentation generates markdown documentation for all rules
func generateRuleDocumentation() error {
	ruleRegistry := registry.NewRuleRegistry()
	registerRules(ruleRegistry)

	fmt.Println("# SpecGrade Rules Documentation")
	fmt.Println("This document describes all available validation rules in SpecGrade.")
	fmt.Println()

	allRules := ruleRegistry.AllRules()
	for _, rule := range allRules {
		fmt.Printf("## %s\n\n", rule.ID())
		fmt.Printf("**Description:** %s\n\n", rule.Description())
		
		// Check which versions this rule applies to
		versions := []string{}
		testVersions := []string{"3.0.0", "3.1.0"}
		for _, version := range testVersions {
			if rule.AppliesTo(version) {
				versions = append(versions, version)
			}
		}
		fmt.Printf("**Applies to:** OpenAPI %s\n\n", strings.Join(versions, ", "))
		fmt.Print("\n---\n")
	}

	return nil
}
