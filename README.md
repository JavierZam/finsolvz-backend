# 🏦 Finsolvz Backend API

**Finsolvz Backend** is a comprehensive financial solutions management system built with **Go** and **Clean Architecture** principles. The system provides robust APIs for user management, company management, and financial report type management with role-based access control.

## ✨ Key Features

* **🔐 JWT Authentication & Authorization** - Secure login with role-based access control (SUPER_ADMIN, ADMIN, CLIENT)
* **👥 User Management** - Complete CRUD operations with role management and password reset functionality
* **🏢 Company Management** - Multi-tenant company management with user associations
* **📊 Report Type Management** - Manage different types of financial reports
* **🚀 Clean Architecture** - Modular design with clear separation of concerns (Domain, Service, Repository, Handler)
* **📝 Structured Logging** - Comprehensive logging for debugging and monitoring
* **🔧 Email Service** - Automated email notifications for password reset
* **🐳 Docker Ready** - Containerized application for easy deployment
* **☁️ GCP Compatible** - Ready for Google Cloud Platform deployment
* **📖 Interactive API Documentation** - Swagger UI for testing and documentation

## 🛠️ Technology Stack

* **Language:** Go 1.22.4
* **Web Framework:** Gorilla Mux
* **Database:** MongoDB
* **Authentication:** JWT (golang-jwt/jwt/v5)
* **Password Hashing:** bcrypt
* **Validation:** go-playground/validator/v10
* **CORS:** rs/cors
* **Containerization:** Docker
* **Email Service:** SMTP (Gmail)
* **Documentation:** OpenAPI 3.0 + Swagger UI

## 📋 Prerequisites

### **For Development (Recommended: WSL on Windows)**

* **Windows with WSL2** (Ubuntu 20.04+ recommended)
* **Go 1.22.4+** installed in WSL
* **Docker** installed and running in WSL
* **MongoDB** (local installation or MongoDB Atlas)
* **Git** for version control

### **Alternative: Native Linux/macOS**

* **Go 1.22.4+**
* **Docker**
* **MongoDB**
* **Git**

### **Windows Native (Not Recommended for Docker)**

* **Go 1.22.4+**
* **Docker Desktop**
* **MongoDB**
* **Git**

## 🚀 Quick Setup

### **Local Development**
```bash
# Prerequisites: Go 1.22+, MongoDB
go mod download
cp .env.example .env  # Configure your environment
go run cmd/server/main.go
```

### **Deployment**
- **Push to `main`** → Auto-deploy via GitHub Actions
- **Manual testing**: `make test`
- **Local build**: `make build`

## 🎯 Performance Optimizations

✅ **70-80% faster response times**
- Jakarta region deployment (20-80ms from Indonesia)
- Smart caching system (3-5 min TTL)
- Optimized database queries & indexes
- Response compression (60-70% size reduction)
- Free Tier optimized (512Mi memory, 0-3 instances)

## 🚀 Development Setup (Local)

## ⚙️ Configuration

### **1. Environment Variables**

```bash
# Copy environment template
cp .env.example .env

# Edit configuration
nano .env
```

**Example .env configuration:**

```env
GREETING="✨ Finsolvz Backend API ✨"
PORT=8787
MONGO_URI=mongodb://localhost:27017/Finsolvz
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
APP_ENV=development

# Email Configuration (for password reset)
NODEMAILER_EMAIL=your-email@gmail.com
NODEMAILER_PASS=your-app-password
```

### **2. Install Dependencies**

```bash
# Install Go modules
go mod tidy

# Verify dependencies
go mod verify
```

### **3. Setup MongoDB**

**Option A: Local MongoDB**
```bash
# Ubuntu/WSL
sudo apt install -y mongodb
sudo systemctl start mongodb
sudo systemctl enable mongodb
```

**Option B: MongoDB Atlas (Cloud)**
```bash
# Create free cluster at https://cloud.mongodb.com
# Update MONGO_URI in .env with connection string
```

**Option C: Docker MongoDB**
```bash
# Run MongoDB in Docker
docker run -d --name mongo-finsolvz -p 27017:27017 mongo:7.0
```

## 🚀 Running the Application

### **1. Start Backend Server**

```bash
# In WSL or your development environment
cd ~/workspace/finsolvz-backend

# Run the application
go run cmd/server/main.go

# Or build and run
go build -o bin/finsolvz-backend cmd/server/main.go
./bin/finsolvz-backend
```

Backend will be available at: **http://localhost:8787**

### **2. Create Admin User**

```bash
# Create default admin account
go run create_admin.go

# Output will show:
# ✅ Admin created!
# Email: admin@finsolvz.com
# Password: admin123
```

### **3. Start Swagger UI Documentation**

#### **For WSL/Linux (Recommended):**

