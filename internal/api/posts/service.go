package posts

import (
	"checkmarks/internal/api/common/access"
	commonModels "checkmarks/internal/api/common/models"
	"checkmarks/internal/api/posts/models"
	"checkmarks/internal/api/posts/repo"
	"context"
	"go.mongodb.org/mongo-driver/bson"
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

func (s *Service) search(ctx context.Context, req *commonModels.SearchReq) ([]models.SearchPostsRes, error) {
	return s.repo.Search(ctx, req)
}

//func (s *Service) search(ctx context.Context, req *commonModels.SearchReq) ([]models.Post, error) {
//	return s.repo.Search(ctx, req)
//}

func (s *Service) get(ctx context.Context, postId string) (*bson.M, error) {
	return s.repo.Get(ctx, postId)
}

func (s *Service) add(ctx context.Context, post *models.Post) (*models.Post, error) {

	now := time.Now()

	post.Updated = &now

	newId, err := s.repo.Add(ctx, post)

	if err != nil {
		return nil, err
	}

	post.Id = newId
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
