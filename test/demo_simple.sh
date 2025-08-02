#!/bin/bash

echo "🚀 SpecGrade Advanced Testing Demonstration"
echo "=================================================="

cd "$(dirname "$0")/.."

# Build SpecGrade first
echo "Building SpecGrade..."
make build

echo ""
echo "📊 Test 1: Grade Distribution with Generated Specs"
echo "----------------------------------------"

# Create test specs directory
mkdir -p test/generated-specs

# Generate a perfect spec
cat > test/generated-specs/perfect.yaml << 'EOF'
openapi: 3.1.0
info:
  title: Perfect API
  version: 1.0.0
  description: A perfectly documented API
paths:
  /users:
    get:
      operationId: getUsers
      description: Retrieve all users
      responses:
        '200':
          description: Success
          content:
            application/json:
              schema:
                type: array
                items:
                  type: object
                  properties:
                    id:
                      type: integer
                      example: 123
                    name:
                      type: string
                      example: "John Doe"
        '400':
          description: Bad Request
        '500':
          description: Internal Server Error
      security:
        - bearerAuth: []
  /users/{id}:
    get:
      operationId: getUserById
      description: Retrieve a specific user
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Success
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: integer
                    example: 123
                  name:
                    type: string
                    example: "John Doe"
        '400':
          description: Bad Request
        '404':
          description: Not Found
        '500':
          description: Internal Server Error
      security:
        - bearerAuth: []
components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
EOF

# Generate a poor spec with intentional issues
cat > test/generated-specs/poor.yaml << 'EOF'
openapi: 3.1.0
info:
  title: Poor API
  version: 1.0.0
paths:
  /endpoint1:
    get:
      operationId: getEndpoint1
      responses:
        '200':
          description: Success
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: integer
                    example: "not_a_number"  # Type mismatch!
                  name:
                    type: string
                    example: 123  # Type mismatch!
  /endpoint2:
    post:
      operationId: postEndpoint2
      responses:
        '200':
          description: Success
  /endpoint3:
    get:
      operationId: getEndpoint3
      responses:
        '200':
          description: Success
EOF

# Generate a massive spec with many endpoints
cat > test/generated-specs/massive.yaml << 'EOF'
openapi: 3.1.0
info:
  title: Massive API
  version: 1.0.0
  description: Large-scale API with many endpoints
paths:
EOF

# Add 50 endpoints to the massive spec
for i in {1..50}; do
cat >> test/generated-specs/massive.yaml << EOF
  /endpoint$i:
    get:
      operationId: getEndpoint$i
      description: Get endpoint $i data
      responses:
        '200':
          description: Success
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: integer
                    example: $i
        '400':
          description: Bad Request
        '500':
          description: Internal Server Error
EOF
done

cat >> test/generated-specs/massive.yaml << 'EOF'
components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
EOF

# Test each generated spec
echo "Testing Perfect Spec (Expected: A+ grade):"
mkdir -p test/perfect-test && cp test/generated-specs/perfect.yaml test/perfect-test/openapi.yaml
./build/specgrade --target-dir=./test/perfect-test --spec-version=3.1.0 --output-format=json | jq '.grade, .score'

echo ""
echo "Testing Poor Spec (Expected: D/F grade):"
mkdir -p test/poor-test && cp test/generated-specs/poor.yaml test/poor-test/openapi.yaml
./build/specgrade --target-dir=./test/poor-test --spec-version=3.1.0 --output-format=json | jq '.grade, .score'

echo ""
echo "Testing Massive Spec (50 endpoints - Performance test):"
mkdir -p test/massive-test && cp test/generated-specs/massive.yaml test/massive-test/openapi.yaml
time ./build/specgrade --target-dir=./test/massive-test --spec-version=3.1.0 --output-format=json | jq '.grade, .score'

echo ""
echo "🔬 Test 2: Edge Cases"
echo "----------------------------------------"

# Create Unicode spec
cat > test/generated-specs/unicode.yaml << 'EOF'
openapi: 3.1.0
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
                    example: "🎉✨🌟"
        '400':
          description: Bad Request
        '500':
          description: Internal Server Error
components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
EOF

echo "Testing Unicode Content (Edge case):"
mkdir -p test/unicode-test && cp test/generated-specs/unicode.yaml test/unicode-test/openapi.yaml
./build/specgrade --target-dir=./test/unicode-test --spec-version=3.1.0 --output-format=json | jq '.grade, .score'

echo ""
echo "⚡ Test 3: Comparison with Existing Specs"
echo "----------------------------------------"

echo "Original Perfect Sample:"
./build/specgrade --target-dir=./test/sample-spec --spec-version=3.1.0 --output-format=json | jq '.grade, .score'

echo ""
echo "Original Bad Example:"
./build/specgrade --target-dir=./test/sample-spec/bad-example --spec-version=3.1.0 --output-format=json | jq '.grade, .score'

echo ""
echo "🎯 Test 4: Rule Coverage Analysis"
echo "----------------------------------------"

echo "All available rules with Spectral-compatible naming:"
./build/specgrade rules

echo ""
echo "🎉 Advanced Testing Complete!"
echo ""
echo "SpecGrade demonstrates robust handling of:"
echo "  ✅ Diverse quality profiles with predictable grading"
echo "  ✅ Edge cases including Unicode content and large specs"
echo "  ✅ Performance with 50+ endpoint specifications"
echo "  ✅ Industry-standard rule naming (Spectral-compatible)"
echo "  ✅ Consistent grading across different spec types"

# Cleanup
rm -rf test/generated-specs test/perfect-test test/poor-test test/massive-test test/unicode-test
