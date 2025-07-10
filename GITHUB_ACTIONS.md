# GitHub Actions CI/CD Setup

Panduan lengkap untuk setup GitHub Actions sebagai pengganti Cloud Build.

## 🏗️ Workflows

### 1. Deploy Workflow (`.github/workflows/deploy.yml`)
- **Trigger**: Push ke branch `main`
- **Jobs**: 
  - Testing (format, lint, unit tests)
  - Build & Deploy ke Cloud Run
- **Coverage**: Upload ke Codecov

### 2. Test Workflow (`.github/workflows/test.yml`)
- **Trigger**: Pull Request ke `main` atau `develop`
- **Jobs**:
  - Matrix testing (Go 1.22, 1.23)
  - Security scan dengan Gosec
- **Artifacts**: Coverage reports

### 3. Health Check Workflow (`.github/workflows/health-check.yml`)
- **Trigger**: Scheduled (setiap 30 menit) + manual
- **Jobs**:
  - Production health check
  - Performance monitoring
- **Alerts**: Notification jika gagal

## 🔐 Setup Secrets

### Otomatis (Recommended)
```bash
# Install GitHub CLI jika belum ada
# https://cli.github.com/

# Login ke GitHub
gh auth login

# Jalankan setup script
./setup-github-secrets.sh
```

### Manual
Tambahkan secrets berikut di GitHub Repository Settings > Secrets and variables > Actions:

| Secret Name | Description | Example |
|-------------|-------------|---------|
| `GCP_PROJECT_ID` | Google Cloud Project ID | `finsolvz-project-123` |
| `GCP_SA_KEY` | Service Account JSON Key | `{"type": "service_account"...}` |

## 📋 Service Account Permissions

Buat service account di GCP dengan permissions:
- **Cloud Run Admin** - Deploy services
- **Artifact Registry Admin** - Push Docker images  
- **Storage Admin** - Access storage
- **Service Account User** - Use service accounts

## 🚀 Deployment Process

### 1. Disable Cloud Build
```bash
# List triggers
gcloud builds triggers list

# Disable trigger (ganti TRIGGER_ID)
gcloud builds triggers delete TRIGGER_ID
```

### 2. Push ke Main Branch
```bash
git add .
git commit -m "Switch to GitHub Actions CI/CD"
git push origin main
```

### 3. Monitor Deployment
- 🌐 GitHub Actions tab: `https://github.com/YOUR_USERNAME/REPO_NAME/actions`
- 📊 Logs real-time: `gh run watch`

## 🔧 Local Testing

Sebelum push, test locally:
```bash
# Format check
gofmt -l .

# Linting  
go vet ./...

# Unit tests
JWT_SECRET=test-secret go test -v ./internal/app/...

# Build test
go build ./cmd/server
```

## 📊 Monitoring

### GitHub Actions
- ✅ Build status badges
- 📈 Test coverage tracking
- 🕐 Performance monitoring
- 🚨 Automated alerts

### Commands
```bash
# List workflows
gh workflow list

# List recent runs
gh run list

# Watch current run
gh run watch

# View logs
gh run view --log

# Re-run failed jobs
gh run rerun --failed
```

## 🔄 Migration Checklist

- [ ] Setup GitHub secrets dengan `./setup-github-secrets.sh`
- [ ] Verify service account permissions
- [ ] Disable Cloud Build triggers  
- [ ] Test GitHub Actions dengan dummy commit
- [ ] Update dokumentasi dengan badges baru
- [ ] Configure branch protection rules
- [ ] Setup Codecov integration (optional)

## 🏷️ Environment Variables

Sama dengan Cloud Build, menggunakan:
- `MONGO_URI` - From Secret Manager
- `JWT_SECRET` - From Secret Manager  
- `APP_ENV=production` - Set in workflow

## 🎯 Benefits vs Cloud Build

| Feature | Cloud Build | GitHub Actions |
|---------|-------------|----------------|
| **Cost** | Pay per build minute | 2000 free minutes/month |
| **Speed** | ~3-5 minutes | ~2-4 minutes |
| **Integration** | GCP native | GitHub native |
| **Matrix Testing** | Manual setup | Built-in support |
| **Secrets** | Secret Manager | GitHub Secrets |
| **Monitoring** | Cloud Console | GitHub UI |
| **Flexibility** | Limited | Highly customizable |

## 🔗 Useful Links

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Google Cloud Run with GitHub Actions](https://cloud.google.com/run/docs/continuous-deployment-with-github-actions)
- [GitHub CLI](https://cli.github.com/)
- [Codecov](https://codecov.io/)