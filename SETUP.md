# 🚀 Finsolvz Backend - GCP Setup Guide

Complete setup guide for deploying Finsolvz Backend with optimal performance in Google Cloud Platform.

## 📋 Prerequisites

- Google Cloud Platform account
- GitHub repository with this code
- `gcloud` CLI installed and authenticated
- Docker installed (optional, for immediate deployment)

## 🎯 One-Command Setup

```bash
./setup-gcp-environment.sh YOUR_PROJECT_ID YOUR_GITHUB_USERNAME
```

**Example:**
```bash
./setup-gcp-environment.sh finsolvz-backend-dev johndoe
```

This script will automatically:
- ✅ Configure GCP project and APIs
- ✅ Setup Artifact Registry in Jakarta (optimal for Indonesia)
- ✅ Configure secrets (MongoDB URI, JWT Secret)
- ✅ Create Cloud Build trigger for auto-deployment
- ✅ Deploy initial version to Cloud Run
- ✅ Configure Free Tier optimizations

## 📊 Performance Testing

After setup, test your deployment:

```bash
./performance-test.sh https://your-service-url.a.run.app
```

## 🌏 Regional Optimization

- **Region**: `asia-southeast2` (Jakarta)
- **Expected Latency**: 20-80ms from Indonesia
- **Performance**: 70-80% faster than US regions

## 💰 Free Tier Configuration

- **Memory**: 512Mi (cost optimized)
- **CPU**: 1 (sufficient for most loads)
- **Max instances**: 3 (Free Tier limit)
- **Min instances**: 0 (scales to zero - no idle costs)

## 🚀 Deployment

### Automatic Deployment
Push to main branch triggers automatic deployment:
```bash
git push origin main
```

### Manual Deployment
```bash
gcloud run deploy finsolvz-backend \
  --source . \
  --region asia-southeast2 \
  --allow-unauthenticated
```

## 📊 Performance Features

✅ **Database Optimizations**
- Connection pooling (50 connections)
- MongoDB indexes for all collections
- Optimized aggregation pipelines

✅ **Caching System**
- In-memory caching (3-5 minute TTL)
- Report caching for faster repeated requests
- Company data caching

✅ **Response Optimization**
- Gzip compression (60-70% size reduction)
- Pagination for large datasets
- Rate limiting (100 req/min)

✅ **Infrastructure**
- Jakarta region deployment
- Free Tier optimized settings
- Auto-scaling (0-3 instances)

## 🔧 Environment Variables

Required secrets (automatically configured by setup script):
- `MONGO_URI`: MongoDB connection string
- `JWT_SECRET`: JWT signing secret

## 📚 API Documentation

After deployment, access:
- **API Docs**: `https://your-service-url/docs`
- **Health Check**: `https://your-service-url/`
- **OpenAPI Spec**: `https://your-service-url/api/openapi.yaml`

## 🎯 Performance Targets

| Endpoint | Target | Optimized |
|----------|--------|-----------|
| Health Check | <50ms | ✅ |
| Companies | <60ms | ✅ |
| Reports (Paginated) | <80ms | ✅ |
| Individual Report | <50ms | ✅ |

## 🔍 Monitoring

Monitor your deployment:
```bash
# View logs
gcloud logs tail finsolvz-backend --region=asia-southeast2

# Check service status
gcloud run services describe finsolvz-backend --region=asia-southeast2

# Performance metrics
./performance-test.sh https://your-service-url.a.run.app
```

## 💡 Troubleshooting

### Common Issues:

1. **Cloud Build fails**: Check quota limits in your region
2. **Secrets not found**: Ensure MongoDB URI and JWT Secret are set
3. **Performance issues**: Run performance test to identify bottlenecks
4. **Free Tier limits**: Max 3 instances, adjust in cloudbuild.yaml if needed

### Support:
- Check logs: `gcloud logs tail SERVICE_NAME`
- Monitor metrics in Cloud Console
- Run performance tests regularly

## ✅ Cleanup

To remove all resources:
```bash
gcloud run services delete finsolvz-backend --region=asia-southeast2
gcloud artifacts repositories delete finsolvz --location=asia-southeast2
gcloud builds triggers delete [TRIGGER_NAME]
```

---

**🎉 Your Finsolvz Backend is now optimized for production with excellent performance in Indonesia!**