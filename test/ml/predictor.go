package ml

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

// QualityFeatures represents extracted features from an OpenAPI specification
type QualityFeatures struct {
	// Structural Features
	TotalPaths      int     `json:"total_paths"`
	TotalOperations int     `json:"total_operations"`
	TotalSchemas    int     `json:"total_schemas"`
	TotalParameters int     `json:"total_parameters"`
	MaxNestingDepth int     `json:"max_nesting_depth"`
	AvgPathLength   float64 `json:"avg_path_length"`

	// Documentation Features
	DescriptionCoverage float64 `json:"description_coverage"`
	ExampleCoverage     float64 `json:"example_coverage"`
	SummaryPresent      bool    `json:"summary_present"`
	TagsPresent         bool    `json:"tags_present"`
	ContactInfoPresent  bool    `json:"contact_info_present"`
	LicensePresent      bool    `json:"license_present"`

	// Schema Quality Features
	TypeConsistency     float64 `json:"type_consistency"`
	RequiredFieldRatio  float64 `json:"required_field_ratio"`
	EnumUsage           float64 `json:"enum_usage"`
	FormatSpecification float64 `json:"format_specification"`

	// API Design Features
	RESTfulnessScore    float64 `json:"restfulness_score"`
	HTTPMethodDiversity float64 `json:"http_method_diversity"`
	StatusCodeCoverage  float64 `json:"status_code_coverage"`
	ErrorResponseRatio  float64 `json:"error_response_ratio"`

	// Security Features
	SecuritySchemeCount int     `json:"security_scheme_count"`
	AuthenticationScore float64 `json:"authentication_score"`
	ScopeDefinitions    int     `json:"scope_definitions"`

	// Complexity Features
	CyclomaticComplexity float64 `json:"cyclomatic_complexity"`
	ReferenceComplexity  float64 `json:"reference_complexity"`
	PolymorphismUsage    float64 `json:"polymorphism_usage"`

	// Consistency Features
	NamingConsistency   float64 `json:"naming_consistency"`
	PatternConsistency  float64 `json:"pattern_consistency"`
	ResponseConsistency float64 `json:"response_consistency"`
}

// QualityPrediction represents a predicted quality assessment
type QualityPrediction struct {
	PredictedGrade     string             `json:"predicted_grade"`
	PredictedScore     int                `json:"predicted_score"`
	Confidence         float64            `json:"confidence"`
	FeatureImportance  map[string]float64 `json:"feature_importance"`
	QualityInsights    []string           `json:"quality_insights"`
	RecommendedActions []string           `json:"recommended_actions"`
	PredictionTime     time.Duration      `json:"prediction_time"`
	ModelVersion       string             `json:"model_version"`
}

// TrainingData represents a labeled example for model training
type TrainingData struct {
	Features    QualityFeatures `json:"features"`
	ActualGrade string          `json:"actual_grade"`
	ActualScore int             `json:"actual_score"`
	APIName     string          `json:"api_name"`
	Category    string          `json:"category"`
}

// QualityPredictor implements ML-based quality prediction
type QualityPredictor struct {
	model        *SimpleMLModel
	featureStats *FeatureStatistics
	modelVersion string
}

// SimpleMLModel implements a basic machine learning model for quality prediction
type SimpleMLModel struct {
	Weights      map[string]float64 `json:"weights"`
	Bias         float64            `json:"bias"`
	FeatureNames []string           `json:"feature_names"`
	TrainedAt    time.Time          `json:"trained_at"`
	TrainingSize int                `json:"training_size"`
}

// FeatureStatistics stores normalization parameters
type FeatureStatistics struct {
	Means   map[string]float64 `json:"means"`
	StdDevs map[string]float64 `json:"std_devs"`
	Mins    map[string]float64 `json:"mins"`
	Maxs    map[string]float64 `json:"maxs"`
}

// NewQualityPredictor creates a new ML-based quality predictor
func NewQualityPredictor() *QualityPredictor {
	return &QualityPredictor{
		model: NewSimpleMLModel(),
		featureStats: &FeatureStatistics{
			Means:   make(map[string]float64),
			StdDevs: make(map[string]float64),
			Mins:    make(map[string]float64),
			Maxs:    make(map[string]float64),
		},
		modelVersion: "1.0.0",
	}
}

