package posts

import (
	"awesomeProject/internal/api/common/access"
	"awesomeProject/internal/api/posts/models"
	"awesomeProject/internal/api/posts/repo"
	"context"
	"time"
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

func (s *Service) getAll(ctx context.Context) ([]models.Post, error) {

	return s.repo.GetAll(ctx)

}

func (s *Service) add(ctx context.Context, post *models.Post) (*models.Post, error) {

	now := time.Now()

	post.Updated = &now

	newId, err := s.repo.Add(ctx, post)

	if err != nil {
		return nil, err
	}

	post.ID = newId
	return post, nil
}

func (s *Service) update(ctx context.Context, post *models.Post) (*models.Post, error) {

	now := time.Now()

	post.Updated = &now

	return post, s.repo.Update(ctx, post)
}

func (s *Service) delete(ctx context.Context, postId string) error {
	return s.repo.Delete(ctx, postId)
}
