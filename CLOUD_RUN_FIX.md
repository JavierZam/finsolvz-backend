# Cloud Run Deployment Fix 🚀

## Issue
Cloud Run deployment was failing with error:
```
ERROR: (gcloud.run.deploy) spec.template.spec.containers[0].env: The following reserved env names were provided: PORT. These values are automatically set by the system.
```

## Root Cause
Google Cloud Run automatically sets the `PORT` environment variable and treats it as a reserved environment variable that cannot be overridden by users.

## Solution Applied

### 1. Removed PORT Environment Variable
**Before:**
```yaml
'--set-env-vars', 'APP_ENV=production,PORT=8787',
```

**After:**
```yaml
'--set-env-vars', 'APP_ENV=production',
```

### 2. Updated Port Configuration
**Cloud Run Deployment:**
- Uses port 8080 (Google Cloud Run default)
- Removed custom PORT environment variable

**Local Development:**
- Continues to use port 8787 (as defined in Go code default)

### 3. Updated Health Check
**Before:**
```dockerfile
CMD wget --no-verbose --tries=1 --spider --timeout=10 http://localhost:8787/ || exit 1
```

**After:**
```dockerfile
CMD wget --no-verbose --tries=1 --spider --timeout=10 http://localhost:${PORT:-8080}/ || exit 1
```

### 4. Files Modified
- `cloudbuild.yaml` - Removed PORT env var, set port to 8080
- `deploy.sh` - Removed PORT env var, set port to 8080
- `Dockerfile` - Updated EXPOSE and health check
- `test-swagger.sh` - Updated documentation
- `SWAGGER_SETUP.md` - Updated documentation

## How It Works Now

### Local Development
1. Go application reads PORT environment variable
2. If not set, defaults to 8787 (hardcoded in main.go)
3. Swagger accessible at `http://localhost:8787/docs`

### Cloud Run Production
1. Google Cloud Run automatically sets PORT=8080
2. Go application reads this and uses port 8080
3. Swagger accessible at `https://[your-cloud-run-url]/docs`

## Key Points
- **PORT environment variable is automatically managed by Cloud Run**
- **Local development uses 8787, Cloud Run uses 8080**
- **Application code correctly reads PORT env var with 8787 fallback**
- **Swagger documentation now works in both environments**

## Deployment Commands
```bash
# Deploy using Cloud Build
./deploy.sh

# Or deploy manually
gcloud builds submit --tag gcr.io/[PROJECT_ID]/finsolvz-backend
gcloud run deploy finsolvz-backend \
  --image gcr.io/[PROJECT_ID]/finsolvz-backend \
  --region asia-southeast2 \
  --platform managed \
  --allow-unauthenticated \
  --port 8080 \
  --set-env-vars="APP_ENV=production"
```

## Testing
```bash
# Test configuration
./test-swagger.sh

# Test deployed service
curl https://[your-cloud-run-url]/
curl https://[your-cloud-run-url]/docs
curl https://[your-cloud-run-url]/api/openapi.yaml
```

The deployment should now work correctly without the PORT environment variable error! 🎉