// NewSimpleMLModel creates a new simple ML model with default weights
func NewSimpleMLModel() *SimpleMLModel {
	// Initialize with reasonable default weights based on domain knowledge
	weights := map[string]float64{
		"description_coverage":  0.25, // High importance for documentation
		"example_coverage":      0.15,
		"type_consistency":      0.20, // Schema quality is crucial
		"status_code_coverage":  0.15, // Good error handling
		"security_scheme_count": 0.10, // Security considerations
		"restfulness_score":     0.10, // API design quality
		"naming_consistency":    0.05, // Consistency matters
	}

	featureNames := make([]string, 0, len(weights))
	for name := range weights {
		featureNames = append(featureNames, name)
	}
	sort.Strings(featureNames)

	return &SimpleMLModel{
		Weights:      weights,
		Bias:         50.0, // Start with neutral bias
		FeatureNames: featureNames,
		TrainedAt:    time.Now(),
		TrainingSize: 0,
	}
}

// ExtractFeatures analyzes an OpenAPI specification and extracts quality features
func (p *QualityPredictor) ExtractFeatures(specContent string) (*QualityFeatures, error) {
	features := &QualityFeatures{}

	// Parse the specification (simplified - would use proper OpenAPI parser in production)
	lines := strings.Split(specContent, "\n")

	// Extract structural features
	features.TotalPaths = p.countMatches(specContent, `^\s*/[^:]*:`, true)
	features.TotalOperations = p.countHTTPMethods(specContent)
	features.TotalSchemas = p.countMatches(specContent, `\s+schemas:`, false)
	features.TotalParameters = p.countMatches(specContent, `parameters:`, false)
	features.MaxNestingDepth = p.calculateMaxNestingDepth(lines)
	features.AvgPathLength = p.calculateAvgPathLength(specContent)

	// Extract documentation features
	features.DescriptionCoverage = p.calculateDescriptionCoverage(specContent)
	features.ExampleCoverage = p.calculateExampleCoverage(specContent)
	features.SummaryPresent = strings.Contains(specContent, "summary:")
	features.TagsPresent = strings.Contains(specContent, "tags:")
	features.ContactInfoPresent = strings.Contains(specContent, "contact:")
	features.LicensePresent = strings.Contains(specContent, "license:")

	// Extract schema quality features
	features.TypeConsistency = p.calculateTypeConsistency(specContent)
	features.RequiredFieldRatio = p.calculateRequiredFieldRatio(specContent)
	features.EnumUsage = p.calculateEnumUsage(specContent)
	features.FormatSpecification = p.calculateFormatSpecification(specContent)

	// Extract API design features
	features.RESTfulnessScore = p.calculateRESTfulnessScore(specContent)
	features.HTTPMethodDiversity = p.calculateHTTPMethodDiversity(specContent)
	features.StatusCodeCoverage = p.calculateStatusCodeCoverage(specContent)
	features.ErrorResponseRatio = p.calculateErrorResponseRatio(specContent)

	// Extract security features
	features.SecuritySchemeCount = p.countMatches(specContent, `securitySchemes:`, false)
	features.AuthenticationScore = p.calculateAuthenticationScore(specContent)
	features.ScopeDefinitions = p.countMatches(specContent, `scopes:`, false)

	// Extract complexity features
	features.CyclomaticComplexity = p.calculateCyclomaticComplexity(specContent)
	features.ReferenceComplexity = p.calculateReferenceComplexity(specContent)
	features.PolymorphismUsage = p.calculatePolymorphismUsage(specContent)

	// Extract consistency features
	features.NamingConsistency = p.calculateNamingConsistency(specContent)
	features.PatternConsistency = p.calculatePatternConsistency(specContent)
	features.ResponseConsistency = p.calculateResponseConsistency(specContent)

	return features, nil
}

