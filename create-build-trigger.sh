#!/bin/bash

# Create Cloud Build Trigger for Finsolvz Backend
set -e

PROJECT_ID="finsolvz-backend-dev"
REPO_OWNER="JavierZam"
REPO_NAME="finsolvz-backend"

echo "🔄 Creating Cloud Build Trigger"
echo "================================"

# Create trigger
gcloud builds triggers create github \
  --name="finsolvz-backend-main" \
  --repo-owner="$REPO_OWNER" \
  --repo-name="$REPO_NAME" \
  --branch-pattern="^main$" \
  --build-config="cloudbuild.yaml" \
  --description="Auto-deploy Finsolvz Backend on main branch push" \
  --substitutions="_SERVICE_NAME=finsolvz-backend,_REGION=asia-southeast2,_REPOSITORY=finsolvz-repo,_MEMORY=1Gi,_CPU=1,_MIN_INSTANCES=0,_MAX_INSTANCES=10"

echo "✅ Build trigger created successfully!"
echo ""

echo "📋 Trigger Details:"
gcloud builds triggers describe finsolvz-backend-main

echo ""
echo "🚀 Next Steps:"
echo "1. Push code to main branch"
echo "2. Monitor build: gcloud builds list --ongoing"
echo "3. Check logs: gcloud builds log [BUILD_ID]"
echo "4. View in console: https://console.cloud.google.com/cloud-build/builds"