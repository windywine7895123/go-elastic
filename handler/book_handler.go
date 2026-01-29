package handler

import (
	"go-elastic/models"
	"go-elastic/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

type BookHandler struct {
	svc service.BookService
}

func NewBookHandler(svc service.BookService) *BookHandler {
	return &BookHandler{svc: svc}
}

func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
	book := new(models.Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Cannot parse JSON",
			"details": err.Error(),
			"hint":    "Ensure Content-Type header is 'application/json' and request body is valid JSON",
		})
	}

	// Validate required fields
	if book.Title == "" || book.Author == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title and Author are required"})
	}

	// Set timestamps if not provided
	if book.CreatedAt.IsZero() {
		book.CreatedAt = time.Now()
	}
	if book.UpdatedAt.IsZero() {
		book.UpdatedAt = time.Now()
	}

	if err := h.svc.CreateBook(c.UserContext(), book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(book)
}

func (h *BookHandler) GetBook(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	book, err := h.svc.GetBookByID(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found"})
	}

	return c.JSON(book)
}

func (h *BookHandler) GetAllBooks(c *fiber.Ctx) error {
	books, err := h.svc.GetAllBooks(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(books)
}

func (h *BookHandler) SearchBooks(c *fiber.Ctx) error {
	searchType := c.Query("type", "") // "title" or "author"
	searchQuery := c.Query("q", "")

	if searchType == "" || searchQuery == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing query parameters: type and q"})
	}

	books, err := h.svc.SearchBooks(c.UserContext(), searchType, searchQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(books)
}
