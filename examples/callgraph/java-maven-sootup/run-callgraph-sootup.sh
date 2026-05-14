#!/bin/bash
# ============================================================
# SootUp Callgraph POC — build & run script
# v2: Uses SootUp instead of classic Soot (SootWrapper.jar)
#
# Steps:
#   1. Check prerequisites (Java, Maven)
#   2. Build the sample app  (mvn package)
#   3. Copy dependencies     (mvn dependency:copy-dependencies)
#   4. Build SootUpWrapper   (mvn package in SootUpWrapper/)
#   5. Run SootUpWrapper.jar directly via java -jar
#   6. Verify output
# ============================================================
set -e
SAMPLE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SAMPLE_DIR"
ALGORITHM="${1:-rta}"   # Pass "cha" as first arg to use CHA instead of RTA
echo "============================================================"
echo " SootUp Callgraph POC  (v2)"
echo " Algorithm: ${ALGORITHM^^}"
echo "============================================================"
echo ""
# ── Step 1: Prerequisites ──────────────────────────────────────
echo "[1/5] Checking prerequisites..."
if ! command -v java &>/dev/null; then
    echo "ERROR: Java not found. Install Java 11+."
    exit 1
fi
JAVA_VER=$(java -version 2>&1 | head -1)
echo "  Java : $JAVA_VER"
if ! command -v mvn &>/dev/null; then
    echo "ERROR: Maven not found."
    exit 1
fi
MVN_VER=$(mvn --version 2>&1 | head -1)
echo "  Maven: $MVN_VER"
echo ""
# ── Step 2: Build the sample app ──────────────────────────────
echo "[2/5] Building sample app (mvn package)..."
mvn package -q -DskipTests
echo "  ✓ Compiled → target/classes/"
echo ""
# ── Step 3: Copy runtime dependencies ─────────────────────────
echo "[3/5] Copying runtime dependencies..."
if [ ! -d ".debrickedTmpFolder" ]; then
    mvn -q -B dependency:copy-dependencies \
        -DoutputDirectory=./.debrickedTmpFolder \
        -DskipTests
    echo "  ✓ Dependencies → .debrickedTmpFolder/"
else
    echo "  ✓ .debrickedTmpFolder/ already present (skipping)"
fi
echo ""
# ── Step 4: Build SootUpWrapper fat JAR ───────────────────────
echo "[4/5] Building SootUpWrapper.jar..."
WRAPPER_DIR="$SAMPLE_DIR/SootUpWrapper"
WRAPPER_JAR="$WRAPPER_DIR/target/sootup-wrapper-1.0.0.jar"
if [ ! -f "$WRAPPER_JAR" ] || [ "$FORCE_REBUILD" = "1" ]; then
    (cd "$WRAPPER_DIR" && mvn package -q -DskipTests)
    echo "  ✓ Built → SootUpWrapper/target/sootup-wrapper-1.0.0.jar"
else
    echo "  ✓ SootUpWrapper.jar already built (set FORCE_REBUILD=1 to rebuild)"
fi
echo ""
# ── Step 5: Run SootUpWrapper ──────────────────────────────────
echo "[5/5] Running callgraph generation with SootUp ($ALGORITHM)..."
OUTPUT_FILE="debricked-call-graph-sootup.java"
java -jar "$WRAPPER_JAR" \
    -u "$SAMPLE_DIR/target/classes" \
    -l "$SAMPLE_DIR/.debrickedTmpFolder" \
    -f "$SAMPLE_DIR/$OUTPUT_FILE" \
    -a "$ALGORITHM"
echo ""
# ── Verify output ──────────────────────────────────────────────
echo "============================================================"
echo " Results"
echo "============================================================"
if [ -f "$OUTPUT_FILE" ]; then
    SIZE=$(wc -c < "$OUTPUT_FILE")
    METHODS=$(grep -o '\[' "$OUTPUT_FILE" | wc -l)
    echo "  ✓ Output: $OUTPUT_FILE"
    echo "    Size   : ${SIZE} bytes"
    echo ""
    echo "  First 300 chars:"
    head -c 300 "$OUTPUT_FILE" | sed 's/^/    /'
    echo ""
else
    echo "  ✗ Output file not found: $OUTPUT_FILE"
    exit 1
fi
echo ""
echo "  ✓ SootUp callgraph generation complete!"
echo ""
echo "  Compare with Soot v1 output:"
echo "    bash compare-callgraphs.sh"