// PredictQuality uses the trained model to predict API quality
func (p *QualityPredictor) PredictQuality(features *QualityFeatures) (*QualityPrediction, error) {
	startTime := time.Now()

	// Normalize features
	normalizedFeatures := p.normalizeFeatures(features)

	// Calculate prediction using simple linear model
	score := p.model.Bias
	featureImportance := make(map[string]float64)

	for _, featureName := range p.model.FeatureNames {
		if weight, exists := p.model.Weights[featureName]; exists {
			featureValue := p.getFeatureValue(normalizedFeatures, featureName)
			contribution := weight * featureValue
			score += contribution
			featureImportance[featureName] = math.Abs(contribution)
		}
	}

	// Clamp score to valid range
	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}

	// Convert score to grade
	grade := p.scoreToGrade(int(score))

	// Calculate confidence (simplified)
	confidence := p.calculateConfidence(features, normalizedFeatures)

	// Generate insights and recommendations
	insights := p.generateQualityInsights(features, featureImportance)
	recommendations := p.generateRecommendations(features, featureImportance)

	prediction := &QualityPrediction{
		PredictedGrade:     grade,
		PredictedScore:     int(score),
		Confidence:         confidence,
		FeatureImportance:  featureImportance,
		QualityInsights:    insights,
		RecommendedActions: recommendations,
		PredictionTime:     time.Since(startTime),
		ModelVersion:       p.modelVersion,
	}

	return prediction, nil
}

// TrainModel trains the ML model using labeled training data
func (p *QualityPredictor) TrainModel(trainingData []TrainingData) error {
	if len(trainingData) == 0 {
		return fmt.Errorf("no training data provided")
	}

	fmt.Printf("🤖 Training ML model with %d examples...\n", len(trainingData))

	// Calculate feature statistics for normalization
	p.calculateFeatureStatistics(trainingData)

	// Simple gradient descent training (simplified implementation)
	learningRate := 0.01
	epochs := 100

	for epoch := 0; epoch < epochs; epoch++ {
		totalError := 0.0

		for _, example := range trainingData {
			// Normalize features
			normalizedFeatures := p.normalizeFeatures(&example.Features)

			// Forward pass
			predicted := p.model.Bias
			for _, featureName := range p.model.FeatureNames {
				if weight, exists := p.model.Weights[featureName]; exists {
					featureValue := p.getFeatureValue(normalizedFeatures, featureName)
					predicted += weight * featureValue
				}
			}

			// Calculate error
			actual := float64(example.ActualScore)
			error := predicted - actual
			totalError += error * error

			// Backward pass (update weights)
			p.model.Bias -= learningRate * error

			for _, featureName := range p.model.FeatureNames {
				if _, exists := p.model.Weights[featureName]; exists {
					featureValue := p.getFeatureValue(normalizedFeatures, featureName)
					p.model.Weights[featureName] -= learningRate * error * featureValue
				}
			}
		}

		// Print progress occasionally
		if epoch%20 == 0 {
			avgError := math.Sqrt(totalError / float64(len(trainingData)))
			fmt.Printf("Epoch %d: Average Error = %.2f\n", epoch, avgError)
		}
	}

	p.model.TrainedAt = time.Now()
	p.model.TrainingSize = len(trainingData)

	fmt.Println("✅ Model training complete!")
	return nil
}

// Helper methods for feature extraction

func (p *QualityPredictor) countMatches(content, pattern string, multiline bool) int {
	var re *regexp.Regexp
	if multiline {
		re = regexp.MustCompile("(?m)" + pattern)
	} else {
		re = regexp.MustCompile(pattern)
	}
	return len(re.FindAllString(content, -1))
}

func (p *QualityPredictor) countHTTPMethods(content string) int {
	methods := []string{"get:", "post:", "put:", "delete:", "patch:", "head:", "options:"}
	total := 0
	for _, method := range methods {
		total += strings.Count(strings.ToLower(content), method)
	}
	return total
}

func (p *QualityPredictor) calculateMaxNestingDepth(lines []string) int {
	maxDepth := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		depth := (len(line) - len(strings.TrimLeft(line, " "))) / 2
		if depth > maxDepth {
			maxDepth = depth
		}
	}
	return maxDepth
}

