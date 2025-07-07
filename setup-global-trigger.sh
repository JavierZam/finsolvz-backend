#!/bin/bash

# Setup GitHub trigger in global region
set -e

echo "🔧 Setting up GitHub trigger in global region"
echo "=============================================="

PROJECT_ID=$(gcloud config get-value project)
echo "📋 Project: $PROJECT_ID"
echo ""

# Step 1: Enable Cloud Build API
echo "🔧 Ensuring Cloud Build API is enabled..."
gcloud services enable cloudbuild.googleapis.com

# Step 2: Check GitHub App installation
echo "🔍 Checking GitHub App installation..."
echo ""
echo "📋 IMPORTANT: You need to install GitHub App first"
echo "================================================"
echo "1. Go to: https://github.com/apps/google-cloud-build"
echo "2. Click 'Install' or 'Configure'"
echo "3. Select your repository: JavierZam/finsolvz-backend"
echo "4. Grant necessary permissions"
echo ""
echo "Have you installed the GitHub App? (y/n):"
read -r GITHUB_APP_INSTALLED

if [[ ! $GITHUB_APP_INSTALLED =~ ^[Yy]$ ]]; then
    echo "❌ Please install GitHub App first, then run this script again"
    exit 1
fi

# Step 3: Create trigger using GitHub App
echo ""
echo "🚀 Creating trigger with GitHub App..."

TRIGGER_NAME="finsolvz-backend-main-global"

# Check if trigger exists
if gcloud builds triggers describe $TRIGGER_NAME --region=global >/dev/null 2>&1; then
    echo "⚠️  Trigger already exists. Deleting..."
    gcloud builds triggers delete $TRIGGER_NAME --region=global --quiet
fi

# Create trigger
gcloud builds triggers create github \
    --region=global \
    --name=$TRIGGER_NAME \
    --repo-name=finsolvz-backend \
    --repo-owner=JavierZam \
    --branch-pattern="^main$" \
    --build-config=cloudbuild.yaml \
    --description="Auto-deploy finsolvz-backend on main branch push - Global region"

echo "✅ Trigger created successfully!"

echo ""
echo "🎯 SUCCESS!"
echo "==========="
echo "✅ Trigger name: $TRIGGER_NAME"
echo "✅ Region: global"
echo "✅ Repository: JavierZam/finsolvz-backend"
echo "✅ Branch: main"
echo ""
echo "📋 Test trigger:"
echo "git push origin main"
echo ""
echo "🔗 Monitor builds: https://console.cloud.google.com/cloud-build/builds"
echo "🔗 Manage triggers: https://console.cloud.google.com/cloud-build/triggers"