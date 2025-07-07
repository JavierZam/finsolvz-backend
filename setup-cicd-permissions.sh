#!/bin/bash

echo "🔄 Setting up CI/CD Pipeline for Finsolvz Backend"
echo "================================================"

PROJECT_ID="finsolvz-backend-dev"
REGION="asia-southeast2"
SERVICE_NAME="finsolvz-backend"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}📋 PROJECT INFO${NC}"
echo "================"
echo "Project ID: $PROJECT_ID"
echo "Region: $REGION" 
echo "Service: $SERVICE_NAME"
echo ""

# 1. Enable required APIs
echo -e "${BLUE}🔧 ENABLING REQUIRED APIS${NC}"
echo "============================"

echo "Enabling Cloud Build API..."
gcloud services enable cloudbuild.googleapis.com

echo "Enabling Cloud Run API..."
gcloud services enable run.googleapis.com

echo "Enabling Container Registry API..."
gcloud services enable containerregistry.googleapis.com

echo "Enabling Secret Manager API..."
gcloud services enable secretmanager.googleapis.com

echo -e "${GREEN}✅ APIs enabled${NC}"
echo ""

# 2. Get project details
echo -e "${BLUE}🔍 GETTING PROJECT DETAILS${NC}"
echo "============================"

PROJECT_NUMBER=$(gcloud projects describe $PROJECT_ID --format="value(projectNumber)")
CLOUD_BUILD_SA="${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com"

echo "Project Number: $PROJECT_NUMBER"
echo "Cloud Build Service Account: $CLOUD_BUILD_SA"
echo ""

# 3. Grant Cloud Build permissions
echo -e "${BLUE}🔐 SETTING UP PERMISSIONS${NC}"
echo "============================"

echo "Granting Cloud Run Admin role..."
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:$CLOUD_BUILD_SA" \
  --role="roles/run.admin"

echo "Granting Service Account User role..."
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:$CLOUD_BUILD_SA" \
  --role="roles/iam.serviceAccountUser"

echo "Granting Secret Manager Secret Accessor role..."
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:$CLOUD_BUILD_SA" \
  --role="roles/secretmanager.secretAccessor"

echo "Granting Storage Admin role (for Container Registry)..."
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:$CLOUD_BUILD_SA" \
  --role="roles/storage.admin"

echo "Granting Cloud Build Service Account role..."
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:$CLOUD_BUILD_SA" \
  --role="roles/cloudbuild.builds.builder"

echo -e "${GREEN}✅ Permissions configured${NC}"
echo ""

# 4. Verify current secrets
echo -e "${BLUE}🔑 VERIFYING SECRETS${NC}"
echo "===================="

echo "Current secrets in project:"
gcloud secrets list --format="table(name,createTime)"

echo ""
echo "Testing secret access..."
if gcloud secrets versions access latest --secret="MONGO_URI" >/dev/null 2>&1; then
    echo -e "${GREEN}✅ MONGO_URI accessible${NC}"
else
    echo -e "${YELLOW}⚠️ MONGO_URI not accessible${NC}"
fi

if gcloud secrets versions access latest --secret="JWT_SECRET" >/dev/null 2>&1; then
    echo -e "${GREEN}✅ JWT_SECRET accessible${NC}"  
else
    echo -e "${YELLOW}⚠️ JWT_SECRET not accessible${NC}"
fi

echo ""

# 5. Summary
echo -e "${BLUE}📋 SETUP SUMMARY${NC}"
echo "================="
echo "✅ Required APIs enabled"
echo "✅ Cloud Build service account configured"
echo "✅ IAM permissions granted"
echo "✅ Secrets verified"
echo ""

echo -e "${BLUE}🚀 NEXT STEPS${NC}"
echo "=============="
echo "1. Push your code to GitHub repository"
echo "2. Connect GitHub to Cloud Build"
echo "3. Create build trigger"
echo "4. Test automated deployment"
echo ""

echo "Ready for GitHub setup!"
echo "======================="