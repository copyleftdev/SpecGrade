#!/bin/bash

# SpecGrade Enhanced Developer Reporting Demo
# Showcases the new actionable, root-cause focused reporting system

echo "🚀 SpecGrade Enhanced Developer Reporting Demo"
echo "=============================================="
echo ""

# Build SpecGrade
echo "📦 Building SpecGrade..."
go build -o specgrade

echo ""
echo "🎯 Demo 1: Enhanced Developer CLI Format"
echo "----------------------------------------"
echo "The new 'developer' format provides rich, actionable insights:"
echo ""

./specgrade --target-dir test/sample-spec --output-format developer

echo ""
echo "🔍 Demo 2: Structured JSON Analytics for CI/CD"
echo "----------------------------------------------"
echo "Enhanced JSON output with developer analytics:"
echo ""

echo "📊 Summary Analytics:"
./specgrade --target-dir test/sample-spec --output-format json | jq '.summary'

echo ""
echo "📈 Complexity & Risk Analytics:"
./specgrade --target-dir test/sample-spec --output-format json | jq '.analytics'

echo ""
echo "🔧 Individual Issue with Actionable Fix:"
./specgrade --target-dir test/sample-spec --output-format json | jq '.rules[] | select(.passed == false) | {ruleID, severity, category, suggestion: .suggestion.title, estimate: .suggestion.estimate, impact: .impact.business_value}'

echo ""
echo "⚡ Demo 3: Comparison - Old vs New Reporting"
echo "-------------------------------------------"

echo "📋 Traditional CLI Format (basic):"
./specgrade --target-dir test/sample-spec --output-format cli

echo ""
echo "🚀 Enhanced Developer Format (actionable):"
echo "   ✅ Step-by-step fix instructions"
echo "   ✅ Code examples and references"
echo "   ✅ Impact analysis (UX, DX, business)"
echo "   ✅ Time estimates and difficulty levels"
echo "   ✅ Location information (JSON paths)"
echo "   ✅ Risk assessment and compliance gaps"
echo "   ✅ Complexity analytics"
echo "   ✅ Maintenance and developer-friendly scores"

echo ""
echo "🎉 Key Benefits of Enhanced Reporting:"
echo "======================================"
echo "🔧 ACTIONABLE: Specific steps to fix each issue"
echo "⏱️  TIME-AWARE: Estimates help prioritize work"
echo "💡 EDUCATIONAL: Learn why issues matter"
echo "🎯 PRIORITIZED: Focus on high-impact, easy wins"
echo "📊 ANALYTICAL: Understand API complexity and risks"
echo "🤖 CI/CD READY: Structured JSON for automation"
echo "👥 DEVELOPER-FOCUSED: Built for real-world workflows"

echo ""
echo "🌟 SpecGrade Enhanced Reporting: Making API quality actionable!"
echo "=============================================================="

# Cleanup
rm -f specgrade
