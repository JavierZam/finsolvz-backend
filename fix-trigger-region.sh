#!/bin/bash

# Fix Cloud Build trigger region quota issue
set -e

echo "🔧 Fixing Cloud Build trigger region quota issue"
echo "==============================================="

PROJECT_ID=$(gcloud config get-value project)
TRIGGER_NAME="finsolvz-backend-trigger"
CURRENT_REGION="asia-southeast2"
NEW_REGION="global"

echo "📋 Project: $PROJECT_ID"
echo "📋 Current region: $CURRENT_REGION"
echo "📋 Target region: $NEW_REGION"
echo ""

# Check current trigger
echo "🔍 Checking current trigger..."
if gcloud builds triggers describe $TRIGGER_NAME --region=$CURRENT_REGION >/dev/null 2>&1; then
    echo "✅ Found trigger in $CURRENT_REGION"
    
    # Get trigger details
    echo "📋 Getting trigger configuration..."
    REPO_INFO=$(gcloud builds triggers describe $TRIGGER_NAME --region=$CURRENT_REGION --format="value(github.name,github.owner)")
    
    if [ -z "$REPO_INFO" ]; then
        echo "❌ Could not get repository info from trigger"
        exit 1
    fi
    
    echo "✅ Repository info retrieved"
    
    # Delete current trigger
    echo "🗑️  Deleting trigger from $CURRENT_REGION..."
    gcloud builds triggers delete $TRIGGER_NAME --region=$CURRENT_REGION --quiet
    
    echo "✅ Trigger deleted from $CURRENT_REGION"
else
    echo "❌ No trigger found in $CURRENT_REGION"
    echo "Let's create a new one in $NEW_REGION"
fi

# Create new trigger in global region
echo ""
echo "🚀 Creating new trigger in $NEW_REGION..."

# Manual input for repository details
echo "📝 Enter repository details:"
echo "GitHub username/organization:"
read -r GITHUB_OWNER
echo "Repository name:"
read -r REPO_NAME

# Create trigger in global region
gcloud builds triggers create github \
    --region=$NEW_REGION \
    --name=$TRIGGER_NAME \
    --repo-name=$REPO_NAME \
    --repo-owner=$GITHUB_OWNER \
    --branch-pattern="^main$" \
    --build-config=cloudbuild.yaml \
    --description="Auto-deploy finsolvz-backend on main branch push"

echo "✅ Trigger created in $NEW_REGION"

echo ""
echo "🎯 SUCCESS!"
echo "==========="
echo "✅ Trigger moved from $CURRENT_REGION to $NEW_REGION"
echo "✅ Should resolve quota restrictions"
echo ""
echo "📋 Test with:"
echo "git push origin main"
echo ""
echo "🔗 Monitor: https://console.cloud.google.com/cloud-build/builds"