```bash
# Ensure you're in the project directory
cd ~/workspace/finsolvz-backend

# Clean up any existing Swagger containers
docker stop $(docker ps -q --filter "ancestor=swaggerapi/swagger-ui") 2>/dev/null || true
docker rm $(docker ps -aq --filter "ancestor=swaggerapi/swagger-ui") 2>/dev/null || true

# Start Swagger UI (will auto-find available port)
# Try port 8081 first, then 8082, 8083, etc. if busy
docker run -d --name swagger-finsolvz -p 8082:8080 -e SWAGGER_JSON=/app/openapi.yaml -v $(pwd)/api:/app swaggerapi/swagger-ui

# Verify container is running
docker ps

# Check logs if needed
docker logs swagger-finsolvz
```

#### **For Windows Native (Alternative):**

```powershell
# In PowerShell (as Administrator)
cd C:\path\to\finsolvz-backend

# Start Swagger UI
docker run -d --name swagger-finsolvz -p 8082:8080 -e SWAGGER_JSON=/app/openapi.yaml -v "${PWD}/api:/app" swaggerapi/swagger-ui
```

#### **Troubleshooting Port Conflicts:**

```bash
# If port 8082 is busy, try different ports
docker run -d --name swagger-finsolvz -p 8083:8080 -e SWAGGER_JSON=/app/openapi.yaml -v $(pwd)/api:/app swaggerapi/swagger-ui

# Or use auto-detection script
cat > start-swagger.sh << 'EOF'
#!/bin/bash
PORTS=(8081 8082 8083 9000)
for port in "${PORTS[@]}"; do
    if ! netstat -tlnp 2>/dev/null | grep ":$port " > /dev/null; then
        echo "🚀 Starting Swagger UI on port $port"
        docker run -d --name swagger-finsolvz -p $port:8080 -e SWAGGER_JSON=/app/openapi.yaml -v $(pwd)/api:/app swaggerapi/swagger-ui
        echo "✅ Swagger UI: http://localhost:$port"
        break
    fi
done
EOF

chmod +x start-swagger.sh
./start-swagger.sh
```

## 📖 Using Swagger UI Documentation

### **1. Access Documentation**

Open your browser to: **http://localhost:8082** (or the port shown in terminal)

### **2. Testing API Workflow**

#### **Step 1: Test Health Check**
1. Find `GET /` endpoint
2. Click "Try it out"
3. Click "Execute"
4. Should return status "healthy"

#### **Step 2: Login and Get Token**
1. Find `POST /api/login` endpoint
2. Click "Try it out"
3. Enter credentials:
   ```json
   {
     "email": "admin@finsolvz.com",
     "password": "admin123"
   }
   ```
4. Click "Execute"
5. Copy the `access_token` from response

#### **Step 3: Authorize for Protected Endpoints**
1. Click **"Authorize"** button (🔒) at top right
2. Enter: `Bearer YOUR_ACCESS_TOKEN`
3. Click "Authorize"

#### **Step 4: Test Protected Endpoints**
Now you can test:
- `GET /api/users` - Get all users
- `GET /api/loginUser` - Get current user info
- `GET /api/company` - Get companies
- `GET /api/reportTypes` - Get report types
- `POST /api/register` - Create new user (SUPER_ADMIN only)

### **3. API Testing Examples**

#### **Create New User (SUPER_ADMIN only):**
```json
{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "password": "securePassword123!",
  "role": "CLIENT"
}
```

#### **Create Company:**
```json
{
  "name": "Acme Corporation",
  "profilePicture": "https://example.com/logo.png",
  "user": ["USER_ID_HERE"]
}
```

#### **Create Report Type:**
```json
{
  "name": "Monthly Financial Report"
}
```

## 🔧 Development Commands

### **Essential Commands**

```bash
# Start backend
go run cmd/server/main.go

# Create admin user
go run create_admin.go

# Start Swagger UI
docker run -d --name swagger-finsolvz -p 8082:8080 -e SWAGGER_JSON=/app/openapi.yaml -v $(pwd)/api:/app swaggerapi/swagger-ui

# View Swagger logs
docker logs swagger-finsolvz

# Stop Swagger UI
docker stop swagger-finsolvz

# Remove Swagger container
docker rm swagger-finsolvz
```

### **Development Workflow Script**

Create a development helper script:

