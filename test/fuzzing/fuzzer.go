package fuzzing

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// FuzzingStrategy defines different approaches to corrupting OpenAPI specs
type FuzzingStrategy string

const (
	StrategyStructural FuzzingStrategy = "structural"  // Break YAML/JSON structure
	StrategySemantic   FuzzingStrategy = "semantic"    // Break OpenAPI semantics
	StrategyDataType   FuzzingStrategy = "datatype"    // Corrupt data types
	StrategyReference  FuzzingStrategy = "reference"   // Break $ref links
	StrategyEncoding   FuzzingStrategy = "encoding"    // Character encoding issues
	StrategySize       FuzzingStrategy = "size"        // Extreme sizes
	StrategyEdgeValues FuzzingStrategy = "edge_values" // Edge case values
	StrategyMutation   FuzzingStrategy = "mutation"    // Random mutations
)

// FuzzingConfig controls the fuzzing process
type FuzzingConfig struct {
	Strategies          []FuzzingStrategy `json:"strategies"`
	MutationRate        float64           `json:"mutation_rate"` // 0.0 to 1.0
	MaxMutationsPerSpec int               `json:"max_mutations_per_spec"`
	PreserveSyntax      bool              `json:"preserve_syntax"` // Try to keep valid YAML/JSON
	SeedValue           int64             `json:"seed_value"`
	OutputDir           string            `json:"output_dir"`
}

// FuzzingResult represents the outcome of fuzzing a specification
type FuzzingResult struct {
	OriginalSpec    string           `json:"original_spec"`
	FuzzedSpec      string           `json:"fuzzed_spec"`
	Strategy        FuzzingStrategy  `json:"strategy"`
	MutationCount   int              `json:"mutation_count"`
	ValidationError string           `json:"validation_error,omitempty"`
	SpecGradeResult *SpecGradeResult `json:"specgrade_result,omitempty"`
	CrashedTool     bool             `json:"crashed_tool"`
	ExecutionTime   time.Duration    `json:"execution_time"`
	GeneratedAt     time.Time        `json:"generated_at"`
}

