#!/bin/bash

# Debricked Java Callgraph Sample - Prerequisites Setup Script
# This script checks and helps install required dependencies:
# - Java 11 or higher
# - Maven
# - Debricked CLI

set -e

SAMPLE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SAMPLE_DIR"

echo "================================"
echo "Prerequisites Check & Setup"
echo "================================"
echo ""

# Detect OS
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    OS="linux"
    PKG_MANAGER=""
    if command -v apt-get &> /dev/null; then
        PKG_MANAGER="apt-get"
    elif command -v yum &> /dev/null; then
        PKG_MANAGER="yum"
    fi
elif [[ "$OSTYPE" == "darwin"* ]]; then
    OS="macos"
    PKG_MANAGER="brew"
else
    OS="unknown"
fi

echo "Detected OS: $OS"
echo ""

# Check Java
echo "[1/3] Checking Java..."
if command -v java &> /dev/null; then
    JAVA_VERSION=$(java -version 2>&1 | grep -oP 'version "\K[0-9]+' | head -1)
    echo "✓ Java is installed (version: $JAVA_VERSION)"
    if [ "$JAVA_VERSION" -lt 11 ]; then
        echo "⚠ WARNING: Java version is less than 11. Callgraph requires Java 11+."
    fi
else
    echo "✗ Java not found"
    echo ""
    if [ "$OS" == "linux" ]; then
        echo "To install Java 11+, run:"
        if [ "$PKG_MANAGER" == "apt-get" ]; then
            echo "  sudo apt-get update && sudo apt-get install -y openjdk-11-jdk"
        elif [ "$PKG_MANAGER" == "yum" ]; then
            echo "  sudo yum install -y java-11-openjdk-devel"
        fi
    elif [ "$OS" == "macos" ]; then
        echo "To install Java 11+, run:"
        echo "  brew install openjdk@11"
        echo "Then add it to your PATH:"
        echo "  sudo ln -sfn /opt/homebrew/opt/openjdk@11/libexec/openjdk.jdk /Library/Java/JavaVirtualMachines/openjdk-11.jdk"
    else
        echo "Please visit https://www.oracle.com/java/technologies/downloads/ to install Java 11+"
    fi
    echo ""
fi

# Check Maven
echo "[2/3] Checking Maven..."
if command -v mvn &> /dev/null; then
    MVN_VERSION=$(mvn -v 2>&1 | grep "Apache Maven" | awk '{print $3}')
    echo "✓ Maven is installed (version: $MVN_VERSION)"
else
    echo "✗ Maven not found"
    echo ""
    if [ "$OS" == "linux" ]; then
        echo "To install Maven, run:"
        if [ "$PKG_MANAGER" == "apt-get" ]; then
            echo "  sudo apt-get update && sudo apt-get install -y maven"
        elif [ "$PKG_MANAGER" == "yum" ]; then
            echo "  sudo yum install -y maven"
        fi
    elif [ "$OS" == "macos" ]; then
        echo "To install Maven, run:"
        echo "  brew install maven"
    else
        echo "Please visit https://maven.apache.org/download.cgi to install Maven"
    fi
    echo ""
fi

# Check Debricked CLI
echo "[3/3] Checking Debricked CLI..."
if command -v debricked &> /dev/null; then
    DEBRICKED_VERSION=$(debricked version 2>&1 || echo "unknown")
    echo "✓ Debricked CLI is installed (version: $DEBRICKED_VERSION)"
else
    echo "✗ Debricked CLI not found"
    echo ""
    echo "To install Debricked CLI, visit:"
    echo "  https://github.com/debricked/cli/releases"
    echo ""
    echo "Or if building from this repo:"
    echo "  cd /home/dritthi/projects/debricked-projects/cli"
    echo "  go build -o debricked ./cmd/debricked"
    echo "  # Then add to PATH or use ./debricked directly"
    echo ""
fi

echo ""
echo "================================"
echo "Setup Check Complete"
echo "================================"
echo ""
echo "Next steps:"
echo "1. Install any missing prerequisites from instructions above"
echo "2. Run: bash run-callgraph.sh"
echo ""