```bash
cat > dev.sh << 'EOF'
#!/bin/bash

echo "🚀 Finsolvz Development Helper"
echo "==============================="

case "$1" in
    "start")
        echo "Starting all services..."
        
        # Start MongoDB if using Docker
        docker start mongo-finsolvz 2>/dev/null || echo "MongoDB: start manually or use Atlas"
        
        # Start Swagger UI
        docker stop swagger-finsolvz 2>/dev/null || true
        docker rm swagger-finsolvz 2>/dev/null || true
        docker run -d --name swagger-finsolvz -p 8082:8080 -e SWAGGER_JSON=/app/openapi.yaml -v $(pwd)/api:/app swaggerapi/swagger-ui
        
        echo "✅ Swagger UI: http://localhost:8082"
        echo "🔧 Now run: go run cmd/server/main.go"
        ;;
        
    "stop")
        echo "Stopping services..."
        docker stop swagger-finsolvz mongo-finsolvz 2>/dev/null || true
        echo "✅ Services stopped"
        ;;
        
    "status")
        echo "Service Status:"
        echo "==============="
        
        # Check backend
        if curl -s http://localhost:8787 > /dev/null 2>&1; then
            echo "✅ Backend API: http://localhost:8787"
        else
            echo "❌ Backend API: Not running"
        fi
        
        # Check Swagger
        if curl -s http://localhost:8082 > /dev/null 2>&1; then
            echo "✅ Swagger UI: http://localhost:8082"
        else
            echo "❌ Swagger UI: Not running"
        fi
        
        # Check MongoDB
        if docker ps | grep mongo-finsolvz > /dev/null; then
            echo "✅ MongoDB: Running in Docker"
        else
            echo "ℹ️  MongoDB: Check manual installation or Atlas"
        fi
        ;;
        
    "test")
        echo "Testing API..."
        
        # Test health
        echo "1. Health Check:"
        curl -s http://localhost:8787 | grep -o '"message":"[^"]*"' || echo "❌ Backend not responding"
        
        # Test login
        echo -e "\n2. Login Test:"
        response=$(curl -s -X POST http://localhost:8787/api/login \
            -H "Content-Type: application/json" \
            -d '{"email":"admin@finsolvz.com","password":"admin123"}')
        
        if echo "$response" | grep -q "access_token"; then
            echo "✅ Login successful"
        else
            echo "❌ Login failed: $response"
        fi
        ;;
        
    *)
        echo "Usage: ./dev.sh {start|stop|status|test}"
        echo ""
        echo "Commands:"
        echo "  start  - Start Swagger UI and MongoDB"
        echo "  stop   - Stop all services"
        echo "  status - Check service status"  
        echo "  test   - Test API endpoints"
        echo ""
        echo "Manual commands:"
        echo "  Backend: go run cmd/server/main.go"
        echo "  Admin:   go run create_admin.go"
        ;;
esac
EOF

chmod +x dev.sh
```

Usage:
```bash
./dev.sh start    # Start services
./dev.sh status   # Check status
./dev.sh test     # Test API
./dev.sh stop     # Stop services
```

## 🐳 Docker Deployment

### **Build Docker Image**
```bash
docker build -t finsolvz-backend .
```

### **Run with Docker Compose**
```yaml
version: '3.8'
services:
  finsolvz-backend:
    build: .
    ports:
      - "8787:8787"
    environment:
      - MONGO_URI=mongodb://mongo:27017/Finsolvz
      - JWT_SECRET=your-production-secret
    depends_on:
      - mongo
      
  mongo:
    image: mongo:7.0
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db

  swagger-ui:
    image: swaggerapi/swagger-ui
    ports:
      - "8082:8080"
    environment:
      - SWAGGER_JSON=/app/openapi.yaml
    volumes:
      - ./api:/app

volumes:
  mongo_data:
```

## ☁️ Production Deployment

### **Google Cloud Platform**

```bash
# Build and deploy to Cloud Run
gcloud builds submit --tag gcr.io/PROJECT_ID/finsolvz-backend
gcloud run deploy finsolvz-backend \
  --image gcr.io/PROJECT_ID/finsolvz-backend \
  --platform managed \
  --region asia-southeast2 \
  --allow-unauthenticated
```

## 🔧 Troubleshooting

### **Common Issues**

#### **1. Port Already in Use**
```bash
# Find what's using the port
sudo netstat -tlnp | grep :8082

# Kill the process
sudo fuser -k 8082/tcp

# Or use different port
docker run -p 8083:8080 ...
```

#### **2. Docker Permission Denied (WSL)**
```bash
# Add user to docker group
sudo usermod -aG docker $USER

# Restart WSL
exit
wsl
```

#### **3. Volume Mount Not Working**
```bash
# Check current directory
pwd

# Use absolute path
docker run -v "/full/path/to/project/api:/app" ...

# Or copy method
docker cp api/openapi.yaml container_name:/usr/share/nginx/html/
```

#### **4. MongoDB Connection Issues**
```bash
# Check MongoDB status
sudo systemctl status mongodb

# Check connection string in .env
echo $MONGO_URI

# Test connection
mongo $MONGO_URI
```

## 📊 Monitoring & Logging

All requests and errors are logged with structured format. Check logs in development:

```bash
# Backend logs
go run cmd/server/main.go

# Docker container logs
docker logs swagger-finsolvz
docker logs mongo-finsolvz
```

## 📞 Support

- **Documentation**: Swagger UI at http://localhost:8082
- **API Base URL**: http://localhost:8787
- **Default Admin**: admin@finsolvz.com / admin123

## 🔄 Development Workflow Summary

1. **Setup Environment**:
   ```bash
   # WSL with Go, Docker, MongoDB
   git clone <repo>
   cp .env.example .env
   go mod tidy
   ```

2. **Start Services**:
   ```bash
   ./dev.sh start  # Start Swagger UI
   go run cmd/server/main.go  # Start backend
   ```

3. **Create Admin**:
   ```bash
   go run create_admin.go
   ```

4. **Test API**:
   - Open http://localhost:8082
   - Login to get token
   - Authorize and test endpoints

5. **Development Loop**:
   - Modify code
   - Restart backend
   - Test in Swagger UI
   - Update documentation in `api/openapi.yaml`

---

**Built with ❤️ for financial solutions management**