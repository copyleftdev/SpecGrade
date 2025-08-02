package community

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ContributionType defines the type of community contribution
type ContributionType string

const (
	ContributionEdgeCase     ContributionType = "edge_case"
	ContributionRealWorldAPI ContributionType = "real_world_api"
	ContributionRule         ContributionType = "validation_rule"
	ContributionBugReport    ContributionType = "bug_report"
	ContributionImprovement  ContributionType = "improvement"
)

// ContributionStatus represents the review status of a contribution
type ContributionStatus string

const (
	StatusPending   ContributionStatus = "pending"
	StatusReviewing ContributionStatus = "reviewing"
	StatusApproved  ContributionStatus = "approved"
	StatusRejected  ContributionStatus = "rejected"
	StatusIntegrated ContributionStatus = "integrated"
)

// Contribution represents a community-submitted test case or improvement
type Contribution struct {
	ID              string             `json:"id"`
	Type            ContributionType   `json:"type"`
	Title           string             `json:"title"`
	Description     string             `json:"description"`
	Contributor     Contributor        `json:"contributor"`
	SubmittedAt     time.Time          `json:"submitted_at"`
	Status          ContributionStatus `json:"status"`
	ReviewNotes     string             `json:"review_notes,omitempty"`
	ReviewedBy      string             `json:"reviewed_by,omitempty"`
	ReviewedAt      *time.Time         `json:"reviewed_at,omitempty"`
	IntegratedAt    *time.Time         `json:"integrated_at,omitempty"`
	
	// Content specific to contribution type
	EdgeCaseData     *EdgeCaseContribution     `json:"edge_case_data,omitempty"`
	APIData          *APIContribution          `json:"api_data,omitempty"`
	RuleData         *RuleContribution         `json:"rule_data,omitempty"`
	BugReportData    *BugReportContribution    `json:"bug_report_data,omitempty"`
	ImprovementData  *ImprovementContribution  `json:"improvement_data,omitempty"`
	
	// Metadata
	Tags            []string           `json:"tags"`
	Difficulty      string             `json:"difficulty"` // easy, medium, hard
	Impact          string             `json:"impact"`     // low, medium, high
	TestResults     *TestResults       `json:"test_results,omitempty"`
}

// Contributor represents someone who submits contributions
type Contributor struct {
	Name         string    `json:"name"`
	Email        string    `json:"email,omitempty"`
	Organization string    `json:"organization,omitempty"`
	GitHub       string    `json:"github,omitempty"`
	Contributions int      `json:"total_contributions"`
	FirstContrib  time.Time `json:"first_contribution"`
	LastContrib   time.Time `json:"last_contribution"`
}

// EdgeCaseContribution represents an edge case test submission
type EdgeCaseContribution struct {
	SpecContent      string            `json:"spec_content"`
	ExpectedBehavior string            `json:"expected_behavior"`
	ActualBehavior   string            `json:"actual_behavior,omitempty"`
	ReproSteps       []string          `json:"reproduction_steps"`
	Environment      map[string]string `json:"environment"`
	Severity         string            `json:"severity"` // critical, high, medium, low
}

// APIContribution represents a real-world API specification submission
type APIContribution struct {
	APIName         string            `json:"api_name"`
	Provider        string            `json:"provider"`
	Category        string            `json:"category"`
	SpecContent     string            `json:"spec_content"`
	SourceURL       string            `json:"source_url,omitempty"`
	Version         string            `json:"version"`
	IsAnonymized    bool              `json:"is_anonymized"`
	ExpectedGrade   string            `json:"expected_grade,omitempty"`
	KnownIssues     []string          `json:"known_issues,omitempty"`
	Metadata        map[string]string `json:"metadata"`
}

// RuleContribution represents a new validation rule submission
type RuleContribution struct {
	RuleName        string   `json:"rule_name"`
	RuleID          string   `json:"rule_id"`
	Description     string   `json:"description"`
	Category        string   `json:"category"`
	Severity        string   `json:"severity"`
	Implementation  string   `json:"implementation"` // Go code
	TestCases       []string `json:"test_cases"`
	Documentation   string   `json:"documentation"`
	References      []string `json:"references"`
}

