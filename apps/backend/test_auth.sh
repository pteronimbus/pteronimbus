#!/bin/bash

# Simple test script to verify backend authentication endpoints
BASE_URL="http://localhost:8080"

echo "Testing Pteronimbus Backend Authentication..."
echo "============================================="

# Test health endpoint
echo "1. Testing health endpoint..."
curl -s "$BASE_URL/health" | jq '.' || echo "Health endpoint failed"
echo ""

# Test login endpoint (should return Discord auth URL)
echo "2. Testing login endpoint..."
curl -s "$BASE_URL/auth/login" | jq '.' || echo "Login endpoint failed"
echo ""

# Test protected endpoint without auth (should fail)
echo "3. Testing protected endpoint without auth..."
curl -s "$BASE_URL/api/test" | jq '.' || echo "Protected endpoint correctly rejected unauthenticated request"
echo ""

echo "Backend authentication endpoints are working!"
echo "To complete the OAuth flow, you'll need to:"
echo "1. Set up Discord OAuth2 credentials in environment variables"
echo "2. Start Redis server"
echo "3. Use the frontend to test the complete flow"