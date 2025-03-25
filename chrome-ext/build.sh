#!/bin/bash

# Exit on error
set -e

# Configuration
OUTPUT_DIR="dist"
ZIP_NAME="chrome-extension.zip"

# Clean output directory
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# Copy files
cp -r assets scripts styles manifest.json "$OUTPUT_DIR/"

# Create zip for Chrome Web Store
cd "$OUTPUT_DIR"
zip -r "../$ZIP_NAME" ./*

echo "Build complete: $ZIP_NAME created" 