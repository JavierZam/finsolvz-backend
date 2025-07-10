#!/bin/bash

# Script untuk disable Cloud Build triggers
# Jalankan setelah GitHub Actions sudah setup

set -e

echo "🔄 Disabling Cloud Build triggers..."

# Check if gcloud is configured
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | head -n 1 > /dev/null; then
    echo "❌ No active gcloud account found"
    echo "Please run: gcloud auth login"
    exit 1
fi

# Get current project
PROJECT_ID=$(gcloud config get-value project 2>/dev/null)
if [ -z "$PROJECT_ID" ]; then
    echo "❌ No project set"
    echo "Please run: gcloud config set project YOUR_PROJECT_ID"
    exit 1
fi

echo "📋 Project: $PROJECT_ID"

# List triggers
echo "📄 Current Cloud Build triggers:"
gcloud builds triggers list --format="table(name,status,github.name)"

echo
echo "🚫 Disabling triggers..."

# Get trigger names and disable them
TRIGGERS=$(gcloud builds triggers list --format="value(name)")

if [ -z "$TRIGGERS" ]; then
    echo "✅ No triggers found to disable"
else
    for trigger in $TRIGGERS; do
        echo "🔄 Disabling trigger: $trigger"
        gcloud builds triggers delete "$trigger" --quiet
        echo "✅ Trigger $trigger disabled"
    done
fi

echo
echo "🎉 Cloud Build triggers disabled successfully!"
echo
echo "📋 Next steps:"
echo "   1. Verify GitHub Actions is working"
echo "   2. Remove cloudbuild.yaml file (optional)"
echo "   3. Update documentation"
echo
echo "🔄 To re-enable if needed:"
echo "   gcloud builds triggers create github --repo-name=REPO --repo-owner=OWNER --branch-pattern=main --build-config=cloudbuild.yaml"