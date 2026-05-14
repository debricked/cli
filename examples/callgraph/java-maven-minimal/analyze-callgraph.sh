#!/bin/bash

# Callgraph Parser and Analyzer Script
# This script helps analyze the generated debricked-call-graph.java file

CALLGRAPH_FILE="${1:-debricked-call-graph.java}"

if [ ! -f "$CALLGRAPH_FILE" ]; then
    echo "ERROR: Callgraph file not found: $CALLGRAPH_FILE"
    exit 1
fi

echo "=========================================="
echo "Callgraph Analysis"
echo "=========================================="
echo ""

# Check for jq (JSON processor)
if ! command -v jq &> /dev/null; then
    echo "⚠ Note: jq is not installed. For JSON parsing, install with:"
    echo "  sudo apt-get install jq   # Linux"
    echo "  brew install jq            # macOS"
    echo ""
    echo "Basic statistics without jq:"
    echo ""
fi

# Count methods
METHOD_COUNT=$(grep -o '"[^"]*"' "$CALLGRAPH_FILE" | grep -c '\[')
echo "Callgraph File: $CALLGRAPH_FILE"
echo "File Size: $(ls -lh "$CALLGRAPH_FILE" | awk '{print $5}')"
echo ""

if command -v jq &> /dev/null; then
    echo "=== JSON Analysis ==="
    echo ""

    # Parse with jq
    TOTAL_METHODS=$(jq '.data | length' "$CALLGRAPH_FILE")
    echo "Total Methods: $TOTAL_METHODS"

    USER_CODE=$(jq '[.data[] | select(.[1] == true)] | length' "$CALLGRAPH_FILE")
    LIBRARY_CODE=$(jq '[.data[] | select(.[1] == false)] | length' "$CALLGRAPH_FILE")

    echo "  - User Code Methods: $USER_CODE"
    echo "  - Library/JDK Methods: $LIBRARY_CODE"
    echo ""

    # Find entry points (methods in user code with no callers or few callers)
    echo "=== User-Code Methods ==="
    jq -r '.data[] | select(.[1] == true) | .[0]' "$CALLGRAPH_FILE" | head -20
    echo ""

    # Show methods that are most called (highest connectivity)
    echo "=== Most Called Methods (Top 10) ==="
    jq '[.data[] | {method: .[0], callerCount: (.[7] | length)}] | sort_by(-.callerCount) | .[0:10] | .[] | "\(.callerCount) callers: \(.method)"' -r "$CALLGRAPH_FILE"
    echo ""

else
    echo "=== Without jq (basic analysis) ==="
    echo ""

    # Simple grep-based analysis
    TOTAL=$(grep -o '\["com\.' "$CALLGRAPH_FILE" | wc -l)
    echo "Methods starting with 'com.': $TOTAL"

    echo ""
    echo "Sample methods found:"
    grep -o '"[^"]*com\.example\.callgraph[^"]*"' "$CALLGRAPH_FILE" | sort -u | head -10

fi

echo ""
echo "=========================================="
echo "For detailed analysis, see CALLGRAPH_OUTPUT.md"
echo "=========================================="

