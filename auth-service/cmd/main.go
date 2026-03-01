package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	appUser "github.com/greeschenko/projectMimir/auth-service/internal/application/user"
	"github.com/greeschenko/projectMimir/auth-service/internal/infrastructure/hasher"
	"github.com/greeschenko/projectMimir/auth-service/internal/infrastructure/migrations"
	"github.com/greeschenko/projectMimir/auth-service/internal/infrastructure/persistence"
	"github.com/greeschenko/projectMimir/auth-service/internal/infrastructure/token"
	grpcHandler "github.com/greeschenko/projectMimir/auth-service/internal/transport/grpc/handler"

	authv1 "github.com/greeschenko/projectMimir/platform/proto/auth/v1"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"google.golang.org/grpc"
)

func main() {
	// =============================
	// Load environment variables
	// =============================
	dbHost := getEnv("DB_HOST")
	dbUser := getEnv("DB_USER")
	dbPassword := getEnv("DB_PASSWORD")
	dbName := getEnv("DB_NAME")
	dbPort := getEnv("DB_PORT")

	jwtSecret := getEnv("JWT_SECRET") // нова змінна для JWT

	// =============================
	// Database
	// =============================
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	if err := migrations.RunMigrations(databaseURL); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// =============================
	// Dependency Injection
	// =============================

	// Репозиторій користувачів
	userRepo := persistence.NewPostgresUserRepository(db)

	// PasswordHasher
	passHasher := hasher.NewBcryptHasher()

	// TokenService
	tokenSvc := token.NewJWTService(jwtSecret, time.Minute*15, time.Hour*24*7)

	// Use-cases
	registerUC := appUser.NewRegisterUseCase(userRepo, passHasher, tokenSvc)
	loginUC := appUser.NewLoginUseCase(userRepo, passHasher, tokenSvc)

	// gRPC handler
	authHandler := grpcHandler.NewAuthHandler(registerUC, loginUC)

	// =============================
	// gRPC Server
	// =============================
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	authv1.RegisterAuthServiceServer(grpcServer, authHandler)

	log.Println("gRPC server started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// getEnv panics if environment variable is not set
func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("environment variable %s is not set", key)
	}
	return value
}
