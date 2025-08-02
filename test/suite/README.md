# SpecGrade Test Suite

## 🎯 Comprehensive OpenAPI Specification Testing

To ensure SpecGrade's rigor and reliability, we need extensive testing across diverse OpenAPI specifications representing real-world scenarios, edge cases, and unknown patterns.

## 📊 Test Categories

### 1. **Grade Distribution Tests**
- `grade-a-plus/` - Perfect specifications (95-100%)
- `grade-a/` - Excellent specifications (90-94%)
- `grade-b/` - Good specifications (75-89%)
- `grade-c/` - Average specifications (60-74%)
- `grade-d/` - Poor specifications (50-59%)
- `grade-f/` - Failing specifications (0-49%)

### 2. **Complexity Tests**
- `minimal/` - Bare minimum valid specs
- `simple/` - Basic CRUD APIs
- `moderate/` - Multi-resource APIs with relationships
- `complex/` - Enterprise APIs with advanced patterns
- `massive/` - Large-scale APIs (100+ endpoints)

### 3. **Industry Domain Tests**
- `healthcare/` - FHIR, medical APIs
- `fintech/` - Banking, payment APIs
- `ecommerce/` - Shopping, inventory APIs
- `iot/` - Device management APIs
- `social/` - Social media, messaging APIs
- `government/` - Public sector APIs

### 4. **Edge Case Tests**
- `edge-cases/circular-refs/` - Circular schema references
- `edge-cases/deep-nesting/` - Deeply nested schemas
- `edge-cases/large-enums/` - Enums with many values
- `edge-cases/unicode/` - International characters
- `edge-cases/empty-objects/` - Empty schemas and responses
- `edge-cases/polymorphism/` - oneOf, anyOf, allOf patterns

### 5. **Version Compatibility Tests**
- `openapi-2.0/` - Swagger 2.0 specifications
- `openapi-3.0/` - OpenAPI 3.0.x specifications
- `openapi-3.1/` - OpenAPI 3.1.x specifications
- `mixed-versions/` - Specs with version inconsistencies

### 6. **Multi-File Architecture Tests**
- `single-file/` - Monolithic specifications
- `multi-file/basic/` - Simple external references
- `multi-file/complex/` - Deep reference hierarchies
- `multi-file/circular/` - Circular file dependencies
- `multi-file/broken/` - Missing or invalid references

### 7. **Security Pattern Tests**
- `security/none/` - No authentication
- `security/basic/` - Basic authentication
- `security/bearer/` - Bearer token authentication
- `security/oauth2/` - OAuth2 flows
- `security/api-key/` - API key authentication
- `security/mixed/` - Multiple security schemes
- `security/malformed/` - Invalid security definitions

### 8. **Documentation Quality Tests**
- `docs/excellent/` - Comprehensive documentation
- `docs/minimal/` - Bare minimum documentation
- `docs/missing/` - Missing descriptions
- `docs/inconsistent/` - Inconsistent documentation style
- `docs/generated/` - Auto-generated documentation

### 9. **Real-World API Tests**
- `real-world/stripe/` - Stripe API patterns
- `real-world/github/` - GitHub API patterns
- `real-world/aws/` - AWS API patterns
- `real-world/google/` - Google API patterns
- `real-world/microsoft/` - Microsoft API patterns

### 10. **Error Condition Tests**
- `malformed/syntax-errors/` - YAML/JSON syntax errors
- `malformed/schema-violations/` - OpenAPI schema violations
- `malformed/type-mismatches/` - Type inconsistencies
- `malformed/missing-required/` - Missing required fields
- `malformed/invalid-refs/` - Invalid $ref references

## 🔬 Advanced Testing Algorithms

### Property-Based Testing
```go
// Generate random OpenAPI specs with known properties
func TestPropertyBased(t *testing.T) {
    quick.Check(func(spec *GeneratedSpec) bool {
        result := specgrade.Validate(spec)
        return result.Grade.IsValid() && result.Score >= 0 && result.Score <= 100
    }, nil)
}
```

### Fuzzing Tests
```go
// Test with malformed/corrupted specifications
func TestFuzzing(t *testing.T) {
    for _, corruptedSpec := range generateCorruptedSpecs() {
        result := specgrade.Validate(corruptedSpec)
        // Should not panic, should return meaningful error
        assert.NotNil(t, result)
    }
}
```

### Regression Tests
```go
// Ensure consistent grading across versions
func TestRegression(t *testing.T) {
    for _, testCase := range regressionSuite {
        result := specgrade.Validate(testCase.Spec)
        assert.Equal(t, testCase.ExpectedGrade, result.Grade)
    }
}
```

### Performance Tests
```go
// Test with large specifications
func BenchmarkLargeSpecs(b *testing.B) {
    largeSpec := loadLargeSpec("massive/1000-endpoints.yaml")
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        specgrade.Validate(largeSpec)
    }
}
```

## 🎯 Expected Outcomes

### Grade Distribution Validation
- A+ specs should consistently score 95-100%
- F specs should consistently score 0-49%
- Grade boundaries should be stable and predictable

### Rule Coverage Analysis
- Each rule should trigger across multiple test cases
- No rule should be untested
- Edge cases should exercise rule boundary conditions

### Performance Benchmarks
- Small specs (< 100 lines): < 100ms
- Medium specs (< 1000 lines): < 500ms
- Large specs (< 10000 lines): < 2s
- Massive specs (> 10000 lines): < 10s

## 🔄 Continuous Testing Strategy

### Automated Test Generation
- Generate specs with known quality issues
- Create synthetic edge cases
- Mutate existing specs to create variants

### Community Contributions
- Accept real-world OpenAPI specs from users
- Build anonymized test cases from production APIs
- Crowdsource edge case discovery

### Machine Learning Integration
- Train models on grading patterns
- Detect anomalous specifications
- Predict quality issues before full validation

## 📈 Success Metrics

1. **Coverage**: 95%+ rule coverage across test suite
2. **Consistency**: < 5% grade variance on repeated runs
3. **Performance**: Sub-second validation for typical specs
4. **Robustness**: Zero crashes on malformed inputs
5. **Accuracy**: Manual review confirms 90%+ of grades