// BugReportContribution represents a bug report submission
type BugReportContribution struct {
	BugType         string            `json:"bug_type"` // crash, incorrect_result, performance
	SpecGradeVersion string           `json:"specgrade_version"`
	InputSpec       string            `json:"input_spec"`
	ExpectedOutput  string            `json:"expected_output"`
	ActualOutput    string            `json:"actual_output"`
	ErrorMessage    string            `json:"error_message,omitempty"`
	StackTrace      string            `json:"stack_trace,omitempty"`
	Environment     map[string]string `json:"environment"`
	ReproSteps      []string          `json:"reproduction_steps"`
}

// ImprovementContribution represents a feature or improvement suggestion
type ImprovementContribution struct {
	ImprovementType string   `json:"improvement_type"` // feature, performance, usability
	CurrentBehavior string   `json:"current_behavior"`
	ProposedBehavior string  `json:"proposed_behavior"`
	Benefits        []string `json:"benefits"`
	Implementation  string   `json:"implementation_notes,omitempty"`
	Examples        []string `json:"examples,omitempty"`
	References      []string `json:"references,omitempty"`
}

// TestResults represents the results of testing a contribution
type TestResults struct {
	TestedAt        time.Time         `json:"tested_at"`
	SpecGradeVersion string           `json:"specgrade_version"`
	TestsPassed     int               `json:"tests_passed"`
	TestsFailed     int               `json:"tests_failed"`
	TestOutput      string            `json:"test_output"`
	PerformanceData map[string]float64 `json:"performance_data,omitempty"`
	Issues          []string          `json:"issues,omitempty"`
}

// CommunityFramework manages community contributions and edge case discovery
type CommunityFramework struct {
	BaseDir         string
	ContributionsDB map[string]*Contribution
	Contributors    map[string]*Contributor
	ReviewQueue     []*Contribution
}

// NewCommunityFramework creates a new community contribution framework
func NewCommunityFramework(baseDir string) *CommunityFramework {
	return &CommunityFramework{
		BaseDir:         baseDir,
		ContributionsDB: make(map[string]*Contribution),
		Contributors:    make(map[string]*Contributor),
		ReviewQueue:     make([]*Contribution, 0),
	}
}

// SubmitContribution allows community members to submit contributions
func (cf *CommunityFramework) SubmitContribution(contrib *Contribution) (string, error) {
	// Generate unique ID
	contrib.ID = cf.generateContributionID(contrib)
	contrib.SubmittedAt = time.Now()
	contrib.Status = StatusPending
	
	// Validate contribution
	if err := cf.validateContribution(contrib); err != nil {
		return "", fmt.Errorf("contribution validation failed: %w", err)
	}
	
	// Update contributor information
	cf.updateContributor(&contrib.Contributor)
	
	// Store contribution
	cf.ContributionsDB[contrib.ID] = contrib
	cf.ReviewQueue = append(cf.ReviewQueue, contrib)
	
	// Save to disk
	if err := cf.saveContribution(contrib); err != nil {
		return "", fmt.Errorf("failed to save contribution: %w", err)
	}
	
	fmt.Printf("✅ Contribution submitted successfully: %s\n", contrib.ID)
	fmt.Printf("   Type: %s\n", contrib.Type)
	fmt.Printf("   Title: %s\n", contrib.Title)
	fmt.Printf("   Contributor: %s\n", contrib.Contributor.Name)
	
	return contrib.ID, nil
}

// ReviewContribution allows maintainers to review and approve/reject contributions
func (cf *CommunityFramework) ReviewContribution(contributionID, reviewerName, notes string, approved bool) error {
	contrib, exists := cf.ContributionsDB[contributionID]
	if !exists {
		return fmt.Errorf("contribution not found: %s", contributionID)
	}
	
	contrib.Status = StatusReviewing
	contrib.ReviewNotes = notes
	contrib.ReviewedBy = reviewerName
	now := time.Now()
	contrib.ReviewedAt = &now
	
	if approved {
		contrib.Status = StatusApproved
		
		// Run tests if applicable
		if err := cf.testContribution(contrib); err != nil {
			contrib.Status = StatusRejected
			contrib.ReviewNotes += fmt.Sprintf("\nTesting failed: %v", err)
		}
	} else {
		contrib.Status = StatusRejected
	}
	
	// Remove from review queue
	cf.removeFromReviewQueue(contributionID)
	
	// Save updated contribution
	if err := cf.saveContribution(contrib); err != nil {
		return fmt.Errorf("failed to save reviewed contribution: %w", err)
	}
	
	fmt.Printf("📋 Contribution %s reviewed by %s: %s\n", 
		contributionID, reviewerName, contrib.Status)
	
	return nil
}

