#!/bin/bash

# Finsolvz Backend Deployment Script
set -e

PROJECT_ID="finsolvz-backend-dev"
REGION="asia-southeast2"
SERVICE_NAME="finsolvz-backend"

echo "🚀 Deploying Finsolvz Backend to GCP Cloud Run"
echo "================================================"

# 1. Set project
echo "📍 Setting GCP project: $PROJECT_ID"
gcloud config set project $PROJECT_ID

# 2. Build and push image
echo "🏗️  Building Docker image..."
gcloud builds submit --tag gcr.io/$PROJECT_ID/$SERVICE_NAME

# 3. Deploy to Cloud Run
echo "🚀 Deploying to Cloud Run..."
gcloud run deploy $SERVICE_NAME \
  --image gcr.io/$PROJECT_ID/$SERVICE_NAME \
  --region $REGION \
  --platform managed \
  --allow-unauthenticated \
  --port 8787 \
  --memory 512Mi \
  --cpu 1 \
  --max-instances 10 \
  --set-env-vars="APP_ENV=development,PORT=8787,GREETING=✨ Finsolvz Backend API ✨" \
  --set-secrets="MONGO_URI=MONGO_URI:latest,JWT_SECRET=JWT_SECRET:latest"

# 4. Get service URL
echo "✅ Deployment complete!"
echo "🌐 Service URL:"
gcloud run services describe $SERVICE_NAME --region $REGION --format "value(status.url)"

echo ""
echo "🔗 Test endpoints:"
echo "Health check: [SERVICE_URL]/"
echo "API docs: [SERVICE_URL]/api (if you add swagger serve)"
echo ""
echo "🔐 Next steps:"
echo "1. Test health endpoint"
echo "2. Create admin user: go run create_admin.go (with production MONGO_URI)"
echo "3. Test authentication endpoints"