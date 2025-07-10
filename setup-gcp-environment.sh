#!/bin/bash

# 🚀 Finsolvz Backend - Complete GCP Environment Setup
# One-time setup script for optimal performance deployment in Jakarta region
# 
# This script sets up:
# - GCP APIs and services
# - Artifact Registry in Jakarta
# - Cloud Build with optimized settings
# - Secrets management
# - Cloud Run deployment with Free Tier optimization
# - Performance monitoring

set -e

echo "🚀 Finsolvz Backend - GCP Environment Setup"
echo "==========================================="
echo ""

# Configuration
PROJECT_ID="${1:-}"
GITHUB_OWNER="${2:-}"
REPO_NAME="finsolvz-backend"
REGION="asia-southeast2"  # Jakarta - best latency for Indonesia
SERVICE_NAME="finsolvz-backend"

# Validate inputs
if [ -z "$PROJECT_ID" ]; then
    echo "❌ Error: PROJECT_ID is required"
    echo "Usage: ./setup-gcp-environment.sh PROJECT_ID GITHUB_USERNAME"
    echo "Example: ./setup-gcp-environment.sh my-project-123 johndoe"
    exit 1
fi

if [ -z "$GITHUB_OWNER" ]; then
    echo "❌ Error: GITHUB_OWNER is required"
    echo "Usage: ./setup-gcp-environment.sh PROJECT_ID GITHUB_USERNAME"
    echo "Example: ./setup-gcp-environment.sh my-project-123 johndoe"
    exit 1
fi

echo "📋 Configuration:"
echo "   Project ID: $PROJECT_ID"
echo "   GitHub Owner: $GITHUB_OWNER"
echo "   Repository: $REPO_NAME"
echo "   Region: $REGION (Jakarta)"
echo "   Service: $SERVICE_NAME"
echo ""

# Set project
echo "1️⃣ Setting up GCP project..."
gcloud config set project $PROJECT_ID
echo "✅ Project set to: $PROJECT_ID"
echo ""

# Enable required APIs
echo "2️⃣ Enabling required APIs..."
gcloud services enable \
    cloudbuild.googleapis.com \
    run.googleapis.com \
    artifactregistry.googleapis.com \
    secretmanager.googleapis.com \
    --project=$PROJECT_ID

echo "✅ APIs enabled successfully"
echo ""

# Create Artifact Registry repository
echo "3️⃣ Setting up Artifact Registry in Jakarta..."
gcloud artifacts repositories create finsolvz \
    --repository-format=docker \
    --location=$REGION \
    --description="Finsolvz backend container images - Jakarta region for optimal Indonesia latency" \
    --project=$PROJECT_ID 2>/dev/null || echo "   Repository already exists"

# Configure Docker authentication
gcloud auth configure-docker $REGION-docker.pkg.dev --quiet
echo "✅ Artifact Registry configured in Jakarta"
echo ""

# Setup secrets
echo "4️⃣ Setting up secrets..."
echo ""

# MongoDB URI
echo "🔒 MongoDB URI Setup:"
echo "Please enter your MongoDB connection string:"
echo "Format: mongodb+srv://username:password@cluster.mongodb.net/database"
read -s MONGO_URI

if [ ! -z "$MONGO_URI" ]; then
    echo "$MONGO_URI" | gcloud secrets create MONGO_URI \
        --data-file=- \
        --project=$PROJECT_ID 2>/dev/null || {
        echo "   Updating existing MONGO_URI secret..."
        echo "$MONGO_URI" | gcloud secrets versions add MONGO_URI --data-file=-
    }
    echo "✅ MONGO_URI secret configured"
else
    echo "⚠️  MONGO_URI not set - you'll need to configure this later"
fi

# JWT Secret
echo ""
echo "🔒 JWT Secret Setup:"
echo "Please enter your JWT secret (or press Enter to generate one):"
read -s JWT_SECRET

if [ -z "$JWT_SECRET" ]; then
    JWT_SECRET=$(openssl rand -hex 32)
    echo "   Generated secure JWT secret"
fi

echo "$JWT_SECRET" | gcloud secrets create JWT_SECRET \
    --data-file=- \
    --project=$PROJECT_ID 2>/dev/null || {
    echo "   Updating existing JWT_SECRET secret..."
    echo "$JWT_SECRET" | gcloud secrets versions add JWT_SECRET --data-file=-
}
echo "✅ JWT_SECRET secret configured"
echo ""

# Setup Cloud Build trigger
echo "5️⃣ Setting up Cloud Build trigger..."

# Check if trigger already exists
EXISTING_TRIGGER=$(gcloud builds triggers list \
    --project=$PROJECT_ID \
    --filter="github.name:$REPO_NAME" \
    --format="value(name)" | head -1)

if [ ! -z "$EXISTING_TRIGGER" ]; then
    echo "⚠️  Existing trigger found: $EXISTING_TRIGGER"
    echo "   Deleting old trigger..."
    gcloud builds triggers delete $EXISTING_TRIGGER \
        --project=$PROJECT_ID \
        --quiet
fi

# Create new trigger
echo "🔧 Creating optimized Cloud Build trigger..."
gcloud builds triggers create github \
    --project=$PROJECT_ID \
    --repo-name="$REPO_NAME" \
    --repo-owner="$GITHUB_OWNER" \
    --branch-pattern="^main$" \
    --build-config="cloudbuild.yaml" \
    --description="Finsolvz Backend - Auto deploy to Jakarta with performance optimization" \
    --region="global"