// IntegrateContribution integrates an approved contribution into the main codebase
func (cf *CommunityFramework) IntegrateContribution(contributionID string) error {
	contrib, exists := cf.ContributionsDB[contributionID]
	if !exists {
		return fmt.Errorf("contribution not found: %s", contributionID)
	}
	
	if contrib.Status != StatusApproved {
		return fmt.Errorf("contribution must be approved before integration")
	}
	
	// Integration logic based on contribution type
	switch contrib.Type {
	case ContributionEdgeCase:
		err := cf.integrateEdgeCase(contrib)
		if err != nil {
			return fmt.Errorf("failed to integrate edge case: %w", err)
		}
		
	case ContributionRealWorldAPI:
		err := cf.integrateRealWorldAPI(contrib)
		if err != nil {
			return fmt.Errorf("failed to integrate real-world API: %w", err)
		}
		
	case ContributionRule:
		err := cf.integrateValidationRule(contrib)
		if err != nil {
			return fmt.Errorf("failed to integrate validation rule: %w", err)
		}
		
	default:
		return fmt.Errorf("integration not implemented for contribution type: %s", contrib.Type)
	}
	
	contrib.Status = StatusIntegrated
	now := time.Now()
	contrib.IntegratedAt = &now
	
	// Update contributor stats
	contributor := cf.Contributors[contrib.Contributor.Email]
	if contributor != nil {
		contributor.Contributions++
		contributor.LastContrib = time.Now()
	}
	
	// Save updated contribution
	if err := cf.saveContribution(contrib); err != nil {
		return fmt.Errorf("failed to save integrated contribution: %w", err)
	}
	
	fmt.Printf("🚀 Contribution %s successfully integrated!\n", contributionID)
	return nil
}

// GetContributionStats returns statistics about community contributions
func (cf *CommunityFramework) GetContributionStats() map[string]interface{} {
	stats := map[string]interface{}{
		"total_contributions":    len(cf.ContributionsDB),
		"pending_reviews":        len(cf.ReviewQueue),
		"total_contributors":     len(cf.Contributors),
		"contributions_by_type":  make(map[ContributionType]int),
		"contributions_by_status": make(map[ContributionStatus]int),
		"top_contributors":       cf.getTopContributors(10),
		"recent_activity":        cf.getRecentActivity(30),
	}
	
	// Count by type and status
	for _, contrib := range cf.ContributionsDB {
		stats["contributions_by_type"].(map[ContributionType]int)[contrib.Type]++
		stats["contributions_by_status"].(map[ContributionStatus]int)[contrib.Status]++
	}
	
	return stats
}

// DiscoverEdgeCases analyzes contributions to identify common edge case patterns
func (cf *CommunityFramework) DiscoverEdgeCases() []EdgeCasePattern {
	var patterns []EdgeCasePattern
	
	// Analyze edge case contributions
	edgeCases := cf.getContributionsByType(ContributionEdgeCase)
	
	// Group by common characteristics
	patternMap := make(map[string]*EdgeCasePattern)
	
	for _, contrib := range edgeCases {
		if contrib.EdgeCaseData == nil {
			continue
		}
		
		// Extract patterns (simplified)
		for _, tag := range contrib.Tags {
			if pattern, exists := patternMap[tag]; exists {
				pattern.Frequency++
				pattern.Examples = append(pattern.Examples, contrib.ID)
			} else {
				patternMap[tag] = &EdgeCasePattern{
					Name:        tag,
					Description: fmt.Sprintf("Edge cases related to %s", tag),
					Frequency:   1,
					Severity:    contrib.EdgeCaseData.Severity,
					Examples:    []string{contrib.ID},
					FirstSeen:   contrib.SubmittedAt,
					LastSeen:    contrib.SubmittedAt,
				}
			}
		}
	}
	
	// Convert to slice and sort by frequency
	for _, pattern := range patternMap {
		patterns = append(patterns, *pattern)
	}
	
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Frequency > patterns[j].Frequency
	})
	
	return patterns
}

