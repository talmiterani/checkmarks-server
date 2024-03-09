package comments

import (
	"checkmarks/internal/api/comments/models"
	"checkmarks/internal/api/comments/repo"
	"checkmarks/internal/api/common/access"
	"context"
	"time"
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

func (s *Service) getComments(ctx context.Context, postId string) ([]models.Comment, error) {

	return s.repo.GetComments(ctx, postId)

}

func (s *Service) add(ctx context.Context, comment *models.Comment) (*models.Comment, error) {

	now := time.Now()

	comment.Updated = &now

	newId, err := s.repo.Add(ctx, comment)

	if err != nil {
		return nil, err
	}

	comment.Id = newId
	return comment, nil
}

func (s *Service) update(ctx context.Context, comment *models.Comment) (*models.Comment, error) {

	now := time.Now()

	comment.Updated = &now

	return comment, s.repo.Update(ctx, comment)
}

func (s *Service) delete(ctx context.Context, commentId string) error {
	return s.repo.Delete(ctx, commentId)
}