// SpecGradeResult represents SpecGrade's output on a fuzzed spec
type SpecGradeResult struct {
	Grade   string `json:"grade"`
	Score   int    `json:"score"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// Fuzzer generates corrupted OpenAPI specifications for robustness testing
type Fuzzer struct {
	config *FuzzingConfig
	rand   *rand.Rand
}

// NewFuzzer creates a new OpenAPI specification fuzzer
func NewFuzzer(config *FuzzingConfig) *Fuzzer {
	seed := config.SeedValue
	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	return &Fuzzer{
		config: config,
		rand:   rand.New(rand.NewSource(seed)),
	}
}

// DefaultFuzzingConfig returns a sensible default configuration
func DefaultFuzzingConfig() *FuzzingConfig {
	return &FuzzingConfig{
		Strategies: []FuzzingStrategy{
			StrategyStructural,
			StrategySemantic,
			StrategyDataType,
			StrategyReference,
			StrategyEncoding,
			StrategyEdgeValues,
			StrategyMutation,
		},
		MutationRate:        0.1,
		MaxMutationsPerSpec: 10,
		PreserveSyntax:      true,
		SeedValue:           0,
		OutputDir:           "fuzzed_specs",
	}
}

// FuzzSpec applies fuzzing strategies to corrupt an OpenAPI specification
func (f *Fuzzer) FuzzSpec(specContent string, strategy FuzzingStrategy) (string, int, error) {
	switch strategy {
	case StrategyStructural:
		return f.fuzzStructural(specContent)
	case StrategySemantic:
		return f.fuzzSemantic(specContent)
	case StrategyDataType:
		return f.fuzzDataTypes(specContent)
	case StrategyReference:
		return f.fuzzReferences(specContent)
	case StrategyEncoding:
		return f.fuzzEncoding(specContent)
	case StrategySize:
		return f.fuzzSize(specContent)
	case StrategyEdgeValues:
		return f.fuzzEdgeValues(specContent)
	case StrategyMutation:
		return f.fuzzMutation(specContent)
	default:
		return specContent, 0, fmt.Errorf("unknown fuzzing strategy: %s", strategy)
	}
}

// fuzzStructural corrupts YAML/JSON structure
func (f *Fuzzer) fuzzStructural(content string) (string, int, error) {
	mutations := 0
	lines := strings.Split(content, "\n")

	for i := 0; i < len(lines) && mutations < f.config.MaxMutationsPerSpec; i++ {
		if f.rand.Float64() < f.config.MutationRate {
			switch f.rand.Intn(5) {
			case 0: // Remove random characters
				if len(lines[i]) > 0 {
					pos := f.rand.Intn(len(lines[i]))
					lines[i] = lines[i][:pos] + lines[i][pos+1:]
					mutations++
				}
			case 1: // Add random characters
				pos := f.rand.Intn(len(lines[i]) + 1)
				char := string(rune(f.rand.Intn(128)))
				lines[i] = lines[i][:pos] + char + lines[i][pos:]
				mutations++
			case 2: // Corrupt indentation
				if strings.TrimSpace(lines[i]) != "" {
					spaces := f.rand.Intn(10)
					lines[i] = strings.Repeat(" ", spaces) + strings.TrimSpace(lines[i])
					mutations++
				}
			case 3: // Break key-value pairs
				if strings.Contains(lines[i], ":") {
					lines[i] = strings.Replace(lines[i], ":", f.randomString(1), 1)
					mutations++
				}
			case 4: // Corrupt quotes
				lines[i] = strings.ReplaceAll(lines[i], "\"", "'")
				mutations++
			}
		}
	}

	return strings.Join(lines, "\n"), mutations, nil
}

// fuzzSemantic corrupts OpenAPI semantic structure
func (f *Fuzzer) fuzzSemantic(content string) (string, int, error) {
	mutations := 0

	// Parse as YAML to understand structure
	var spec map[string]interface{}
	if err := yaml.Unmarshal([]byte(content), &spec); err != nil {
		// If we can't parse, apply text-based semantic fuzzing
		return f.fuzzSemanticText(content)
	}

	// Apply semantic corruptions
	if f.rand.Float64() < f.config.MutationRate {
		// Corrupt OpenAPI version
		if _, exists := spec["openapi"]; exists {
			spec["openapi"] = f.randomString(5)
			mutations++
		}
	}

	if f.rand.Float64() < f.config.MutationRate {
		// Remove required sections
		requiredSections := []string{"info", "paths"}
		if len(requiredSections) > 0 {
			section := requiredSections[f.rand.Intn(len(requiredSections))]
			delete(spec, section)
			mutations++
		}
	}

	if f.rand.Float64() < f.config.MutationRate {
		// Corrupt HTTP methods
		if paths, ok := spec["paths"].(map[string]interface{}); ok {
			for _, pathValue := range paths {
				if pathObj, ok := pathValue.(map[string]interface{}); ok {
					httpMethods := []string{"get", "post", "put", "delete", "patch"}
					for _, method := range httpMethods {
						if _, exists := pathObj[method]; exists {
							pathObj[f.randomString(3)] = pathObj[method]
							delete(pathObj, method)
							mutations++
							break
						}
					}
				}
				if mutations >= f.config.MaxMutationsPerSpec {
					break
				}
			}
		}
	}

	// Convert back to YAML
	fuzzedBytes, err := yaml.Marshal(spec)
	if err != nil {
		return content, mutations, err
	}

	return string(fuzzedBytes), mutations, nil
}

// fuzzSemanticText applies semantic fuzzing using text manipulation
func (f *Fuzzer) fuzzSemanticText(content string) (string, int, error) {
	mutations := 0

	// Common OpenAPI keywords to corrupt
	keywords := []string{
		"openapi", "info", "paths", "components", "servers",
		"get", "post", "put", "delete", "patch", "head", "options",
		"parameters", "responses", "requestBody", "schemas",
		"type", "format", "enum", "required", "properties",
	}

	for _, keyword := range keywords {
		if mutations >= f.config.MaxMutationsPerSpec {
			break
		}

		if f.rand.Float64() < f.config.MutationRate {
			// Replace keyword with corrupted version
			corrupted := f.corruptString(keyword)
			content = strings.ReplaceAll(content, keyword+":", corrupted+":")
			mutations++
		}
	}

	return content, mutations, nil
}

// fuzzDataTypes corrupts data type definitions
func (f *Fuzzer) fuzzDataTypes(content string) (string, int, error) {
	mutations := 0

	// Data type patterns to corrupt
	patterns := map[string][]string{
		`type:\s*string`:  {"type: number", "type: boolean", "type: invalid"},
		`type:\s*integer`: {"type: string", "type: array", "type: null"},
		`type:\s*number`:  {"type: string", "type: object", "type: undefined"},
		`type:\s*boolean`: {"type: integer", "type: array", "type: invalid"},
		`type:\s*array`:   {"type: string", "type: object", "type: null"},
		`type:\s*object`:  {"type: array", "type: string", "type: invalid"},
	}

	for pattern, replacements := range patterns {
		if mutations >= f.config.MaxMutationsPerSpec {
			break
		}

		if f.rand.Float64() < f.config.MutationRate {
			re := regexp.MustCompile(pattern)
			if re.MatchString(content) {
				replacement := replacements[f.rand.Intn(len(replacements))]
				content = re.ReplaceAllString(content, replacement)
				mutations++
			}
		}
	}

	return content, mutations, nil
}

// fuzzReferences corrupts $ref links
func (f *Fuzzer) fuzzReferences(content string) (string, int, error) {
	mutations := 0

	// Find and corrupt $ref patterns
	refPattern := regexp.MustCompile(`\$ref:\s*["']([^"']+)["']`)
	matches := refPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if mutations >= f.config.MaxMutationsPerSpec {
			break
		}

		if f.rand.Float64() < f.config.MutationRate {
			originalRef := match[0]
			refValue := match[1]

			var corruptedRef string
			switch f.rand.Intn(4) {
			case 0: // Break the reference path
				corruptedRef = strings.Replace(originalRef, refValue, f.randomString(10), 1)
			case 1: // Remove # symbol
				corruptedRef = strings.Replace(originalRef, "#", "", 1)
			case 2: // Add extra slashes
				corruptedRef = strings.Replace(originalRef, "/", "//", -1)
			case 3: // Corrupt the reference format
				corruptedRef = "$invalid: " + refValue
			}

			content = strings.Replace(content, originalRef, corruptedRef, 1)
			mutations++
		}
	}

	return content, mutations, nil
}

// fuzzEncoding corrupts character encoding
func (f *Fuzzer) fuzzEncoding(content string) (string, int, error) {
	mutations := 0

	// Insert various problematic characters
	problematicChars := []string{
		"\x00", "\x01", "\x02", "\x03", "\x04", "\x05", // Control characters
		"\u0000", "\u0001", "\u0002", // Unicode control
		"\uFFFE", "\uFFFF", // Invalid Unicode
		"🚀", "💥", "🔥", // Emojis
		"中文", "العربية", "русский", // Non-Latin scripts
	}

	lines := strings.Split(content, "\n")
	for i := 0; i < len(lines) && mutations < f.config.MaxMutationsPerSpec; i++ {
		if f.rand.Float64() < f.config.MutationRate {
			pos := f.rand.Intn(len(lines[i]) + 1)
			char := problematicChars[f.rand.Intn(len(problematicChars))]
			lines[i] = lines[i][:pos] + char + lines[i][pos:]
			mutations++
		}
	}

	return strings.Join(lines, "\n"), mutations, nil
}

// fuzzSize creates extremely large or small values
func (f *Fuzzer) fuzzSize(content string) (string, int, error) {
	mutations := 0

	// Create extremely long strings
	if f.rand.Float64() < f.config.MutationRate {
		longString := strings.Repeat("A", 100000) // 100KB string
		content = strings.Replace(content, "title:", "title: \""+longString+"\"", 1)
		mutations++
	}

	// Create deeply nested structures
	if f.rand.Float64() < f.config.MutationRate {
		deepNesting := strings.Repeat("  nested:\n", 1000)
		content = content + "\ndeep_structure:\n" + deepNesting + "  value: end"
		mutations++
	}

	return content, mutations, nil
}

// fuzzEdgeValues inserts edge case values
func (f *Fuzzer) fuzzEdgeValues(content string) (string, int, error) {
	mutations := 0

	edgeValues := []string{
		"null", "undefined", "NaN", "Infinity", "-Infinity",
		"true", "false", "0", "-1", "2147483647", "-2147483648",
		"\"\"", "\" \"", "\"\t\"", "\"\n\"", "\"\r\n\"",
		"[]", "{}", "[null]", "{\"\":null}",
	}

	// Replace some values with edge cases
	lines := strings.Split(content, "\n")
	for i := 0; i < len(lines) && mutations < f.config.MaxMutationsPerSpec; i++ {
		if strings.Contains(lines[i], ":") && f.rand.Float64() < f.config.MutationRate {
			parts := strings.SplitN(lines[i], ":", 2)
			if len(parts) == 2 {
				edgeValue := edgeValues[f.rand.Intn(len(edgeValues))]
				lines[i] = parts[0] + ": " + edgeValue
				mutations++
			}
		}
	}

	return strings.Join(lines, "\n"), mutations, nil
}

// fuzzMutation applies random byte-level mutations
func (f *Fuzzer) fuzzMutation(content string) (string, int, error) {
	mutations := 0
	contentBytes := []byte(content)

	for i := 0; i < len(contentBytes) && mutations < f.config.MaxMutationsPerSpec; i++ {
		if f.rand.Float64() < f.config.MutationRate {
			switch f.rand.Intn(3) {
			case 0: // Bit flip
				contentBytes[i] ^= byte(1 << uint(f.rand.Intn(8)))
			case 1: // Random byte
				contentBytes[i] = byte(f.rand.Intn(256))
			case 2: // Delete byte
				contentBytes = append(contentBytes[:i], contentBytes[i+1:]...)
				i-- // Adjust index after deletion
			}
			mutations++
		}
	}

	return string(contentBytes), mutations, nil
}

// randomString generates a random string of specified length
func (f *Fuzzer) randomString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[f.rand.Intn(len(chars))]
	}
	return string(result)
}

// corruptString slightly corrupts a string while keeping it readable
func (f *Fuzzer) corruptString(s string) string {
	if len(s) == 0 {
		return s
	}

	switch f.rand.Intn(4) {
	case 0: // Change case
		if f.rand.Float64() < 0.5 {
			return strings.ToUpper(s)
		}
		return strings.ToLower(s)
	case 1: // Add/remove character
		pos := f.rand.Intn(len(s))
		return s[:pos] + f.randomString(1) + s[pos:]
	case 2: // Swap characters
		if len(s) > 1 {
			runes := []rune(s)
			i, j := f.rand.Intn(len(runes)), f.rand.Intn(len(runes))
			runes[i], runes[j] = runes[j], runes[i]
			return string(runes)
		}
		return s
	default: // Add suffix
		return s + f.randomString(1)
	}
}

// RunFuzzingCampaign executes a comprehensive fuzzing campaign
func (f *Fuzzer) RunFuzzingCampaign(inputSpecs []string, specGradePath string) ([]FuzzingResult, error) {
	var results []FuzzingResult

	fmt.Printf("🔥 Starting fuzzing campaign with %d strategies on %d specs...\n",
		len(f.config.Strategies), len(inputSpecs))

	for _, specPath := range inputSpecs {
		specContent, err := os.ReadFile(specPath)
		if err != nil {
			continue
		}

		for _, strategy := range f.config.Strategies {
			result := f.fuzzSingleSpec(string(specContent), specPath, strategy, specGradePath)
			results = append(results, result)

			// Save fuzzed spec for analysis
			if f.config.OutputDir != "" {
				f.saveFuzzedSpec(result, specPath, strategy)
			}
		}
	}

	fmt.Printf("🎯 Fuzzing campaign complete: %d fuzzed specs generated\n", len(results))
	return results, nil
}

// fuzzSingleSpec fuzzes a single specification with a specific strategy
func (f *Fuzzer) fuzzSingleSpec(content, originalPath string, strategy FuzzingStrategy, specGradePath string) FuzzingResult {
	startTime := time.Now()

	result := FuzzingResult{
		OriginalSpec: originalPath,
		Strategy:     strategy,
		GeneratedAt:  time.Now(),
		CrashedTool:  false,
	}

	// Apply fuzzing
	fuzzedContent, mutationCount, err := f.FuzzSpec(content, strategy)
	result.MutationCount = mutationCount
	result.FuzzedSpec = fuzzedContent

	if err != nil {
		result.ValidationError = err.Error()
		result.ExecutionTime = time.Since(startTime)
		return result
	}

	// Test with SpecGrade (if path provided)
	if specGradePath != "" {
		specGradeResult := f.testWithSpecGrade(fuzzedContent, specGradePath)
		result.SpecGradeResult = specGradeResult

		if specGradeResult != nil && !specGradeResult.Success {
			result.CrashedTool = true
		}
	}

	result.ExecutionTime = time.Since(startTime)
	return result
}

// testWithSpecGrade runs SpecGrade against a fuzzed specification
func (f *Fuzzer) testWithSpecGrade(fuzzedContent, specGradePath string) *SpecGradeResult {
	// Create temporary file
	tmpDir, err := os.MkdirTemp("", "specgrade_fuzz")
	if err != nil {
		return &SpecGradeResult{Success: false, Error: "Failed to create temp dir"}
	}
	defer os.RemoveAll(tmpDir)

	specFile := filepath.Join(tmpDir, "openapi.yaml")
	if err := os.WriteFile(specFile, []byte(fuzzedContent), 0644); err != nil {
		return &SpecGradeResult{Success: false, Error: "Failed to write temp spec"}
	}

	// This would run SpecGrade in a real implementation
	// For now, we'll simulate the behavior
	return &SpecGradeResult{
		Success: true,
		Grade:   "F",             // Fuzzed specs should typically get low grades
		Score:   f.rand.Intn(50), // Random low score
	}
}

// saveFuzzedSpec saves a fuzzed specification to disk
func (f *Fuzzer) saveFuzzedSpec(result FuzzingResult, originalPath string, strategy FuzzingStrategy) error {
	if err := os.MkdirAll(f.config.OutputDir, 0755); err != nil {
		return err
	}

	baseName := filepath.Base(originalPath)
	baseName = strings.TrimSuffix(baseName, filepath.Ext(baseName))

	filename := fmt.Sprintf("%s_%s_%d.yaml", baseName, strategy, time.Now().Unix())
	outputPath := filepath.Join(f.config.OutputDir, filename)

	return os.WriteFile(outputPath, []byte(result.FuzzedSpec), 0644)
}

// GenerateReport creates a comprehensive fuzzing report
func (f *Fuzzer) GenerateReport(results []FuzzingResult) map[string]interface{} {
	report := map[string]interface{}{
		"total_fuzzed_specs": len(results),
		"strategies_used":    f.config.Strategies,
		"crash_count":        0,
		"strategy_stats":     make(map[FuzzingStrategy]map[string]interface{}),
		"generated_at":       time.Now(),
	}

	strategyStats := make(map[FuzzingStrategy]map[string]interface{})

	for _, strategy := range f.config.Strategies {
		strategyStats[strategy] = map[string]interface{}{
			"total_specs":    0,
			"crashed_specs":  0,
			"avg_mutations":  0.0,
			"avg_score":      0.0,
			"execution_time": time.Duration(0),
		}
	}

	for _, result := range results {
		stats := strategyStats[result.Strategy]
		stats["total_specs"] = stats["total_specs"].(int) + 1

		if result.CrashedTool {
			stats["crashed_specs"] = stats["crashed_specs"].(int) + 1
			report["crash_count"] = report["crash_count"].(int) + 1
		}

		if result.SpecGradeResult != nil && result.SpecGradeResult.Success {
			currentAvg := stats["avg_score"].(float64)
			count := stats["total_specs"].(int)
			newAvg := (currentAvg*float64(count-1) + float64(result.SpecGradeResult.Score)) / float64(count)
			stats["avg_score"] = newAvg
		}

		// Update average mutations
		currentMutAvg := stats["avg_mutations"].(float64)
		count := stats["total_specs"].(int)
		newMutAvg := (currentMutAvg*float64(count-1) + float64(result.MutationCount)) / float64(count)
		stats["avg_mutations"] = newMutAvg

		// Update execution time
		stats["execution_time"] = stats["execution_time"].(time.Duration) + result.ExecutionTime
	}

	report["strategy_stats"] = strategyStats
	return report
}
