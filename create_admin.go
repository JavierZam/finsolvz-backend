package main

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"finsolvz-backend/internal/config"
	"finsolvz-backend/internal/domain"
	"finsolvz-backend/internal/repository"
	"finsolvz-backend/internal/utils"
)

func main() {
	godotenv.Load()

	db, err := config.ConnectMongoDB(context.Background())
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}

	userRepo := repository.NewUserMongoRepository(db)

	hashedPassword, _ := utils.HashPassword("admin123")

	admin := &domain.User{
		Name:      "Super Admin",
		Email:     "admin@finsolvz.com",
		Password:  hashedPassword,
		Role:      domain.RoleSuperAdmin,
		Company:   []primitive.ObjectID{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = userRepo.Create(context.Background(), admin)
	if err != nil {
		log.Fatal("Failed to create admin:", err)
	}

	log.Println("✅ Admin created!")
	log.Println("Email: admin@finsolvz.com")
	log.Println("Password: admin123")
}
