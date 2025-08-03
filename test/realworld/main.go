package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/copyleftdev/specgrade/test/realworld/tools"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <command>")
		fmt.Println("Commands:")
		fmt.Println("  list-sources    - List all available API sources")
		fmt.Println("  collect-all     - Download all API specifications")
		fmt.Println("  collect <name>  - Download a specific API specification")
		fmt.Println("  validate-all    - Validate all collected APIs")
		fmt.Println("  stats           - Show collection statistics")
		return
	}

	command := os.Args[1]
	baseDir := "./collected-apis"

	switch command {
	case "list-sources":
		listSources()
	case "collect-all":
		collectAll(baseDir)
	case "collect":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run main.go collect <api-name>")
			return
		}
		collectSpecific(baseDir, os.Args[2])
	case "validate-all":
		validateAll(baseDir)
	case "stats":
		showStats(baseDir)
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}

func listSources() {
	fmt.Println("🌍 Available Real-World API Sources:")
	fmt.Println("=====================================")
	
	sources := tools.GetRealWorldAPISources()
	for _, source := range sources {
		fmt.Printf("\n📋 %s (%s)\n", source.Name, source.Provider)
		fmt.Printf("   Category: %s\n", source.Category)
		fmt.Printf("   Expected Grade: %s\n", source.ExpectedGrade)
		fmt.Printf("   Complexity: %s\n", source.Complexity)
		fmt.Printf("   Description: %s\n", source.Description)
		fmt.Printf("   URL: %s\n", source.URL)
		fmt.Printf("   Tags: %v\n", source.Tags)
	}
	
	fmt.Printf("\n📊 Total APIs: %d\n", len(sources))
}

func collectAll(baseDir string) {
	fmt.Println("🔄 Collecting all real-world API specifications...")
	
	collector := tools.NewCollector(baseDir)
	err := collector.CollectAll()
	if err != nil {
		log.Fatalf("Failed to collect APIs: %v", err)
	}
	
	fmt.Println("✅ Successfully collected all API specifications!")
	showStats(baseDir)
}

func collectSpecific(baseDir, apiName string) {
	fmt.Printf("🔄 Collecting API specification for: %s\n", apiName)
	
	sources := tools.GetRealWorldAPISources()
	var targetSource *tools.APISource
	
	for _, source := range sources {
		if source.Name == apiName {
			targetSource = &source
			break
		}
	}
	
	if targetSource == nil {
		fmt.Printf("❌ API '%s' not found in sources\n", apiName)
		fmt.Println("Available APIs:")
		for _, source := range sources {
			fmt.Printf("  - %s\n", source.Name)
		}
		return
	}
	
	collector := tools.NewCollector(baseDir)
	err := collector.CollectAPI(*targetSource)
	if err != nil {
		log.Fatalf("Failed to collect API %s: %v", apiName, err)
	}
	
	fmt.Printf("✅ Successfully collected %s API specification!\n", apiName)
}

func validateAll(baseDir string) {
	fmt.Println("🔍 Validating all collected API specifications...")
	
	// Find the SpecGrade binary
	specGradePath, err := findSpecGradeBinary()
	if err != nil {
		log.Fatalf("Failed to find SpecGrade binary: %v", err)
	}
	
	validator := tools.NewBatchValidator(baseDir, specGradePath)
	report, err := validator.ValidateAll()
	if err != nil {
		log.Fatalf("Failed to validate APIs: %v", err)
	}
	
	fmt.Printf("📊 Validation Results:\n")
	fmt.Printf("   Total APIs: %d\n", report.TotalAPIs)
	fmt.Printf("   Successful: %d\n", report.SuccessfulAPIs)
	fmt.Printf("   Failed: %d\n", report.FailedAPIs)
	fmt.Printf("   Total Time: %v\n", report.ValidationTime)
	
	fmt.Println("\n📋 Individual Results:")
	for _, result := range report.Results {
		status := "✅"
		if !result.Success {
			status = "❌"
		}
		fmt.Printf("   %s %s (%s): %s (Score: %d)\n", 
			status, result.APIName, result.Provider, result.ActualGrade, result.Score)
	}
	
	// Save detailed report
	reportFile := filepath.Join(baseDir, "validation-report.json")
	err = validator.SaveReport(report, reportFile)
	if err != nil {
		log.Printf("Failed to save report: %v", err)
	} else {
		fmt.Printf("\n💾 Detailed report saved to: %s\n", reportFile)
	}
}

func showStats(baseDir string) {
	fmt.Println("📊 Collection Statistics:")
	fmt.Println("========================")
	
	collector := tools.NewCollector(baseDir)
	stats, err := collector.GetAPIStats()
	if err != nil {
		log.Fatalf("Failed to get stats: %v", err)
	}
	
	for key, value := range stats {
		fmt.Printf("   %s: %v\n", key, value)
	}
	
	// List collected APIs
	apis, err := collector.ListCollectedAPIs()
	if err != nil {
		log.Printf("Failed to list APIs: %v", err)
		return
	}
	
	if len(apis) > 0 {
		fmt.Println("\n📋 Collected APIs:")
		for _, api := range apis {
			fmt.Printf("   - %s (%s) - %s [%s]\n", 
				api.Name, api.Provider, api.ExpectedGrade, api.Category)
		}
	}
}

func findSpecGradeBinary() (string, error) {
	// Try different possible locations
	candidates := []string{
		"../../../build/specgrade",
		"../../build/specgrade", 
		"../build/specgrade",
		"./build/specgrade",
		"specgrade", // In PATH
	}
	
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			abs, err := filepath.Abs(candidate)
			if err != nil {
				return candidate, nil
			}
			return abs, nil
		}
	}
	
	return "", fmt.Errorf("SpecGrade binary not found in any of the expected locations")
}
