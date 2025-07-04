package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"finsolvz-backend/internal/app/auth"
	"finsolvz-backend/internal/app/company"
	"finsolvz-backend/internal/app/reporttype"
	"finsolvz-backend/internal/app/user"
	"finsolvz-backend/internal/config"
	"finsolvz-backend/internal/platform/http/middleware"
	"finsolvz-backend/internal/repository"
	"finsolvz-backend/internal/utils"
	"finsolvz-backend/internal/utils/log"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Warnf(context.Background(), "No .env file found: %v", err)
	}

	ctx := context.Background()

	// Connect to MongoDB
	db, err := config.ConnectMongoDB(ctx)
	if err != nil {
		log.Fatalf(ctx, "Failed to connect to database: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserMongoRepository(db)
	reportTypeRepo := repository.NewReportTypeMongoRepository(db)
	companyRepo := repository.NewCompanyMongoRepository(db) // ✅ NEW

	// Initialize services
	emailService := utils.NewEmailService()
	authService := auth.NewService(userRepo, emailService)
	userService := user.NewService(userRepo)
	reportTypeService := reporttype.NewService(reportTypeRepo)
	companyService := company.NewService(companyRepo, userRepo) // ✅ Simplified

	// Initialize handlers
	authHandler := auth.NewHandler(authService)
	userHandler := user.NewHandler(userService, authService)
	reportTypeHandler := reporttype.NewHandler(reportTypeService)
	companyHandler := company.NewHandler(companyService) // ✅ NEW

	// Setup router
	router := mux.NewRouter()

	// Apply global middleware
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.RecoveryMiddleware)

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Register routes
	authHandler.RegisterRoutes(router)
	userHandler.RegisterRoutes(router, middleware.AuthMiddleware)
	reportTypeHandler.RegisterRoutes(router, middleware.AuthMiddleware)
	companyHandler.RegisterRoutes(router, middleware.AuthMiddleware) // ✅ NEW

	// Health check endpoint
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		greeting := os.Getenv("GREETING")
		if greeting == "" {
			greeting = "✨ Finsolvz Backend API ✨"
		}
		utils.RespondJSON(w, http.StatusOK, map[string]string{
			"message": greeting,
			"status":  "healthy",
		})
	}).Methods("GET")

	// Apply CORS middleware
	handler := c.Handler(router)

	// Server configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8787"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Infof(ctx, "🚀🚀 Server running on http://localhost:%s 🚀🚀", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(ctx, "Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info(ctx, "Shutting down server...")

	// Graceful shutdown
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatalf(ctx, "Server forced to shutdown: %v", err)
	}

	log.Info(ctx, "Server exited")
}