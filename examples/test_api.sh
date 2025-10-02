#!/bin/bash

# MinIO API Test Script
# This script demonstrates how to test the upload and download endpoints

BASE_URL="http://localhost:8080"
TEST_FILE="test-document.txt"

echo "=== MinIO API Test Script ==="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Create a test file if it doesn't exist
if [ ! -f "$TEST_FILE" ]; then
    echo "Creating test file..."
    echo "This is a test document for MinIO upload/download testing." > "$TEST_FILE"
    echo "Created at: $(date)" >> "$TEST_FILE"
fi

# Test 1: Health Check
echo -e "${YELLOW}Test 1: Health Check${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/health")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✓ Health check passed${NC}"
else
    echo -e "${RED}✗ Health check failed (HTTP $http_code)${NC}"
    exit 1
fi
echo ""

# Test 2: Upload file with original name
echo -e "${YELLOW}Test 2: Upload file (original name)${NC}"
echo "Uploading $TEST_FILE..."
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/upload" \
    -F "file=@$TEST_FILE")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✓ Upload successful${NC}"
    echo "Response: $body" | jq '.' 2>/dev/null || echo "$body"
else
    echo -e "${RED}✗ Upload failed (HTTP $http_code)${NC}"
    echo "Response: $body"
fi
echo ""

# Test 3: Upload file with custom name
echo -e "${YELLOW}Test 3: Upload file (custom name)${NC}"
CUSTOM_NAME="custom-$(date +%s).txt"
echo "Uploading with custom name: $CUSTOM_NAME"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/upload?objectName=$CUSTOM_NAME" \
    -F "file=@$TEST_FILE")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✓ Upload with custom name successful${NC}"
    echo "Response: $body" | jq '.' 2>/dev/null || echo "$body"
else
    echo -e "${RED}✗ Upload with custom name failed (HTTP $http_code)${NC}"
    echo "Response: $body"
fi
echo ""

# Test 4: Download file
echo -e "${YELLOW}Test 4: Download file${NC}"
DOWNLOAD_FILE="downloaded-$TEST_FILE"
echo "Downloading to $DOWNLOAD_FILE..."
http_code=$(curl -s -w "%{http_code}" "$BASE_URL/api/download?objectName=$TEST_FILE" \
    -o "$DOWNLOAD_FILE")

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✓ Download successful${NC}"
    echo "File saved to: $DOWNLOAD_FILE"
    echo "File content:"
    cat "$DOWNLOAD_FILE"
else
    echo -e "${RED}✗ Download failed (HTTP $http_code)${NC}"
fi
echo ""

# Test 5: Download custom named file
echo -e "${YELLOW}Test 5: Download custom named file${NC}"
DOWNLOAD_CUSTOM="downloaded-$CUSTOM_NAME"
echo "Downloading $CUSTOM_NAME to $DOWNLOAD_CUSTOM..."
http_code=$(curl -s -w "%{http_code}" "$BASE_URL/api/download?objectName=$CUSTOM_NAME" \
    -o "$DOWNLOAD_CUSTOM")

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✓ Download successful${NC}"
    echo "File saved to: $DOWNLOAD_CUSTOM"
else
    echo -e "${RED}✗ Download failed (HTTP $http_code)${NC}"
fi
echo ""

# Test 6: Download non-existent file (should fail)
echo -e "${YELLOW}Test 6: Download non-existent file (should fail)${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/download?objectName=nonexistent-file.txt")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "404" ]; then
    echo -e "${GREEN}✓ Correctly returned 404 for non-existent file${NC}"
    echo "Response: $body" | jq '.' 2>/dev/null || echo "$body"
else
    echo -e "${RED}✗ Expected 404, got HTTP $http_code${NC}"
    echo "Response: $body"
fi
echo ""

# Test 7: Upload without file (should fail)
echo -e "${YELLOW}Test 7: Upload without file (should fail)${NC}"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/upload")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "400" ]; then
    echo -e "${GREEN}✓ Correctly returned 400 for missing file${NC}"
    echo "Response: $body" | jq '.' 2>/dev/null || echo "$body"
else
    echo -e "${RED}✗ Expected 400, got HTTP $http_code${NC}"
    echo "Response: $body"
fi
echo ""

# Test 8: Get file info (headers)
echo -e "${YELLOW}Test 8: Get file info (headers only)${NC}"
echo "Getting headers for $TEST_FILE..."
curl -s -I "$BASE_URL/api/download?objectName=$TEST_FILE"
echo ""

# Cleanup
echo -e "${YELLOW}Cleanup${NC}"
read -p "Do you want to delete test files? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    rm -f "$DOWNLOAD_FILE" "$DOWNLOAD_CUSTOM" 2>/dev/null
    echo "Cleaned up downloaded files"
fi

echo ""
echo -e "${GREEN}=== Test completed ===${NC}"
