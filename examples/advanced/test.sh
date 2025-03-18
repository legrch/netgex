#!/bin/bash

# Test script for the V2 example

# Set the base URL
BASE_URL="http://localhost:8080"

# Function to test an endpoint
test_endpoint() {
  local endpoint=$1
  local expected_status=$2
  local description=$3

  echo "Testing $description..."
  response=$(curl -s -o /dev/null -w "%{http_code}" $endpoint)
  
  if [ "$response" -eq "$expected_status" ]; then
    echo "✅ $description: OK (Status: $response)"
  else
    echo "❌ $description: FAILED (Expected: $expected_status, Got: $response)"
    exit 1
  fi
}

# Function to test an endpoint and check the response body
test_endpoint_body() {
  local endpoint=$1
  local expected_status=$2
  local expected_content=$3
  local description=$4

  echo "Testing $description..."
  response_status=$(curl -s -o response.txt -w "%{http_code}" $endpoint)
  
  if [ "$response_status" -ne "$expected_status" ]; then
    echo "❌ $description: FAILED (Expected status: $expected_status, Got: $response_status)"
    exit 1
  fi

  if grep -q "$expected_content" response.txt; then
    echo "✅ $description: OK (Status: $response_status, Content contains: $expected_content)"
  else
    echo "❌ $description: FAILED (Content does not contain: $expected_content)"
    cat response.txt
    exit 1
  fi
}

# Wait for the server to start
echo "Waiting for the server to start..."
sleep 2

# Test HTTP endpoint with path parameter
test_endpoint_body "$BASE_URL/api/v1/hello/World" 200 "Hello, World!" "HTTP endpoint with path parameter"

# Test HTTP endpoint with query parameter
test_endpoint_body "$BASE_URL/api/v1/hello?name=World" 200 "Hello, World!" "HTTP endpoint with query parameter"

# Test health endpoint
test_endpoint_body "$BASE_URL/health" 200 "OK" "Health endpoint"

# Test Swagger UI
test_endpoint_body "$BASE_URL/swagger/index.html" 200 "Swagger UI" "Swagger UI"

echo "All tests passed! 🎉" 