#!/bin/bash

# Setup Cloud Build Logs using gcloud commands (Windows compatible)
set -e

PROJECT_ID=$(gcloud config get-value project)
REGION="asia-southeast2"
BUCKET_NAME="${PROJECT_ID}-build-logs"

echo "📋 Setting up Cloud Build Logs using gcloud commands"
echo "===================================================="
echo "Project: $PROJECT_ID"
echo "Bucket: $BUCKET_NAME"
echo "Region: $REGION"
echo ""

# Check if Storage API is enabled
echo "🔧 Enabling Storage API..."
gcloud services enable storage.googleapis.com

# Create bucket using gcloud instead of gsutil
echo "🪣 Creating logs bucket using gcloud..."

# Check if bucket exists first
if gcloud storage buckets describe gs://$BUCKET_NAME >/dev/null 2>&1; then
    echo "✅ Bucket already exists: gs://$BUCKET_NAME"
else
    echo "Creating bucket..."
    
    # Create bucket using gcloud storage
    gcloud storage buckets create gs://$BUCKET_NAME \
        --location=$REGION \
        --uniform-bucket-level-access
    
    echo "✅ Bucket created: gs://$BUCKET_NAME"
fi

# Set bucket permissions for Cloud Build
echo "🔐 Setting bucket permissions..."

# Get Cloud Build service account
PROJECT_NUMBER=$(gcloud projects describe $PROJECT_ID --format="value(projectNumber)")
CLOUD_BUILD_SA="${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com"

echo "Cloud Build SA: $CLOUD_BUILD_SA"

# Grant permissions using gcloud instead of gsutil
gcloud storage buckets add-iam-policy-binding gs://$BUCKET_NAME \
    --member="serviceAccount:$CLOUD_BUILD_SA" \
    --role="roles/storage.admin"

echo "✅ Permissions granted"

# Create lifecycle policy using gcloud
echo "🗂️ Setting lifecycle policy..."

# Create lifecycle JSON
cat > lifecycle.json << EOF
{
  "rule": [
    {
      "action": {"type": "Delete"},
      "condition": {"age": 30}
    }
  ]
}
EOF

# Apply lifecycle policy
gcloud storage buckets update gs://$BUCKET_NAME \
    --lifecycle-file=lifecycle.json

rm lifecycle.json

echo "✅ Lifecycle policy set (30 days retention)"

# Create cloudbuild.yaml with logs bucket
echo "📝 Creating cloudbuild.yaml with logs bucket..."

cat > cloudbuild-with-bucket.yaml << EOF
# Cloud Build with Custom Logs Bucket
steps:
  # Build steps
  - name: 'gcr.io/cloud-builders/go'
    entrypoint: 'go'
    args: ['mod', 'download']
    
  - name: 'gcr.io/cloud-builders/go'
    entrypoint: 'go'
    args: ['test', './...']
    
  - name: 'gcr.io/cloud-builders/go'
    entrypoint: 'go'
    args: ['build', '-o', 'main', './cmd/server']
    env: ['CGO_ENABLED=0', 'GOOS=linux']

  - name: 'gcr.io/cloud-builders/docker'
    args: ['build', '-t', 'gcr.io/\$PROJECT_ID/finsolvz-backend:\$COMMIT_SHA', '.']
    
  - name: 'gcr.io/cloud-builders/docker'
    args: ['push', 'gcr.io/\$PROJECT_ID/finsolvz-backend:\$COMMIT_SHA']

  - name: 'gcr.io/cloud-builders/gcloud'
    args: [
      'run', 'deploy', 'finsolvz-backend',
      '--image', 'gcr.io/\$PROJECT_ID/finsolvz-backend:\$COMMIT_SHA',
      '--region', 'asia-southeast2',
      '--platform', 'managed',
      '--allow-unauthenticated',
      '--set-secrets', 'MONGO_URI=MONGO_URI:latest,JWT_SECRET=JWT_SECRET:latest',
      '--quiet'
    ]

images:
  - 'gcr.io/\$PROJECT_ID/finsolvz-backend:\$COMMIT_SHA'

timeout: '1200s'

# Use custom logs bucket
logsBucket: 'gs://$BUCKET_NAME'

options:
  machineType: 'E2_HIGHCPU_8'
  diskSizeGb: 20
  substitutionOption: 'ALLOW_LOOSE'
EOF

echo "✅ cloudbuild-with-bucket.yaml created"

echo ""
echo "🎯 SUMMARY"
echo "=========="
echo "✅ Logs bucket created: gs://$BUCKET_NAME"
echo "✅ Permissions configured"
echo "✅ Lifecycle policy set (30 days)"
echo "✅ cloudbuild-with-bucket.yaml generated"
echo ""
echo "📋 Next Steps:"
echo "1. Test with: gcloud builds submit --config=cloudbuild-with-bucket.yaml ."
echo "2. Or copy config to your main cloudbuild.yaml"
echo ""
echo "🧪 Verify bucket access:"
echo "   gcloud storage ls gs://$BUCKET_NAME"