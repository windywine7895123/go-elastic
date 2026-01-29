package repository

import (
	"context"
	"go-elastic/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	client *mongo.Client
	db     *mongo.Database
	repo   UserRepository
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:password@localhost:27017"))
	s.Require().NoError(err)

	s.client = client
	s.db = client.Database("test_db")
}

func (s *UserRepositoryTestSuite) SetupTest() {
	// Clear collection before each test
	s.db.Collection("users").DeleteMany(context.Background(), map[string]interface{}{})
	s.repo = NewUserRepository(s.db.Collection("users"))
}

func (s *UserRepositoryTestSuite) TearDownSuite() {
	s.client.Disconnect(context.Background())
}

func (s *UserRepositoryTestSuite) TestCreateSuccess() {
	user := &models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.repo.Create(context.Background(), user)

	assert.NoError(s.T(), err)
	assert.NotEqual(s.T(), primitive.NilObjectID, user.ID)
}

func (s *UserRepositoryTestSuite) TestFindByIDSuccess() {
	// Create a user first
	user := &models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.repo.Create(context.Background(), user)
	s.Require().NoError(err)

	// Find the user
	foundUser, err := s.repo.FindByID(context.Background(), user.ID.Hex())

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), foundUser)
	assert.Equal(s.T(), user.Username, foundUser.Username)
	assert.Equal(s.T(), user.Email, foundUser.Email)
}

func (s *UserRepositoryTestSuite) TestFindAllSuccess() {
	// Create multiple users
	users := []*models.User{
		{
			Username:  "user1",
			Email:     "user1@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Username:  "user2",
			Email:     "user2@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, u := range users {
		err := s.repo.Create(context.Background(), u)
		s.Require().NoError(err)
	}

	// Find all users
	foundUsers, err := s.repo.FindAll(context.Background())

	assert.NoError(s.T(), err)
	assert.Len(s.T(), foundUsers, 2)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
