package repo

import (
	"checkmarks/internal/api/comments/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentsRepo interface {
	GetByPostId(ctx context.Context, postId string) ([]models.Comment, error)
	Add(ctx context.Context, comment *models.Comment) (*primitive.ObjectID, error)
	Update(ctx context.Context, comment *models.Comment) error
	Delete(ctx context.Context, commentId string) error
}