func (p *QualityPredictor) calculateAvgPathLength(content string) float64 {
	pathPattern := regexp.MustCompile(`^\s*(/[^:]*):`)
	matches := pathPattern.FindAllStringSubmatch(content, -1)

	if len(matches) == 0 {
		return 0
	}

	totalLength := 0
	for _, match := range matches {
		totalLength += len(match[1])
	}

	return float64(totalLength) / float64(len(matches))
}

func (p *QualityPredictor) calculateDescriptionCoverage(content string) float64 {
	totalItems := p.countHTTPMethods(content) + p.countMatches(content, `parameters:`, false)
	if totalItems == 0 {
		return 0
	}

	descriptions := p.countMatches(content, `description:`, false)
	return math.Min(float64(descriptions)/float64(totalItems), 1.0)
}

func (p *QualityPredictor) calculateExampleCoverage(content string) float64 {
	schemas := p.countMatches(content, `type:`, false)
	if schemas == 0 {
		return 0
	}

	examples := p.countMatches(content, `example:`, false)
	return math.Min(float64(examples)/float64(schemas), 1.0)
}

func (p *QualityPredictor) calculateTypeConsistency(content string) float64 {
	// Simplified: check for common type inconsistencies
	typePattern := regexp.MustCompile(`type:\s*(\w+)`)
	matches := typePattern.FindAllStringSubmatch(content, -1)

	if len(matches) == 0 {
		return 1.0
	}

	validTypes := map[string]bool{
		"string": true, "number": true, "integer": true,
		"boolean": true, "array": true, "object": true,
	}

	validCount := 0
	for _, match := range matches {
		if validTypes[match[1]] {
			validCount++
		}
	}

	return float64(validCount) / float64(len(matches))
}

func (p *QualityPredictor) calculateRequiredFieldRatio(content string) float64 {
	properties := p.countMatches(content, `properties:`, false)
	if properties == 0 {
		return 0
	}

	required := p.countMatches(content, `required:`, false)
	return math.Min(float64(required)/float64(properties), 1.0)
}

func (p *QualityPredictor) calculateEnumUsage(content string) float64 {
	strings := p.countMatches(content, `type:\s*string`, false)
	if strings == 0 {
		return 0
	}

	enums := p.countMatches(content, `enum:`, false)
	return math.Min(float64(enums)/float64(strings), 1.0)
}

func (p *QualityPredictor) calculateFormatSpecification(content string) float64 {
	strings := p.countMatches(content, `type:\s*string`, false)
	if strings == 0 {
		return 0
	}

	formats := p.countMatches(content, `format:`, false)
	return math.Min(float64(formats)/float64(strings), 1.0)
}

func (p *QualityPredictor) calculateRESTfulnessScore(content string) float64 {
	// Simplified RESTfulness scoring
	score := 0.0

	// Check for proper HTTP method usage
	if strings.Contains(content, "get:") {
		score += 0.2
	}
	if strings.Contains(content, "post:") {
		score += 0.2
	}
	if strings.Contains(content, "put:") {
		score += 0.2
	}
	if strings.Contains(content, "delete:") {
		score += 0.2
	}

	// Check for resource-based paths
	resourcePaths := p.countMatches(content, `^\s*/\w+(/\{\w+\})?:`, true)
	totalPaths := p.countMatches(content, `^\s*/[^:]*:`, true)

	if totalPaths > 0 {
		score += 0.2 * (float64(resourcePaths) / float64(totalPaths))
	}

	return score
}

func (p *QualityPredictor) calculateHTTPMethodDiversity(content string) float64 {
	methods := []string{"get:", "post:", "put:", "delete:", "patch:", "head:", "options:"}
	presentMethods := 0

	for _, method := range methods {
		if strings.Contains(strings.ToLower(content), method) {
			presentMethods++
		}
	}

	return float64(presentMethods) / float64(len(methods))
}

