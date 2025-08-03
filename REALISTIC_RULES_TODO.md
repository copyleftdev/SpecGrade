# Realistic Rules TODO

## Overview
This document outlines new validation rules to be implemented based on real-world API validation findings from major providers (Stripe, Twilio, DigitalOcean, SendGrid, etc.). These rules address common quality issues found in production APIs.

## Real-World Findings Summary

### ✅ Successfully Validated APIs
- **Twilio API**: Grade B (75) - Missing error responses, some documentation gaps
- **Bad-Example Spec**: Grade C (62) - Type mismatches, missing descriptions

### ❌ APIs with Critical Issues
- **DigitalOcean API**: Circular schema references (`apiAgent` → `apiWorkspace` → `apiAgent`)
- **Stripe API**: Complex circular references across multiple schemas
- **SendGrid API**: Schema definition errors (`bad data in definitions`)

---

## Proposed New Rules

### 🔄 **Circular Reference Detection Rules**

#### 1. `circular-schema-references`
- **Priority**: HIGH
- **Category**: `schema_integrity`
- **Description**: Detect and prevent circular schema references that cause parsing failures
- **Real-world Impact**: Found in Stripe, DigitalOcean APIs
- **Implementation**: 
  - Track schema reference chains during validation
  - Detect cycles using visited set/map
  - Report circular paths with full reference chain
- **Grade Impact**: Should be CRITICAL (automatic F if present)
- **Example**: `#/components/schemas/User → #/components/schemas/Profile → #/components/schemas/User`

#### 2. `max-schema-depth`
- **Priority**: MEDIUM
- **Category**: `schema_complexity`
- **Description**: Limit maximum nesting depth in schema references
- **Rationale**: Prevents overly complex schemas that are hard to understand/implement
- **Threshold**: Warn at depth 5, error at depth 10
- **Grade Impact**: Minor deduction for warnings, major for errors

### 📝 **Enhanced Documentation Rules**

#### 3. `comprehensive-error-responses`
- **Priority**: HIGH
- **Category**: `error_handling`
- **Description**: Ensure all operations define proper error responses (400, 401, 403, 404, 500)
- **Real-world Impact**: Missing in Twilio (197 operations), common pattern
- **Implementation**:
  - Check each operation for standard HTTP error codes
  - Validate error response schemas are defined
  - Ensure error responses have meaningful descriptions
- **Grade Impact**: Significant deduction for missing error handling

#### 4. `operation-description-quality`
- **Priority**: MEDIUM
- **Category**: `documentation`
- **Description**: Validate quality of operation descriptions beyond just presence
- **Criteria**:
  - Minimum length (e.g., 20 characters)
  - Contains action verbs
  - Describes what the operation does, not just restates the path
  - No placeholder text ("TODO", "TBD", etc.)
- **Grade Impact**: Minor deduction for poor quality descriptions

#### 5. `parameter-description-completeness`
- **Priority**: MEDIUM
- **Category**: `documentation`
- **Description**: Ensure all parameters have meaningful descriptions
- **Criteria**:
  - All path, query, header parameters have descriptions
  - Descriptions explain parameter purpose and format
  - Required vs optional clearly indicated
- **Real-world Impact**: Often missing in production APIs

### 🔒 **Security and Best Practices Rules**

#### 6. `security-scheme-completeness`
- **Priority**: HIGH
- **Category**: `security`
- **Description**: Validate security schemes are properly defined and applied
- **Checks**:
  - Security schemes have proper type and configuration
  - Operations reference appropriate security schemes
  - No operations without security (unless explicitly public)
- **Real-world Impact**: Security often incomplete or inconsistent

#### 7. `sensitive-data-exposure`
- **Priority**: HIGH
- **Category**: `security`
- **Description**: Detect potential sensitive data exposure in examples/schemas
- **Patterns to detect**:
  - API keys, tokens, passwords in examples
  - Email addresses, phone numbers in examples
  - Credit card numbers, SSNs in examples
- **Implementation**: Regex patterns + heuristics
- **Grade Impact**: Critical security issue

### 📊 **Schema Quality Rules**

#### 8. `schema-example-consistency-enhanced`
- **Priority**: MEDIUM
- **Category**: `schema_validation`
- **Description**: Enhanced version of existing rule with better type checking
- **Improvements**:
  - Better handling of numeric types (int vs float)
  - Date/time format validation
  - Enum value validation
  - Null handling for nullable fields
- **Real-world Impact**: Type mismatches found in test specs

#### 9. `schema-property-naming`
- **Priority**: LOW
- **Category**: `consistency`
- **Description**: Enforce consistent property naming conventions
- **Options**:
  - camelCase vs snake_case consistency
  - Avoid reserved keywords
  - Meaningful property names (not x, y, z, etc.)
