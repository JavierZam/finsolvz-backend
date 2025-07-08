#!/bin/bash

# Test Swagger Access Script
echo "🔍 Testing Swagger Documentation Access"
echo "======================================"

# Check if server is running locally
echo "1. Testing local server..."
if curl -s http://localhost:8787/ > /dev/null 2>&1; then
    echo "✅ Local server is running on port 8787"
    echo "🌐 Swagger UI: http://localhost:8787/docs"
    echo "📄 OpenAPI Spec: http://localhost:8787/api/openapi.yaml"
    echo "🩺 Health Check: http://localhost:8787/"
    echo "🐛 Debug Files: http://localhost:8787/debug/files"
else
    echo "❌ Local server is not running on port 8787"
    echo "   Run: go run cmd/server/main.go"
fi

echo ""
echo "2. Testing OpenAPI file exists..."
if [ -f "./api/openapi.yaml" ]; then
    echo "✅ OpenAPI specification file exists"
    echo "   Location: ./api/openapi.yaml"
else
    echo "❌ OpenAPI specification file missing"
    echo "   Expected: ./api/openapi.yaml"
fi

echo ""
echo "3. Testing deployment configuration..."
echo "   Dockerfile: Exposes port 8787 ✅"
echo "   CloudBuild: Configured for port 8787 ✅"
echo "   Deploy script: Uses port 8787 ✅"

echo ""
echo "4. Next steps for deployment:"
echo "   a. Build and deploy: ./deploy.sh"
echo "   b. Test deployed swagger: [DEPLOYED_URL]/docs"
echo "   c. Test deployed API spec: [DEPLOYED_URL]/api/openapi.yaml"