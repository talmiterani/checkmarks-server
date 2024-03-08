package repo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentsRepo interface {
	GetComments(ctx context.Context, postId string) ([]primitive.M, error)
}
