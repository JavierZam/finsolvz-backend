# Service Account Setup untuk GitHub Actions

## 🔐 Required Permissions

Service account untuk GitHub Actions membutuhkan roles berikut:

### Core Roles:
1. **Cloud Run Admin** (`roles/run.admin`)
   - Deploy services ke Cloud Run
   - Manage traffic allocation
   - Update service configurations

2. **Artifact Registry Writer** (`roles/artifactregistry.writer`)
   - Push Docker images ke registry
   - Upload artifacts to repositories

3. **Service Account User** (`roles/iam.serviceAccountUser`)
   - Act as other service accounts
   - Required for Cloud Run deployments

### Additional Permissions:
4. **Storage Admin** (`roles/storage.admin`) - Optional
   - Access to Cloud Storage buckets

## 🛠️ Setup Commands

### 1. Create Service Account
```bash
PROJECT_ID="finsolvz-backend-dev"
SA_NAME="github-serviceaccount"
SA_EMAIL="$SA_NAME@$PROJECT_ID.iam.gserviceaccount.com"

gcloud iam service-accounts create $SA_NAME \
    --description="Service account for GitHub Actions" \
    --display-name="GitHub Actions Service Account" \
    --project=$PROJECT_ID
```

### 2. Add Required Roles
```bash
# Cloud Run Admin
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:$SA_EMAIL" \
    --role="roles/run.admin"

# Artifact Registry Writer
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:$SA_EMAIL" \
    --role="roles/artifactregistry.writer"

# Service Account User  
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:$SA_EMAIL" \
    --role="roles/iam.serviceAccountUser"

# Storage Admin (optional)
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:$SA_EMAIL" \
    --role="roles/storage.admin"
```

### 3. Create JSON Key
```bash
gcloud iam service-accounts keys create key.json \
    --iam-account=$SA_EMAIL \
    --project=$PROJECT_ID
```

### 4. Verify Permissions
```bash
gcloud projects get-iam-policy $PROJECT_ID \
    --flatten="bindings[].members" \
    --format="table(bindings.role)" \
    --filter="bindings.members:$SA_EMAIL"
```

## 🔍 Troubleshooting

### Error: "Unauthenticated request"
**Cause**: Docker tidak ter-authenticate dengan Artifact Registry

**Solution**:
1. Pastikan `gcloud auth configure-docker` telah dijalankan
2. Verify service account memiliki `roles/artifactregistry.writer`
3. Check repository Artifact Registry sudah ada

### Error: "Permission denied"
**Cause**: Service account tidak memiliki permissions yang cukup

**Solution**:
1. Verify semua required roles sudah di-assign
2. Wait 1-2 menit untuk permission propagation
3. Re-run deployment

### Error: "Repository not found"
**Cause**: Artifact Registry repository belum dibuat

**Solution**:
```bash
gcloud artifacts repositories create finsolvz \
    --repository-format=docker \
    --location=asia-southeast2 \
    --description="Docker repository for Finsolvz Backend"
```

## 📋 Verification Checklist

- [ ] Service account created
- [ ] All 4 roles assigned (run.admin, artifactregistry.writer, iam.serviceAccountUser, storage.admin)
- [ ] JSON key downloaded
- [ ] GitHub secret `GCP_SA_KEY` updated with new key content
- [ ] Artifact Registry repository exists
- [ ] Service account can authenticate

## 🔗 Useful Commands

```bash
# List all service accounts
gcloud iam service-accounts list

# Check service account permissions
gcloud projects get-iam-policy $PROJECT_ID

# List Artifact Registry repositories
gcloud artifacts repositories list --location=asia-southeast2

# Test Docker authentication
echo "asia-southeast2-docker.pkg.dev" | docker-credential-gcloud get
```