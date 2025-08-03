package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// APIMetadata contains information about a collected API specification
type APIMetadata struct {
	Name            string    `json:"name"`
	Provider        string    `json:"provider"`
	Category        string    `json:"category"`
	Version         string    `json:"version"`
	SourceURL       string    `json:"source_url"`
	CollectedAt     time.Time `json:"collected_at"`
	ExpectedGrade   string    `json:"expected_grade"`
	ComplexityLevel string    `json:"complexity_level"`
	Description     string    `json:"description"`
	Tags            []string  `json:"tags"`
}

// APISource defines where to collect an API specification from
type APISource struct {
	Name          string
	Provider      string
	Category      string
	URL           string
	ExpectedGrade string
	Complexity    string
	Description   string
	Tags          []string
}

// Collector handles downloading and organizing real-world API specifications
type Collector struct {
	BaseDir    string
	HTTPClient *http.Client
}

// NewCollector creates a new API specification collector
func NewCollector(baseDir string) *Collector {
	return &Collector{
		BaseDir: baseDir,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetRealWorldAPISources returns a curated list of high-quality real-world APIs
func GetRealWorldAPISources() []APISource {
	return []APISource{
		// Fintech APIs
		{
			Name:          "stripe",
			Provider:      "Stripe",
			Category:      "fintech",
			URL:           "https://raw.githubusercontent.com/stripe/openapi/master/openapi/spec3.yaml",
			ExpectedGrade: "A+",
			Complexity:    "high",
			Description:   "Payment processing with complex webhooks and subscriptions",
			Tags:          []string{"payments", "webhooks", "subscriptions", "complex"},
		},
		{
			Name:          "plaid",
			Provider:      "Plaid",
			Category:      "fintech",
			URL:           "https://raw.githubusercontent.com/plaid/plaid-openapi/master/2020-09-14.yml",
			ExpectedGrade: "A",
			Complexity:    "high",
			Description:   "Financial data aggregation and banking integrations",
			Tags:          []string{"banking", "financial-data", "aggregation"},
		},

		// Developer Platform APIs
		{
			Name:          "github",
			Provider:      "GitHub",
			Category:      "developer",
			URL:           "https://raw.githubusercontent.com/github/rest-api-description/main/descriptions/api.github.com/api.github.com.yaml",
			ExpectedGrade: "A+",
			Complexity:    "very-high",
			Description:   "Version control, repositories, and collaboration platform",
			Tags:          []string{"git", "repositories", "collaboration", "massive"},
		},
		{
			Name:          "gitlab",
			Provider:      "GitLab",
			Category:      "developer",
			URL:           "https://docs.gitlab.com/ee/api/openapi/openapi.yaml",
			ExpectedGrade: "A",
			Complexity:    "high",
			Description:   "DevOps lifecycle and CI/CD pipelines",
			Tags:          []string{"devops", "cicd", "git", "pipelines"},
		},

		// Cloud APIs
		{
			Name:          "digitalocean",
			Provider:      "DigitalOcean",
			Category:      "cloud",
			URL:           "https://api-engineering.nyc3.digitaloceanspaces.com/spec-ci/DigitalOcean-public.v2.yaml",
			ExpectedGrade: "A",
			Complexity:    "medium",
			Description:   "Simple cloud infrastructure management",
			Tags:          []string{"cloud", "infrastructure", "droplets", "kubernetes"},
		},

		// Communication APIs
		{
			Name:          "twilio",
			Provider:      "Twilio",
			Category:      "communication",
			URL:           "https://raw.githubusercontent.com/twilio/twilio-oai/main/spec/yaml/twilio_api_v2010.yaml",
			ExpectedGrade: "B+",
			Complexity:    "medium",
			Description:   "SMS, voice, and communication services",
			Tags:          []string{"sms", "voice", "communication", "telephony"},
		},
		{
			Name:          "sendgrid",
			Provider:      "SendGrid",
			Category:      "communication",
			URL:           "https://raw.githubusercontent.com/sendgrid/sendgrid-oai/main/oai.yaml",
			ExpectedGrade: "B",
			Complexity:    "medium",
			Description:   "Email delivery and marketing automation",
			Tags:          []string{"email", "marketing", "automation", "delivery"},
		},

		// E-commerce APIs
		{
			Name:          "shopify-admin",
			Provider:      "Shopify",
			Category:      "ecommerce",
			URL:           "https://shopify.dev/admin-api-reference.yaml",
			ExpectedGrade: "B+",
			Complexity:    "high",
			Description:   "E-commerce platform and store management",
			Tags:          []string{"ecommerce", "retail", "stores", "products"},
		},

		// Analytics APIs
		{
			Name:          "mixpanel",
			Provider:      "Mixpanel",
			Category:      "analytics",
			URL:           "https://raw.githubusercontent.com/mixpanel/docs/main/reference/openapi.yaml",
			ExpectedGrade: "B",
			Complexity:    "medium",
			Description:   "Product analytics and user behavior tracking",
			Tags:          []string{"analytics", "tracking", "events", "user-behavior"},
		},

		// Social APIs (Note: Many social APIs don't provide public OpenAPI specs)
		// We'll need to create representative examples or find community versions
	}
}

// CollectAPI downloads and stores an API specification with metadata
func (c *Collector) CollectAPI(source APISource) error {
	fmt.Printf("📥 Collecting %s API from %s...\n", source.Name, source.Provider)

	// Create directory structure
	categoryDir := filepath.Join(c.BaseDir, source.Category)
	apiDir := filepath.Join(categoryDir, source.Name)

	if err := os.MkdirAll(apiDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", apiDir, err)
	}

	// Download the OpenAPI specification
	resp, err := c.HTTPClient.Get(source.URL)
	if err != nil {
		return fmt.Errorf("failed to download API spec: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error %d when downloading %s", resp.StatusCode, source.URL)
	}

	// Save the specification file
	specPath := filepath.Join(apiDir, "openapi.yaml")
	specFile, err := os.Create(specPath)
	if err != nil {
		return fmt.Errorf("failed to create spec file: %w", err)
	}
	defer specFile.Close()

	_, err = io.Copy(specFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save spec file: %w", err)
	}

	// Create metadata
	metadata := APIMetadata{
		Name:            source.Name,
		Provider:        source.Provider,
		Category:        source.Category,
		Version:         "latest",
		SourceURL:       source.URL,
		CollectedAt:     time.Now(),
		ExpectedGrade:   source.ExpectedGrade,
		ComplexityLevel: source.Complexity,
		Description:     source.Description,
		Tags:            source.Tags,
	}

	// Save metadata
	metadataPath := filepath.Join(apiDir, "metadata.json")
	metadataFile, err := os.Create(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}
	defer metadataFile.Close()

	encoder := json.NewEncoder(metadataFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(metadata); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	// Create README
	readmePath := filepath.Join(apiDir, "README.md")
	readmeContent := fmt.Sprintf(`# %s API

**Provider**: %s  
**Category**: %s  
**Expected Grade**: %s  
**Complexity**: %s  

## Description
%s

## Tags
%s

## Source
- **URL**: %s
- **Collected**: %s

## Files
- openapi.yaml - The OpenAPI specification
- metadata.json - Collection metadata
- README.md - This file

## Usage

Validate with SpecGrade:
specgrade --target-dir=. --spec-version=3.1.0 --output-format=json
`,
		source.Provider,
		source.Provider,
		source.Category,
		source.ExpectedGrade,
		source.Complexity,
		source.Description,
		strings.Join(source.Tags, ", "),
		source.URL,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to create README: %w", err)
	}

	fmt.Printf("✅ Successfully collected %s API\n", source.Name)
	return nil
}

// CollectAll downloads all APIs from the curated source list
func (c *Collector) CollectAll() error {
	sources := GetRealWorldAPISources()

	fmt.Printf("🌍 Collecting %d real-world API specifications...\n", len(sources))

	successCount := 0
	for _, source := range sources {
		if err := c.CollectAPI(source); err != nil {
			fmt.Printf("❌ Failed to collect %s: %v\n", source.Name, err)
			continue
		}
		successCount++

		// Be respectful to API providers
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("🎉 Successfully collected %d/%d APIs\n", successCount, len(sources))
	return nil
}

// UpdateAPI re-downloads a specific API to get the latest version
func (c *Collector) UpdateAPI(apiName string) error {
	sources := GetRealWorldAPISources()

	for _, source := range sources {
		if source.Name == apiName {
			return c.CollectAPI(source)
		}
	}

	return fmt.Errorf("API '%s' not found in source list", apiName)
}

// ListCollectedAPIs returns information about all collected APIs
func (c *Collector) ListCollectedAPIs() ([]APIMetadata, error) {
	var apis []APIMetadata

	categories := []string{"fintech", "developer", "cloud", "communication", "ecommerce", "analytics", "social"}

	for _, category := range categories {
		categoryDir := filepath.Join(c.BaseDir, category)

		if _, err := os.Stat(categoryDir); os.IsNotExist(err) {
			continue
		}

		entries, err := os.ReadDir(categoryDir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			metadataPath := filepath.Join(categoryDir, entry.Name(), "metadata.json")
			metadataFile, err := os.Open(metadataPath)
			if err != nil {
				continue
			}

			var metadata APIMetadata
			if err := json.NewDecoder(metadataFile).Decode(&metadata); err != nil {
				metadataFile.Close()
				continue
			}
			metadataFile.Close()

			apis = append(apis, metadata)
		}
	}

	return apis, nil
}

// GetAPIStats returns statistics about the collected API collection
func (c *Collector) GetAPIStats() (map[string]interface{}, error) {
	apis, err := c.ListCollectedAPIs()
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_apis": len(apis),
		"categories": make(map[string]int),
		"providers":  make(map[string]int),
		"grades":     make(map[string]int),
		"complexity": make(map[string]int),
	}

	for _, api := range apis {
		stats["categories"].(map[string]int)[api.Category]++
		stats["providers"].(map[string]int)[api.Provider]++
		stats["grades"].(map[string]int)[api.ExpectedGrade]++
		stats["complexity"].(map[string]int)[api.ComplexityLevel]++
	}

	return stats, nil
}
