package main

import (
	"context"
	"go-elastic/database"
	"go-elastic/handler"
	"go-elastic/repository"
	"go-elastic/service"
	"log"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// ---- Load .env file ----
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults")
	}

	// ---- Logger ----
	logger := NewLogger()

	// ---- Database ----
	database.InitDB()
	defer database.CloseDB()

	// ---- Elasticsearch ----
	database.InitElasticsearch()
	// Create books index
	bookMapping := `{
		"mappings": {
			"properties": {
				"title": {"type": "text"},
				"author": {"type": "text"},
				"isbn": {"type": "keyword"},
				"description": {"type": "text"},
				"publisher": {"type": "text"},
				"publish_date": {"type": "date"},
				"pages": {"type": "integer"},
				"language": {"type": "keyword"},
				"created_at": {"type": "date"},
				"updated_at": {"type": "date"}
			}
		}
	}`
	database.CreateIndexIfNotExists("books", bookMapping)

	// ---- Dependency Injection (Three-tier) ----
	userCollection := database.DB.Collection("users")
	userRepo := repository.NewUserRepository(userCollection)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	bookCollection := database.DB.Collection("books")
	bookRepo := repository.NewBookRepository(bookCollection)
	bookSvc := service.NewBookService(bookRepo)
	bookHandler := handler.NewBookHandler(bookSvc)

	// ---- Tracer ----
	shutdown := InitTracer()
	defer shutdown(context.Background())

	// ---- Fiber ----
	app := fiber.New()

	// ---- Middlewares (ORDER MATTERS) ----
	app.Use(RecoveryMiddleware(logger))
	app.Use(CORSMiddleware())
	app.Use(otelfiber.Middleware(
		otelfiber.WithServerName("fiber-logger"),
	))
	app.Use(RequestIDMiddleware())
	app.Use(LoggerMiddleware(logger))

	// ---- Routes ----
	SetupRoutes(app, logger, userHandler, bookHandler)

	logger.Info("server starting on :8080")
	logger.Fatal(app.Listen(":8080"))
}
