package main

import (
	"go-elastic/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func SetupRoutes(app *fiber.App, logger *logrus.Logger, userHandler *handler.UserHandler, bookHandler *handler.BookHandler) {

	// logger test
	app.Get("/hello", func(c *fiber.Ctx) error {
		logger.WithField("request_id", c.Locals("request_id")).Info("hello_called")
		return c.JSON(fiber.Map{"message": "hello"})
	})

	// panic use case
	app.Get("/panic", func(c *fiber.Ctx) error {
		panic("lab panic test ")
	})
	// tracing with jaeger use case
	app.Get("/v2/hello", handler.HelloHandler)

	// User Routes (Three-tier pattern)
	api := app.Group("/api")
	users := api.Group("/users")
	users.Post("/", userHandler.CreateUser)
	users.Get("/", userHandler.GetAllUsers)
	users.Get("/:id", userHandler.GetUser)

	// Book Routes (Three-tier pattern)
	books := api.Group("/books")
	books.Post("/", bookHandler.CreateBook)
	books.Get("/", bookHandler.GetAllBooks)
	books.Get("/search", bookHandler.SearchBooks)
	books.Get("/:id", bookHandler.GetBook)
}