func (p *QualityPredictor) calculateStatusCodeCoverage(content string) float64 {
	standardCodes := []string{"200", "201", "400", "401", "403", "404", "500"}
	presentCodes := 0

	for _, code := range standardCodes {
		if strings.Contains(content, "'"+code+"'") || strings.Contains(content, "\""+code+"\"") {
			presentCodes++
		}
	}

	return float64(presentCodes) / float64(len(standardCodes))
}

func (p *QualityPredictor) calculateErrorResponseRatio(content string) float64 {
	totalResponses := p.countMatches(content, `responses:`, false)
	if totalResponses == 0 {
		return 0
	}

	errorResponses := p.countMatches(content, `'[45]\d\d':`, false) +
		p.countMatches(content, `"[45]\d\d":`, false)

	return math.Min(float64(errorResponses)/float64(totalResponses), 1.0)
}

func (p *QualityPredictor) calculateAuthenticationScore(content string) float64 {
	score := 0.0

	if strings.Contains(content, "securitySchemes:") {
		score += 0.5
	}
	if strings.Contains(content, "security:") {
		score += 0.3
	}
	if strings.Contains(content, "oauth2") || strings.Contains(content, "bearer") {
		score += 0.2
	}

	return math.Min(score, 1.0)
}

func (p *QualityPredictor) calculateCyclomaticComplexity(content string) float64 {
	// Simplified complexity based on conditional structures
	complexity := float64(p.countMatches(content, `oneOf:`, false))
	complexity += float64(p.countMatches(content, `anyOf:`, false))
	complexity += float64(p.countMatches(content, `allOf:`, false))

	totalOperations := float64(p.countHTTPMethods(content))
	if totalOperations == 0 {
		return 0
	}

	return complexity / totalOperations
}

func (p *QualityPredictor) calculateReferenceComplexity(content string) float64 {
	references := float64(p.countMatches(content, `\$ref:`, false))
	totalItems := float64(p.countHTTPMethods(content) + p.countMatches(content, `schemas:`, false))

	if totalItems == 0 {
		return 0
	}

	return references / totalItems
}

func (p *QualityPredictor) calculatePolymorphismUsage(content string) float64 {
	polymorphism := float64(p.countMatches(content, `oneOf:`, false))
	polymorphism += float64(p.countMatches(content, `anyOf:`, false))
	polymorphism += float64(p.countMatches(content, `discriminator:`, false))

	schemas := float64(p.countMatches(content, `schemas:`, false))
	if schemas == 0 {
		return 0
	}

	return math.Min(polymorphism/schemas, 1.0)
}

func (p *QualityPredictor) calculateNamingConsistency(content string) float64 {
	// Simplified: check for consistent naming patterns
	// This would be more sophisticated in a real implementation
	return 0.8 // Placeholder
}

func (p *QualityPredictor) calculatePatternConsistency(content string) float64 {
	// Simplified: check for consistent patterns
	return 0.7 // Placeholder
}

func (p *QualityPredictor) calculateResponseConsistency(content string) float64 {
	// Simplified: check for consistent response structures
	return 0.75 // Placeholder
}

// Helper methods for model operations

func (p *QualityPredictor) normalizeFeatures(features *QualityFeatures) *QualityFeatures {
	// Create a copy and normalize
	normalized := *features

	// Apply z-score normalization where statistics are available
	// This is a simplified version - would be more comprehensive in production

	return &normalized
}

func (p *QualityPredictor) getFeatureValue(features *QualityFeatures, featureName string) float64 {
	// Use reflection or a map to get feature values
	// Simplified implementation
	switch featureName {
	case "description_coverage":
		return features.DescriptionCoverage
	case "example_coverage":
		return features.ExampleCoverage
	case "type_consistency":
		return features.TypeConsistency
	case "status_code_coverage":
		return features.StatusCodeCoverage
	case "security_scheme_count":
		return float64(features.SecuritySchemeCount)
	case "restfulness_score":
		return features.RESTfulnessScore
	case "naming_consistency":
		return features.NamingConsistency
	default:
		return 0.0
	}
}

