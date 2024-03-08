package comments

import (
	"awesomeProject/internal/api/comments/repo"
	"awesomeProject/internal/api/common/access"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	*access.DbConnections
	repo repo.CommentsRepo
}

func NewService(sdc *access.DbConnections) *Service {
	return &Service{
		sdc,
		repo.New(sdc),
	}
}

func (s *Service) getComments(ctx context.Context, postId string) ([]primitive.M, error) {

	return s.repo.GetComments(ctx, postId)

}