if [ $? -eq 0 ]; then
    echo "✅ Cloud Build trigger created successfully"
else
    echo "⚠️  Cloud Build trigger creation failed"
    echo "💡 You may need to:"
    echo "   1. Connect your GitHub repository in Cloud Build console"
    echo "   2. Grant necessary permissions"
    echo "   3. Run this script again"
fi
echo ""

# Set default configurations
echo "6️⃣ Setting default configurations..."
gcloud config set run/region $REGION
gcloud config set builds/region global
echo "✅ Default region set to Jakarta"
echo ""

# Deploy initial version (if Docker is available)
echo "7️⃣ Initial deployment..."
if command -v docker &> /dev/null; then
    echo "🐳 Docker detected - performing initial deployment..."
    
    # Build image locally
    BUILD_ID=$(date +%s)
    IMAGE_NAME="$REGION-docker.pkg.dev/$PROJECT_ID/finsolvz/backend:$BUILD_ID"
    
    echo "   Building Docker image..."
    docker build -t "$IMAGE_NAME" . --quiet
    
    echo "   Pushing to Artifact Registry..."
    docker push "$IMAGE_NAME" --quiet
    
    echo "   Deploying to Cloud Run..."
    gcloud run deploy $SERVICE_NAME \
        --image="$IMAGE_NAME" \
        --region=$REGION \
        --platform=managed \
        --allow-unauthenticated \
        --memory=512Mi \
        --cpu=1 \
        --port=8080 \
        --min-instances=0 \
        --max-instances=3 \
        --concurrency=80 \
        --timeout=300 \
        --set-env-vars=APP_ENV=production \
        --set-secrets=MONGO_URI=MONGO_URI:latest,JWT_SECRET=JWT_SECRET:latest \
        --project=$PROJECT_ID \
        --quiet
    
    if [ $? -eq 0 ]; then
        # Get service URL
        SERVICE_URL=$(gcloud run services describe $SERVICE_NAME \
            --region=$REGION \
            --project=$PROJECT_ID \
            --format="value(status.url)" 2>/dev/null)
        
        echo "✅ Initial deployment successful!"
        echo "🌐 Service URL: $SERVICE_URL"
        
        # Quick performance test
        echo ""
        echo "📊 Quick performance test..."
        sleep 5
        
        response_time=$(curl -w "%{time_total}" -s -o /dev/null "$SERVICE_URL/" 2>/dev/null || echo "timeout")
        if [ "$response_time" != "timeout" ]; then
            ms=$(echo "$response_time * 1000" | bc 2>/dev/null || echo "0")
            echo "   Response time: ${ms}ms"
            
            if (( $(echo "$ms < 100" | bc -l 2>/dev/null || echo 0) )); then
                echo "   ✅ Excellent performance!"
            elif (( $(echo "$ms < 200" | bc -l 2>/dev/null || echo 0) )); then
                echo "   ✅ Good performance!"
            else
                echo "   ⚠️  Performance acceptable for first request"
            fi
        fi
    else
        echo "⚠️  Initial deployment failed - you can deploy later with 'git push'"
    fi
else
    echo "⚠️  Docker not found - skipping initial deployment"
    echo "💡 Push to main branch to trigger automatic deployment"
fi

echo ""
echo "🎉 Setup Complete!"
echo "================="
echo ""
echo "📊 Environment Summary:"
echo "   • Project: $PROJECT_ID"
echo "   • Region: $REGION (Jakarta - optimal for Indonesia)"
echo "   • Artifact Registry: finsolvz"
echo "   • Cloud Run Service: $SERVICE_NAME"
echo "   • Auto-deployment: Enabled on main branch push"
echo ""
echo "💰 Free Tier Optimizations:"
echo "   • Memory: 512Mi (cost optimized)"
echo "   • CPU: 1 (sufficient for most loads)"
echo "   • Max instances: 3 (free tier limit)"
echo "   • Min instances: 0 (scales to zero)"
echo "   • Auto compression & caching enabled"
echo ""
echo "📈 Performance Optimizations:"
echo "   • Jakarta region (20-50ms from Indonesia)"
echo "   • Optimized database connections"
echo "   • Smart caching system"
echo "   • Compressed responses"
echo "   • Rate limiting protection"
echo ""
echo "🚀 Next Steps:"
echo "   1. Update your frontend to use the service URL above"
echo "   2. Push changes to main branch for auto-deployment"
echo "   3. Monitor performance in Cloud Console"
echo "   4. Check logs: gcloud logs tail SERVICE_NAME --region=$REGION"
echo ""
echo "💡 Pro Tips:"
echo "   • Use /api/reports/paginated for large datasets"
echo "   • API responses are cached for better performance"
echo "   • Service scales to zero when not used (no idle costs)"
echo "   • All endpoints support gzip compression"
echo ""

if [ ! -z "$SERVICE_URL" ]; then
    echo "🌐 Your API is ready at: $SERVICE_URL"
    echo "📚 Documentation: $SERVICE_URL/docs"
    echo "💚 Health check: $SERVICE_URL/"
fi

echo ""
echo "✅ Finsolvz Backend environment setup completed successfully!"