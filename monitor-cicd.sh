#!/bin/bash

# Finsolvz Backend CI/CD Monitoring & Debugging Script
set -e

PROJECT_ID="finsolvz-backend-dev"
SERVICE_NAME="finsolvz-backend"
REGION="asia-southeast2"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

show_help() {
    echo "🔧 Finsolvz Backend CI/CD Monitor"
    echo "=================================="
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  status      - Show overall CI/CD status"
    echo "  builds      - Show recent builds"
    echo "  logs        - Show logs for latest build"
    echo "  service     - Show Cloud Run service status"
    echo "  test        - Test deployed service"
    echo "  trigger     - Show trigger configuration"
    echo "  secrets     - Check secrets status"
    echo "  debug       - Full debugging information"
    echo "  help        - Show this help"
    echo ""
}

check_status() {
    echo -e "${BLUE}📊 CI/CD STATUS OVERVIEW${NC}"
    echo "=========================="
    
    # Check triggers
    echo "🔄 Build Triggers:"
    gcloud builds triggers list --format="table(name,status,github.name,github.push.branch)" --filter="name:finsolvz*"
    
    echo ""
    echo "🏗️ Recent Builds:"
    gcloud builds list --limit=5 --format="table(id,status,createTime,duration,substitutions.TRIGGER_NAME)"
    
    echo ""
    echo "🚀 Cloud Run Service:"
    gcloud run services describe $SERVICE_NAME --region=$REGION --format="table(metadata.name,status.url,status.conditions[0].status,spec.template.spec.containers[0].image)" 2>/dev/null || echo "Service not found"
}

show_builds() {
    echo -e "${BLUE}🏗️ RECENT BUILDS${NC}"
    echo "=================="
    gcloud builds list --limit=10 --format="table(id,status,createTime,duration,sourceProvenance.resolvedRepoSource.commitSha.slice(0:7),substitutions.TRIGGER_NAME)"
}

show_logs() {
    echo -e "${BLUE}📝 LATEST BUILD LOGS${NC}"
    echo "===================="
    
    LATEST_BUILD=$(gcloud builds list --limit=1 --format="value(id)")
    if [ -z "$LATEST_BUILD" ]; then
        echo "No builds found"
        return
    fi
    
    echo "Build ID: $LATEST_BUILD"
    echo "Logs:"
    gcloud builds log $LATEST_BUILD --stream
}

check_service() {
    echo -e "${BLUE}🚀 CLOUD RUN SERVICE STATUS${NC}"
    echo "============================"
    
    # Service details
    if gcloud run services describe $SERVICE_NAME --region=$REGION >/dev/null 2>&1; then
        echo "✅ Service exists"
        
        SERVICE_URL=$(gcloud run services describe $SERVICE_NAME --region=$REGION --format="value(status.url)")
        echo "🌐 URL: $SERVICE_URL"
        
        # Check traffic allocation
        echo "📊 Traffic Allocation:"
        gcloud run services describe $SERVICE_NAME --region=$REGION --format="table(status.traffic[].revisionName,status.traffic[].percent)"
        
        # Check latest revision
        echo "📋 Latest Revision:"
        gcloud run revisions list --service=$SERVICE_NAME --region=$REGION --limit=1 --format="table(metadata.name,status.conditions[0].status,spec.containers[0].image.slice(0:80))"
        
    else
        echo "❌ Service not found"
    fi
}