func (p *QualityPredictor) scoreToGrade(score int) string {
	switch {
	case score >= 95:
		return "A+"
	case score >= 90:
		return "A"
	case score >= 85:
		return "A-"
	case score >= 80:
		return "B+"
	case score >= 75:
		return "B"
	case score >= 70:
		return "B-"
	case score >= 65:
		return "C+"
	case score >= 60:
		return "C"
	case score >= 55:
		return "C-"
	case score >= 50:
		return "D"
	default:
		return "F"
	}
}

func (p *QualityPredictor) calculateConfidence(original, normalized *QualityFeatures) float64 {
	// Simplified confidence calculation
	// In practice, this would consider model uncertainty, feature reliability, etc.
	return 0.85 // Placeholder
}

func (p *QualityPredictor) generateQualityInsights(features *QualityFeatures, importance map[string]float64) []string {
	var insights []string

	if features.DescriptionCoverage < 0.5 {
		insights = append(insights, "Low documentation coverage detected - many operations lack descriptions")
	}

	if features.StatusCodeCoverage < 0.4 {
		insights = append(insights, "Limited HTTP status code coverage - missing error response definitions")
	}

	if features.SecuritySchemeCount == 0 {
		insights = append(insights, "No security schemes defined - API may lack proper authentication")
	}

	if features.TypeConsistency < 0.8 {
		insights = append(insights, "Type consistency issues detected - some invalid or inconsistent type definitions")
	}

	if len(insights) == 0 {
		insights = append(insights, "Overall API structure appears well-designed")
	}

	return insights
}

func (p *QualityPredictor) generateRecommendations(features *QualityFeatures, importance map[string]float64) []string {
	var recommendations []string

	if features.DescriptionCoverage < 0.7 {
		recommendations = append(recommendations, "Add descriptions to operations and parameters for better documentation")
	}

	if features.ExampleCoverage < 0.5 {
		recommendations = append(recommendations, "Include examples in schema definitions to improve API usability")
	}

	if features.ErrorResponseRatio < 0.3 {
		recommendations = append(recommendations, "Define error responses (4xx, 5xx) for better error handling")
	}

	if features.SecuritySchemeCount == 0 {
		recommendations = append(recommendations, "Implement security schemes (OAuth2, API keys, etc.)")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Consider adding more detailed examples and edge case documentation")
	}

	return recommendations
}

func (p *QualityPredictor) calculateFeatureStatistics(trainingData []TrainingData) {
	// Calculate means, standard deviations, mins, maxs for normalization
	// Simplified implementation

	featureValues := make(map[string][]float64)

	for _, example := range trainingData {
		featureValues["description_coverage"] = append(featureValues["description_coverage"], example.Features.DescriptionCoverage)
		featureValues["example_coverage"] = append(featureValues["example_coverage"], example.Features.ExampleCoverage)
		// ... add other features
	}

	for featureName, values := range featureValues {
		if len(values) == 0 {
			continue
		}

		// Calculate mean
		sum := 0.0
		for _, v := range values {
			sum += v
		}
		mean := sum / float64(len(values))
		p.featureStats.Means[featureName] = mean

		// Calculate standard deviation
		sumSquares := 0.0
		for _, v := range values {
			diff := v - mean
			sumSquares += diff * diff
		}
		stdDev := math.Sqrt(sumSquares / float64(len(values)))
		p.featureStats.StdDevs[featureName] = stdDev

		// Find min and max
		min, max := values[0], values[0]
		for _, v := range values {
			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
		}
		p.featureStats.Mins[featureName] = min
		p.featureStats.Maxs[featureName] = max
	}
}

// SaveModel saves the trained model to disk
func (p *QualityPredictor) SaveModel(filename string) error {
	data := map[string]interface{}{
		"model":         p.model,
		"feature_stats": p.featureStats,
		"model_version": p.modelVersion,
		"saved_at":      time.Now(),
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// LoadModel loads a trained model from disk
func (p *QualityPredictor) LoadModel(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return err
	}

	// This would need proper deserialization in a real implementation
	fmt.Println("Model loading not fully implemented - would deserialize from JSON")

	return nil
}
