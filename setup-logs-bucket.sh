#!/bin/bash

# Setup Cloud Build Logs Bucket
set -e

PROJECT_ID=$(gcloud config get-value project)
REGION="asia-southeast2"
BUCKET_NAME="${PROJECT_ID}-build-logs"

echo "📋 Setting up Cloud Build Logs Bucket"
echo "====================================="
echo "Project: $PROJECT_ID"
echo "Bucket: $BUCKET_NAME"
echo "Region: $REGION"
echo ""

# Check if bucket exists
if gsutil ls -b gs://$BUCKET_NAME >/dev/null 2>&1; then
    echo "✅ Bucket already exists: gs://$BUCKET_NAME"
else
    echo "🪣 Creating logs bucket..."
    
    # Create bucket
    gsutil mb -l $REGION gs://$BUCKET_NAME
    
    echo "✅ Bucket created: gs://$BUCKET_NAME"
fi

# Set bucket permissions for Cloud Build
echo "🔐 Setting bucket permissions..."

# Get Cloud Build service account
PROJECT_NUMBER=$(gcloud projects describe $PROJECT_ID --format="value(projectNumber)")
CLOUD_BUILD_SA="${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com"

echo "Cloud Build SA: $CLOUD_BUILD_SA"

# Grant permissions
gsutil iam ch serviceAccount:$CLOUD_BUILD_SA:roles/storage.admin gs://$BUCKET_NAME

echo "✅ Permissions granted"

# Set lifecycle policy to auto-delete old logs
echo "🗂️ Setting lifecycle policy..."

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

gsutil lifecycle set lifecycle.json gs://$BUCKET_NAME
rm lifecycle.json

echo "✅ Lifecycle policy set (30 days retention)"

# Create cloudbuild.yaml with logs bucket
echo "📝 Creating cloudbuild.yaml with logs bucket..."

cat > cloudbuild-with-bucket.yaml << 'EOF'
steps:
  # Basic build steps
  - name: 'gcr.io/cloud-builders/go'
    entrypoint: 'go'
    args: ['mod', 'download']
    
  - name: 'gcr.io/cloud-builders/go'
    entrypoint: 'go'
    args: ['build', '-o', 'main', './cmd/server']

  - name: 'gcr.io/cloud-builders/docker'
    args: [
      'build',
      '-t', 'gcr.io/$PROJECT_ID/${_SERVICE_NAME}:$COMMIT_SHA',
      '.'
    ]

  - name: 'gcr.io/cloud-builders/docker'
    args: [
      'push',
      'gcr.io/$PROJECT_ID/${_SERVICE_NAME}:$COMMIT_SHA'
    ]

  - name: 'gcr.io/cloud-builders/gcloud'
    args: [
      'run', 'deploy', '${_SERVICE_NAME}',
      '--image', 'gcr.io/$PROJECT_ID/${_SERVICE_NAME}:$COMMIT_SHA',
      '--region', '${_REGION}',
      '--platform', 'managed',
      '--allow-unauthenticated',
      '--set-secrets', 'MONGO_URI=MONGO_URI:latest,JWT_SECRET=JWT_SECRET:latest'
    ]

images:
  - 'gcr.io/$PROJECT_ID/${_SERVICE_NAME}:$COMMIT_SHA'

timeout: '1200s'

# Use custom logs bucket
logsBucket: 'gs://PROJECT_ID_PLACEHOLDER-build-logs'

options:
  machineType: 'E2_HIGHCPU_8'
  diskSizeGb: 20
  substitutionOption: 'ALLOW_LOOSE'

substitutions:
  _SERVICE_NAME: 'finsolvz-backend'
  _REGION: 'asia-southeast2'
EOF

# Replace placeholder with actual project ID
sed "s/PROJECT_ID_PLACEHOLDER/$PROJECT_ID/g" cloudbuild-with-bucket.yaml > cloudbuild-bucket-ready.yaml
rm cloudbuild-with-bucket.yaml

echo "✅ cloudbuild-bucket-ready.yaml created"

echo ""
echo "🎯 SUMMARY"
echo "=========="
echo "✅ Logs bucket created: gs://$BUCKET_NAME"
echo "✅ Permissions configured"
echo "✅ Lifecycle policy set (30 days)"
echo "✅ cloudbuild-bucket-ready.yaml generated"
echo ""
echo "📋 Next Steps:"
echo "1. Use cloudbuild-bucket-ready.yaml instead of cloudbuild.yaml"
echo "2. Or add this to your existing cloudbuild.yaml:"
echo "   logsBucket: 'gs://$BUCKET_NAME'"
echo ""
echo "🧪 Test the setup:"
echo "   gcloud builds submit --config=cloudbuild-bucket-ready.yaml ."