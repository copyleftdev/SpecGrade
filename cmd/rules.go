package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/copyleftdev/specgrade/registry"
	"github.com/copyleftdev/specgrade/rules"
)

var rulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "Manage validation rules",
	Long:  "Commands for managing and discovering validation rules",
}

var rulesListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all available validation rules",
	Long:  "List all available validation rules with their descriptions and applicable versions",
	RunE:  listRules,
}

func init() {
	rootCmd.AddCommand(rulesCmd)
	rulesCmd.AddCommand(rulesListCmd)
}

func listRules(cmd *cobra.Command, args []string) error {
	ruleRegistry := registry.NewRuleRegistry()
	
	// Register all rules
	ruleRegistry.Register(&rules.InfoTitleRule{})
	ruleRegistry.Register(&rules.InfoVersionRule{})
	ruleRegistry.Register(&rules.PathsExistRule{})
	ruleRegistry.Register(&rules.OperationIDRule{})

	allRules := ruleRegistry.AllRules()
	
	fmt.Printf("📋 Available Rules (%d total)\n\n", len(allRules))
	
	for _, rule := range allRules {
		fmt.Printf("🔍 %s\n", rule.ID())
		fmt.Printf("   Description: %s\n", rule.Description())
		
		// Check which versions this rule applies to
		versions := []string{}
		testVersions := []string{"3.0.0", "3.1.0"}
		for _, version := range testVersions {
			if rule.AppliesTo(version) {
				versions = append(versions, version)
			}
		}
		fmt.Printf("   Applies to: OpenAPI %s\n\n", strings.Join(versions, ", "))
	}
	
	return nil
}
