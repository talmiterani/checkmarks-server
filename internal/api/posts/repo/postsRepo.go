package repo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostsRepo interface {
	GetMovies(ctx context.Context) ([]primitive.M, error)
}
