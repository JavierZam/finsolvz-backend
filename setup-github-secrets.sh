#!/bin/bash

# Setup GitHub Secrets untuk GitHub Actions
# Jalankan script ini di terminal setelah menginstall GitHub CLI

set -e

echo "🔐 Setting up GitHub Secrets for GitHub Actions..."

# Check if gh CLI is installed
if ! command -v gh &> /dev/null; then
    echo "❌ GitHub CLI not found. Please install it first:"
    echo "   https://cli.github.com/"
    exit 1
fi

# Check if logged in to GitHub
if ! gh auth status &> /dev/null; then
    echo "🔑 Please login to GitHub first:"
    gh auth login
fi

# Get current repository
REPO=$(gh repo view --json nameWithOwner -q .nameWithOwner)
echo "📋 Repository: $REPO"

# Function to set secret
set_secret() {
    local secret_name=$1
    local secret_description=$2
    
    echo -n "🔑 Enter $secret_description: "
    read -s secret_value
    echo
    
    if [ -n "$secret_value" ]; then
        echo "$secret_value" | gh secret set "$secret_name"
        echo "✅ Secret $secret_name set successfully"
    else
        echo "⚠️  Skipping empty secret $secret_name"
    fi
}

echo
echo "📝 Please provide the following secrets:"
echo

# GCP Project ID
set_secret "GCP_PROJECT_ID" "GCP Project ID"

# GCP Service Account Key
echo
echo "🔑 For GCP_SA_KEY, you need a service account JSON key with these permissions:"
echo "   - Cloud Run Admin"
echo "   - Storage Admin" 
echo "   - Artifact Registry Admin"
echo "   - Service Account User"
echo
echo "💡 To create service account key:"
echo "   1. Go to IAM & Admin > Service Accounts in GCP Console"
echo "   2. Create or select service account"
echo "   3. Add the required roles above"
echo "   4. Create JSON key and copy entire content"
echo
echo -n "🔑 Paste the entire JSON key content: "
read -s gcp_sa_key
echo

if [ -n "$gcp_sa_key" ]; then
    echo "$gcp_sa_key" | gh secret set "GCP_SA_KEY"
    echo "✅ Secret GCP_SA_KEY set successfully"
else
    echo "❌ GCP_SA_KEY is required for deployment"
    exit 1
fi

echo
echo "🎉 GitHub Secrets setup completed!"
echo
echo "📋 Secrets configured:"
echo "   - GCP_PROJECT_ID"
echo "   - GCP_SA_KEY"
echo
echo "🚀 Next steps:"
echo "   1. Disable Cloud Build trigger in GCP Console"
echo "   2. Push code to trigger GitHub Actions"
echo "   3. Check Actions tab in GitHub for deployment status"
echo
echo "🔗 Useful commands:"
echo "   gh secret list                    # List all secrets"
echo "   gh workflow list                  # List workflows"
echo "   gh run list                       # List workflow runs"
echo "   gh run watch                      # Watch current workflow run"