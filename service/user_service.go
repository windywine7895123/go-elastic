package service

import (
	"context"
	"go-elastic/models"
	"go-elastic/repository"

	"go.opentelemetry.io/otel"
)

const tracerName = "user-service"

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *models.User) error {
	tr := otel.Tracer(tracerName)
	ctx, span := tr.Start(ctx, "CreateUser")
	defer span.End()

	return s.repo.Create(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	tr := otel.Tracer(tracerName)
	ctx, span := tr.Start(ctx, "GetUserByID")
	defer span.End()

	return s.repo.FindByID(ctx, id)
}

func (s *userService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	tr := otel.Tracer(tracerName)
	ctx, span := tr.Start(ctx, "GetAllUsers")
	defer span.End()

	return s.repo.FindAll(ctx)
}