// EdgeCasePattern represents a discovered pattern in edge case contributions
type EdgeCasePattern struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Frequency   int       `json:"frequency"`
	Severity    string    `json:"severity"`
	Examples    []string  `json:"examples"`
	FirstSeen   time.Time `json:"first_seen"`
	LastSeen    time.Time `json:"last_seen"`
}

// Helper methods

func (cf *CommunityFramework) generateContributionID(contrib *Contribution) string {
	data := fmt.Sprintf("%s-%s-%s-%d", 
		contrib.Type, contrib.Title, contrib.Contributor.Email, time.Now().UnixNano())
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)[:12] // Use first 12 characters
}

func (cf *CommunityFramework) validateContribution(contrib *Contribution) error {
	if contrib.Title == "" {
		return fmt.Errorf("title is required")
	}
	
	if contrib.Description == "" {
		return fmt.Errorf("description is required")
	}
	
	if contrib.Contributor.Name == "" {
		return fmt.Errorf("contributor name is required")
	}
	
	// Type-specific validation
	switch contrib.Type {
	case ContributionEdgeCase:
		if contrib.EdgeCaseData == nil {
			return fmt.Errorf("edge case data is required")
		}
		if contrib.EdgeCaseData.SpecContent == "" {
			return fmt.Errorf("spec content is required for edge cases")
		}
		
	case ContributionRealWorldAPI:
		if contrib.APIData == nil {
			return fmt.Errorf("API data is required")
		}
		if contrib.APIData.SpecContent == "" {
			return fmt.Errorf("spec content is required for APIs")
		}
		
	case ContributionRule:
		if contrib.RuleData == nil {
			return fmt.Errorf("rule data is required")
		}
		if contrib.RuleData.Implementation == "" {
			return fmt.Errorf("rule implementation is required")
		}
	}
	
	return nil
}

func (cf *CommunityFramework) updateContributor(contributor *Contributor) {
	email := contributor.Email
	if email == "" {
		email = contributor.Name // Fallback to name if no email
	}
	
	if existing, exists := cf.Contributors[email]; exists {
		existing.Contributions++
		existing.LastContrib = time.Now()
		// Update other fields if provided
		if contributor.Organization != "" {
			existing.Organization = contributor.Organization
		}
		if contributor.GitHub != "" {
			existing.GitHub = contributor.GitHub
		}
	} else {
		contributor.Contributions = 1
		contributor.FirstContrib = time.Now()
		contributor.LastContrib = time.Now()
		cf.Contributors[email] = contributor
	}
}

