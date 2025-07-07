#!/bin/bash

# Setup GitHub Trigger for Cloud Build CI/CD
set -e

echo "🔗 Setting up GitHub Trigger for Cloud Build"
echo "============================================="

# Get project info
PROJECT_ID=$(gcloud config get-value project)
if [ -z "$PROJECT_ID" ]; then
    echo "❌ Error: No project configured. Run 'gcloud config set project YOUR_PROJECT_ID'"
    exit 1
fi

echo "📋 Project: $PROJECT_ID"
echo ""

# Check connection type
echo "🔍 Checking GitHub connections..."

# Check for Developer Connect (generation 2)
DEV_CONNECTIONS=$(gcloud builds connections list --region=asia-southeast2 --format="value(name)" 2>/dev/null || echo "")
# Check for GitHub App (generation 1)
GITHUB_APP_REPOS=$(gcloud source repos list --format="value(name)" 2>/dev/null | grep -E "(github|bitbucket)" || echo "")

echo ""
if [ -n "$DEV_CONNECTIONS" ]; then
    echo "✅ Developer Connect connection found (Generation 2)"
    USE_DEV_CONNECT=true
    REGION="asia-southeast2"
elif [ -n "$GITHUB_APP_REPOS" ]; then
    echo "✅ GitHub App connection found (Generation 1)"
    USE_DEV_CONNECT=false
    REGION="global"
else
    echo "❌ No GitHub connection found."
    echo ""
    echo "📋 Choose connection method:"
    echo "1. Developer Connect (Generation 2) - Recommended for new projects"
    echo "2. GitHub App (Generation 1) - Legacy method"
    echo ""
    echo "For Developer Connect:"
    echo "- Go to: https://console.cloud.google.com/cloud-build/triggers"
    echo "- Use region: asia-southeast2"
    echo "- Connect via Developer Connect"
    echo ""
    echo "For GitHub App:"
    echo "- Go to: https://console.cloud.google.com/cloud-build/triggers"
    echo "- Use region: global"
    echo "- Connect via GitHub App"
    echo ""
    echo "After connecting, run this script again."
    exit 1
fi

# Get GitHub repository details
echo ""
echo "📝 GitHub Repository Setup"
echo "========================="
echo "Enter your GitHub username/organization:"
read -r GITHUB_OWNER
echo "Enter your repository name (e.g., finsolvz-backend):"
read -r REPO_NAME

# Create trigger
echo ""
echo "🚀 Creating Cloud Build trigger..."

TRIGGER_NAME="finsolvz-backend-main"

# Check if trigger already exists
if gcloud builds triggers describe $TRIGGER_NAME --region=$REGION >/dev/null 2>&1; then
    echo "⚠️  Trigger '$TRIGGER_NAME' already exists"
    echo "Do you want to delete and recreate it? (y/n):"
    read -r RECREATE
    if [[ $RECREATE =~ ^[Yy]$ ]]; then
        echo "🗑️  Deleting existing trigger..."
        gcloud builds triggers delete $TRIGGER_NAME --region=$REGION --quiet
    else
        echo "❌ Keeping existing trigger. Exiting."
        exit 0
    fi
fi

# Create the trigger based on connection type
if [ "$USE_DEV_CONNECT" = true ]; then
    echo "🔧 Creating trigger with Developer Connect..."
    
    # List available connections
    echo "Available connections:"
    gcloud builds connections list --region=$REGION --format="table(name,installationState)"
    
    echo "Enter connection name from above:"
    read -r CONNECTION_NAME
    
    # Create trigger with Developer Connect
    gcloud builds triggers create github \
        --region=$REGION \
        --name=$TRIGGER_NAME \
        --repository="projects/$PROJECT_ID/locations/$REGION/connections/$CONNECTION_NAME/repositories/$REPO_NAME" \
        --branch-pattern="^main$" \
        --build-config=cloudbuild.yaml \
        --description="Auto-deploy finsolvz-backend on main branch push"
else
    echo "🔧 Creating trigger with GitHub App..."
    
    # Create trigger with GitHub App
    gcloud builds triggers create github \
        --region=$REGION \
        --name=$TRIGGER_NAME \
        --repo-name=$REPO_NAME \
        --repo-owner=$GITHUB_OWNER \
        --branch-pattern="^main$" \
        --build-config=cloudbuild.yaml \
        --description="Auto-deploy finsolvz-backend on main branch push"
fi

echo "✅ Trigger created successfully!"

# Create PR trigger (optional)
echo ""
echo "📝 Create Pull Request trigger for testing? (y/n):"
read -r CREATE_PR_TRIGGER

if [[ $CREATE_PR_TRIGGER =~ ^[Yy]$ ]]; then
    echo "🔄 Creating PR trigger..."
    
    # Create PR-specific cloudbuild config
    cat > cloudbuild-pr.yaml << 'EOF'
steps:
  # Step 1: Download dependencies
  - name: 'golang:1.23-alpine'
    entrypoint: 'go'
    args: ['mod', 'download']

  # Step 2: Run tests
  - name: 'golang:1.23-alpine'
    entrypoint: 'go'
    args: ['test', './...']

  # Step 3: Build application (no deploy)
  - name: 'golang:1.23-alpine'
    entrypoint: 'go'
    args: ['build', '-o', 'main', './cmd/server']
    env:
      - 'CGO_ENABLED=0'
      - 'GOOS=linux'

  # Step 4: Build Docker image (no push)
  - name: 'gcr.io/cloud-builders/docker'
    args: [
      'build',
      '-t', 'gcr.io/$PROJECT_ID/finsolvz-backend:pr-$_PR_NUMBER',
      '.'
    ]

timeout: '600s'

options:
  machineType: 'E2_HIGHCPU_8'
  diskSizeGb: 20
EOF

    if [ "$USE_DEV_CONNECT" = true ]; then
        gcloud builds triggers create github \
            --region=$REGION \
            --name="finsolvz-backend-pr" \
            --repository="projects/$PROJECT_ID/locations/$REGION/connections/$CONNECTION_NAME/repositories/$REPO_NAME" \
            --pull-request-pattern=".*" \
            --build-config=cloudbuild-pr.yaml \
            --description="Test finsolvz-backend on pull requests"
    else
        gcloud builds triggers create github \
            --region=$REGION \
            --name="finsolvz-backend-pr" \
            --repo-name=$REPO_NAME \
            --repo-owner=$GITHUB_OWNER \
            --pull-request-pattern=".*" \
            --build-config=cloudbuild-pr.yaml \
            --description="Test finsolvz-backend on pull requests"
    fi
    
    echo "✅ PR trigger created!"
fi

echo ""
echo "🎯 SETUP COMPLETE!"
echo "=================="
echo "✅ GitHub repository connected"
echo "✅ Main branch trigger created"
echo "$([ "$CREATE_PR_TRIGGER" = "y" ] && echo "✅ PR trigger created" || echo "⏭️  PR trigger skipped")"
echo ""
echo "📋 Trigger Details:"
echo "- Name: $TRIGGER_NAME"
echo "- Repository: $GITHUB_OWNER/$REPO_NAME"
echo "- Branch: main"
echo "- Config: cloudbuild.yaml"
echo ""
echo "🧪 Next Steps:"
echo "1. Push a commit to main branch"
echo "2. Monitor build: https://console.cloud.google.com/cloud-build/builds"
echo "3. Check deployment: https://console.cloud.google.com/run"
echo ""
echo "🔗 Manage triggers: https://console.cloud.google.com/cloud-build/triggers"