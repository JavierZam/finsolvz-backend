#!/bin/bash

# Setup Cloud Build CI/CD untuk finsolvz-backend
set -e

echo "🚀 Setting up Cloud Build CI/CD for finsolvz-backend"
echo "=================================================="

# Get project info
PROJECT_ID=$(gcloud config get-value project)
PROJECT_NUMBER=$(gcloud projects describe $PROJECT_ID --format="value(projectNumber)")

if [ -z "$PROJECT_ID" ]; then
    echo "❌ Error: No project configured. Run 'gcloud config set project YOUR_PROJECT_ID'"
    exit 1
fi

echo "📋 Project ID: $PROJECT_ID"
echo "📋 Project Number: $PROJECT_NUMBER"
echo ""

# Step 1: Enable required APIs
echo "🔧 Enabling required APIs..."
gcloud services enable cloudbuild.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable containerregistry.googleapis.com
gcloud services enable storage.googleapis.com
gcloud services enable secretmanager.googleapis.com
echo "✅ APIs enabled"

# Step 2: Create secrets for environment variables
echo "🔐 Creating Secret Manager secrets..."

# Create MONGO_URI secret if not exists
if ! gcloud secrets describe MONGO_URI >/dev/null 2>&1; then
    echo "Creating MONGO_URI secret..."
    echo "Enter your MongoDB URI:"
    read -s MONGO_URI
    echo -n "$MONGO_URI" | gcloud secrets create MONGO_URI --data-file=-
    echo "✅ MONGO_URI secret created"
else
    echo "✅ MONGO_URI secret already exists"
fi

# Create JWT_SECRET secret if not exists
if ! gcloud secrets describe JWT_SECRET >/dev/null 2>&1; then
    echo "Creating JWT_SECRET secret..."
    echo "Enter your JWT secret:"
    read -s JWT_SECRET
    echo -n "$JWT_SECRET" | gcloud secrets create JWT_SECRET --data-file=-
    echo "✅ JWT_SECRET secret created"
else
    echo "✅ JWT_SECRET secret already exists"
fi

# Step 3: Grant Cloud Build service account permissions
echo "🔑 Granting Cloud Build service account permissions..."

CLOUD_BUILD_SA="${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com"
echo "Cloud Build SA: $CLOUD_BUILD_SA"

# Grant necessary roles
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:$CLOUD_BUILD_SA" \
    --role="roles/run.admin"

gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:$CLOUD_BUILD_SA" \
    --role="roles/storage.admin"

gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:$CLOUD_BUILD_SA" \
    --role="roles/secretmanager.secretAccessor"

gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:$CLOUD_BUILD_SA" \
    --role="roles/iam.serviceAccountUser"

echo "✅ Permissions granted"

# Step 4: Test the build
echo "🧪 Testing Cloud Build configuration..."
echo "Running: gcloud builds submit --config=cloudbuild.yaml ."
echo ""

read -p "Do you want to test the build now? (y/n): " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    gcloud builds submit --config=cloudbuild.yaml .
    echo "✅ Build test completed"
else
    echo "⏭️ Build test skipped"
fi

# Step 5: Setup GitHub trigger (optional)
echo ""
echo "🔗 GitHub Repository Setup"
echo "========================="
echo "To setup automatic builds from GitHub:"
echo ""
echo "1. Go to Cloud Build Triggers: https://console.cloud.google.com/cloud-build/triggers"
echo "2. Click 'Create Trigger'"
echo "3. Connect your GitHub repository"
echo "4. Set trigger configuration:"
echo "   - Name: finsolvz-backend-main"
echo "   - Event: Push to branch"
echo "   - Branch: ^main$"
echo "   - Configuration: Cloud Build configuration file"
echo "   - Location: /cloudbuild.yaml"
echo ""

echo "🎯 SETUP COMPLETE!"
echo "=================="
echo "✅ APIs enabled"
echo "✅ Secrets created (MONGO_URI, JWT_SECRET)"
echo "✅ Service account permissions configured"
echo "✅ cloudbuild.yaml ready"
echo ""
echo "📋 Next Steps:"
echo "1. Setup GitHub trigger (see instructions above)"
echo "2. Test manual build: gcloud builds submit --config=cloudbuild.yaml ."
echo "3. Push to GitHub to trigger automatic builds"
echo ""
echo "🌐 Monitor builds: https://console.cloud.google.com/cloud-build/builds"
echo "🚀 View deployments: https://console.cloud.google.com/run"