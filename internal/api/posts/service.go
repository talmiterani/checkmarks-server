package posts

import (
	"awesomeProject/internal/api/common/access"
	"awesomeProject/internal/api/posts/repo"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	*access.DbConnections
	repo repo.PostsRepo
}

func NewService(sdc *access.DbConnections) *Service {
	return &Service{
		sdc,
		repo.New(sdc),
	}
}

func (s *Service) getMovies(ctx context.Context) ([]primitive.M, error) {

	return s.repo.GetMovies(ctx)

}