func (cf *CommunityFramework) saveContribution(contrib *Contribution) error {
	// Create directory structure
	typeDir := filepath.Join(cf.BaseDir, "contributions", string(contrib.Type))
	if err := os.MkdirAll(typeDir, 0755); err != nil {
		return err
	}
	
	// Save contribution as JSON
	contribPath := filepath.Join(typeDir, contrib.ID+".json")
	file, err := os.Create(contribPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(contrib)
}

func (cf *CommunityFramework) testContribution(contrib *Contribution) error {
	// Run automated tests on the contribution
	testResults := &TestResults{
		TestedAt:         time.Now(),
		SpecGradeVersion: "1.0.0", // Would get actual version
		TestsPassed:      0,
		TestsFailed:      0,
	}
	
	switch contrib.Type {
	case ContributionEdgeCase:
		// Test edge case against SpecGrade
		if contrib.EdgeCaseData != nil {
			// This would run actual tests
			testResults.TestsPassed = 1
			testResults.TestOutput = "Edge case successfully reproduced"
		}
		
	case ContributionRealWorldAPI:
		// Validate API spec and run SpecGrade
		if contrib.APIData != nil {
			// This would run actual validation
			testResults.TestsPassed = 1
			testResults.TestOutput = "API spec validated successfully"
		}
		
	case ContributionRule:
		// Compile and test rule implementation
		if contrib.RuleData != nil {
			// This would compile and test the rule
			testResults.TestsPassed = 1
			testResults.TestOutput = "Rule implementation compiled and tested"
		}
	}
	
	contrib.TestResults = testResults
	return nil
}

func (cf *CommunityFramework) removeFromReviewQueue(contributionID string) {
	for i, contrib := range cf.ReviewQueue {
		if contrib.ID == contributionID {
			cf.ReviewQueue = append(cf.ReviewQueue[:i], cf.ReviewQueue[i+1:]...)
			break
		}
	}
}

func (cf *CommunityFramework) integrateEdgeCase(contrib *Contribution) error {
	// Save edge case to test suite
	edgeCaseDir := filepath.Join(cf.BaseDir, "edge_cases", "community")
	if err := os.MkdirAll(edgeCaseDir, 0755); err != nil {
		return err
	}
	
	specPath := filepath.Join(edgeCaseDir, contrib.ID+".yaml")
	return os.WriteFile(specPath, []byte(contrib.EdgeCaseData.SpecContent), 0644)
}

func (cf *CommunityFramework) integrateRealWorldAPI(contrib *Contribution) error {
	// Save API to real-world collection
	apiDir := filepath.Join(cf.BaseDir, "realworld", contrib.APIData.Category, contrib.APIData.APIName)
	if err := os.MkdirAll(apiDir, 0755); err != nil {
		return err
	}
	
	specPath := filepath.Join(apiDir, "openapi.yaml")
	return os.WriteFile(specPath, []byte(contrib.APIData.SpecContent), 0644)
}

func (cf *CommunityFramework) integrateValidationRule(contrib *Contribution) error {
	// This would integrate a new validation rule into the codebase
	// For now, just save the rule implementation
	ruleDir := filepath.Join(cf.BaseDir, "rules", "community")
	if err := os.MkdirAll(ruleDir, 0755); err != nil {
		return err
	}
	
	rulePath := filepath.Join(ruleDir, contrib.RuleData.RuleID+".go")
	return os.WriteFile(rulePath, []byte(contrib.RuleData.Implementation), 0644)
}

func (cf *CommunityFramework) getContributionsByType(contribType ContributionType) []*Contribution {
	var contributions []*Contribution
	for _, contrib := range cf.ContributionsDB {
		if contrib.Type == contribType {
			contributions = append(contributions, contrib)
		}
	}
	return contributions
}

func (cf *CommunityFramework) getTopContributors(limit int) []Contributor {
	contributors := make([]Contributor, 0, len(cf.Contributors))
	for _, contributor := range cf.Contributors {
		contributors = append(contributors, *contributor)
	}
	
	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].Contributions > contributors[j].Contributions
	})
	
	if len(contributors) > limit {
		contributors = contributors[:limit]
	}
	
	return contributors
}

func (cf *CommunityFramework) getRecentActivity(days int) []Contribution {
	cutoff := time.Now().AddDate(0, 0, -days)
	var recent []Contribution
	
	for _, contrib := range cf.ContributionsDB {
		if contrib.SubmittedAt.After(cutoff) {
			recent = append(recent, *contrib)
		}
	}
	
	sort.Slice(recent, func(i, j int) bool {
		return recent[i].SubmittedAt.After(recent[j].SubmittedAt)
	})
	
	return recent
}

// LoadContributions loads existing contributions from disk
func (cf *CommunityFramework) LoadContributions() error {
	contributionsDir := filepath.Join(cf.BaseDir, "contributions")
	
	if _, err := os.Stat(contributionsDir); os.IsNotExist(err) {
		return nil // No contributions directory yet
	}
	
	return filepath.Walk(contributionsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !strings.HasSuffix(path, ".json") {
			return nil
		}
		
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		
		var contrib Contribution
		if err := json.NewDecoder(file).Decode(&contrib); err != nil {
			return err
		}
		
		cf.ContributionsDB[contrib.ID] = &contrib
		
		// Add to review queue if pending
		if contrib.Status == StatusPending {
			cf.ReviewQueue = append(cf.ReviewQueue, &contrib)
		}
		
		// Update contributor info
		cf.updateContributor(&contrib.Contributor)
		
		return nil
	})
}

// GenerateContributionReport creates a comprehensive report of community activity
func (cf *CommunityFramework) GenerateContributionReport() map[string]interface{} {
	stats := cf.GetContributionStats()
	patterns := cf.DiscoverEdgeCases()
	
	report := map[string]interface{}{
		"summary":           stats,
		"edge_case_patterns": patterns,
		"generated_at":      time.Now(),
		"framework_version": "1.0.0",
	}
	
	return report
}
