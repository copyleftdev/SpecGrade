package fetcher

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
)

// LocalSpecLoader loads OpenAPI specs from local files
type LocalSpecLoader struct {
	targetDir string
}

// NewLocalSpecLoader creates a new local spec loader
func NewLocalSpecLoader(targetDir string) *LocalSpecLoader {
	return &LocalSpecLoader{
		targetDir: targetDir,
	}
}

// Load loads an OpenAPI spec from the target directory
func (l *LocalSpecLoader) Load(version string) (*openapi3.T, error) {
	// Look for common OpenAPI spec file names
	possibleFiles := []string{
		"openapi.yaml", "openapi.yml", "openapi.json",
		"swagger.yaml", "swagger.yml", "swagger.json",
		"api.yaml", "api.yml", "api.json",
	}

	var specFile string
	for _, file := range possibleFiles {
		fullPath := filepath.Join(l.targetDir, file)
		if fileExists(fullPath) {
			specFile = fullPath
			break
		}
	}

	if specFile == "" {
		return nil, fmt.Errorf("no OpenAPI spec file found in directory: %s", l.targetDir)
	}

	// Create a loader that can resolve external references
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	// Load the spec from file (this will automatically resolve external $ref)
	spec, err := loader.LoadFromFile(specFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	// Note: We skip strict validation to allow grading of specs with minor issues
	// This allows SpecGrade to provide feedback on specs that have validation errors
	// but are still structurally sound enough to analyze

	return spec, nil
}

// fileExists checks if a file exists
func fileExists(filename string) bool {
	_, err := ioutil.ReadFile(filename)
	return err == nil
}
