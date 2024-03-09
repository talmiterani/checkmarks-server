package repo

import (
	"awesomeProject/internal/api/posts/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostsRepo interface {
	GetAll(ctx context.Context) ([]models.Post, error)
	Add(ctx context.Context, post *models.Post) (*primitive.ObjectID, error)
	Update(ctx context.Context, post *models.Post) error
	Delete(ctx context.Context, postId string) error
	DeleteAll(ctx context.Context)
}