- **Grade Impact**: Minor deduction for inconsistency

#### 10. `required-fields-validation`
- **Priority**: MEDIUM
- **Category**: `schema_integrity`
- **Description**: Validate that required fields are properly defined and used
- **Checks**:
  - Required fields exist in schema properties
  - Required fields are not nullable (unless explicitly allowed)
  - Examples include all required fields

### 🌐 **API Design Best Practices Rules**

#### 11. `http-method-semantics`
- **Priority**: MEDIUM
- **Category**: `api_design`
- **Description**: Ensure HTTP methods are used semantically correctly
- **Checks**:
  - GET operations don't have request bodies
  - POST/PUT/PATCH have appropriate request bodies
  - DELETE operations return appropriate status codes
  - Idempotency considerations for PUT vs POST

#### 12. `resource-naming-consistency`
- **Priority**: LOW
- **Category**: `api_design`
- **Description**: Validate consistent resource naming in paths
- **Checks**:
  - Plural vs singular consistency
  - Consistent path parameter naming
  - RESTful path structure
- **Example**: `/users/{userId}` vs `/user/{id}` inconsistency

#### 13. `response-schema-completeness`
- **Priority**: MEDIUM
- **Category**: `response_validation`
- **Description**: Ensure all successful responses have proper schemas
- **Checks**:
  - 200/201 responses have content schemas
  - Response schemas match expected data types
  - Pagination metadata for list endpoints

### 🔧 **Implementation Quality Rules**

#### 14. `openapi-version-compliance`
- **Priority**: HIGH
- **Category**: `specification_compliance`
- **Description**: Validate strict compliance with OpenAPI specification version
- **Checks**:
  - Use of deprecated features
  - Version-specific syntax validation
  - Required fields for the specified version

#### 15. `external-reference-validation`
- **Priority**: MEDIUM
- **Category**: `reference_integrity`
- **Description**: Validate external references are accessible and valid
- **Implementation**:
  - Check HTTP references return 200
  - Validate referenced schemas are valid
  - Warn about unreachable external references

---

## Implementation Priority

### Phase 1 (Critical Issues) 🚨
1. `circular-schema-references` - Prevents parsing failures
2. `comprehensive-error-responses` - Major usability issue
3. `security-scheme-completeness` - Security critical
4. `sensitive-data-exposure` - Security critical

### Phase 2 (Quality Improvements) 📈
5. `operation-description-quality` - Developer experience
6. `schema-example-consistency-enhanced` - Data integrity
7. `required-fields-validation` - Schema integrity
8. `openapi-version-compliance` - Specification compliance

### Phase 3 (Polish & Consistency) ✨
9. `max-schema-depth` - Complexity management
10. `parameter-description-completeness` - Documentation
11. `http-method-semantics` - API design
12. `response-schema-completeness` - Response validation

### Phase 4 (Nice-to-Have) 🎯
13. `schema-property-naming` - Consistency
14. `resource-naming-consistency` - API design
15. `external-reference-validation` - Reference integrity

---

## Testing Strategy

### Unit Tests
- Each rule should have comprehensive unit tests
- Test against known good/bad examples
- Edge case coverage (empty specs, malformed data)

### Integration Tests
- Test rules against real-world API specs
- Validate rule interactions don't conflict
- Performance testing with large specs

### Real-World Validation
- Test against collected APIs in `test/realworld/collected-apis/`
- Compare results with expected grades
- Validate rule effectiveness on production APIs

---

## Configuration Options

### Rule Severity Levels
- `CRITICAL`: Automatic grade reduction to F
- `ERROR`: Major grade impact (10-20 points)
- `WARNING`: Moderate grade impact (5-10 points)
- `INFO`: Minor grade impact (1-5 points)

### Customization
- Allow rules to be disabled via config
- Configurable thresholds (e.g., max schema depth)
- Industry-specific rule sets (fintech, healthcare, etc.)

---

## Success Metrics

### Quality Improvements
- Increase in average grades for real-world APIs after fixes
- Reduction in parsing failures
- Better developer experience scores

### Adoption Metrics
- Rule usage in CI/CD pipelines
- Community feedback on rule effectiveness
- Contribution of new rules from users

---

## Next Steps

1. **Create rule implementation framework** - Standardize rule creation process
2. **Implement Phase 1 rules** - Focus on critical issues first
3. **Add comprehensive tests** - Ensure rule reliability
4. **Update documentation** - Document new rules and usage
5. **Validate against real-world APIs** - Test effectiveness
6. **Gather community feedback** - Iterate based on usage

---

## Notes

- Rules should be backward compatible with existing validation
- Consider performance impact of complex rules
- Provide clear error messages with fix suggestions
- Document rule rationale for community understanding
- Consider rule interactions and conflicts
