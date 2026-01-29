package service

import (
	"context"
	"go-elastic/models"
	"go-elastic/repository"

	"go.opentelemetry.io/otel"
)

const bookTracerName = "book-service"

type BookService interface {
	CreateBook(ctx context.Context, book *models.Book) error
	GetBookByID(ctx context.Context, id string) (*models.Book, error)
	GetAllBooks(ctx context.Context) ([]models.Book, error)
	SearchBooks(ctx context.Context, searchType string, query string) ([]models.Book, error)
}

type bookService struct {
	repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) BookService {
	return &bookService{
		repo: repo,
	}
}

func (s *bookService) CreateBook(ctx context.Context, book *models.Book) error {
	tr := otel.Tracer(bookTracerName)
	ctx, span := tr.Start(ctx, "CreateBook")
	defer span.End()

	return s.repo.Create(ctx, book)
}

func (s *bookService) GetBookByID(ctx context.Context, id string) (*models.Book, error) {
	tr := otel.Tracer(bookTracerName)
	ctx, span := tr.Start(ctx, "GetBookByID")
	defer span.End()

	return s.repo.FindByID(ctx, id)
}

func (s *bookService) GetAllBooks(ctx context.Context) ([]models.Book, error) {
	tr := otel.Tracer(bookTracerName)
	ctx, span := tr.Start(ctx, "GetAllBooks")
	defer span.End()

	return s.repo.FindAll(ctx)
}

func (s *bookService) SearchBooks(ctx context.Context, searchType string, query string) ([]models.Book, error) {
	tr := otel.Tracer(bookTracerName)
	ctx, span := tr.Start(ctx, "SearchBooks")
	defer span.End()

	switch searchType {
	case "title":
		return s.repo.SearchByTitle(ctx, query)
	case "author":
		return s.repo.SearchByAuthor(ctx, query)
	default:
		return s.repo.FindAll(ctx)
	}
}
