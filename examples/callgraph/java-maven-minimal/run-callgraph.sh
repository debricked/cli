#!/bin/bash

# Debricked Java Callgraph Sample - Automated Build & Generation Script
# This script runs the full flow described in the README:
# 1. Verify Java version
# 2. Build with Maven
# 3. Copy dependencies
# 4. Generate callgraph

set -e

SAMPLE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SAMPLE_DIR"

echo "================================"
echo "Debricked Java Callgraph Sample"
echo "================================"
echo ""

# Step 1: Verify Java version
echo "[1/4] Checking Java version..."
if ! command -v java &> /dev/null; then
    echo "ERROR: Java not found. Please install Java 11 or higher."
    exit 1
fi
java -version
echo ""

# Step 2: Build Maven project
echo "[2/4] Building Maven project..."
if ! command -v mvn &> /dev/null; then
    echo "ERROR: Maven not found. Please install Maven."
    exit 1
fi
mvn package -q -DskipTests -e
echo "✓ Maven build complete"
echo ""

# Step 3: Copy dependencies
echo "[3/4] Copying external dependencies..."
mvn -q -B dependency:copy-dependencies -DoutputDirectory=./.debrickedTmpFolder -DskipTests -e
echo "✓ Dependencies copied to .debrickedTmpFolder/"
echo ""

# Step 3b: Seed .debricked/soot-wrapper.jar so CLI does not attempt a network download
# The embedded jar in the repo supports Java 21+; for Java 11/17 the CLI would normally
# download a versioned jar from GitHub. We pre-seed it with the embedded jar to keep
# this sample fully offline.
SOOT_JAR_DIR="$SAMPLE_DIR/.debricked"
SOOT_JAR_PATH="$SOOT_JAR_DIR/soot-wrapper.jar"
REPO_ROOT_JAR="$(cd "$SAMPLE_DIR/../../.." && pwd)/internal/callgraph/language/java/soot-wrapper.jar"
if [ ! -f "$SOOT_JAR_PATH" ]; then
    echo "Seeding .debricked/soot-wrapper.jar from repo..."
    mkdir -p "$SOOT_JAR_DIR"
    if [ -f "$REPO_ROOT_JAR" ]; then
        cp "$REPO_ROOT_JAR" "$SOOT_JAR_PATH"
        echo "✓ soot-wrapper.jar seeded from repo"
    else
        echo "⚠ Could not find soot-wrapper.jar at $REPO_ROOT_JAR"
        echo "  The CLI will attempt to download it from GitHub (requires internet access)"
    fi
fi
echo ""

# Step 4: Generate callgraph
echo "[4/4] Generating callgraph..."

# Try to find debricked in PATH, or use local repo binary
DEBRICKED_CMD="debricked"
if ! command -v debricked &> /dev/null; then
    # Try to find debricked binary in repo root
    # Sample is at: examples/callgraph/java-maven-minimal
    # Repo root is: ../../..
    REPO_ROOT="$(cd "$SAMPLE_DIR/../../.." && pwd)"
    if [ -f "$REPO_ROOT/debricked" ]; then
        DEBRICKED_CMD="$REPO_ROOT/debricked"
        echo "Using debricked from: $DEBRICKED_CMD"
    else
        echo "ERROR: debricked CLI not found in PATH or at $REPO_ROOT/debricked"
        echo "Please either:"
        echo "  1. Add debricked to your PATH"
        echo "  2. Build from repo root: cd $REPO_ROOT && go build -o debricked ./cmd/debricked"
        exit 1
    fi
fi

# Try with --no-build flag (assumes already built)
$DEBRICKED_CMD callgraph --no-build

echo ""
echo "================================"
echo "Callgraph generation complete!"
echo "================================"
echo ""

# Verify artifacts
if [ -f debricked-call-graph.java ]; then
    FILE_SIZE=$(wc -c < debricked-call-graph.java)
    echo "✓ Artifact created: debricked-call-graph.java"
    echo "  Size: $(ls -lh debricked-call-graph.java | awk '{print $5}')"
    echo "  Format: JSON callgraph data ($(echo "($FILE_SIZE / 1024)" | bc)KB)"
    echo ""
    echo "  Sample output (first 200 chars):"
    head -c 200 debricked-call-graph.java | sed 's/^/    /'
    echo ""
else
    echo "⚠ Expected artifact not found: debricked-call-graph.java"
    echo "  This may indicate an error in callgraph generation."
    exit 1
fi

if [ -d .debrickedTmpFolder ]; then
    DEPS_COUNT=$(find .debrickedTmpFolder -name "*.jar" 2>/dev/null | wc -l)
    echo "✓ Dependencies folder: .debrickedTmpFolder/ (contains $DEPS_COUNT jars)"
fi

echo ""
echo "✓ Sample complete! Callgraph successfully generated."
echo ""
echo "Next steps:"
echo "  - Review the generated debricked-call-graph.java"
echo "  - Parse the JSON to analyze call relationships"
echo "  - See README.md for details"

