package repo

import (
	commonModels "checkmarks/internal/api/common/models"
	"checkmarks/internal/api/posts/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostsRepo interface {
	Search(ctx context.Context, req *commonModels.SearchReq) ([]models.SearchPosts, int, error)
	Get(ctx context.Context, postId string) (*bson.M, error)
	Add(ctx context.Context, post *models.Post) (*primitive.ObjectID, error)
	Update(ctx context.Context, post *models.Post) error
	Delete(ctx context.Context, postId string) error
}