test_service() {
    echo -e "${BLUE}🧪 SERVICE TESTING${NC}"
    echo "==================="
    
    SERVICE_URL=$(gcloud run services describe $SERVICE_NAME --region=$REGION --format="value(status.url)" 2>/dev/null)
    
    if [ -z "$SERVICE_URL" ]; then
        echo "❌ Service not found"
        return
    fi
    
    echo "🌐 Testing: $SERVICE_URL"
    echo ""
    
    # Health check
    echo "1. Health Check:"
    if curl -f -s --max-time 10 "$SERVICE_URL/" | jq -r '.message' 2>/dev/null; then
        echo -e "${GREEN}✅ Health check passed${NC}"
    else
        echo -e "${RED}❌ Health check failed${NC}"
    fi
    
    echo ""
    
    # API endpoints test
    echo "2. API Endpoints:"
    
    # Test login
    echo -n "   Login endpoint: "
    LOGIN_RESPONSE=$(curl -s --max-time 10 -X POST "$SERVICE_URL/api/login" \
        -H "Content-Type: application/json" \
        -d '{"email":"admin@finsolvz.com","password":"admin123"}' 2>/dev/null)
    
    if echo "$LOGIN_RESPONSE" | grep -q "access_token" 2>/dev/null; then
        echo -e "${GREEN}✅ Working${NC}"
        
        # Extract token and test protected endpoint
        TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.access_token' 2>/dev/null)
        if [ "$TOKEN" != "null" ] && [ ! -z "$TOKEN" ]; then
            echo -n "   Protected endpoint: "
            if curl -f -s --max-time 10 "$SERVICE_URL/api/loginUser" \
                -H "Authorization: Bearer $TOKEN" >/dev/null 2>&1; then
                echo -e "${GREEN}✅ Working${NC}"
            else
                echo -e "${RED}❌ Failed${NC}"
            fi
        fi
    else
        echo -e "${RED}❌ Failed${NC}"
        echo "Response: $LOGIN_RESPONSE"
    fi
    
    echo ""
    echo "3. Documentation:"
    echo -n "   Swagger docs: "
    if curl -f -s --max-time 10 "$SERVICE_URL/docs" >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Available${NC}"
        echo "   📖 Docs URL: $SERVICE_URL/docs"
    else
        echo -e "${RED}❌ Not available${NC}"
    fi
}

check_trigger() {
    echo -e "${BLUE}🔄 BUILD TRIGGER CONFIGURATION${NC}"
    echo "================================"
    
    # List all triggers
    gcloud builds triggers list --format="yaml" --filter="name:finsolvz*"
}

check_secrets() {
    echo -e "${BLUE}🔐 SECRETS STATUS${NC}"
    echo "=================="
    
    echo "📋 Available Secrets:"
    gcloud secrets list --format="table(name,createTime)"
    
    echo ""
    echo "🔍 Secret Access Test:"
    
    secrets=("MONGO_URI" "JWT_SECRET" "NODEMAILER_EMAIL" "NODEMAILER_PASS")
    for secret in "${secrets[@]}"; do
        echo -n "   $secret: "
        if gcloud secrets versions access latest --secret="$secret" >/dev/null 2>&1; then
            echo -e "${GREEN}✅ Accessible${NC}"
        else
            echo -e "${RED}❌ Not accessible${NC}"
        fi
    done
}

debug_full() {
    echo -e "${BLUE}🔍 FULL DEBUG INFORMATION${NC}"
    echo "=========================="
    
    check_status
    echo ""
    check_service
    echo ""
    check_secrets
    echo ""
    check_trigger
    echo ""
    
    echo -e "${BLUE}📋 TROUBLESHOOTING CHECKLIST${NC}"
    echo "============================="
    echo "□ GitHub repository connected to Cloud Build"
    echo "□ Build trigger configured for main branch"
    echo "□ All required secrets created and accessible"
    echo "□ Cloud Build service account has proper permissions"
    echo "□ Artifact Registry repository exists"
    echo "□ cloudbuild.yaml file exists in repository root"
    echo "□ Latest build completed successfully"
    echo "□ Cloud Run service is deployed and healthy"
    echo ""
    
    echo -e "${BLUE}🔗 USEFUL LINKS${NC}"
    echo "==============="
    echo "Cloud Build Console: https://console.cloud.google.com/cloud-build/builds?project=$PROJECT_ID"
    echo "Cloud Run Console: https://console.cloud.google.com/run?project=$PROJECT_ID"
    echo "Secret Manager Console: https://console.cloud.google.com/security/secret-manager?project=$PROJECT_ID"
    
    if [ ! -z "$(gcloud run services describe $SERVICE_NAME --region=$REGION --format='value(status.url)' 2>/dev/null)" ]; then
        SERVICE_URL=$(gcloud run services describe $SERVICE_NAME --region=$REGION --format="value(status.url)")
        echo "Service URL: $SERVICE_URL"
        echo "API Docs: $SERVICE_URL/docs"
    fi
}

# Main script logic
case "${1:-help}" in
    "status")
        check_status
        ;;
    "builds")
        show_builds
        ;;
    "logs")
        show_logs
        ;;
    "service")
        check_service
        ;;
    "test")
        test_service
        ;;
    "trigger")
        check_trigger
        ;;
    "secrets")
        check_secrets
        ;;
    "debug")
        debug_full
        ;;
    "help")
        show_help
        ;;
    *)
        echo "Unknown command: $1"
        show_help
        exit 1
        ;;
esac