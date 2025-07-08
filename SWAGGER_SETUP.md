# Swagger Documentation Setup - Fixed ✅

## Problem Analysis
The Swagger documentation was inaccessible in deployment due to several configuration issues:

1. **Missing OpenAPI file in container** - The `api/openapi.yaml` file was not copied to the Docker container
2. **Port mismatch** - Container exposed port 8080 but app ran on port 8787
3. **Cloud Build port configuration** - Deployment used wrong port
4. **Documentation inconsistencies** - OpenAPI spec didn't match actual API implementation

## Fixes Applied

### 1. Docker Configuration (`Dockerfile`)
```dockerfile
# Copy the OpenAPI specification file
COPY --from=builder /app/api ./api

# Expose port (default 8787, configurable via PORT env var)
EXPOSE 8787

# Health check updated to use correct port
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider --timeout=10 http://localhost:8787/ || exit 1
```

### 2. Cloud Build Configuration (`cloudbuild.yaml`)
```yaml
'--port', '8787',
'--set-env-vars', 'APP_ENV=production,PORT=8787',
```

### 3. OpenAPI Documentation Updates (`api/openapi.yaml`)
- Fixed server URLs to reflect actual deployment
- Updated company endpoints to use single `/api/company/{idOrName}` route
- Updated report type endpoints to use single `/api/reportTypes/{idOrName}` route
- Added documentation for intelligent ID/name detection
- Added Swagger UI access information

### 4. Testing Script (`test-swagger.sh`)
Created automated testing script to verify:
- Local server running on correct port
- OpenAPI file exists
- Deployment configuration
- Access URLs

## API Endpoints Overview

### Authentication
- `POST /api/login` - User authentication
- `POST /api/forgot-password` - Password reset request
- `POST /api/reset-password` - Password reset with token

### User Management
- `POST /api/register` - Register new user (SUPER_ADMIN)
- `GET /api/users` - Get all users (ADMIN+)
- `GET /api/users/{id}` - Get user by ID
- `GET /api/loginUser` - Get current user
- `PUT /api/users/{id}` - Update user (SUPER_ADMIN)
- `DELETE /api/users/{id}` - Delete user (SUPER_ADMIN)
- `PUT /api/updateRole` - Update user role (SUPER_ADMIN)
- `PATCH /api/change-password` - Change password

### Company Management
- `GET /api/company` - Get all companies
- `POST /api/company` - Create new company
- `GET /api/company/{idOrName}` - Get company by ID or name (smart routing)
- `PUT /api/company/{id}` - Update company (SUPER_ADMIN)
- `DELETE /api/company/{id}` - Delete company (SUPER_ADMIN)
- `GET /api/user/companies` - Get current user's companies

### Report Types
- `GET /api/reportTypes` - Get all report types
- `POST /api/reportTypes` - Create new report type
- `GET /api/reportTypes/{idOrName}` - Get report type by ID or name (smart routing)
- `PUT /api/reportTypes/{id}` - Update report type
- `DELETE /api/reportTypes/{id}` - Delete report type

### Reports
- `GET /api/reports` - Get all reports
- `POST /api/reports` - Create new report
- `GET /api/reports/{id}` - Get report by ID
- `PUT /api/reports/{id}` - Update report
- `DELETE /api/reports/{id}` - Delete report
- `GET /api/reports/name/{name}` - Get report by name
- `GET /api/reports/company/{companyId}` - Get reports by company
- `POST /api/reports/companies` - Get reports by multiple companies
- `GET /api/reports/reportType/{reportType}` - Get reports by type
- `GET /api/reports/userAccess/{id}` - Get reports by user access
- `GET /api/reports/createdBy/{id}` - Get reports by creator

### Documentation & Health
- `GET /` - Health check
- `GET /docs` - Swagger UI interface
- `GET /api/openapi.yaml` - OpenAPI specification
- `GET /debug/files` - Debug endpoint

## How to Access Swagger Documentation

### Local Development
1. Start the server: `go run cmd/server/main.go`
2. Access Swagger UI: `http://localhost:8787/docs`
3. View OpenAPI spec: `http://localhost:8787/api/openapi.yaml`

### Production Deployment
1. Deploy using: `./deploy.sh`
2. Access Swagger UI: `https://[your-cloud-run-url]/docs`
3. View OpenAPI spec: `https://[your-cloud-run-url]/api/openapi.yaml`

### Testing
Run the test script: `./test-swagger.sh`

## Key Features

### Smart Routing
- Company and Report Type endpoints support both ID and name parameters
- Automatic detection: 24-character hex strings are treated as ObjectIDs, others as names
- Single endpoint handles both cases for cleaner API design

### Authentication & Authorization
- JWT-based authentication
- Role-based access control (SUPER_ADMIN, ADMIN, CLIENT)
- Middleware-based security enforcement

### Data Population
- Reports include populated company, user, and report type data
- Comprehensive filtering and querying options
- Consistent error handling and validation

## Security Notes
- All endpoints except authentication require valid JWT tokens
- Role-based permissions enforced at controller level
- CORS properly configured for cross-origin requests
- Input validation and sanitization implemented

## Next Steps
1. Deploy the fixed configuration
2. Test Swagger access in production
3. Verify all endpoints work correctly
4. Monitor API usage and performance

The Swagger documentation should now be fully accessible both locally and in production! 🚀