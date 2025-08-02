package generator

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// SpecGenerator creates OpenAPI specifications with controlled quality issues
type SpecGenerator struct {
	rand *rand.Rand
}

// NewSpecGenerator creates a new specification generator
func NewSpecGenerator() *SpecGenerator {
	return &SpecGenerator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// QualityProfile defines the expected quality characteristics
type QualityProfile struct {
	TargetGrade        string  // A+, A, B, C, D, F
	MissingDescriptions float64 // 0.0 = none missing, 1.0 = all missing
	TypeMismatches     int     // Number of intentional type/example mismatches
	MissingErrorCodes  float64 // 0.0 = all present, 1.0 = all missing
	SecurityIssues     bool    // Whether to include security problems
	ComplexityLevel    int     // 1-5, affects number of paths/schemas
}

// GenerateSpec creates an OpenAPI spec matching the quality profile
func (g *SpecGenerator) GenerateSpec(profile QualityProfile) string {
	spec := strings.Builder{}
	
	// Header
	spec.WriteString("openapi: 3.1.0\n")
	spec.WriteString("info:\n")
	spec.WriteString("  title: Generated Test API\n")
	spec.WriteString("  version: 1.0.0\n")
	
	// Conditionally add description based on quality profile
	if g.rand.Float64() > profile.MissingDescriptions {
		spec.WriteString("  description: A generated API for testing SpecGrade\n")
	}
	
	// Security schemes (or lack thereof)
	if !profile.SecurityIssues {
		spec.WriteString("components:\n")
		spec.WriteString("  securitySchemes:\n")
		spec.WriteString("    bearerAuth:\n")
		spec.WriteString("      type: http\n")
		spec.WriteString("      scheme: bearer\n")
	}
	
	// Generate paths based on complexity
	spec.WriteString("paths:\n")
	numPaths := profile.ComplexityLevel * 3
	
	for i := 0; i < numPaths; i++ {
		pathName := fmt.Sprintf("/resource%d", i+1)
		spec.WriteString(fmt.Sprintf("  %s:\n", pathName))
		
		// GET operation
		spec.WriteString("    get:\n")
		spec.WriteString(fmt.Sprintf("      operationId: getResource%d\n", i+1))
		
		// Conditionally add description
		if g.rand.Float64() > profile.MissingDescriptions {
			spec.WriteString(fmt.Sprintf("      description: Retrieve resource %d\n", i+1))
		}
		
		// Add responses
		spec.WriteString("      responses:\n")
		spec.WriteString("        '200':\n")
		spec.WriteString("          description: Success\n")
		spec.WriteString("          content:\n")
		spec.WriteString("            application/json:\n")
		spec.WriteString("              schema:\n")
		spec.WriteString("                type: object\n")
		spec.WriteString("                properties:\n")
		spec.WriteString("                  id:\n")
		spec.WriteString("                    type: integer\n")
		
		// Intentionally add type mismatches
		if profile.TypeMismatches > 0 && i < profile.TypeMismatches {
			spec.WriteString("                    example: \"not_a_number\"\n") // String example for integer type
		} else {
			spec.WriteString("                    example: 123\n")
		}
		
		spec.WriteString("                  name:\n")
		spec.WriteString("                    type: string\n")
		spec.WriteString("                    example: \"Resource Name\"\n")
		
		// Conditionally add error responses
		if g.rand.Float64() > profile.MissingErrorCodes {
			spec.WriteString("        '400':\n")
			spec.WriteString("          description: Bad Request\n")
			spec.WriteString("        '500':\n")
			spec.WriteString("          description: Internal Server Error\n")
		}
		
		// Add security if not problematic
		if !profile.SecurityIssues {
			spec.WriteString("      security:\n")
			spec.WriteString("        - bearerAuth: []\n")
		}
	}
	
	return spec.String()
}

// PredefinedProfiles returns common quality profiles for testing
func PredefinedProfiles() map[string]QualityProfile {
	return map[string]QualityProfile{
		"perfect": {
			TargetGrade:         "A+",
			MissingDescriptions: 0.0,
			TypeMismatches:      0,
			MissingErrorCodes:   0.0,
			SecurityIssues:      false,
			ComplexityLevel:     3,
		},
		"excellent": {
			TargetGrade:         "A",
			MissingDescriptions: 0.1,
			TypeMismatches:      0,
			MissingErrorCodes:   0.1,
			SecurityIssues:      false,
			ComplexityLevel:     4,
		},
		"good": {
			TargetGrade:         "B",
			MissingDescriptions: 0.3,
			TypeMismatches:      1,
			MissingErrorCodes:   0.3,
			SecurityIssues:      false,
			ComplexityLevel:     5,
		},
		"average": {
			TargetGrade:         "C",
			MissingDescriptions: 0.5,
			TypeMismatches:      2,
			MissingErrorCodes:   0.5,
			SecurityIssues:      false,
			ComplexityLevel:     4,
		},
		"poor": {
			TargetGrade:         "D",
			MissingDescriptions: 0.7,
			TypeMismatches:      3,
			MissingErrorCodes:   0.7,
			SecurityIssues:      true,
			ComplexityLevel:     3,
		},
		"failing": {
			TargetGrade:         "F",
			MissingDescriptions: 0.9,
			TypeMismatches:      5,
			MissingErrorCodes:   0.9,
			SecurityIssues:      true,
			ComplexityLevel:     2,
		},
	}
}

// EdgeCaseGenerator creates specifications with specific edge cases
type EdgeCaseGenerator struct {
	rand *rand.Rand
}

// NewEdgeCaseGenerator creates a new edge case generator
func NewEdgeCaseGenerator() *EdgeCaseGenerator {
	return &EdgeCaseGenerator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateCircularRef creates a spec with circular schema references
func (g *EdgeCaseGenerator) GenerateCircularRef() string {
	return `openapi: 3.1.0
info:
  title: Circular Reference Test
  version: 1.0.0
paths:
  /test:
    get:
      operationId: testCircular
      responses:
        '200':
          description: Success
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/NodeA'
components:
  schemas:
    NodeA:
      type: object
      properties:
        id:
          type: string
        child:
          $ref: '#/components/schemas/NodeB'
    NodeB:
      type: object
      properties:
        id:
          type: string
        parent:
          $ref: '#/components/schemas/NodeA'`
}

// GenerateDeepNesting creates a spec with deeply nested schemas
func (g *EdgeCaseGenerator) GenerateDeepNesting(depth int) string {
	spec := strings.Builder{}
	spec.WriteString("openapi: 3.1.0\n")
	spec.WriteString("info:\n")
	spec.WriteString("  title: Deep Nesting Test\n")
	spec.WriteString("  version: 1.0.0\n")
	spec.WriteString("paths:\n")
	spec.WriteString("  /test:\n")
	spec.WriteString("    get:\n")
	spec.WriteString("      operationId: testDeepNesting\n")
	spec.WriteString("      responses:\n")
	spec.WriteString("        '200':\n")
	spec.WriteString("          description: Success\n")
	spec.WriteString("          content:\n")
	spec.WriteString("            application/json:\n")
	spec.WriteString("              schema:\n")
	
	// Create deeply nested object
	for i := 0; i < depth; i++ {
		indent := strings.Repeat("  ", 8+i*2)
		if i == 0 {
			spec.WriteString("                type: object\n")
			spec.WriteString("                properties:\n")
			spec.WriteString("                  level0:\n")
		} else {
			spec.WriteString(fmt.Sprintf("%stype: object\n", indent))
			spec.WriteString(fmt.Sprintf("%sproperties:\n", indent))
			spec.WriteString(fmt.Sprintf("%s  level%d:\n", indent, i))
		}
	}
	
	// Final property
	finalIndent := strings.Repeat("  ", 8+depth*2)
	spec.WriteString(fmt.Sprintf("%stype: string\n", finalIndent))
	spec.WriteString(fmt.Sprintf("%sexample: \"Deep value at level %d\"\n", finalIndent, depth))
	
	return spec.String()
}

// GenerateUnicodeContent creates a spec with international characters
func (g *EdgeCaseGenerator) GenerateUnicodeContent() string {
	return `openapi: 3.1.0
info:
  title: 国际化API测试 (Internationalization API Test)
  version: 1.0.0
  description: API с поддержкой Unicode и эмодзи 🌍🚀
paths:
  /用户:
    get:
      operationId: getUsers
      description: Получить список пользователей 👥
      responses:
        '200':
          description: Успешный ответ
          content:
            application/json:
              schema:
                type: object
                properties:
                  名前:
                    type: string
                    example: "田中太郎"
                  email:
                    type: string
                    example: "user@例え.テスト"
                  статус:
                    type: string
                    enum: ["активный", "неактивный"]
                    example: "активный"
                  emoji:
                    type: string
                    example: "🎉✨🌟"`
}

// GenerateMassiveSpec creates a specification with many endpoints
func (g *EdgeCaseGenerator) GenerateMassiveSpec(numEndpoints int) string {
	spec := strings.Builder{}
	spec.WriteString("openapi: 3.1.0\n")
	spec.WriteString("info:\n")
	spec.WriteString("  title: Massive API Test\n")
	spec.WriteString("  version: 1.0.0\n")
	spec.WriteString("  description: Large-scale API with many endpoints\n")
	spec.WriteString("paths:\n")
	
	for i := 0; i < numEndpoints; i++ {
		spec.WriteString(fmt.Sprintf("  /endpoint%d:\n", i))
		spec.WriteString("    get:\n")
		spec.WriteString(fmt.Sprintf("      operationId: getEndpoint%d\n", i))
		spec.WriteString(fmt.Sprintf("      description: Get endpoint %d data\n", i))
		spec.WriteString("      responses:\n")
		spec.WriteString("        '200':\n")
		spec.WriteString("          description: Success\n")
		spec.WriteString("          content:\n")
		spec.WriteString("            application/json:\n")
		spec.WriteString("              schema:\n")
		spec.WriteString("                type: object\n")
		spec.WriteString("                properties:\n")
		spec.WriteString("                  id:\n")
		spec.WriteString("                    type: integer\n")
		spec.WriteString(fmt.Sprintf("                    example: %d\n", i))
		spec.WriteString("        '400':\n")
		spec.WriteString("          description: Bad Request\n")
		spec.WriteString("        '500':\n")
		spec.WriteString("          description: Internal Server Error\n")
	}
	
	return spec.String()
}
